package blogcms

import (
	"fmt"

	"github.com/jomei/notionapi"
)

type MarkdownParagraph struct {
	Markdown string
}

func parseAnnotations(annotations *notionapi.Annotations) (string, string) {
	if annotations == nil {
		return "", ""
	}
	var prefix = ""
	var suffix = ""

	if annotations.Bold {
		prefix = prefix + "**"
		suffix = "**" + suffix
	}

	if annotations.Italic {
		prefix = prefix + "_"
		suffix = "_" + suffix
	}

	if annotations.Strikethrough {
		prefix = prefix + "~~"
		suffix = "~~" + suffix
	}

	if annotations.Code {
		prefix = prefix + "`"
		suffix = "`" + suffix
	}

	return prefix, suffix
}

func (cms *Blogcms) parseParagraph(block notionapi.Block) (*MarkdownParagraph, error) {
	paragraphBlock := block.(*notionapi.ParagraphBlock)
	fmt.Printf(">>paragraph block>%#v\n", paragraphBlock)

	var md = ""
	for _, richText := range paragraphBlock.Paragraph.Text {
		prefix, suffix := parseAnnotations(richText.Annotations)
		md = md + prefix + richText.PlainText + suffix
	}
	return &MarkdownParagraph{
		Markdown: md,
	}, nil
}
