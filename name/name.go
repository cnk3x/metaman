package name

import (
	"path/filepath"
	"regexp"
	"strings"
)

func repl(src, pattern, repl string) string {
	return regexp.MustCompile(pattern).ReplaceAllString(src, repl)
}

func redel(src, pattern string) string {
	return repl(src, pattern, "")
}

func Clean(name string) string {
	// name = strings.ReplaceAll(name, "dy.ygdy8.com", "")
	// name = strings.ReplaceAll(name, "www.ygdy8.com", "")
	// name = strings.ReplaceAll(name, "ygdy8.com", "")
	// name = strings.ReplaceAll(name, "dygod.org", "")
	// name = strings.ReplaceAll(name, "www.dy2018.com", "")
	// name = strings.ReplaceAll(name, "dy2018.com", "")
	// name = strings.ReplaceAll(name, "dy2018", "")
	// name = strings.ReplaceAll(name, "电影天堂", "")
	// name = strings.ReplaceAll(name, "阳光电影", "")
	name = redel(name, `(www|dy|[\s\.-_])(\w+\.)?(\w+\.)?\w+\.(com|org|net|cn|vip|cc)`)
	name = redel(name, `(电影天堂|阳光电影)`)
	name = redel(name, `(蓝光|高清)`)

	name = redel(name, `([中英国粤韩英日双三]+语)?[中英国粤韩英日双三]+字`)
	name = redel(name, `[中英双国粤韩英日]+[语文](字幕)?`)
	name = redel(name, `中字`)

	name = repl(name, `[\s\.\-_【\(\[](20[012]\d|19[5-9]\d)[s\]\)】\s\.\-_]`, " ($1)")
	name = repl(name, `[\.\s](20[012]\d|19[5-9]\d)$`, " ($1)")

	name = redel(name, `(?i)(BD)?\d{3,4}(p|P)`)
	name = redel(name, `(?i)(DTSHD-MA|E?AAC|E?AC3|DTS|IMAX)`)
	name = redel(name, `(?i)(IQY|WEB-DL|H265|60fps|-Dream)`)
	name = redel(name, `(?i)(iTunes|DDP5\.1|Atmos|DV|H\.?26[54])`)
	name = redel(name, `(?i)(HD|BD)`)
	name = redel(name, `\[[^\]]*\]`)
	name = repl(name, `[\s\.\-\_\[\]]+`, " ")

	return strings.TrimSpace(name)
}

func CleanPath(path string) string {
	dir, name := filepath.Split(path)
	ext := filepath.Ext(name)
	name = strings.TrimSuffix(name, ext)
	name = Clean(name)
	return filepath.Join(dir, name+strings.ToLower(ext))
}
