package n2m

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jomei/notionapi"
)

// references:
// https://developers.notion.com/reference/post-database-query

func (cms *Notion2Markdown) extractMetaData(block notionapi.Block, page *CmsPageDescription) (*MetaDataInformation, error) {
	childDatabase := block.(*notionapi.ChildDatabaseBlock)
	database, err := cms.client.Database.Query(context.Background(), notionapi.DatabaseID(childDatabase.ID), nil)
	if err != nil {
		return nil, err
	}

	// extract all the properties from the inline database block
	props := map[string]string{}
	for _, result := range database.Results {
		var propName = ""
		var propValue = ""
		for key, element := range result.Properties {
			switch key {
			case "Tag":
				if element.GetType() != "title" {
					return nil, fmt.Errorf("invalid property type (%s) in meta table:%s", element.GetType(), key)
				}
				title := element.(*notionapi.TitleProperty)
				if len(title.Title) != 1 {
					return nil, fmt.Errorf("invalid type data (%s) in meta table:%s", element.GetType(), key)
				}
				propName = title.Title[0].PlainText
			case "Value":
				if element.GetType() != "rich_text" {
					return nil, fmt.Errorf("invalid property type (%s) in meta table:%s", element.GetType(), key)
				}
				richText := element.(*notionapi.RichTextProperty)
				if len(richText.RichText) != 1 {
					return nil, fmt.Errorf("invalid type data (%s) in meta table:%s", element.GetType(), key)
				}
				propValue = richText.RichText[0].PlainText
			default:
				return nil, fmt.Errorf("invalid property in meta table:%s", key)

			}
		}
		props[propName] = propValue
	}

	var medataData = MetaDataInformation{
		Title:                page.Title,
		NotionPageId:         page.Id,
		NotionLastEditedTime: page.LastEditedTime,
	}
	// check those properties
	for propName, propValue := range props {
		switch strings.ToLower(propName) {
		case "description":
			medataData.Description = propValue
		case "slug":
			medataData.Slug = propValue
		case "date":
			const layout = "2006/01/02"
			tm, err := time.Parse(layout, propValue)
			if err != nil {
				return nil, fmt.Errorf("invalid date format (%s) in meta table", propValue)
			}
			medataData.Date = tm
		case "tags":
			medataData.Tags = make([]string, 0)
			tags := strings.Split(propValue, ",")
			for _, tag := range tags {
				medataData.Tags = append(medataData.Tags, strings.Trim(tag, " "))
			}
		case "draft":
			b, err := strconv.ParseBool(propValue)
			if err != nil {
				return nil, fmt.Errorf("invalid boolean format (%s) in meta table", propValue)
			}
			medataData.IsDraft = b

		default:
			return nil, fmt.Errorf("invalid property name (%s) in meta table", propName)
		}
	}

	return &medataData, nil
}
