package n2m

import (
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
	return prefix + richText.PlainText + suffix
}

func (cms *Notion2Markdown) mdFromRichTexts(richTexts []notionapi.RichText) string {
	var md = ""
	for _, richText := range richTexts {
		md = md + cms.mdFromRichText(richText)
	}
	return md
}
