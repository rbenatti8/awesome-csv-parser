package employer

import (
	"awesome-csv-parser/internal/concurrency"
	"awesome-csv-parser/pkg/schema"
	"awesome-csv-parser/pkg/schema/formatters"
	"awesome-csv-parser/pkg/schema/sanitizers"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestJsonRepository_GetByID(t *testing.T) {
	testCases := []struct {
		expectedErr  error
		buildReaders func() (io.Reader, io.Reader)
		expected     *Employer
		name         string
		arg          string
	}{
		{
			name: "success",
			arg:  "1",
			buildReaders: func() (io.Reader, io.Reader) {
				schemas := `[{
		 "id": "1",
		 "name": "schema for vendor A",
		 "target_fields": [
		   {
		     "name": "employee_name",
		     "is_required": true,
		     "is_unique": false,
		     "source": {
		       "fields": ["Name"],
		       "sanitizations": [
		         {
		           "type": "remove_special_characters",
		           "options": {
		             "exclude": []
		           }
		         },
		         {
		           "type": "capitalize",
		           "options": {}
		         }
		       ]
		     }
		   },
		   {
		     "name": "employee_salary",
		     "is_required": true,
		     "is_unique": false,
		     "source": {
		       "fields": ["Wage"],
		       "sanitizations": [
		         {
		           "type": "remove_special_characters",
		           "options": {
		             "exclude": ["."]
		           }
		         },
		         {
		           "type": "remove_letters",
		           "options": {}
		         }
		       ],
		       "format": {
		         "type": "decimal",
		         "options": {
		           "places": 2
		         }
		       }
		     }
		   },
		   {
		     "name": "employee_email",
		     "is_required": true,
		     "is_unique": true,
		     "source": {
		       "fields": ["Email"],
		       "sanitizations": [
		         {
		           "type": "lowercase",
		           "options": {}
		         }
		       ],
		       "format": {
		         "type": "email",
		         "options": {}
		       }
		     }
		   },
		   {
		     "name": "employee_id",
		     "is_required": true,
		     "is_unique": false,
		     "source": {
		       "fields": ["Number"],
		       "sanitizations": []
		     }
		   },
		   {
		     "name": "employee_phone",
		     "is_required": false,
		     "is_unique": false,
		     "source": {
		       "fields": [""]
		     }
		   }
		 ]
		}]`
				employers := `[{"employer_id":"1","schema_id":"1"},{"employer_id":"1","schema_id":"2"},{"employer_id":"1","schema_id":"3"},{"employer_id":"1","schema_id":"4"},{"employer_id":"5","schema_id":"5"}]`

				return strings.NewReader(schemas), strings.NewReader(employers)
			},
			expected: &Employer{
				ID: "1",
				Schema: &schema.V1{
					ID:   "1",
					Name: "schema for vendor A",
					TargetFields: []schema.Field{
						{
							Name:         "employee_name",
							SourceFields: []string{"Name"},
							IsRequired:   true,
							IsUnique:     false,
							Sanitizers:   []sanitizers.Sanitizer{&sanitizers.RemoveSpecialCharacters{}, &sanitizers.Capitalize{}},
							Aggregator:   nil,
							Formatter:    nil,
						},
						{
							Name:         "employee_salary",
							SourceFields: []string{"Wage"},
							IsRequired:   true,
							IsUnique:     false,
							Sanitizers:   []sanitizers.Sanitizer{&sanitizers.RemoveSpecialCharacters{Exclude: []rune{'.'}}, &sanitizers.RemoveLetters{}},
							Formatter:    &formatters.Decimal{Places: 2},
						},
						{
							Name:         "employee_email",
							SourceFields: []string{"Email"},
							IsRequired:   true,
							IsUnique:     true,
							Sanitizers:   []sanitizers.Sanitizer{&sanitizers.LowerCase{}},
							Formatter:    &formatters.Email{},
						},
						{
							Name:         "employee_id",
							SourceFields: []string{"Number"},
							IsRequired:   true,
							IsUnique:     false,
						},
						{
							Name:         "employee_phone",
							SourceFields: []string{""},
							IsRequired:   false,
							IsUnique:     false,
						},
					},
					ShardedMap: concurrency.NewShardedMap[bool](1024),
				},
			},
		},
		{
			name: "when schema not found",
			arg:  "1",
			buildReaders: func() (io.Reader, io.Reader) {
				schemas := `[{
		 "id": "2",
		 "name": "schema for vendor A",
		 "target_fields": [
		   {
		     "name": "employee_name",
		     "is_required": true,
		     "is_unique": false,
		     "source": {
		       "fields": ["Name"],
		       "sanitizations": [
		         {
		           "type": "remove_special_characters",
		           "options": {
		             "exclude": []
		           }
		         },
		         {
		           "type": "capitalize",
		           "options": {}
		         }
		       ]
		     }
		   },
		   {
		     "name": "employee_salary",
		     "is_required": true,
		     "is_unique": false,
		     "source": {
		       "fields": ["Wage"],
		       "sanitizations": [
		         {
		           "type": "remove_special_characters",
		           "options": {
		             "exclude": ["."]
		           }
		         },
		         {
		           "type": "remove_letters",
		           "options": {}
		         }
		       ],
		       "format": {
		         "type": "decimal",
		         "options": {
		           "places": 2
		         }
		       }
		     }
		   },
		   {
		     "name": "employee_email",
		     "is_required": true,
		     "is_unique": true,
		     "source": {
		       "fields": ["Email"],
		       "sanitizations": [
		         {
		           "type": "lowercase",
		           "options": {}
		         }
		       ],
		       "format": {
		         "type": "email",
		         "options": {}
		       }
		     }
		   },
		   {
		     "name": "employee_id",
		     "is_required": true,
		     "is_unique": false,
		     "source": {
		       "fields": ["Number"],
		       "sanitizations": []
		     }
		   },
		   {
		     "name": "employee_phone",
		     "is_required": false,
		     "is_unique": false,
		     "source": {
		       "fields": [""]
		     }
		   }
		 ]
		}]`
				employers := `[{"employer_id":"1","schema_id":"1"},{"employer_id":"1","schema_id":"2"},{"employer_id":"1","schema_id":"3"},{"employer_id":"1","schema_id":"4"},{"employer_id":"5","schema_id":"5"}]`

				return strings.NewReader(schemas), strings.NewReader(employers)
			},
			expected: nil,
			expectedErr: ErrSchemaNotFound{
				ID: "1",
			},
		},
		{
			name: "when employer not found",
			arg:  "10",
			buildReaders: func() (io.Reader, io.Reader) {
				schemas := `[{
		 "id": "1",
		 "name": "schema for vendor A",
		 "target_fields": [
		   {
		     "name": "employee_name",
		     "is_required": true,
		     "is_unique": false,
		     "source": {
		       "fields": ["Name"],
		       "sanitizations": [
		         {
		           "type": "remove_special_characters",
		           "options": {
		             "exclude": []
		           }
		         },
		         {
		           "type": "capitalize",
		           "options": {}
		         }
		       ]
		     }
		   },
		   {
		     "name": "employee_salary",
		     "is_required": true,
		     "is_unique": false,
		     "source": {
		       "fields": ["Wage"],
		       "sanitizations": [
		         {
		           "type": "remove_special_characters",
		           "options": {
		             "exclude": ["."]
		           }
		         },
		         {
		           "type": "remove_letters",
		           "options": {}
		         }
		       ],
		       "format": {
		         "type": "decimal",
		         "options": {
		           "places": 2
		         }
		       }
		     }
		   },
		   {
		     "name": "employee_email",
		     "is_required": true,
		     "is_unique": true,
		     "source": {
		       "fields": ["Email"],
		       "sanitizations": [
		         {
		           "type": "lowercase",
		           "options": {}
		         }
		       ],
		       "format": {
		         "type": "email",
		         "options": {}
		       }
		     }
		   },
		   {
		     "name": "employee_id",
		     "is_required": true,
		     "is_unique": false,
		     "source": {
		       "fields": ["Number"],
		       "sanitizations": []
		     }
		   },
		   {
		     "name": "employee_phone",
		     "is_required": false,
		     "is_unique": false,
		     "source": {
		       "fields": [""]
		     }
		   }
		 ]
		}]`
				employers := `[{"employer_id":"1","schema_id":"1"},{"employer_id":"1","schema_id":"2"},{"employer_id":"1","schema_id":"3"},{"employer_id":"1","schema_id":"4"},{"employer_id":"5","schema_id":"5"}]`

				return strings.NewReader(schemas), strings.NewReader(employers)
			},
			expected: nil,
			expectedErr: ErrEmployerNotFound{
				ID: "10",
			},
		},
		{
			name: "when error getting schemas",
			arg:  "1",
			buildReaders: func() (io.Reader, io.Reader) {
				schemas := `[`
				employers := `[{"employer_id":"1","schema_id":"1"},{"employer_id":"1","schema_id":"2"},{"employer_id":"1","schema_id":"3"},{"employer_id":"1","schema_id":"4"},{"employer_id":"5","schema_id":"5"}]`

				return strings.NewReader(schemas), strings.NewReader(employers)
			},
			expected:    nil,
			expectedErr: ErrInvalidReader{InternalErr: io.EOF},
		},
		{
			name: "when error getting employers",
			arg:  "1",
			buildReaders: func() (io.Reader, io.Reader) {
				schemas := `[{
    "id": "1",
    "name": "schema for vendor A",
    "target_fields": [
      {
        "name": "employee_name",
        "is_required": true,
        "is_unique": false,
        "source": {
          "fields": ["Name"],
          "sanitizations": [
            {
              "type": "remove_special_characters",
              "options": {
                "exclude": []
              }
            },
            {
              "type": "capitalize",
              "options": {}
            }
          ]
        }
      },
      {
        "name": "employee_salary",
        "is_required": true,
        "is_unique": false,
        "source": {
          "fields": ["Wage"],
          "sanitizations": [
            {
              "type": "remove_special_characters",
              "options": {
                "exclude": ["."]
              }
            },
            {
              "type": "remove_letters",
              "options": {}
            }
          ],
          "format": {
            "type": "decimal",
            "options": {
              "places": 2
            }
          }
        }
      },
      {
        "name": "employee_email",
        "is_required": true,
        "is_unique": true,
        "source": {
          "fields": ["Email"],
          "sanitizations": [
            {
              "type": "lowercase",
              "options": {}
            }
          ],
          "format": {
            "type": "email",
            "options": {}
          }
        }
      },
      {
        "name": "employee_id",
        "is_required": true,
        "is_unique": false,
        "source": {
          "fields": ["Number"],
          "sanitizations": []
        }
      },
      {
        "name": "employee_phone",
        "is_required": false,
        "is_unique": false,
        "source": {
          "fields": [""]
        }
      }
    ]
  }]`
				employers := `[`

				return strings.NewReader(schemas), strings.NewReader(employers)
			},
			expected:    nil,
			expectedErr: ErrInvalidReader{InternalErr: io.EOF},
		},
	}

	for _, tc := range testCases {
		tc = tc

		t.Run(tc.name, func(t *testing.T) {
			schemas, employers := tc.buildReaders()

			repo := NewJsonRepository(schemas, employers)

			employer, err := repo.GetByID(tc.arg)
			assert.Equal(t, tc.expected, employer)

			if tc.expectedErr != nil {
				assert.ErrorAs(t, err, &tc.expectedErr)
			}
		})
	}
}
