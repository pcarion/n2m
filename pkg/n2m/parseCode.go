package n2m

import (
	"fmt"

	"github.com/jomei/notionapi"
)

func (cms *Notion2Markdown) parseCode(block notionapi.Block) (*MarkdownCode, error) {
	codeBlock := block.(*notionapi.CodeBlock)

	md := cms.mdFromRichTexts(codeBlock.Code.Text)
	return &MarkdownCode{
		markdown: fmt.Sprintf("```\n%s\n```", md),
		language: codeBlock.Code.Language,
	}, nil
}
