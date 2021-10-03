package blogcms

import (
	"context"
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
	var metaData *MetaDataInformation

	// https://developers.notion.com/changelog/api-support-for-code-blocks-and-inline-databases
	for _, b := range block.Results {
		switch b.GetType().String() {
		case "child_database":
			// meta information
			metaData, err = cms.extractMetaData(b)
			if err != nil {
				return err
			}
			fmt.Printf(">metaData>%#v\n", metaData)

		case "paragraph":
			paragraph, err := cms.parseParagraph(b)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf(">paragraph>%s>\n\n%v\n\n", b.GetType().String(), paragraph)

		default:
			return fmt.Errorf("block type parsing not implemented for:%s", b.GetType().String())
		}
	}

	return nil
}
