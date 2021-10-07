package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pcarion/gonotionapi/pkg/blogcms"
)

type JsonConfig struct {
	RootPageId      string `json:"rootPageId"`
	OutputDirectory string `json:"outputDirectory"`
}

func parseJSONConfig() (*JsonConfig, error) {
	content, err := os.ReadFile("n2b.json")
	if err != nil {
		return nil, err
	}
	config := JsonConfig{}
	json.Unmarshal(content, &config)
	return &config, nil
}

func main() {
	notionIntegrationToken := os.Getenv("NOTION_INTEGRATION_TOKEN")

	if notionIntegrationToken == "" {
		fmt.Printf("Missing environment variable: NOTION_INTEGRATION_TOKEN\n")
		os.Exit(1)
	}

	config, err := parseJSONConfig()
	if err != nil {
		fmt.Printf("Error reading configuration file\n")
		os.Exit(1)
	}

	cms, err := blogcms.NewBlocgCms(notionIntegrationToken)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	err = cms.GenerateContent(config.RootPageId, config.OutputDirectory)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
