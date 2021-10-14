package n2m

import (
	"github.com/jomei/notionapi"
)

type Notion2Markdown struct {
	client    *notionapi.Client
	debugMode bool
}

func NewNotionToMarkdown(notionIntegrationToken string, debugMode bool) (*Notion2Markdown, error) {
	client := notionapi.NewClient(notionapi.Token(notionIntegrationToken))

	return &Notion2Markdown{
		client:    client,
		debugMode: debugMode,
	}, nil
}

func (cms *Notion2Markdown) GenerateMardown(rootPageId string, outputDirectory string, pageIndex int) error {
	// create result directory
	err := ensureDir(outputDirectory)
	if err != nil {
		return err
	}
	pages, err := cms.extractListPages(rootPageId)

	if err != nil {
		return err
	}

	for ix, page := range pages {
		// test if we limit to one page
		if pageIndex >= 0 && ix != pageIndex {
			continue
		}
		err = cms.convertPageToMarkdown(page.Id, outputDirectory)
		if err != nil {
			return err
		}
	}
	return nil
}
