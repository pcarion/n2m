package n2m

import (
	"fmt"
	"os"
	"path/filepath"
)

func (cms *Notion2Markdown) writeMarkdownFile(outputDirectory string, metaData *MetaDataInformation, blocks []MarkdownBlock) error {
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
	f.WriteString(fmt.Sprintf("title: %q\n", metaData.Title))
	f.WriteString(fmt.Sprintf("slug: %q\n", metaData.Slug))
	f.WriteString(fmt.Sprintf("date: %q\n", metaData.Date.Format("2006-01-02")))
	if len(metaData.Tags) > 0 {
		f.WriteString("tags:\n")
		for _, t := range metaData.Tags {
			f.WriteString(fmt.Sprintf("  - %q\n", t))
		}
	}
	f.WriteString(fmt.Sprintf("description: %q\n", metaData.Description))
	if metaData.HasToc {
		f.WriteString("toc: true\n")
	} else {
		f.WriteString("toc: false\n")
	}
	if metaData.IsDraft {
		f.WriteString("draft: true\n")
	} else {
		f.WriteString("draft: false\n")
	}
	f.WriteString("\n")
	// notion meta information
	f.WriteString(fmt.Sprintf("notionPageId: %q\n", metaData.NotionPageId))
	f.WriteString(fmt.Sprintf("notionLastEditedTime: %q\n", metaData.NotionLastEditedTime))
	f.WriteString("---\n\n")
	for ix, block := range blocks {
		f.WriteString(block.markdown)
		f.WriteString("\n")
		// check if we need to add an extra line
		if block.mdType != MdTypeListItem || (ix < len(blocks)-1 && blocks[ix+1].mdType != MdTypeListItem) {
			f.WriteString("\n")
		}
	}
	return nil
}
