package sanitizers

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"unicode"
)

func New(name string, opts Opts) (Sanitizer, error) {
	switch name {
	case "trim_spaces":
		return &TrimSpace{}, nil
	case "lowercase":
		return &LowerCase{}, nil
	case "capitalize":
		return &Capitalize{}, nil
	case "remove_special_characters":
		if len(opts.Exclude) == 0 {
			return &RemoveSpecialCharacters{}, nil
		}

		exclude := make([]rune, 0, len(opts.Exclude))

		for _, e := range opts.Exclude {
			exclude = append(exclude, []rune(e)...)
		}

		return &RemoveSpecialCharacters{
			Exclude: exclude,
		}, nil
	case "remove_letters":
		if len(opts.Exclude) == 0 {
			return &RemoveLetters{}, nil
		}

		exclude := make([]rune, 0, len(opts.Exclude))

		for _, e := range opts.Exclude {
			exclude = append(exclude, []rune(e)...)
		}

		return &RemoveLetters{
			Exclude: exclude,
		}, nil
	default:
		return nil, &ErrSanitizerNotFound{Name: name}
	}
}

type Opts struct {
	Exclude []string `json:"exclude"`
}

type Sanitizer interface {
	Sanitize(string) string
}

type ErrSanitizerNotFound struct {
	Name string
}

func (e *ErrSanitizerNotFound) Error() string {
	return "sanitizer not found: " + e.Name
}

type Capitalize struct {
}

func (f *Capitalize) Sanitize(input string) string {
	return cases.Title(language.English).String(input)
}

type LowerCase struct{}

func (l *LowerCase) Sanitize(input string) string {
	return strings.ToLower(input)
}

type RemoveLetters struct {
	Exclude []rune
}

func (r *RemoveLetters) Sanitize(input string) string {
	var resultBuilder []rune

	for _, char := range input {
		if !unicode.IsLetter(char) || r.isExcluded(char) {
			resultBuilder = append(resultBuilder, char)
		}
	}

	return string(resultBuilder)
}

func (r *RemoveLetters) isExcluded(char rune) bool {
	for _, excludedChar := range r.Exclude {
		if excludedChar == char {
			return true
		}
	}

	return false
}

type RemoveSpecialCharacters struct {
	Exclude []rune
}

func (r *RemoveSpecialCharacters) Sanitize(input string) string {
	var resultBuilder []rune

	for _, char := range input {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || unicode.IsSpace(char) || r.isExcluded(char) {
			resultBuilder = append(resultBuilder, char)
		}
	}

	return string(resultBuilder)
}

func (r *RemoveSpecialCharacters) isExcluded(char rune) bool {
	for _, excludedChar := range r.Exclude {
		if excludedChar == char {
			return true
		}
	}

	return false
}

type TrimSpace struct{}

func (r *TrimSpace) Sanitize(value string) string {
	return strings.TrimSpace(value)
}
