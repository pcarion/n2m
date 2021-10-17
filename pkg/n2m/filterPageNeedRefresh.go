package n2m

import "fmt"

func (cms *Notion2Markdown) filterPageNeedRefresh(pages []CmsPageDescription, outputDirectory string) []CmsPageDescription {
	fmt.Printf("\n\n## TODO: filter page that need a refresh\n\n")
	return pages
}
