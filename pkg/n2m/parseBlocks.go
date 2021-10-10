package n2m

import (
	"github.com/jomei/notionapi"
)

type MarkdownParagraph struct {
	markdown string
}

func (cms *Notion2Markdown) parseParagraphBlock(block notionapi.Block) (*MarkdownParagraph, error) {
	paragraphBlock := block.(*notionapi.ParagraphBlock)

	md := cms.mdFromRichTexts(paragraphBlock.Paragraph.Text)
	return &MarkdownParagraph{
		markdown: md,
	}, nil
}

func (cms *Notion2Markdown) parseBulletedListItemBlock(block notionapi.Block) (*MarkdownParagraph, error) {
	bulletedListItemBlock := block.(*notionapi.BulletedListItemBlock)
	listItem := bulletedListItemBlock.BulletedListItem
	md := cms.mdFromRichTexts(listItem.Text)
	return &MarkdownParagraph{
		markdown: "* " + md,
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
