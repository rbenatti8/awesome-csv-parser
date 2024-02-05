package schema

type OptionsDTO struct {
	Delimiter string   `json:"delimiter"`
	Exclude   []string `json:"exclude"`
	Places    int      `json:"places"`
}

type AggregationDTO struct {
	Type    string     `json:"type"`
	Options OptionsDTO `json:"options"`
}

type SanitizationDTO struct {
	Type    string     `json:"type"`
	Options OptionsDTO `json:"options"`
}

type FormatterDTO struct {
	Type    string     `json:"type"`
	Options OptionsDTO `json:"options"`
}

type SourceDTO struct {
	Aggregation   *AggregationDTO   `json:"aggregation"`
	Formatter     *FormatterDTO     `json:"format"`
	Fields        []string          `json:"fields"`
	Sanitizations []SanitizationDTO `json:"sanitizations"`
}

type TargetFieldDTO struct {
	Source     SourceDTO `json:"source"`
	Name       string    `json:"name"`
	IsRequired bool      `json:"is_required"`
	IsUnique   bool      `json:"is_unique"`
}

type V1DTO struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	TargetFields []TargetFieldDTO `json:"target_fields"`
}
