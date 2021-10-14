package n2m

import (
	"context"

	"github.com/jomei/notionapi"
)

type withChidren struct {
	hasChildren bool
	blockID     string
}

type BlockVisitor func(notionapi.Block, int) error

func (cms *Notion2Markdown) visitBlockChildren(blockId string, visitor BlockVisitor, level int) error {
	var pagination = notionapi.Pagination{
		PageSize: 10,
	}
	var doContinueLoop = true
	var returnError error = nil

	for doContinueLoop {

		children, err := cms.client.Block.GetChildren(context.Background(), notionapi.BlockID(blockId), &pagination)
		if err != nil {
			returnError = err
			doContinueLoop = false
		} else {
			if children.HasMore {
				pagination.StartCursor = notionapi.Cursor(children.NextCursor)
			} else {
				doContinueLoop = false
			}
			for _, block := range children.Results {
				err = visitor(block, level)
				if err != nil {
					returnError = err
					doContinueLoop = false
				} else {
					// check if there are children nodes
					var withBlockChildren *withChidren = nil
					if block.GetType().String() == "paragraph" {
						paragraphBlock := block.(*notionapi.ParagraphBlock)
						withBlockChildren = &withChidren{
							hasChildren: paragraphBlock.HasChildren,
							blockID:     paragraphBlock.ID.String(),
						}
					} else if block.GetType().String() == "bulleted_list_item" {
						bulletedListItemBlock := block.(*notionapi.BulletedListItemBlock)
						withBlockChildren = &withChidren{
							hasChildren: bulletedListItemBlock.HasChildren,
							blockID:     bulletedListItemBlock.ID.String(),
						}
					} else if block.GetType().String() == "heading_1" {
						headingBlock := block.(*notionapi.Heading1Block)
						withBlockChildren = &withChidren{
							hasChildren: headingBlock.HasChildren,
							blockID:     headingBlock.ID.String(),
						}
					} else if block.GetType().String() == "heading_2" {
						headingBlock := block.(*notionapi.Heading2Block)
						withBlockChildren = &withChidren{
							hasChildren: headingBlock.HasChildren,
							blockID:     headingBlock.ID.String(),
						}
					} else if block.GetType().String() == "heading_3" {
						headingBlock := block.(*notionapi.Heading3Block)
						withBlockChildren = &withChidren{
							hasChildren: headingBlock.HasChildren,
							blockID:     headingBlock.ID.String(),
						}
					}

					if withBlockChildren != nil {
						// we iterate through children blocks
						err = cms.visitBlockChildren(withBlockChildren.blockID, visitor, level+1)
						if err != nil {
							returnError = err
							doContinueLoop = false
						}
					}
				}
			}
		}
	}
	return returnError
}
