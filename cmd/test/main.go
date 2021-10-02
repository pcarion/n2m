package main

import (
	"fmt"
	"os"

	"github.com/pcarion/gonotionapi/pkg/blogcms"
)

func main() {
	notionIntegrationToken := os.Getenv("NOTION_INTEGRATION_TOKEN")

	if notionIntegrationToken == "" {
		fmt.Printf("Missing environment variable: NOTION_INTEGRATION_TOKEN\n")
		os.Exit(1)
	}

	cms, err := blogcms.NewBlocgCms(notionIntegrationToken)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	pages, err := cms.ExtractListPages("8ebca155ffda45f7b5d49b0892672dea")

	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Pages: \n%#v", pages)
}
