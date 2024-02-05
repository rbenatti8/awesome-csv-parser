package schema

import (
	"awesome-csv-parser/internal/concurrency"
	"errors"
	"fmt"
	"sync"
)

type V1 struct {
	ID           string
	Name         string
	TargetFields []Field
	ShardedMap   concurrency.ShardedMap[bool]
	headers      []string
	locker       sync.Mutex
}

func (v *V1) Headers() []string {
	v.locker.Lock()
	defer v.locker.Unlock()

	if v.headers != nil {
		return v.headers
	}

	headers := make([]string, 0, len(v.TargetFields))

	for _, f := range v.TargetFields {
		headers = append(headers, f.Name)
	}

	v.headers = headers

	return headers
}

func (v *V1) Build(record []string, headerMap HeaderMap) ([]string, error) {
	targetRecords := make([]string, 0, len(v.TargetFields))

	for _, f := range v.TargetFields {
		targetRecord, err := f.Build(record, headerMap)
		if err != nil {
			return nil, err
		}

		if f.IsUnique && !v.isUnique(targetRecord) {
			return nil, errors.New(fmt.Sprintf("record is not unique: %s", targetRecord))
		}

		targetRecords = append(targetRecords, targetRecord)
	}

	return targetRecords, nil
}

func NewFromDTO(dto V1DTO) (*V1, error) {
	v := &V1{
		ID:           dto.ID,
		Name:         dto.Name,
		TargetFields: make([]Field, 0, len(dto.TargetFields)),
		ShardedMap:   concurrency.NewShardedMap[bool](1024),
	}

	for _, f := range dto.TargetFields {
		field, err := buildField(f)
		if err != nil {
			return nil, err
		}

		v.TargetFields = append(v.TargetFields, *field)
	}

	return v, nil
}

func (v *V1) isUnique(targetRecord string) bool {

	_, ok := v.ShardedMap.Get(targetRecord)

	if ok {
		return false
	}

	v.ShardedMap.Set(targetRecord, true)

	return true
}
