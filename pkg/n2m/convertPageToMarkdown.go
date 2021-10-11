package n2m

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jomei/notionapi"
)

const (
	MdTypePara = iota
	MdTypeHeader
	MdTypeImage
	MdTypeListItem
	MdTypeCode
)

// description of a generated markdown block
type MarkdownBlock struct {
	mdType   int
	level    int
	markdown string
}

type ImageDescription struct {
	imageUrl           string
	imageLocalFileName string
}

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
	var mdBlocks []MarkdownBlock = make([]MarkdownBlock, 0, 20)
	var mdImages []ImageDescription = make([]ImageDescription, 0, 5)
	var err error
	var imagesCount = 0

	var addLine = func(md string, mdType int, level int) {
		mdBlocks = append(mdBlocks, MarkdownBlock{
			mdType:   mdType,
			level:    level,
			markdown: md,
		})
	}

	page, err := cms.client.Page.Get(context.Background(), notionapi.PageID(pageId))
	if err != nil {
		return err
	}
	pageTitle = getPageTitle(page)

	fmt.Printf("ConvertPageToMarkdown: pageId=%s title=%s\n", pageId, pageTitle)
	err = cms.visitBlockChildren(pageId, func(blocks []notionapi.Block) error {
		// https://developers.notion.com/changelog/api-support-for-code-blocks-and-inline-databases
		for _, b := range blocks {
			var blockType = b.GetType().String()

			switch blockType {
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
				fmt.Printf(">paragraph>%s>\n\n%v\n\n", blockType, paragraph)
				// get lines
				addLine(paragraph.markdown, MdTypePara, 0)

			case "bulleted_list_item":
				paragraph, err := cms.parseBulletedListItemBlock(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">bulleted_list>%s>\n\n%v\n\n", blockType, paragraph)
				// get lines
				addLine(paragraph.markdown, MdTypeListItem, 0)

			case "heading_1":
				paragraph, err := cms.parseParagraphHeading1(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">heading 1>%s>\n\n%v\n\n", blockType, paragraph)
				// get lines
				addLine(paragraph.markdown, MdTypeHeader, 1)

			case "heading_2":
				paragraph, err := cms.parseParagraphHeading2(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">heading 2>%s>\n\n%v\n\n", blockType, paragraph)
				// get lines
				addLine(paragraph.markdown, MdTypeHeader, 2)

			case "code":
				paragraph, err := cms.parseCode(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">paragraph>%s>\n\n%v\n\n", blockType, paragraph)
				// get lines
				addLine(paragraph.markdown, MdTypeCode, 0)

			case "image":
				imagesCount++
				paragraph, err := cms.parseImageBlock(b, metaData.Slug, imagesCount)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">paragraph (image)>%s>\n\n%v\n\n", blockType, paragraph)
				// get lines
				addLine(paragraph.markdown, MdTypeImage, 0)

				// store image
				mdImages = append(mdImages, ImageDescription{
					imageUrl:           paragraph.imageToDownloadUrl,
					imageLocalFileName: paragraph.imageLocalFileName,
				})

			default:
				blockData, err := json.Marshal(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf("blockData:%s \n\n%s\n\n", blockType, string(blockData))
				return fmt.Errorf("block type parsing not implemented for:%s", blockType)
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
	err = cms.writeMarkdownFile(outputDirectory, metaData, mdBlocks)
	if err != nil {
		return err
	}

	downloadImages(outputDirectory, metaData.Slug, mdImages)

	return nil
}
