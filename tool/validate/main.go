package main

import (
	"log"
	"os"

	"github.com/TheThingsNetwork/lorawan-webhook-templates/tool/internal/shared"
	"gopkg.in/yaml.v2"
)

var check struct{}

func main() {
	indexBytes, err := os.ReadFile("templates.yml")
	if err != nil {
		log.Fatal(err)
	}
	var index []string
	err = yaml.Unmarshal(indexBytes, &index)
	if err != nil {
		log.Fatal(err)
	}
	for _, templateName := range index {
		templateBytes, err := os.ReadFile(templateName + ".yml")
		if err != nil {
			log.Fatalf("failed to read %s.yml: %v", templateName, err)
		}
		var template shared.WebhookTemplate
		err = yaml.Unmarshal(templateBytes, &template)
		if err != nil {
			log.Fatalf("failed to parse %s.yml: %v", templateName, err)
		}

		err = template.Validate()
		if err != nil {
			log.Fatalf("%s.yml is invalid: %v", templateName, err)
		}
	}
}
