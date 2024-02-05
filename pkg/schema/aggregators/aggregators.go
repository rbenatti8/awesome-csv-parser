package aggregators

import "strings"

type Opts struct {
	Delimiter string
}

type Aggregator interface {
	Aggregate([]string) (string, error)
}

type ErrAggregatorNotFound struct {
	Name string
}

func (e *ErrAggregatorNotFound) Error() string {
	return "aggregator not found: " + e.Name
}

func New(name string, opts Opts) (Aggregator, error) {
	switch name {
	case "concat":
		return &Concat{
			Delimiter: opts.Delimiter,
		}, nil
	default:
		return nil, &ErrAggregatorNotFound{Name: name}
	}
}

type Concat struct {
	Delimiter string
}

func (c *Concat) Aggregate(values []string) (string, error) {
	return strings.Join(values, c.Delimiter), nil
}
