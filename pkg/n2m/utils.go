package n2m

import (
	"fmt"
	"os"
	"path/filepath"
)

func ensureDir(dirName string) error {
	newpath := filepath.Join(".", dirName)
	fmt.Printf("Creating directory:%s\n", newpath)
	err := os.MkdirAll(newpath, os.ModePerm)

	return err
}
