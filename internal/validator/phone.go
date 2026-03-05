package validator

import (
	"regexp"
	"strings"
)

var tajikPhoneRegex = regexp.MustCompile(`^\d{9}$`)

func NormalizeTajikPhone(phone string) (string, bool) {
	phone = strings.TrimSpace(phone)

	phone = strings.TrimPrefix(phone, "+")

	if after, ok := strings.CutPrefix(phone, "992"); ok {
		phone = after
	}

	if !tajikPhoneRegex.MatchString(phone) {
		return "", false
	}

	return "+992" + phone, true
}
