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

type VisitorContext struct {
	metaData    *MetaDataInformation
	pageTitle   string
	mdBlocks    []MarkdownBlock
	mdImages    []ImageDescription
	imagesCount int
	cms         *Notion2Markdown
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
	var err error

	var visitorContext = VisitorContext{
		metaData:    nil,
		pageTitle:   "",
		mdBlocks:    make([]MarkdownBlock, 0, 20),
		mdImages:    make([]ImageDescription, 0, 5),
		imagesCount: 0,
		cms:         cms,
	}

	page, err := cms.client.Page.Get(context.Background(), notionapi.PageID(pageId))
	if err != nil {
		return err
	}
	visitorContext.pageTitle = getPageTitle(page)

	var visitorFunction = mkVisitor(&visitorContext)

	fmt.Printf("ConvertPageToMarkdown: pageId=%s title=%s\n", pageId, visitorContext.pageTitle)
	err = cms.visitBlockChildren(pageId, visitorFunction)
	// test result of "visitor"
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}
	fmt.Printf(">metaData>%#v\n", visitorContext.metaData)

	// generate file
	err = cms.writeMarkdownFile(outputDirectory, visitorContext.metaData, visitorContext.mdBlocks)
	if err != nil {
		return err
	}

	downloadImages(outputDirectory, visitorContext.metaData.Slug, visitorContext.mdImages)

	return nil
}

func mkVisitor(context *VisitorContext) BlockChildrenVisitor {
	var err error

	var addLine = func(md string, mdType int, level int) {
		context.mdBlocks = append(context.mdBlocks, MarkdownBlock{
			mdType:   mdType,
			level:    level,
			markdown: md,
		})
	}

	var visitorFunction = func(blocks []notionapi.Block) error {
		// https://developers.notion.com/changelog/api-support-for-code-blocks-and-inline-databases
		for _, b := range blocks {
			var blockType = b.GetType().String()

			switch blockType {
			case "child_database":
				// meta information
				context.metaData, err = context.cms.extractMetaData(b, context.pageTitle)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">metaData>%#v\n", context.metaData)

			case "paragraph":
				paragraph, err := context.cms.parseParagraphBlock(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">paragraph>%s>\n\n%v\n\n", blockType, paragraph)
				// get lines
				addLine(paragraph.markdown, MdTypePara, 0)

			case "bulleted_list_item":
				paragraph, err := context.cms.parseBulletedListItemBlock(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">bulleted_list>%s>\n\n%v\n\n", blockType, paragraph)
				// get lines
				addLine(paragraph.markdown, MdTypeListItem, 0)

				LogAsJson(b, "@@ bulleted_list_item")

			case "heading_1":
				paragraph, err := context.cms.parseParagraphHeading1(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">heading 1>%s>\n\n%v\n\n", blockType, paragraph)
				// get lines
				addLine(paragraph.markdown, MdTypeHeader, 1)

			case "heading_2":
				paragraph, err := context.cms.parseParagraphHeading2(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">heading 2>%s>\n\n%v\n\n", blockType, paragraph)
				// get lines
				addLine(paragraph.markdown, MdTypeHeader, 2)

			case "code":
				paragraph, err := context.cms.parseCode(b)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">paragraph>%s>\n\n%v\n\n", blockType, paragraph)
				// get lines
				addLine(paragraph.markdown, MdTypeCode, 0)

			case "image":
				context.imagesCount++
				paragraph, err := context.cms.parseImageBlock(b, context.metaData.Slug, context.imagesCount)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return err
				}
				fmt.Printf(">paragraph (image)>%s>\n\n%v\n\n", blockType, paragraph)
				// get lines
				addLine(paragraph.markdown, MdTypeImage, 0)

				// store image
				context.mdImages = append(context.mdImages, ImageDescription{
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
	}
	return visitorFunction
}
