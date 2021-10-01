package blogcms

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
)

func (cms *Blogcms) ExtractListPages(pageId string) error {
	page, err := cms.client.Page.Get(context.Background(), notionapi.PageID(pageId))
	if err != nil {
		return err
	}

	block, err := cms.client.Block.GetChildren(context.Background(), notionapi.BlockID(page.ID), nil)
	if err != nil {
		return err
	}

	for _, b := range block.Results {
		if b.GetType().String() == "child_page" {
			childPage := b.(*notionapi.ChildPageBlock)
			fmt.Printf(">child page>:%#v\n", childPage.ChildPage)
		}
	}
	return nil
}
