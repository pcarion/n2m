package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jomei/notionapi"
)

func main() {
	notionIntegrationToken := os.Getenv("NOTION_INTEGRATION_TOKEN")

	if notionIntegrationToken == "" {
		fmt.Printf("Missing environment variable: NOTION_INTEGRATION_TOKEN\n")
		os.Exit(1)
	}

	client := notionapi.NewClient(notionapi.Token(notionIntegrationToken))

	page, err := client.Page.Get(context.Background(), "8ebca155ffda45f7b5d49b0892672dea")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("page: %v\n", page)

	pageData, err := json.Marshal(page)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("pageData: %s\n", string(pageData))

	block, err := client.Block.GetChildren(context.Background(), notionapi.BlockID(page.ID), nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	blockData, err := json.Marshal(block)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("blockData: %s\n", string(blockData))

}
