package validator

import "unicode"

func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasLetter, hasDigit bool

	for _, r := range password {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}
