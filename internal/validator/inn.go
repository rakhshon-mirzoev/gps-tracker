package validator

import (
	"regexp"
	"strings"
)

var innRegex = regexp.MustCompile(`^\d{9}$`)

func NormalizeAndValidateINN(inn string) (string, bool) {
	inn = strings.TrimSpace(inn)

	if !innRegex.MatchString(inn) {
		return "", false
	}

	return inn, true
}
