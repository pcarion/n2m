package n2m

import (
	"github.com/jomei/notionapi"
)

type Notion2Markdown struct {
	client *notionapi.Client
}

func NewNotionToMarkdown(notionIntegrationToken string) (*Notion2Markdown, error) {
	client := notionapi.NewClient(notionapi.Token(notionIntegrationToken))

	return &Notion2Markdown{
		client: client,
	}, nil
}

func (cms *Notion2Markdown) GenerateContent(rootPageId string, outputDirectory string) error {
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
		// TODO: test to limit
		if ix != 0 {
			continue
		}
		err = cms.ConvertPageToMarkdown(page.Id)
		if err != nil {
			return err
		}
	}
	return nil
}
