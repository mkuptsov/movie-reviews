package validation

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/mkuptsov/movie-reviews/internal/modules/users"
	"gopkg.in/validator.v2"
)

var (
	passwordMinLength         = 8
	passwordMaxLenth          = 72
	emailMaxLegth             = 127
	passwordSpecialCharacters = "!%$#()[]{}?+*~@^&-_"
	passwordRequiredEntries   = []struct {
		name  string
		chars string
	}{
		{"lowercase character", "abcdefghijklmnopqrstuvwxyz"},
		{"uppercase character", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{"digit", "0123456789"},
		{"special character ( " + passwordSpecialCharacters + ")", passwordSpecialCharacters},
	}
)

func SetupValidators() {
	validators := []struct {
		name string
		fn   validator.ValidationFunc
	}{
		{"password", password},
		{"email", email},
		{"role", role},
		{"sort", sort},
	}

	for _, v := range validators {
		_ = validator.SetValidationFunc(v.name, v.fn)
	}
}

//nolint:revive // function requires param
func password(v interface{}, param string) error {
	// ...
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("password only validates strings")
	}

	if len(s) < passwordMinLength || len(s) > passwordMaxLenth {
		return fmt.Errorf("password must be at least %d and not more than %d characters long", passwordMinLength, passwordMaxLenth)
	}

	for _, required := range passwordRequiredEntries {
		if !strings.ContainsAny(s, required.chars) {
			return fmt.Errorf("password must contain at least one %s", required.name)
		}
	}
	return nil
}

//nolint:revive // function requires param
func email(v interface{}, param string) error {
	// ...
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("email only validates strings")
	}

	if len(s) > emailMaxLegth {
		return fmt.Errorf("email must be at most %d characters long", emailMaxLegth)
	}

	_, err := mail.ParseAddress(s)
	if err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	return nil
}

//nolint:revive // function requires param
func role(v interface{}, param string) error {
	// ...
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("role only validates strings")
	}
	if s != "" && s != users.AdminRole && s != users.EditorRole && s != users.UserRole {
		return fmt.Errorf("invalid role")
	}
	return nil
}

func sort(v interface{}, _ string) error {
	validate := func(s *string) error {
		if s == nil {
			return nil
		}
		switch *s {
		case "desc", "asc":
			return nil
		}
		return fmt.Errorf("sort must be one of desc or asc")
	}

	switch s := v.(type) {
	case string:
		return validate(&s)
	case *string:
		return validate(s)
	default:
		return fmt.Errorf("sort only validates string or pointer to string")
	}
}
