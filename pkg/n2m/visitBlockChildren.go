package n2m

import (
	"context"

	"github.com/jomei/notionapi"
)

type BlockChildrenVisitor func([]notionapi.Block) error

func (cms *Notion2Markdown) visitBlockChildren(pageId string, visitor BlockChildrenVisitor) error {
	var pagination = notionapi.Pagination{
		PageSize: 10,
	}
	var doContinueLoop = true
	var returnError error = nil

	for doContinueLoop {

		block, err := cms.client.Block.GetChildren(context.Background(), notionapi.BlockID(pageId), &pagination)
		if err != nil {
			returnError = err
			doContinueLoop = false
		}
		if block.HasMore {
			pagination.StartCursor = notionapi.Cursor(block.NextCursor)
		} else {
			doContinueLoop = false
		}
		err = visitor(block.Results)
		if err != nil {
			returnError = err
			doContinueLoop = false
		}
	}
	return returnError
}
