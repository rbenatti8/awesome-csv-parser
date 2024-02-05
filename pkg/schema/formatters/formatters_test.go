package formatters

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecimal_Format(t *testing.T) {
	type fields struct {
		DecimalPlaces int32
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		fields  fields
		wantErr bool
	}{
		{name: "test 1", fields: fields{2}, args: args{"1.234"}, want: "1.23"},
		{name: "test 2", fields: fields{2}, args: args{"1.2"}, want: "1.20"},
		{name: "test 3", fields: fields{2}, args: args{"aa"}, wantErr: true},
	}

	for _, tt := range tests {
		tt = tt
		t.Run(tt.name, func(t *testing.T) {
			d := &Decimal{
				Places: tt.fields.DecimalPlaces,
			}

			r, err := d.Format(tt.args.value)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, r)
		})
	}
}

func TestEmail_Format(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "test 1", args: args{"joao@gmail.com"}, want: "joao@gmail.com"},
		{name: "test 2", args: args{"joao@gmail"}, wantErr: true},
	}

	for _, tt := range tests {
		tt = tt
		t.Run(tt.name, func(t *testing.T) {
			e := &Email{}

			r, err := e.Format(tt.args.value)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, r)
		})
	}
}
