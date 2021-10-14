package n2m

import (
	"strings"

	"github.com/jomei/notionapi"
)

func (cms *Notion2Markdown) parseParagraphBlock(block notionapi.Block) (*MarkdownParagraph, error) {
	paragraphBlock := block.(*notionapi.ParagraphBlock)

	md := cms.mdFromRichTexts(paragraphBlock.Paragraph.Text)
	return &MarkdownParagraph{
		markdown: md,
	}, nil
}

func (cms *Notion2Markdown) parseBulletedListItemBlock(block notionapi.Block, level int) (*MarkdownParagraph, error) {
	bulletedListItemBlock := block.(*notionapi.BulletedListItemBlock)
	var prefix = strings.Repeat("  ", level)

	listItem := bulletedListItemBlock.BulletedListItem
	md := cms.mdFromRichTexts(listItem.Text)
	return &MarkdownParagraph{
		markdown: prefix + "* " + md,
	}, nil
}

func (cms *Notion2Markdown) parseParagraphHeading1(block notionapi.Block) (*MarkdownParagraph, error) {
	headingBlock := block.(*notionapi.Heading1Block)

	md := cms.mdFromRichTexts(headingBlock.Heading1.Text)
	return &MarkdownParagraph{
		markdown: "# " + md,
	}, nil
}

func (cms *Notion2Markdown) parseParagraphHeading2(block notionapi.Block) (*MarkdownParagraph, error) {
	headingBlock := block.(*notionapi.Heading2Block)

	md := cms.mdFromRichTexts(headingBlock.Heading2.Text)
	return &MarkdownParagraph{
		markdown: "## " + md,
	}, nil
}

func (cms *Notion2Markdown) parseParagraphHeading3(block notionapi.Block) (*MarkdownParagraph, error) {
	headingBlock := block.(*notionapi.Heading3Block)

	md := cms.mdFromRichTexts(headingBlock.Heading3.Text)
	return &MarkdownParagraph{
		markdown: "## " + md,
	}, nil
}
