package utils

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func FormatFieldName(field string) string {
	re := regexp.MustCompile(`([a-z])([A-Z])`)

	return cases.Title(language.English).String(
		re.ReplaceAllString(
			strings.ReplaceAll(field, "_", " "),
			"$1 $2",
		),
	)
}