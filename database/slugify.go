package database

import (
	"strings"
	"unicode"
)

func Slugify(str string) string {
	var slugBuilder strings.Builder

	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			slugBuilder.WriteRune(unicode.ToLower(r))
		} else if unicode.IsSpace(r) {
			slugBuilder.WriteRune('-')
		}
	}

	slug := strings.Trim(slugBuilder.String(), "-")

	return slug
}
