package blogcms

import (
	"github.com/jomei/notionapi"
)

type CmsPageDescription struct {
	Id    string
	Title string
}

func (cms *Blogcms) extractListPages(pageId string) ([]CmsPageDescription, error) {
	result := make([]CmsPageDescription, 0, 10)

	cms.visitBlockChildren(pageId, func(blocks []notionapi.Block) error {
		for _, b := range blocks {
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
		return nil
	})
	return result, nil
}
