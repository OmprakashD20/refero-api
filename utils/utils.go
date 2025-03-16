package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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

func GenerateShortURL(url string) string {
	hash := sha256.Sum256([]byte(url + time.Now().String()))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return strings.TrimRight(encoded[:10], "=")
}
