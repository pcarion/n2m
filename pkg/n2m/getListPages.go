package n2m

import (
	"time"

	"github.com/jomei/notionapi"
)

func (cms *Notion2Markdown) extractListPages(pageId string) ([]CmsPageDescription, error) {
	result := make([]CmsPageDescription, 0, 10)
	var index = 0

	cms.visitBlockChildren(pageId, func(block notionapi.Block, level int) error {
		if block.GetType().String() == notionapi.BlockTypeChildPage.String() {
			childPage := block.(*notionapi.ChildPageBlock)
			id := childPage.ID.String()
			title := childPage.ChildPage.Title
			lastEditedTime := "none"
			if childPage.LastEditedTime != nil {
				lastEditedTime = childPage.LastEditedTime.Format(time.RFC3339)
			}
			result = append(result, CmsPageDescription{
				Id:             id,
				Title:          title,
				LastEditedTime: lastEditedTime,
				Index:          index,
			})
			index++
		}
		return nil
	}, 0)
	return result, nil
}
