package blogcms

import (
	"context"

	"github.com/jomei/notionapi"
)

type CmsPageDescription struct {
	Id    string
	Title string
}

func (cms *Blogcms) ExtractListPages(pageId string) ([]CmsPageDescription, error) {
	page, err := cms.client.Page.Get(context.Background(), notionapi.PageID(pageId))
	if err != nil {
		return nil, err
	}

	block, err := cms.client.Block.GetChildren(context.Background(), notionapi.BlockID(page.ID), nil)
	if err != nil {
		return nil, err
	}

	result := make([]CmsPageDescription, 0, len(block.Results))
	for _, b := range block.Results {
		if b.GetType().String() == "child_page" {
			childPage := b.(*notionapi.ChildPageBlock)
			id := childPage.ID.String()
			title := childPage.ChildPage.Title
			result = append(result, CmsPageDescription{
				Id:    id,
				Title: title,
			})
		}
	}
	return result, nil
}
