package strutil

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func NormalizeTitle(s string) string {
	caser := cases.Title(language.Indonesian)
	return caser.String(strings.ToLower(strings.TrimSpace(s)))
}
