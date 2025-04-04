package utils

import (
	gonanoid "github.com/matoous/go-nanoid"
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
