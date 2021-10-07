package blogcms

import (
	"github.com/jomei/notionapi"
)

type Blogcms struct {
	client *notionapi.Client
}

func NewBlocgCms(notionIntegrationToken string) (*Blogcms, error) {
	client := notionapi.NewClient(notionapi.Token(notionIntegrationToken))

	return &Blogcms{
		client: client,
	}, nil
}

func (cms *Blogcms) GenerateContent(rootPageId string, outputDirectory string) error {
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
