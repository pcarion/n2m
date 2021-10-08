package n2m

import (
	"fmt"

	"github.com/jomei/notionapi"
)

func (cms *Notion2Markdown) parseImageBlock(block notionapi.Block) (*MarkdownParagraph, error) {
	imageBlock := block.(*notionapi.ImageBlock)
	image := imageBlock.Image

	// TODO: handle external URLs
	fmt.Printf(">image block>%#v\n", imageBlock)
	fmt.Printf(">image URL>\n%s\n", image.File.URL)
	return &MarkdownParagraph{
		Markdown: "",
	}, nil
}
