package n2m

import (
	"fmt"

	"github.com/jomei/notionapi"
)

type Notion2Markdown struct {
	client          *notionapi.Client
	debugMode       bool
	forceGeneration bool
}

func NewNotionToMarkdown(notionIntegrationToken string, debugMode, forceGeneration bool) (*Notion2Markdown, error) {
	client := notionapi.NewClient(notionapi.Token(notionIntegrationToken))

	return &Notion2Markdown{
		client:          client,
		debugMode:       debugMode,
		forceGeneration: forceGeneration,
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

	// check list of pages that actually need a refresh
	if !cms.forceGeneration {
		err = cms.markPagesToSkip(pages, outputDirectory)
		if err != nil {
			return err
		}
	}

	// loop through all the pages
	for ix, page := range pages {
		// test if we limit to one page
		if pageIndex >= 0 && ix != pageIndex {
			continue
		}
		// skip generation is marked as such
		if page.Skip {
			fmt.Printf("[%03d]: pageId=%s title=%s  *** SKIPPED ***\n", page.Index, page.Id, page.Title)
			continue
		}
		fmt.Printf("[%03d]: pageId=%s title=%s\n", page.Index, page.Id, page.Title)

		err = cms.convertPageToMarkdown(page, outputDirectory)
		if err != nil {
			return err
		}
	}
	return nil
}
