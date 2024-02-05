package employer

import (
	"errors"
	"io"
	"testing"
)

func TestErrEmployerNotFound_Error(t *testing.T) {
	type fields struct {
		ID string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "returns error message",
			fields: fields{
				ID: "1",
			},
			want: "employer with id (1) not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrEmployerNotFound{
				ID: tt.fields.ID,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ErrEmployerNotFound.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrInvalidReader_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "returns error message",

			want: "invalid reader: EOF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrInvalidReader{
				InternalErr: io.EOF,
			}

			if got := e.Error(); got != tt.want {
				t.Errorf("ErrEmployerNotFound.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrInvalidReader_Unwrap(t *testing.T) {
	tests := []struct {
		want error
		name string
	}{
		{
			name: "returns internal error",
			want: io.EOF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrInvalidReader{
				InternalErr: io.EOF,
			}
			if got := e.Unwrap(); !errors.Is(got, tt.want) {
				t.Errorf("ErrInvalidReader.Unwrap() = %v, want %v", got, tt.want)
			}
		})
	}
}
