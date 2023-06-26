package strs

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Find(src string, pattern string, format string) string {
	re := regexp.MustCompile(pattern)
	dst := re.FindString(src)
	if format == "" || format == "$0" || dst == "" {
		return dst
	}
	return re.ReplaceAllString(dst, format)
}

func Json(n any) string {
	data, _ := json.MarshalIndent(n, "", "  ")
	return string(data)
}

func Unwrap(s string, prefix, suffix string) string {
	if prefix != "" {
		s = strings.TrimPrefix(s, prefix)
	}
	if suffix != "" {
		s = strings.TrimSuffix(s, suffix)
	}
	return s
}

func Clean(s string) string {
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}

func Text(s *goquery.Selection) string {
	return Clean(s.Text())
}

func Sub(src string, segs ...string) string {
	if len(segs) == 0 {
		return ""
	}

	var suffix string
	if len(segs) > 1 {
		suffix = segs[len(segs)-1]
		segs = segs[:len(segs)-1]
	}

	for _, prefix := range segs {
		if i := strings.Index(src, prefix); i != -1 {
			src = src[i+len(prefix):]
		} else {
			return ""
		}
	}

	if suffix != "" {
		if i := strings.Index(src, suffix); i != -1 {
			src = src[:i]
		} else {
			return ""
		}
	}

	return src
}
