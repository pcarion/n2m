package n2m

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jomei/notionapi"
)

func (cms *Notion2Markdown) ConvertPageToMarkdown(pageId string, outputDirectory string) error {
	var metaData *MetaDataInformation
	var lines []string = make([]string, 0, 20)
	var err error

	fmt.Printf("ConvertPageToMarkdown: pageId=%s\n", pageId)
	cms.visitBlockChildren(pageId, func(blocks []notionapi.Block) error {
		// https://developers.notion.com/changelog/api-support-for-code-blocks-and-inline-databases
		for _, b := range blocks {
			switch b.GetType().String() {
			case "child_database":
				// meta information
				metaData, err = cms.extractMetaData(b)
				if err != nil {
					return err
				}
				fmt.Printf(">metaData>%#v\n", metaData)

			case "paragraph":
				paragraph, err := cms.parseParagraphBlock(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf(">paragraph>%s>\n\n%v\n\n", b.GetType().String(), paragraph)
				// get lines
				lines = append(lines, paragraph.Markdown)

			case "bulleted_list_item":
				paragraph, err := cms.parseBulletedListItemBlock(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf(">paragraph (bulleted)>%s>\n\n%v\n\n", b.GetType().String(), paragraph)
				// get lines
				lines = append(lines, paragraph.Markdown)

			case "image":
				paragraph, err := cms.parseImageBlock(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf(">paragraph (image)>%s>\n\n%v\n\n", b.GetType().String(), paragraph)

			default:
				blockData, err := json.Marshal(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("blockData:%s \n\n%s\n\n", b.GetType().String(), string(blockData))
				return fmt.Errorf("block type parsing not implemented for:%s", b.GetType().String())
			}
		}
		return nil
	})
	fmt.Printf(">metaData>%#v\n", metaData)

	return nil
}
