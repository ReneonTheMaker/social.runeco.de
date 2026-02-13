package web

import (
	"fmt"
	"time"

	"github.com/gofiber/template/html/v2"
)

func defaultValue(v any, d any) any {
	if v == nil {
		return d
	}
	return v
}

func dict(values ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i < len(values); i += 2 {
		m[values[i].(string)] = values[i+1]
	}
	return m
}

func printf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func printtime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

func printstatus(status int) string {
	switch status {
	case 0:
		return "Pending"
	case 1:
		return "In Progress"
	case 2:
		return "Synced"
	default:
		return "Unknown"
	}
}

func registerRenderFunctions(engine *html.Engine) {
	// Register custom template functions
	engine.AddFunc("default", defaultValue)
	engine.AddFunc("dict", dict)
	engine.AddFunc("printf", printf)
	engine.AddFunc("printtime", printtime)
	engine.AddFunc("printstatus", printstatus)
}
