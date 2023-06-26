package name

import (
	"os"
	"path/filepath"
)

func Link(dstRoot, srcPath string) (dstPath string, err error) {
	dstPath = CleanPath(filepath.Join(dstRoot, filepath.Base(srcPath)))
	if err = os.MkdirAll(dstRoot, 0755); err != nil {
		return
	}
	err = os.Link(srcPath, dstPath)
	return
}
