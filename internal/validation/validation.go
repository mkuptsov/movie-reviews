package validation

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
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
	}

	for _, v := range validators {
		_ = validator.SetValidationFunc(v.name, v.fn)
	}
}

func password(v interface{}, param string) error {
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

func email(v interface{}, param string) error {
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

func role(v interface{}, param string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("role only validates strings")
	}
	if s != "" && s != users.AdminRole && s != users.EditorRole && s != users.UserRole {
		return fmt.Errorf("invalid role")
	}
	return nil
}
