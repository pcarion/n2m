package n2m

import (
	"fmt"

	"github.com/jomei/notionapi"
)

func (cms *Notion2Markdown) parseCode(block notionapi.Block) (*MarkdownCode, error) {
	codeBlock := block.(*notionapi.CodeBlock)
	language := codeBlock.Code.Language

	md := cms.mdFromRichTexts(codeBlock.Code.Text)
	return &MarkdownCode{
		markdown: fmt.Sprintf("```%s\n%s\n```", language, md),
		language: language,
	}, nil
}
