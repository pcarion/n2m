package n2m

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gernest/front"
)

func getListOfIndexMdFiles(outputDirectory string) ([]string, error) {
	var result = make([]string, 0, 10)

	// serach for index.md files
	err := filepath.Walk(filepath.Join(".", outputDirectory),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == "index.md" {
				result = append(result, path)
			}
			return nil
		})

	return result, err
}

func getExistingMdInformation(files []string) (map[string]string, error) {
	var result = make(map[string]string)
	m := front.NewMatter()
	m.Handle("+++", front.JSONHandler)
	m.Handle("---", front.YAMLHandler)

	for _, indexMdFile := range files {
		b, err := ioutil.ReadFile(indexMdFile)
		if err != nil {
			fmt.Printf("error: %#v", err)
			return nil, err
		}
		front, _, err := m.Parse(strings.NewReader(string(b)))
		if err != nil {
			fmt.Printf("[%s] error: %#v", indexMdFile, err)
			return nil, err
		}
		frontNotionPageId := front["notionPageId"]
		if frontNotionPageId != nil {
			pageId := fmt.Sprintf("%v", frontNotionPageId)
			if pageId != "" {
				lastEditedTime := fmt.Sprintf("%v", front["notionLastEditedTime"])
				result[pageId] = lastEditedTime
			}
		}
	}
	return result, nil
}

func (cms *Notion2Markdown) markPagesToSkip(pages []*CmsPageDescription, outputDirectory string) error {
	fmt.Printf("\n\n## TODO: filter page that need a refresh\n\n")

	indexMdFiles, err := getListOfIndexMdFiles(outputDirectory)
	if err != nil {
		fmt.Printf("error: %#v", err)
		return err
	}

	existing, err := getExistingMdInformation(indexMdFiles)
	if err != nil {
		fmt.Printf("error: %#v", err)
		return err
	}

	// fmt.Printf("@@ existing notion ids: %#v", existing)

	// check pages that don't need to be regenerated
	for _, page := range pages {
		//		fmt.Printf("\n@@ check skip pageId=[%s], editedTime=[%s] existing[%s] title=[%s]\n", page.Id, page.LastEditedTime, existing[page.Id], page.Title)
		if existing[page.Id] != "" && existing[page.Id] == page.LastEditedTime {
			// fmt.Printf("@@@ skipping\n")
			page.Skip = true
		}
	}
	return nil
}
