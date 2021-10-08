package n2m

import (
	"os"
	"path/filepath"
)

func ensureDir(dirName string) error {
	newpath := filepath.Join(".", dirName)
	err := os.MkdirAll(newpath, os.ModePerm)
	return err
}
