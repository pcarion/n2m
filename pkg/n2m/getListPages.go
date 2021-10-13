package n2m

import (
	"github.com/jomei/notionapi"
)

type CmsPageDescription struct {
	Id    string
	Title string
}

func (cms *Notion2Markdown) extractListPages(pageId string) ([]CmsPageDescription, error) {
	result := make([]CmsPageDescription, 0, 10)

	cms.visitBlockChildren(pageId, func(block notionapi.Block, level int) error {
		if block.GetType().String() == "child_page" {
			childPage := block.(*notionapi.ChildPageBlock)
			id := childPage.ID.String()
			title := childPage.ChildPage.Title
			result = append(result, CmsPageDescription{
				Id:    id,
				Title: title,
			})
		}
		return nil
	}, 0)
	return result, nil
}
