package n2m

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func downloadImage(imgUrl string, outFileName string) error {
	fmt.Printf("@@ downloadImage url=%s, path=%s\n", imgUrl, outFileName)
	resp, err := http.Get(imgUrl)
	if err != nil {
		return fmt.Errorf("couldn't download image: %s", err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(outFileName)
	if err != nil {
		return fmt.Errorf("couldn't create image file: %s", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return err
	}

	return nil
}

func downloadImages(outputDirectory string, slug string, mdImages []ImageDescription) error {
	for _, image := range mdImages {
		err := downloadImage(image.imageUrl, filepath.Join(outputDirectory, slug, image.imageLocalFileName))
		if err != nil {
			return err
		}
	}
	return nil
}
