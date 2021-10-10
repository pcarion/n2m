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
	f.WriteString(fmt.Sprintf("date: %s\n", metaData.Date.Format("2006-01-02")))
	if len(metaData.Tags) > 0 {
		f.WriteString("tags:\n")
		for _, t := range metaData.Tags {
			f.WriteString(fmt.Sprintf("  - %s\n", t))
		}
	}
	f.WriteString(fmt.Sprintf("description: %s\n", metaData.Description))
	f.WriteString("toc: true\n")
	if metaData.IsDraft {
		f.WriteString("draft: true\n")
	} else {
		f.WriteString("draft: false\n")
	}
	f.WriteString("---\n\n")
	for _, line := range lines {
		f.WriteString(line)
		f.WriteString("\n\n")
	}
	return nil
}
