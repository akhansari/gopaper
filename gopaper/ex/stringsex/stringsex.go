package stringsex

import (
	"regexp"
	"strings"
)

func SplitAndTrimSpace(str string, sep string) []string {
	a := strings.Split(str, sep)
	for key, value := range a {
		a[key] = strings.TrimSpace(value)
		a[key] = strings.Replace(a[key], " ", "-", -1)
	}
	return a
}

func FormatUrl(str string) string {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)
	str = strings.Replace(str, " ", "-", -1)
	rgx, _ := regexp.Compile("[^\\w-]")
	str = rgx.ReplaceAllString(str, "")
	return str
}
