package schema

import (
	"awesome-csv-parser/pkg/schema/aggregators"
	"awesome-csv-parser/pkg/schema/formatters"
	"awesome-csv-parser/pkg/schema/sanitizers"
	"fmt"
)

type (
	HeaderMap struct {
		m    map[string]int
		keys []string
	}
	Field struct {
		Aggregator   aggregators.Aggregator
		Formatter    formatters.Formatter
		Name         string
		SourceFields []string
		Sanitizers   []sanitizers.Sanitizer
		IsRequired   bool
		IsUnique     bool
	}
)

func buildField(fieldDTO TargetFieldDTO) (*Field, error) {
	if len(fieldDTO.Source.Fields) == 0 {
		return nil, fmt.Errorf("field %s has no source fields", fieldDTO.Name)
	}

	f := Field{
		Name:         fieldDTO.Name,
		SourceFields: fieldDTO.Source.Fields,
		IsRequired:   fieldDTO.IsRequired,
		IsUnique:     fieldDTO.IsUnique,
	}

	if len(fieldDTO.Source.Fields) > 1 && fieldDTO.Source.Aggregation == nil {
		return nil, fmt.Errorf("field %s has multiple source fields but no aggregation", fieldDTO.Name)
	}

	aggregator, err := buildAggregator(fieldDTO.Source.Aggregation)
	if err != nil {
		return nil, err
	}

	f.Aggregator = aggregator

	formatter, err := buildFormatter(fieldDTO.Source.Formatter)
	if err != nil {
		return nil, err
	}

	f.Formatter = formatter

	_sanitizers, err := buildSanitizers(fieldDTO.Source.Sanitizations)
	if err != nil {
		return nil, err
	}

	f.Sanitizers = _sanitizers

	return &f, nil

}

func (f *Field) Build(record []string, headerMap HeaderMap) (string, error) {
	if len(f.SourceFields) > 1 {
		return f.buildMultiple(record, headerMap)
	}

	return f.buildSingle(record, headerMap)
}

func (f *Field) buildSingle(record []string, headerMap HeaderMap) (string, error) {
	ok, input := f.getValue(record, headerMap)

	if err := f.checkRequired(ok, input); err != nil {
		return "", err
	}

	value := f.sanitize(input)

	if f.Formatter != nil {
		return f.Formatter.Format(value)
	}

	return value, nil
}

func (f *Field) getValue(record []string, headerMap HeaderMap) (bool, string) {
	headerIndex, ok := headerMap.Get(f.SourceFields[0])

	if !ok {
		return ok, ""
	}

	return ok, record[headerIndex]
}

func (f *Field) buildMultiple(record []string, headerMap HeaderMap) (string, error) {
	sanitizedFields := make([]string, 0, len(f.SourceFields))

	for i := 0; i < len(f.SourceFields); i++ {
		headerIndex, ok := headerMap.Get(f.SourceFields[i])
		input := record[headerIndex]

		if err := f.checkRequired(ok, input); err != nil {
			return "", err
		}

		sanitizedFields = append(sanitizedFields, f.sanitize(input))
	}

	aggregatedValue, err := f.Aggregator.Aggregate(sanitizedFields)
	if err != nil {
		return "", err
	}

	if f.Formatter != nil {
		return f.Formatter.Format(aggregatedValue)
	}

	return aggregatedValue, nil
}

func (f *Field) sanitize(value string) string {
	for _, sanitizer := range f.Sanitizers {
		value = sanitizer.Sanitize(value)
	}

	return value
}

func (f *Field) checkRequired(found bool, value string) error {
	if !found && f.IsRequired {
		return &ErrFieldNotFound{Name: f.SourceFields[0]}
	}

	if found && f.IsRequired && value == "" {
		return &ErrEmptyField{Name: f.SourceFields[0]}
	}

	return nil
}

func NewHeaderMap(size int) HeaderMap {
	return HeaderMap{
		m:    make(map[string]int, size),
		keys: make([]string, 0, size),
	}
}

func (hm *HeaderMap) Keys() []string {
	return hm.keys
}

func (hm *HeaderMap) Get(key string) (int, bool) {
	index, ok := hm.m[key]
	return index, ok
}

func (hm *HeaderMap) Set(key string, index int) {
	hm.m[key] = index
	hm.keys = append(hm.keys, key)
}

func buildAggregator(aggregatorDTO *AggregationDTO) (aggregators.Aggregator, error) {
	if aggregatorDTO == nil {
		return nil, nil
	}

	return aggregators.New(aggregatorDTO.Type, aggregators.Opts{Delimiter: aggregatorDTO.Options.Delimiter})
}

func buildFormatter(formatterDTO *FormatterDTO) (formatters.Formatter, error) {
	if formatterDTO == nil {
		return nil, nil
	}

	return formatters.New(formatterDTO.Type, formatters.Opts{Places: int32(formatterDTO.Options.Places)})
}

func buildSanitizers(sanitizersDTO []SanitizationDTO) ([]sanitizers.Sanitizer, error) {
	if len(sanitizersDTO) == 0 {
		return nil, nil
	}

	sanitizerList := make([]sanitizers.Sanitizer, 0, len(sanitizersDTO))

	for _, sanitizerDTO := range sanitizersDTO {
		s, err := sanitizers.New(sanitizerDTO.Type, sanitizers.Opts{
			Exclude: sanitizerDTO.Options.Exclude,
		})

		if err != nil {
			return nil, err
		}

		sanitizerList = append(sanitizerList, s)
	}

	return sanitizerList, nil
}
