package n2m

import (
	"github.com/jomei/notionapi"
)

type MarkdownParagraph struct {
	Markdown string
}

func (cms *Notion2Markdown) parseParagraphBlock(block notionapi.Block) (*MarkdownParagraph, error) {
	paragraphBlock := block.(*notionapi.ParagraphBlock)

	md := cms.mdFromRichTexts(paragraphBlock.Paragraph.Text)
	return &MarkdownParagraph{
		Markdown: md,
	}, nil
}

func (cms *Notion2Markdown) parseBulletedListItemBlock(block notionapi.Block) (*MarkdownParagraph, error) {
	bulletedListItemBlock := block.(*notionapi.BulletedListItemBlock)
	listItem := bulletedListItemBlock.BulletedListItem
	md := cms.mdFromRichTexts(listItem.Text)
	return &MarkdownParagraph{
		Markdown: "* " + md,
	}, nil
}
