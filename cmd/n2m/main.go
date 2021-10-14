package main

import (
	"encoding/json"
	"flag"
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
	// cli flags
	var debugMode bool
	var pageIndex int
	flag.BoolVar(&debugMode, "debug", false, "debug mode")
	flag.IntVar(&pageIndex, "page", -1, "the page index to generate")
	flag.Parse()

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

	cms, err := n2m.NewNotionToMarkdown(notionIntegrationToken, debugMode)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	err = cms.GenerateMardown(config.RootPageId, config.OutputDirectory, pageIndex)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
