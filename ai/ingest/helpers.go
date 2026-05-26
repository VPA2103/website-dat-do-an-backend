package ingest

import (
	"fmt"
	"strings"
)

func toString(v any) string {
	s, ok := v.(string)
	if ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func normalizeAnyList(v any) []string {
	if v == nil {
		return nil
	}
	switch t := v.(type) {
	case []any:
		out := make([]string, 0, len(t))
		for _, x := range t {
			s := strings.TrimSpace(toString(x))
			if s != "" {
				out = append(out, s)
			}
		}
		if len(out) == 0 {
			return nil
		}
		return out
	case []string:
		out := make([]string, 0, len(t))
		for _, x := range t {
			s := strings.TrimSpace(x)
			if s != "" {
				out = append(out, s)
			}
		}
		if len(out) == 0 {
			return nil
		}
		return out
	case string:
		return toList(t)
	default:
		return toList(toString(t))
	}
}
