package validation

import (
	"fmt"
	"net/mail"
	"strings"

	"gopkg.in/validator.v2"
)

var (
	passwordMinLength         = 8
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
	}

	for _, v := range validators {
		_ = validator.SetValidationFunc(v.name, v.fn)
	}
}

func password(v interface{}, param string) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("passord only valifates strings")
	}

	if len(s) < passwordMinLength {
		return fmt.Errorf("password must be at least %d characters long", passwordMinLength)
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
		return fmt.Errorf("passord only valifates strings")
	}

	if len(s) > emailMaxLegth {
		return fmt.Errorf("email must at most %d characters long", emailMaxLegth)
	}

	_, err := mail.ParseAddress(s)
	return err
}
