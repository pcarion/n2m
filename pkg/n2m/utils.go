package n2m

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ensureDir(dirName string) error {
	newpath := filepath.Join(".", dirName)
	err := os.MkdirAll(newpath, os.ModePerm)
	return err
}

func LogAsJson(b interface{}, what string) {
	blockData, err := json.Marshal(b)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("data as json:%s \n\n%s\n\n", what, string(blockData))

}
