package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jomei/notionapi"
)

func main() {
	notionIntegrationToken := os.Getenv("NOTION_INTEGRATION_TOKEN")

	if notionIntegrationToken == "" {
		fmt.Println("Missing environment variable: NOTION_INTEGRATION_TOKEN")
		os.Exit(1)
	}

	fmt.Println("Hello from test")
	client := notionapi.NewClient(notionapi.Token(notionIntegrationToken))

	page, err := client.Page.Get(context.Background(), "8ebca155ffda45f7b5d49b0892672dea")
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("page: %v", page)
}
