package blogcms

import "github.com/jomei/notionapi"

type Blogcms struct {
	client *notionapi.Client
}

func NewBlocgCms(notionIntegrationToken string) (*Blogcms, error) {
	client := notionapi.NewClient(notionapi.Token(notionIntegrationToken))

	return &Blogcms{
		client: client,
	}, nil
}
