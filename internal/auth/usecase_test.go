package auth

import (
	"balance-ledger-database-design/pkg/response"
	"errors"
	"strings"
	"testing"
)

func TestValidateRegister(t *testing.T) {
	valid := RegisterRequest{Email: "a@example.com", FullName: "A", Password: "secret123"}
	cases := []struct {
		name string
		edit func(r *RegisterRequest)
		want error
	}{
		{"valid", func(r *RegisterRequest) {}, nil},
		{"not an email", func(r *RegisterRequest) { r.Email = "not-an-email" }, ErrInvalidEmail},
		{"display name", func(r *RegisterRequest) { r.Email = "A <a@example.com>" }, ErrInvalidEmail},
		{"blank name", func(r *RegisterRequest) { r.FullName = "  " }, ErrInvalidFullName},
		{"7-byte password", func(r *RegisterRequest) { r.Password = "1234567" }, ErrPasswordLength},
		{"72-byte password", func(r *RegisterRequest) { r.Password = strings.Repeat("a", 72) }, nil},
		{"73-byte password", func(r *RegisterRequest) { r.Password = strings.Repeat("a", 73) }, ErrPasswordLength},
	}
	for _, c := range cases {
		r := valid
		c.edit(&r)
		err := validateRegister(r)
		var re *response.Error
		if errors.As(err, &re) {
			err = re.Err
		}
		if !errors.Is(err, c.want) {
			t.Errorf("%s: got %v, want %v", c.name, err, c.want)
		}
	}
}
