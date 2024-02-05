package schema

import (
	"awesome-csv-parser/internal/concurrency"
	"awesome-csv-parser/pkg/schema/formatters"
	"awesome-csv-parser/pkg/schema/sanitizers"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestV1_UnmarshalJSON(t *testing.T) {
//	tests := []struct {
//		name          string
//		args          string
//		runAssertions func(*testing.T, *V1, error)
//	}{
//		{
//			name: "Schema 1 - Valid",
//			args: `{"id":"1","name":"schema for vendor A","target_fields":[{"name":"employee_name","is_required":true,"is_unique":false,"source":{"fields":["Name"],"sanitizations":[{"type":"remove_special_characters","options":{"exclude":[]}},{"type":"capitalize","options":{}}]}},{"name":"employee_salary","is_required":true,"is_unique":false,"source":{"fields":["Wage"],"sanitizations":[{"type":"remove_special_characters","options":{"exclude":["."]}},{"type":"remove_letters","options":{}}],"format":{"type":"decimal","options":{"places":2}}}},{"name":"employee_email","is_required":true,"is_unique":true,"source":{"fields":["Email"],"sanitizations":[{"type":"lowercase","options":{}}],"format":{"type":"email","options":{}}}},{"name":"employee_id","is_required":true,"is_unique":false,"source":{"fields":["Number"],"sanitizations":[]}},{"name":"employee_phone","is_required":false,"is_unique":false,"source":{"fields":[""]}}]}`,
//			runAssertions: func(t *testing.T, schema *V1, err error) {
//				expected := &V1{
//					ID:   "1",
//					Name: "schema for vendor A",
//					TargetFields: []Field{
//						{
//							Name:         "employee_name",
//							SourceFields: []string{"Name"},
//							IsRequired:   true,
//							IsUnique:     false,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.RemoveSpecialCharacters{},
//								&sanitizers.Capitalize{},
//							},
//						},
//						{
//							Name:         "employee_salary",
//							SourceFields: []string{"Wage"},
//							IsRequired:   true,
//							IsUnique:     false,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.RemoveSpecialCharacters{Exclude: []rune{'.'}},
//								&sanitizers.RemoveLetters{},
//							},
//							Formatter: &formatters.Decimal{
//								Places: 2,
//							},
//						},
//						{
//							Name:         "employee_email",
//							SourceFields: []string{"Email"},
//							IsRequired:   true,
//							IsUnique:     true,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.LowerCase{},
//							},
//							Formatter: &formatters.Email{},
//						},
//						{
//							Name:         "employee_id",
//							SourceFields: []string{"Number"},
//							IsRequired:   true,
//							IsUnique:     false,
//						},
//						{
//							Name:         "employee_phone",
//							SourceFields: []string{""},
//							IsRequired:   false,
//							IsUnique:     false,
//						},
//					},
//				}
//
//				assert.Equal(t, expected, schema)
//				assert.Nil(t, err)
//			},
//		},
//		{
//			name: "Schema 2 - Valid",
//			args: `{"id":"2","name":"schema for vendor B","target_fields":[{"name":"employee_name","is_required":true,"is_unique":false,"source":{"fields":["First","Last"],"sanitizations":[{"type":"remove_special_characters","options":{"exclude":[]}},{"type":"capitalize","options":{}}],"aggregation":{"type":"concat","options":{"delimiter":" "}}}},{"name":"employee_salary","is_required":true,"is_unique":false,"source":{"fields":["Salary"],"sanitizations":[{"type":"remove_special_characters","options":{"exclude":["."]}},{"type":"remove_letters","options":{}}],"format":{"type":"decimal","options":{"places":2}}}},{"name":"employee_email","is_required":true,"is_unique":true,"source":{"fields":["E-mail"],"sanitizations":[{"type":"lowercase","options":{}}],"format":{"type":"email","options":{}}}},{"name":"employee_id","is_required":true,"is_unique":false,"source":{"fields":["ID"],"sanitizations":[]}},{"name":"employee_phone","is_required":false,"is_unique":false,"source":{"fields":[""]}}]}`,
//			runAssertions: func(t *testing.T, schema *V1, err error) {
//				expected := &V1{
//					ID:   "2",
//					Name: "schema for vendor B",
//					TargetFields: []Field{
//						{
//							Name:         "employee_name",
//							SourceFields: []string{"First", "Last"},
//							IsRequired:   true,
//							IsUnique:     false,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.RemoveSpecialCharacters{},
//								&sanitizers.Capitalize{},
//							},
//							Aggregator: &aggregators.Concat{
//								Delimiter: " ",
//							},
//						},
//						{
//							Name:         "employee_salary",
//							SourceFields: []string{"Salary"},
//							IsRequired:   true,
//							IsUnique:     false,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.RemoveSpecialCharacters{Exclude: []rune{'.'}},
//								&sanitizers.RemoveLetters{},
//							},
//							Formatter: &formatters.Decimal{
//								Places: 2,
//							},
//						},
//						{
//							Name:         "employee_email",
//							SourceFields: []string{"E-mail"},
//							IsRequired:   true,
//							IsUnique:     true,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.LowerCase{},
//							},
//							Formatter: &formatters.Email{},
//						},
//						{
//							Name:         "employee_id",
//							SourceFields: []string{"ID"},
//							IsRequired:   true,
//							IsUnique:     false,
//						},
//						{
//							Name:         "employee_phone",
//							SourceFields: []string{""},
//							IsRequired:   false,
//							IsUnique:     false,
//						},
//					},
//				}
//
//				assert.Equal(t, expected, schema)
//				assert.Nil(t, err)
//			},
//		},
//		{
//			name: "Schema 3 - Valid",
//			args: `{"id":"3","name":"schema for vendor C","target_fields":[{"name":"employee_name","is_required":true,"is_unique":false,"source":{"fields":["first name","last name"],"sanitizations":[{"type":"remove_special_characters","options":{"exclude":[]}},{"type":"capitalize","options":{}}],"aggregation":{"type":"concat","options":{"delimiter":" "}}}},{"name":"employee_salary","is_required":true,"is_unique":false,"source":{"fields":["Rate"],"sanitizations":[{"type":"remove_special_characters","options":{"exclude":["."]}},{"type":"remove_letters","options":{}}],"format":{"type":"decimal","options":{"places":2}}}},{"name":"employee_email","is_required":true,"is_unique":true,"source":{"fields":["e-mail"],"sanitizations":[{"type":"lowercase","options":{}}],"format":{"type":"email","options":{}}}},{"name":"employee_id","is_required":true,"is_unique":false,"source":{"fields":["Employee Number"],"sanitizations":[]}},{"name":"employee_phone","is_required":false,"is_unique":false,"source":{"fields":["Mobile"],"sanitizations":[{"type":"remove_special_characters","options":{"exclude":["+","(",")","-"]}},{"type":"trim_spaces","options":{}}]}}]}`,
//			runAssertions: func(t *testing.T, schema *V1, err error) {
//				expected := &V1{
//					ID:   "3",
//					Name: "schema for vendor C",
//					TargetFields: []Field{
//						{
//							Name:         "employee_name",
//							SourceFields: []string{"first name", "last name"},
//							IsRequired:   true,
//							IsUnique:     false,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.RemoveSpecialCharacters{},
//								&sanitizers.Capitalize{},
//							},
//							Aggregator: &aggregators.Concat{
//								Delimiter: " ",
//							},
//						},
//						{
//							Name:         "employee_salary",
//							SourceFields: []string{"Rate"},
//							IsRequired:   true,
//							IsUnique:     false,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.RemoveSpecialCharacters{Exclude: []rune{'.'}},
//								&sanitizers.RemoveLetters{},
//							},
//							Formatter: &formatters.Decimal{
//								Places: 2,
//							},
//						},
//						{
//							Name:         "employee_email",
//							SourceFields: []string{"e-mail"},
//							IsRequired:   true,
//							IsUnique:     true,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.LowerCase{},
//							},
//							Formatter: &formatters.Email{},
//						},
//						{
//							Name:         "employee_id",
//							SourceFields: []string{"Employee Number"},
//							IsRequired:   true,
//							IsUnique:     false,
//						},
//						{
//							Name:         "employee_phone",
//							SourceFields: []string{"Mobile"},
//							IsRequired:   false,
//							IsUnique:     false,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.RemoveSpecialCharacters{Exclude: []rune{'+', '(', ')', '-'}},
//								&sanitizers.TrimSpace{},
//							},
//						},
//					},
//				}
//
//				assert.Equal(t, expected, schema)
//				assert.Nil(t, err)
//			},
//		},
//		{
//			name: "Schema 4 - Valid",
//			args: `{"id":"4","name":"schema for vendor D","target_fields":[{"name":"employee_name","is_required":true,"is_unique":false,"source":{"fields":["f. name","l. name"],"sanitizations":[{"type":"remove_special_characters","options":{"exclude":[]}},{"type":"capitalize","options":{}}],"aggregation":{"type":"concat","options":{"delimiter":" "}}}},{"name":"employee_salary","is_required":true,"is_unique":false,"source":{"fields":["wage"],"sanitizations":[{"type":"remove_special_characters","options":{"exclude":["."]}},{"type":"remove_letters","options":{}}],"format":{"type":"decimal","options":{"places":2}}}},{"name":"employee_email","is_required":true,"is_unique":true,"source":{"fields":["email"],"sanitizations":[{"type":"lowercase","options":{}}],"format":{"type":"email","options":{}}}},{"name":"employee_id","is_required":true,"is_unique":false,"source":{"fields":["emp id"],"sanitizations":[]}},{"name":"employee_phone","is_required":false,"is_unique":false,"source":{"fields":["phone"],"sanitizations":[{"type":"remove_special_characters","options":{"exclude":["+","(",")","-"]}},{"type":"trim_spaces","options":{}}]}}]}`,
//			runAssertions: func(t *testing.T, schema *V1, err error) {
//				expected := &V1{
//					ID:   "4",
//					Name: "schema for vendor D",
//					TargetFields: []Field{
//						{
//							Name:         "employee_name",
//							SourceFields: []string{"f. name", "l. name"},
//							IsRequired:   true,
//							IsUnique:     false,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.RemoveSpecialCharacters{},
//								&sanitizers.Capitalize{},
//							},
//							Aggregator: &aggregators.Concat{
//								Delimiter: " ",
//							},
//						},
//						{
//							Name:         "employee_salary",
//							SourceFields: []string{"wage"},
//							IsRequired:   true,
//							IsUnique:     false,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.RemoveSpecialCharacters{Exclude: []rune{'.'}},
//								&sanitizers.RemoveLetters{},
//							},
//							Formatter: &formatters.Decimal{
//								Places: 2,
//							},
//						},
//						{
//							Name:         "employee_email",
//							SourceFields: []string{"email"},
//							IsRequired:   true,
//							IsUnique:     true,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.LowerCase{},
//							},
//							Formatter: &formatters.Email{},
//						},
//						{
//							Name:         "employee_id",
//							SourceFields: []string{"emp id"},
//							IsRequired:   true,
//							IsUnique:     false,
//						},
//						{
//							Name:         "employee_phone",
//							SourceFields: []string{"phone"},
//							IsRequired:   false,
//							IsUnique:     false,
//							Sanitizers: []sanitizers.Sanitizer{
//								&sanitizers.RemoveSpecialCharacters{Exclude: []rune{'+', '(', ')', '-'}},
//								&sanitizers.TrimSpace{},
//							},
//						},
//					},
//				}
//
//				assert.Equal(t, expected, schema)
//				assert.Nil(t, err)
//			},
//		},
//	}
//	for _, testCase := range tests {
//		tc := testCase
//
//		t.Run(tc.name, func(t *testing.T) {
//			got := new(V1)
//
//			err := got.UnmarshalJSON([]byte(tc.args))
//			tc.runAssertions(t, got, err)
//		})
//	}
//}

func TestV1_Headers(t *testing.T) {
	tests := []struct {
		runAssertions func(*testing.T, *V1, error)
		name          string
		args          string
	}{
		{
			name: "Schema 1 - Valid",
			args: `{"id":"1","name":"schema for vendor A","target_fields":[{"name":"employee_name","is_required":true,"is_unique":false,"source":{"fields":["Name"],"sanitizations":[{"type":"remove_special_characters","options":{"exclude":[]}},{"type":"capitalize","options":{}}]}},{"name":"employee_salary","is_required":true,"is_unique":false,"source":{"fields":["Wage"],"sanitizations":[{"type":"remove_special_characters","options":{"exclude":["."]}},{"type":"remove_letters","options":{}}],"format":{"type":"decimal","options":{"places":2}}}},{"name":"employee_email","is_required":true,"is_unique":true,"source":{"fields":["Email"],"sanitizations":[{"type":"lowercase","options":{}}],"format":{"type":"email","options":{}}}},{"name":"employee_id","is_required":true,"is_unique":false,"source":{"fields":["Number"],"sanitizations":[]}},{"name":"employee_phone","is_required":false,"is_unique":false,"source":{"fields":[""]}}]}`,
			runAssertions: func(t *testing.T, schema *V1, err error) {
				expected := []string{
					"employee_name",
					"employee_salary",
					"employee_email",
					"employee_id",
					"employee_phone",
				}

				assert.Equal(t, expected, schema.Headers())
				assert.Nil(t, err)
			},
		},
	}

	for _, testCase := range tests {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			var dto V1DTO

			err := json.Unmarshal([]byte(tc.args), &dto)
			assert.NoError(t, err)

			got, err := NewFromDTO(dto)
			tc.runAssertions(t, got, err)
		})
	}
}

func TestV1_Build(t *testing.T) {
	type args struct {
		schema    *V1
		headerMap HeaderMap
		record    []string
	}

	tests := []struct {
		args          args
		runAssertions func(*testing.T, []string, error)
		name          string
	}{
		{
			name: "Build - Valid",
			args: args{
				record: []string{
					"John doe!!!",
					"doe@Test.com",
					"$1,000.00",
				},
				headerMap: HeaderMap{
					m: map[string]int{
						"Name":  0,
						"Email": 1,
						"Wage":  2,
					},
					keys: []string{
						"Name",
						"Email",
						"Wage",
					},
				},
				schema: &V1{
					ID:   "1",
					Name: "schema for vendor A",
					TargetFields: []Field{
						{
							Name:         "employee_name",
							SourceFields: []string{"Name"},
							IsRequired:   true,
							IsUnique:     false,
							Sanitizers: []sanitizers.Sanitizer{
								&sanitizers.RemoveSpecialCharacters{},
								&sanitizers.Capitalize{},
							},
						},
						{
							Name:         "employee_salary",
							SourceFields: []string{"Wage"},
							IsRequired:   true,
							IsUnique:     false,
							Sanitizers: []sanitizers.Sanitizer{
								&sanitizers.RemoveSpecialCharacters{Exclude: []rune{'.'}},
								&sanitizers.RemoveLetters{},
							},
							Formatter: &formatters.Decimal{
								Places: 2,
							},
						},
						{
							Name:         "employee_email",
							SourceFields: []string{"Email"},
							IsRequired:   true,
							IsUnique:     true,
							Sanitizers: []sanitizers.Sanitizer{
								&sanitizers.LowerCase{},
							},
							Formatter: &formatters.Email{},
						},
					},
					ShardedMap: concurrency.NewShardedMap[bool](1024),
				},
			},
			runAssertions: func(t *testing.T, got []string, err error) {
				expected := []string{
					"John Doe",
					"1000.00",
					"doe@test.com",
				}

				assert.Equal(t, expected, got)
			},
		},
		{
			name: "Should return error when required field is missing",
			args: args{
				record: []string{
					"John doe!!!",
					"doe@Test.com",
				},
				headerMap: HeaderMap{
					m: map[string]int{
						"Name":  0,
						"Email": 1,
					},
					keys: []string{
						"Name",
						"Email",
					},
				},
				schema: &V1{
					ID:   "1",
					Name: "schema for vendor A",
					TargetFields: []Field{
						{
							Name:         "employee_name",
							SourceFields: []string{"Name"},
							IsRequired:   true,
							IsUnique:     false,
							Sanitizers: []sanitizers.Sanitizer{
								&sanitizers.RemoveSpecialCharacters{},
								&sanitizers.Capitalize{},
							},
						},
						{
							Name:         "employee_salary",
							SourceFields: []string{"Wage"},
							IsRequired:   true,
							IsUnique:     false,
							Sanitizers: []sanitizers.Sanitizer{
								&sanitizers.RemoveSpecialCharacters{Exclude: []rune{'.'}},
								&sanitizers.RemoveLetters{},
							},
							Formatter: &formatters.Decimal{
								Places: 2,
							},
						},
						{
							Name:         "employee_email",
							SourceFields: []string{"Email"},
							IsRequired:   true,
							IsUnique:     true,
							Sanitizers: []sanitizers.Sanitizer{
								&sanitizers.LowerCase{},
							},
							Formatter: &formatters.Email{},
						},
					},
				},
			},
			runAssertions: func(t *testing.T, got []string, err error) {
				assert.Equal(t, "field Wage not found", err.Error())
			},
		},
		{
			name: "Should return error when required field is empty",
			args: args{
				record: []string{
					"John doe!!!",
					"doe@Test.com",
					"",
				},
				headerMap: HeaderMap{
					m: map[string]int{
						"Name":  0,
						"Email": 1,
						"Wage":  2,
					},
					keys: []string{
						"Name",
						"Email",
						"Wage",
					},
				},
				schema: &V1{
					ID:   "1",
					Name: "schema for vendor A",
					TargetFields: []Field{
						{
							Name:         "employee_name",
							SourceFields: []string{"Name"},
							IsRequired:   true,
							IsUnique:     false,
							Sanitizers: []sanitizers.Sanitizer{
								&sanitizers.RemoveSpecialCharacters{},
								&sanitizers.Capitalize{},
							},
						},
						{
							Name:         "employee_salary",
							SourceFields: []string{"Wage"},
							IsRequired:   true,
							IsUnique:     false,
							Sanitizers: []sanitizers.Sanitizer{
								&sanitizers.RemoveSpecialCharacters{Exclude: []rune{'.'}},
								&sanitizers.RemoveLetters{},
							},
							Formatter: &formatters.Decimal{
								Places: 2,
							},
						},
						{
							Name:         "employee_email",
							SourceFields: []string{"Email"},
							IsRequired:   true,
							IsUnique:     true,
							Sanitizers: []sanitizers.Sanitizer{
								&sanitizers.LowerCase{},
							},
							Formatter: &formatters.Email{},
						},
					},
				},
			},
			runAssertions: func(t *testing.T, got []string, err error) {
				assert.Equal(t, "field Wage is empty", err.Error())
			},
		},
	}

	for _, testCase := range tests {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.args.schema.Build(tc.args.record, tc.args.headerMap)
			tc.runAssertions(t, got, err)
		})
	}
}
