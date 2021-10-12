package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pcarion/gonotionapi/pkg/n2m"
)

type JsonConfig struct {
	RootPageId      string `json:"rootPageId"`
	OutputDirectory string `json:"outputDirectory"`
}

func parseJSONConfig() (*JsonConfig, error) {
	content, err := os.ReadFile("n2m.json")
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Printf("Cannot read config file n2m.json in %s\n", pwd)
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

	cms, err := n2m.NewNotionToMarkdown(notionIntegrationToken)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	err = cms.GenerateMardown(config.RootPageId, config.OutputDirectory)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
