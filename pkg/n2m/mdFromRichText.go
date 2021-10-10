package n2m

import (
	"fmt"

	"github.com/jomei/notionapi"
)

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

func (cms *Notion2Markdown) mdFromRichText(richText notionapi.RichText) string {
	prefix, suffix := parseAnnotations(richText.Annotations)
	var md = prefix + richText.PlainText + suffix
	if richText.Href != "" {
		return fmt.Sprintf("[%s](%s)", md, richText.Href)
	}
	return md
}

func (cms *Notion2Markdown) mdFromRichTexts(richTexts []notionapi.RichText) string {
	var md = ""
	for _, richText := range richTexts {
		md = md + cms.mdFromRichText(richText)
	}
	return md
}
