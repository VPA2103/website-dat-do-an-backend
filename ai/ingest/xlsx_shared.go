package ingest

import (
	"path/filepath"
	"regexp"
	"strings"
)

var canonicalHeaders = map[string][]string{
	"id":          {"id", "ma", "mã", "code"},
	"name":        {"name", "ten", "tên", "dish", "món", "mon"},
	"description": {"description", "mo ta", "mô tả", "desc"},
	"price":       {"price", "gia", "giá"},
	"tags":        {"tags", "tag"},
	"allergens":   {"allergens", "di ung", "dị ứng"},
	"ingredients": {"ingredients", "nguyen lieu", "nguyên liệu"},
}

var wsRe = regexp.MustCompile(`\s+`)

func HasXLSXExtension(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".xlsx"
}

func normHeader(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = wsRe.ReplaceAllString(s, " ")
	return s
}

func headerMap(headers []string) map[string]int {
	idx := map[string]int{}
	for canonical, variants := range canonicalHeaders {
		for i, h := range headers {
			n := normHeader(h)
			for _, v := range variants {
				if n == v {
					idx[canonical] = i
					goto nextCanonical
				}
			}
		}
	nextCanonical:
	}
	return idx
}

func cell(row []string, i int, ok bool) string {
	if !ok {
		return ""
	}
	if i < 0 || i >= len(row) {
		return ""
	}
	return row[i]
}

func toList(v string) []string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	parts := strings.FieldsFunc(v, func(r rune) bool { return r == ',' || r == ';' })
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func has(m map[string]int, k string) bool {
	_, ok := m[k]
	return ok
}

func isBlankRow(row []string) bool {
	for _, x := range row {
		if strings.TrimSpace(x) != "" {
			return false
		}
	}
	return true
}
