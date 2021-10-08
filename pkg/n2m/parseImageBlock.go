package n2m

import (
	"fmt"

	"github.com/jomei/notionapi"
)

type MarkdownImage struct {
	caption  string
	imageUrl string
}

func (cms *Notion2Markdown) parseImageBlock(block notionapi.Block) (*MarkdownImage, error) {
	imageBlock := block.(*notionapi.ImageBlock)
	image := imageBlock.Image

	// TODO: handle external URLs
	fmt.Printf(">image block>%#v\n", imageBlock)
	fmt.Printf(">image URL>\n%s\n", image.File.URL)
	var caption = ""
	if len(imageBlock.Image.Caption) > 0 {
		caption = imageBlock.Image.Caption[0].PlainText
	}
	return &MarkdownImage{
		imageUrl: image.File.URL,
		caption:  caption,
	}, nil
}
