package n2m

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jomei/notionapi"
)

func getPageTitle(page *notionapi.Page) string {
	for _, prop := range page.Properties {
		if prop.GetType() == "title" {
			titleProp := prop.(*notionapi.TitleProperty)
			title := ""
			for _, t := range titleProp.Title {
				title += t.PlainText
			}
			return title
		}
	}
	return ""
}

func (cms *Notion2Markdown) convertPageToMarkdown(pageId string, outputDirectory string) error {
	var metaData *MetaDataInformation
	var pageTitle = ""
	var lines []string = make([]string, 0, 20)
	var err error
	var imagesCount = 0

	page, err := cms.client.Page.Get(context.Background(), notionapi.PageID(pageId))
	if err != nil {
		return err
	}
	pageTitle = getPageTitle(page)

	fmt.Printf("ConvertPageToMarkdown: pageId=%s title=%s\n", pageId, pageTitle)
	err = cms.visitBlockChildren(pageId, func(blocks []notionapi.Block) error {
		fmt.Printf("@@@ in visitor [%d]... \n", len(blocks))
		// https://developers.notion.com/changelog/api-support-for-code-blocks-and-inline-databases
		for _, b := range blocks {
			fmt.Printf("@@@ block [%s]... \n", b.GetType().String())
			switch b.GetType().String() {
			case "child_database":
				// meta information
				metaData, err = cms.extractMetaData(b, pageTitle)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">metaData>%#v\n", metaData)

			case "paragraph":
				paragraph, err := cms.parseParagraphBlock(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">paragraph>%s>\n\n%v\n\n", b.GetType().String(), paragraph)
				// get lines
				lines = append(lines, paragraph.markdown)

			case "bulleted_list_item":
				paragraph, err := cms.parseBulletedListItemBlock(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">bulleted_list>%s>\n\n%v\n\n", b.GetType().String(), paragraph)
				// get lines
				lines = append(lines, paragraph.markdown)

			case "heading_1":
				paragraph, err := cms.parseParagraphHeading1(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">heading 1>%s>\n\n%v\n\n", b.GetType().String(), paragraph)
				// get lines
				lines = append(lines, paragraph.markdown)

			case "code":
				paragraph, err := cms.parseCode(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">paragraph>%s>\n\n%v\n\n", b.GetType().String(), paragraph)
				// get lines
				lines = append(lines, paragraph.markdown)

			case "image":
				imagesCount++
				paragraph, err := cms.parseImageBlock(b, metaData.Slug, imagesCount)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">paragraph (image)>%s>\n\n%v\n\n", b.GetType().String(), paragraph)
				// get lines
				lines = append(lines, paragraph.markdown)

			default:
				blockData, err := json.Marshal(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf("blockData:%s \n\n%s\n\n", b.GetType().String(), string(blockData))
				return fmt.Errorf("block type parsing not implemented for:%s", b.GetType().String())
			}
		}
		return nil
	})
	// test result of "visitor"
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}
	fmt.Printf(">metaData>%#v\n", metaData)

	// generate file
	err = cms.writeMarkdownFile(outputDirectory, metaData, lines)
	if err != nil {
		return err
	}

	return nil
}
