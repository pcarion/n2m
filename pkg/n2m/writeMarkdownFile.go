package n2m

import (
	"fmt"
	"os"
	"path/filepath"
)

func (cms *Notion2Markdown) writeMarkdownFile(outputDirectory string, metaData *MetaDataInformation, lines []string) error {
	dir := filepath.Join(outputDirectory, metaData.Slug)
	err := ensureDir(dir)
	if err != nil {
		return err
	}
	fileName := filepath.Join(dir, "index.md")

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString("---\n")
	f.WriteString(fmt.Sprintf("title: %s\n", metaData.Title))
	f.WriteString(fmt.Sprintf("slug: %s\n", metaData.Slug))
	f.WriteString(fmt.Sprintf("date: %s\n", metaData.Date.Format("2006/01/02")))
	f.WriteString("tags: ")
	for ix, t := range metaData.Tags {
		if ix > 0 {
			f.WriteString(", ")
		}
		f.WriteString(t)
	}
	f.WriteString("\n")
	f.WriteString(fmt.Sprintf("excerpt: %s\n", metaData.Excerpt))
	f.WriteString("---\n\n")
	for _, line := range lines {
		f.WriteString(line)
		f.WriteString("\n\n")
	}
	return nil
}
