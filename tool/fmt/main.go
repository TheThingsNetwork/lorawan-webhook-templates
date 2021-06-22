package main

import (
	"bytes"
	"flag"
	"log"
	"os"

	"github.com/TheThingsNetwork/lorawan-webhook-templates/tool/internal/shared"
	"gopkg.in/yaml.v2"
)

var exitCode = flag.Bool("exit-code", false, "Set exit code if there is a difference.")

func main() {
	flag.Parse()

	var madeChanges bool

	indexBytes, err := os.ReadFile("templates.yml")
	if err != nil {
		log.Fatal(err)
	}
	var index []string
	err = yaml.Unmarshal(indexBytes, &index)
	if err != nil {
		log.Fatal(err)
	}
	formattedIndexBytes, err := yaml.Marshal(index)
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Equal(indexBytes, formattedIndexBytes) {
		err = os.WriteFile("templates.yml", formattedIndexBytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Formatted templates.yml")
		madeChanges = true
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
		formattedTemplateBytes, err := yaml.Marshal(template)
		if err != nil {
			log.Fatal(err)
		}
		if !bytes.Equal(templateBytes, formattedTemplateBytes) {
			err = os.WriteFile(templateName+".yml", formattedTemplateBytes, 0644)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Formatted %s.yml", templateName)
			madeChanges = true
		}
	}

	if *exitCode && madeChanges {
		os.Exit(1)
	}
}
