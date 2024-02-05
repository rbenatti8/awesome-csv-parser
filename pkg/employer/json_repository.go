package employer

import (
	"awesome-csv-parser/pkg/schema"
	"encoding/json"
	"fmt"
	"io"
)

type employerDTO struct {
	ID       string `json:"employer_id"`
	SchemaID string `json:"schema_id"`
}

type JsonRepository struct {
	schemasReader   io.Reader
	employersReader io.Reader
}

func NewJsonRepository(schemasReader, employersReader io.Reader) *JsonRepository {
	return &JsonRepository{
		schemasReader:   schemasReader,
		employersReader: employersReader,
	}
}

func (v *JsonRepository) getEmployers() (employers []employerDTO, err error) {
	decoder := json.NewDecoder(v.employersReader)
	err = decoder.Decode(&employers)
	if err != nil {
		return
	}

	return
}

func (v *JsonRepository) getEmployer(id string) (*employerDTO, error) {
	employers, err := v.getEmployers()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(employers); i++ {
		if employers[i].ID == id {
			return &employers[i], nil
		}
	}

	return nil, fmt.Errorf("employer with id %s not found", id)
}

func (v *JsonRepository) getSchemas() ([]*schema.V1, error) {
	var dtos []schema.V1DTO

	decoder := json.NewDecoder(v.schemasReader)
	err := decoder.Decode(&dtos)
	if err != nil {
		return nil, err
	}

	schemas := make([]*schema.V1, 0, len(dtos))

	var s *schema.V1

	for _, dto := range dtos {
		s, err = schema.NewFromDTO(dto)
		if err != nil {
			return nil, err
		}

		schemas = append(schemas, s)
	}

	return schemas, nil
}

func (v *JsonRepository) getSchema(id string) (*schema.V1, error) {
	schemas, err := v.getSchemas()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(schemas); i++ {
		if schemas[i].ID == id {
			return schemas[i], nil
		}
	}

	return nil, ErrSchemaNotFound{ID: id}
}

func (v *JsonRepository) GetByID(id string) (*Employer, error) {
	employer, err := v.getEmployer(id)
	if err != nil {
		return nil, err
	}

	s, err := v.getSchema(employer.SchemaID)
	if err != nil {
		return nil, err
	}

	return &Employer{
		ID:     employer.ID,
		Schema: s,
	}, nil
}
