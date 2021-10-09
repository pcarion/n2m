package n2m

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/jomei/notionapi"
)

type MarkdownImage struct {
	caption            string
	imageToDownloadUrl string
	imageLocalFileName string
	markdown           string
}

func extractExtensionFromFileUrl(urlToParse string) (string, error) {
	u, err := url.Parse(urlToParse)
	if err != nil {
		return "", err
	}
	extension := filepath.Ext(u.Path)
	return extension, nil
}

func (cms *Notion2Markdown) parseImageBlock(block notionapi.Block, slug string, imageIndex int) (*MarkdownImage, error) {
	imageBlock := block.(*notionapi.ImageBlock)
	image := imageBlock.Image

	// TODO: handle external URLs
	fmt.Printf(">image block>%#v\n", imageBlock)
	fmt.Printf(">image URL>\n%s\n", image.File.URL)
	var caption = ""
	if len(imageBlock.Image.Caption) > 0 {
		caption = imageBlock.Image.Caption[0].PlainText
	}
	extension, err := extractExtensionFromFileUrl(image.File.URL)
	if err != nil {
		return nil, err
	}
	imageLocalFileName := fmt.Sprintf("%s-%03d%s", slug, imageIndex, extension)
	return &MarkdownImage{
		imageToDownloadUrl: image.File.URL,
		imageLocalFileName: imageLocalFileName,
		caption:            caption,
		markdown:           fmt.Sprintf("![](%s)", imageLocalFileName),
	}, nil
}
