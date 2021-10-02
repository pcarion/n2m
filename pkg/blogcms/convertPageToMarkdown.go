package blogcms

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
)

func (cms *Blogcms) ConvertPageToMarkdown(pageId string) error {
	fmt.Printf("ConvertPageToMarkdown: pageId=%s\n", pageId)
	block, err := cms.client.Block.GetChildren(context.Background(), notionapi.BlockID(pageId), nil)
	if err != nil {
		return err
	}
	fmt.Printf("@@ block: %#v", block)

	// blockData, err := json.Marshal(block)
	// if err != nil {
	// 	fmt.Printf("error: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("blockData: %s\n", string(blockData))
	return nil
}
