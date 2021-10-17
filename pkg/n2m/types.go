package n2m

import "time"

const (
	MdTypePara = iota
	MdTypeHeader
	MdTypeImage
	MdTypeListItem
	MdTypeCode
)

// description of a generated markdown block
type MarkdownBlock struct {
	mdType   int
	level    int
	markdown string
}

type ImageDescription struct {
	imageUrl           string
	imageLocalFileName string
}

type VisitorContext struct {
	metaData    *MetaDataInformation
	page        *CmsPageDescription
	mdBlocks    []MarkdownBlock
	mdImages    []ImageDescription
	imagesCount int
	cms         *Notion2Markdown
}

type MetaDataInformation struct {
	Title                string
	Slug                 string
	Date                 time.Time
	Tags                 []string
	Description          string
	IsDraft              bool
	HasToc               bool
	NotionPageId         string
	NotionLastEditedTime string
}

type MarkdownParagraph struct {
	markdown string
}

type MarkdownCode struct {
	language string
	markdown string
}

type MarkdownImage struct {
	caption            string
	imageToDownloadUrl string
	imageLocalFileName string
	markdown           string
}

type CmsPageDescription struct {
	Id             string
	Title          string
	LastEditedTime string
	Index          int
}
