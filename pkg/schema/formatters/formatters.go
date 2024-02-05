package formatters

import (
	"fmt"
	"github.com/shopspring/decimal"
	"regexp"
)

type ErrFormatterNotFound struct {
	Name string
}

func (e *ErrFormatterNotFound) Error() string {
	return fmt.Sprintf("formatter %s not found", e.Name)
}

type Opts struct {
	Places int32
}

type Formatter interface {
	Format(value string) (string, error)
}

func New(name string, opts Opts) (Formatter, error) {
	switch name {
	case "decimal":
		return &Decimal{
			Places: opts.Places,
		}, nil
	case "email":
		return &Email{}, nil
	default:
		return nil, &ErrFormatterNotFound{Name: name}
	}
}

type Decimal struct {
	Places int32
}

func (d *Decimal) Format(value string) (string, error) {
	v, err := decimal.NewFromString(value)
	if err != nil {
		return "", fmt.Errorf("field %s is not a valid decimal", value)
	}

	return v.Truncate(d.Places).StringFixed(d.Places), nil
}

var regex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

type Email struct {
}

func (e *Email) Format(value string) (string, error) {
	isValid, err := regexp.MatchString(regex, value)
	if err != nil || !isValid {
		return "", fmt.Errorf("%s is not a valid email", value)
	}

	return value, nil
}
