package templateex

import (
	"html"
	"html/template"
	"net/url"
	"reflect"
	"time"
)

func Equal(args ...interface{}) bool {
	if len(args) == 0 {
		return false
	}
	x := args[0]
	switch x := x.(type) {
	case string, int, int64, byte, float32, float64:
		for _, y := range args[1:] {
			if x == y {
				return true
			}
		}
		return false
	}

	for _, y := range args[1:] {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}

func Addition(num1 int, num2 int) int {
	return num1 + num2
}

func FormatDate(datetime time.Time) string {
	if datetime.IsZero() {
		return ""
	}
	if time.Now().Format("020106") == datetime.Format("020106") {
		return "Today"
	}
	if time.Now().AddDate(0, 0, -1).Format("020106") == datetime.Format("020106") {
		return "Yesterday"
	}
	return datetime.Format("02 Jan 2006 15h04")
}

func FormatBool(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

func HtmlSafe(text string) template.HTML {
	return template.HTML(text)
}

func HtmlEscape(text string) string {
	return html.EscapeString(text)
}

func QueryEscape(text string) string {
	return url.QueryEscape(text)
}
