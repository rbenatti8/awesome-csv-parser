package schema

import (
	"awesome-csv-parser/pkg/schema/aggregators"
	"awesome-csv-parser/pkg/schema/formatters"
	"awesome-csv-parser/pkg/schema/sanitizers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestField_Build(t *testing.T) {
	type fields struct {
		Sanitizers   func() []sanitizers.Sanitizer
		Aggregator   func() aggregators.Aggregator
		Formatter    func() formatters.Formatter
		Name         string
		SourceFields []string
		IsRequired   bool
		Unique       bool
	}
	type args struct {
		record    []string
		headerMap HeaderMap
	}

	tests := []struct {
		fields  fields
		name    string
		want    string
		args    args
		wantErr bool
	}{
		{
			name: "Build name from single source field",
			fields: fields{
				Name:         "employee_name",
				SourceFields: []string{"name"},
				IsRequired:   true,
				Unique:       false,
				Sanitizers: func() []sanitizers.Sanitizer {
					return []sanitizers.Sanitizer{
						&sanitizers.RemoveSpecialCharacters{},
						&sanitizers.Capitalize{},
					}
				},
				Aggregator: func() aggregators.Aggregator {
					return nil
				},
				Formatter: func() formatters.Formatter {
					return nil
				},
			},
			args: args{
				record: []string{"John doe!!"},
				headerMap: HeaderMap{
					m: map[string]int{
						"name": 0,
					},
					keys: []string{"name"},
				},
			},
			want:    "John Doe",
			wantErr: false,
		},
		{
			name: "Build name from multiple source fields",
			fields: fields{
				Name:         "employee_name",
				SourceFields: []string{"first name", "last name"},
				IsRequired:   true,
				Unique:       false,
				Sanitizers: func() []sanitizers.Sanitizer {
					return []sanitizers.Sanitizer{
						&sanitizers.RemoveSpecialCharacters{},
						&sanitizers.Capitalize{},
					}
				},
				Aggregator: func() aggregators.Aggregator {
					return &aggregators.Concat{
						Delimiter: " ",
					}
				},
				Formatter: func() formatters.Formatter {
					return nil
				},
			},
			args: args{
				record: []string{"John", "doe!!"},
				headerMap: HeaderMap{
					m: map[string]int{
						"first name": 0,
						"last name":  1,
					},
					keys: []string{"first name", "last name"},
				},
			},
			want:    "John Doe",
			wantErr: false,
		},
		{
			name: "Build salary from single source field",
			fields: fields{
				Name:         "employee_salary",
				SourceFields: []string{"salary"},
				IsRequired:   true,
				Unique:       false,
				Sanitizers: func() []sanitizers.Sanitizer {
					return []sanitizers.Sanitizer{
						&sanitizers.RemoveSpecialCharacters{
							Exclude: []rune{'.'},
						},
						&sanitizers.RemoveLetters{},
						&sanitizers.TrimSpace{},
					}
				},
				Aggregator: func() aggregators.Aggregator {
					return nil
				},
				Formatter: func() formatters.Formatter {
					return &formatters.Decimal{Places: 2}
				},
			},
			args: args{
				record: []string{"R$ 1000.678"},
				headerMap: HeaderMap{
					m: map[string]int{
						"salary": 0,
					},
					keys: []string{"salary"},
				},
			},
			want:    "1000.67",
			wantErr: false,
		},
		{
			name: "Build e-mail from single source field",
			fields: fields{
				Name:         "employee_email",
				SourceFields: []string{"email"},
				IsRequired:   true,
				Unique:       false,
				Sanitizers: func() []sanitizers.Sanitizer {
					return []sanitizers.Sanitizer{
						&sanitizers.LowerCase{},
					}
				},
				Aggregator: func() aggregators.Aggregator {
					return nil
				},
				Formatter: func() formatters.Formatter {
					return &formatters.Email{}
				},
			},
			args: args{
				record: []string{"Tester_123@gmail.com"},
				headerMap: HeaderMap{
					m: map[string]int{
						"email": 0,
					},
					keys: []string{"email"},
				},
			},
			want:    "tester_123@gmail.com",
			wantErr: false,
		},
		{
			name: "Build an invalid e-mail should return an error",
			fields: fields{
				Name:         "employee_email",
				SourceFields: []string{"email"},
				IsRequired:   true,
				Unique:       false,
				Sanitizers: func() []sanitizers.Sanitizer {
					return []sanitizers.Sanitizer{
						&sanitizers.LowerCase{},
					}
				},
				Aggregator: func() aggregators.Aggregator {
					return nil
				},
				Formatter: func() formatters.Formatter {
					return &formatters.Email{}
				},
			},
			args: args{
				record: []string{"Tester_123gmail.com"},
				headerMap: HeaderMap{
					m: map[string]int{
						"email": 0,
					},
					keys: []string{"email"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt = tt

		t.Run(tt.name, func(t *testing.T) {
			f := &Field{
				Name:         tt.fields.Name,
				SourceFields: tt.fields.SourceFields,
				IsRequired:   tt.fields.IsRequired,
				IsUnique:     tt.fields.Unique,
				Sanitizers:   tt.fields.Sanitizers(),
				Aggregator:   tt.fields.Aggregator(),
				Formatter:    tt.fields.Formatter(),
			}

			got, err := f.Build(tt.args.record, tt.args.headerMap)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
