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

func ToPgUUID(uuidStr string) pgtype.UUID {
	var uuid pgtype.UUID
	if err := uuid.Scan(uuidStr); err != nil {
		return pgtype.UUID{}
	}
	return uuid
}

func PgUUIDToStringPtr(uuid pgtype.UUID) *string {
	if uuid.Valid {
		str := uuid.String()
		return &str
	}
	return nil
}
