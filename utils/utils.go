package utils

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func FormatTag(tag string) string {
	tag = strings.ReplaceAll(strings.TrimSpace(tag), ".", " ")

	tag = strings.ToLower(tag)

	tag = cases.Title(language.English).String(tag)

	re := regexp.MustCompile(`\b(js|db|Js|Db)\b`)
	tag = re.ReplaceAllStringFunc(tag, func(s string) string {
		return strings.ToUpper(s)
	})

	return tag
}