package sanitizers

import "testing"

func TestCapitalize_Sanitize(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test 1", args: args{"joao"}, want: "Joao"},
		{name: "test 2", args: args{"joao silva"}, want: "Joao Silva"},
	}
	for _, tt := range tests {
		tt = tt
		t.Run(tt.name, func(t *testing.T) {
			c := &Capitalize{}
			if got := c.Sanitize(tt.args.value); got != tt.want {
				t.Errorf("Capitalize.Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLowerCase_Sanitize(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test 1", args: args{"Joao"}, want: "joao"},
		{name: "test 2", args: args{"Joao Silva"}, want: "joao silva"},
	}
	for _, tt := range tests {
		tt = tt
		t.Run(tt.name, func(t *testing.T) {
			l := &LowerCase{}
			if got := l.Sanitize(tt.args.value); got != tt.want {
				t.Errorf("LowerCase.Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveLetters_Sanitize(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test 1", args: args{"Joao"}, want: ""},
		{name: "test 2", args: args{"Joao Silva"}, want: " "},
	}
	for _, tt := range tests {
		tt = tt
		t.Run(tt.name, func(t *testing.T) {
			r := &RemoveLetters{}
			if got := r.Sanitize(tt.args.value); got != tt.want {
				t.Errorf("RemoveLetters.Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveSpecialCharacters_Sanitize(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test 1", args: args{"Joao"}, want: "Joao"},
		{name: "test 2", args: args{"Joao Silva"}, want: "Joao Silva"},
		{name: "test 3", args: args{"Joao Silva 123"}, want: "Joao Silva 123"},
		{name: "test 4", args: args{"Joao Silva 123 !@#$%"}, want: "Joao Silva 123 "},
	}
	for _, tt := range tests {
		tt = tt
		t.Run(tt.name, func(t *testing.T) {
			r := &RemoveSpecialCharacters{}
			if got := r.Sanitize(tt.args.value); got != tt.want {
				t.Errorf("RemoveSpecialCharacters.Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimSpace_Sanitize(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test 1", args: args{" Joao "}, want: "Joao"},
		{name: "test 2", args: args{" Joao Silva "}, want: "Joao Silva"},
	}
	for _, tt := range tests {
		tt = tt
		t.Run(tt.name, func(t *testing.T) {
			ts := &TrimSpace{}
			if got := ts.Sanitize(tt.args.value); got != tt.want {
				t.Errorf("TrimSpace.Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		name string
		opts Opts
	}
	tests := []struct {
		want    Sanitizer
		name    string
		args    args
		wantErr bool
	}{
		{name: "test 1", args: args{"capitalize", Opts{}}, want: &Capitalize{}, wantErr: false},
		{name: "test 2", args: args{"lowercase", Opts{}}, want: &LowerCase{}, wantErr: false},
		{name: "test 3", args: args{"remove_letters", Opts{}}, want: &RemoveLetters{}, wantErr: false},
		{name: "test 4", args: args{"remove_special_characters", Opts{}}, want: &RemoveSpecialCharacters{}, wantErr: false},
		{name: "test 5", args: args{"trim_spaces", Opts{}}, want: &TrimSpace{}, wantErr: false},
		{name: "test 6", args: args{"notfound", Opts{}}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		tt = tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && got.Sanitize("Joao") != tt.want.Sanitize("Joao") {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
