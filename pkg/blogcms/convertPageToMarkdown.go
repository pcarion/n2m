package blogcms

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jomei/notionapi"
)

func (cms *Blogcms) ConvertPageToMarkdown(pageId string) error {
	fmt.Printf("ConvertPageToMarkdown: pageId=%s\n", pageId)
	block, err := cms.client.Block.GetChildren(context.Background(), notionapi.BlockID(pageId), nil)
	if err != nil {
		return err
	}

	// https://developers.notion.com/changelog/api-support-for-code-blocks-and-inline-databases
	for _, b := range block.Results {
		if b.GetType().String() == "child_database" {
			metaData, err := cms.extractMetaData(b)
			if err != nil {
				return err
			}
			fmt.Printf(">metaData>%#v\n", metaData)
		}
	}

	blockData, err := json.Marshal(block)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("blockData: \n\n%s\n\n", string(blockData))
	return nil
}
