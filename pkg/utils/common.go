package utils

import (
	gonanoid "github.com/matoous/go-nanoid"
	"strconv"
	"strings"
	"time"
)

// RandStringBytes Generate random string by length n
func RandStringBytes(n int, upper bool) string {
	var LetterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	if upper {
		LetterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	}

	rand, err := gonanoid.Generate(LetterBytes, n)
	if err != nil {
		return ""
	}

	return rand
}

func ConvertUnixMilliToTime(unixMilli int64) string {
	// Convert Unix milliseconds to time.Time
	seconds := unixMilli / 1000
	nanoseconds := (unixMilli % 1000) * int64(time.Millisecond)
	t := time.Unix(seconds, nanoseconds).UTC()

	// Format the time.Time object to a string with the specified format
	return t.Format("2006-01-02 15:04:05.000000 -07:00")
}

func SafeFloatPointer(f *float64, defaultValue float64) *float64 {
	if f != nil {
		return f
	}
	return &defaultValue
}

func SafeIntPointer(i *int, defaultValue int) *int {
	if i != nil {
		return i
	}
	return &defaultValue
}

func SafeStringPointer(s *string, defaultValue string) *string {
	if s != nil {
		return s
	}
	return &defaultValue
}

func ConvertTimeToMillisString(t *time.Time) string {
	return strconv.FormatInt(t.UnixMilli(), 10)
}

var bannedWords = []string{
	"cặc",
	"lồn",
	"địt",
	"đù",
	"má mày",
	"mẹ mày",
	"đù mẹ",
	"đụ mẹ",
	"đù má mày",
	"đụ mẹ mày",
	"địt mẹ mày",
	"cac",
	"lon",
	"dit",
	"du",
	"ma may",
	"du me",
	"du ma may",
	"du me may",
	"me may",
	"dit me may",
}

func SanitizeComment(input string) string {
	result := input
	lowered := strings.ToLower(input)

	for _, banned := range bannedWords {
		index := strings.Index(lowered, banned)
		if index != -1 {
			mask := strings.Repeat("*", len(banned))
			result = strings.ReplaceAll(result, banned, mask)
		}
	}
	return result
}
