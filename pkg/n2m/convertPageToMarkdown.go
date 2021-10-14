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
	err = cms.visitBlockChildren(pageId, visitorFunction, 0)
	// test result of "visitor"
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}
	if cms.debugMode {
		fmt.Printf(">metaData>%#v\n", visitorContext.metaData)
	}

	// generate file
	err = cms.writeMarkdownFile(outputDirectory, visitorContext.metaData, visitorContext.mdBlocks)
	if err != nil {
		return err
	}

	downloadImages(outputDirectory, visitorContext.metaData.Slug, visitorContext.mdImages)

	return nil
}

func mkVisitor(context *VisitorContext) BlockVisitor {
	var err error

	var addLine = func(md string, mdType int, level int) {
		context.mdBlocks = append(context.mdBlocks, MarkdownBlock{
			mdType:   mdType,
			level:    level,
			markdown: md,
		})
	}

	var visitorFunction = func(block notionapi.Block, level int) error {
		// https://developers.notion.com/changelog/api-support-for-code-blocks-and-inline-databases
		var blockType = block.GetType().String()
		var debugMode = context.cms.debugMode

		switch blockType {
		case "child_database":
			// meta information
			context.metaData, err = context.cms.extractMetaData(block, context.pageTitle)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return err
			}
			if debugMode {
				fmt.Printf(">metaData>%#v\n", context.metaData)
			}

		case "paragraph":
			paragraph, err := context.cms.parseParagraphBlock(block)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return err
			}
			if debugMode {
				fmt.Printf(">paragraph>%s>\n\n%v\n\n", blockType, paragraph)
			}
			// get lines
			addLine(paragraph.markdown, MdTypePara, 0)

		case "bulleted_list_item":
			paragraph, err := context.cms.parseBulletedListItemBlock(block, level)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return err
			}
			if debugMode {
				fmt.Printf(">bulleted_list>%s>\n\n%v\n\n", blockType, paragraph)
			}
			// get lines
			addLine(paragraph.markdown, MdTypeListItem, 0)

		case "heading_1":
			paragraph, err := context.cms.parseParagraphHeading1(block)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return err
			}
			if debugMode {
				fmt.Printf(">heading 1>%s>\n\n%v\n\n", blockType, paragraph)
			}
			// get lines
			addLine(paragraph.markdown, MdTypeHeader, 1)

		case "heading_2":
			paragraph, err := context.cms.parseParagraphHeading2(block)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return err
			}
			if debugMode {
				fmt.Printf(">heading 2>%s>\n\n%v\n\n", blockType, paragraph)
			}
			// get lines
			addLine(paragraph.markdown, MdTypeHeader, 2)

		case "heading_3":
			paragraph, err := context.cms.parseParagraphHeading3(block)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return err
			}
			if debugMode {
				fmt.Printf(">heading 2>%s>\n\n%v\n\n", blockType, paragraph)
			}
			// get lines
			addLine(paragraph.markdown, MdTypeHeader, 2)

		case "code":
			paragraph, err := context.cms.parseCode(block)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return err
			}
			if debugMode {
				fmt.Printf(">paragraph>%s>\n\n%v\n\n", blockType, paragraph)
			}
			// get lines
			addLine(paragraph.markdown, MdTypeCode, 0)

		case "image":
			context.imagesCount++
			paragraph, err := context.cms.parseImageBlock(block, context.metaData.Slug, context.imagesCount, debugMode)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return err
			}
			if debugMode {
				fmt.Printf(">paragraph (image)>%s>\n\n%v\n\n", blockType, paragraph)
			}
			// get lines
			addLine(paragraph.markdown, MdTypeImage, 0)

			// store image
			context.mdImages = append(context.mdImages, ImageDescription{
				imageUrl:           paragraph.imageToDownloadUrl,
				imageLocalFileName: paragraph.imageLocalFileName,
			})

		default:
			blockData, err := json.Marshal(block)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return err
			}
			fmt.Printf("blockData:%s \n\n%s\n\n", blockType, string(blockData))
			return fmt.Errorf("block type parsing not implemented for:%s", blockType)
		}
		return nil
	}
	return visitorFunction
}
