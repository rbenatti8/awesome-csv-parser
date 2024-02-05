package aggregators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	r, err := New("concat", Opts{Delimiter: ","})
	assert.NoError(t, err)

	assert.Equal(t, &Concat{Delimiter: ","}, r)
}

func TestConcat_Aggregate(t *testing.T) {
	c := Concat{Delimiter: ","}

	result, err := c.Aggregate([]string{"a", "b", "c"})
	assert.NoError(t, err)

	assert.Equal(t, "a,b,c", result)
}
