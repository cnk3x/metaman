package name

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func Walk(srcRoot string, walkFn func(path string, info fs.FileInfo) (err error), minSize int) (err error) {
	if srcRoot, err = filepath.Abs(srcRoot); err != nil {
		return
	}

	err = filepath.Walk(srcRoot, func(path string, info fs.FileInfo, err error) error {
		if err != nil || !allowFile(path, info, minSize) {
			return err
		}
		return walkFn(path, info)
	})

	return
}

func allowFile(path string, info fs.FileInfo, minSize int) bool {
	allowExts := []string{".mkv", ".mp4", ".rmvb"}
	minBytes := int64(minSize) << 20
	allowExt := func(path string) bool {
		ext := strings.ToLower(filepath.Ext(path))
		for _, allowExt := range allowExts {
			if ext == allowExt {
				return true
			}
		}
		return false
	}
	return info != nil && info.Mode().IsRegular() && allowExt(path) && info.Size() >= minBytes
}
