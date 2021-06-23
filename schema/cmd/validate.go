// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/TheThingsNetwork/lorawan-webhook-templates/schema"
	"gopkg.in/yaml.v2"
)

// numWorkers is the number of worker goroutines to spawn for template validation.
const numWorkers = 12

func validateTemplate(baseDir, templateID string) error {
	templateYml := filepath.Join(baseDir, templateID+".yml")
	b, err := os.ReadFile(templateYml)
	if err != nil {
		return err
	}
	template := schema.WebhookTemplate{}
	if err := yaml.Unmarshal(b, &template); err != nil {
		return err
	}
	return template.Validate()
}

func readTemplatesIndex(baseDir string) ([]string, error) {
	templatesYml := filepath.Join(baseDir, "templates.yml")
	b, err := os.ReadFile(templatesYml)
	if err != nil {
		return nil, err
	}
	var templateIDs []string
	if err := yaml.Unmarshal(b, &templateIDs); err != nil {
		return nil, err
	}
	if len(templateIDs) == 0 {
		return nil, fmt.Errorf("no template IDs found in %s", templatesYml)
	}
	return templateIDs, nil
}

func init() {
	http.DefaultClient.Timeout = 5 * time.Second
}

func main() {
	baseDir := os.Getenv("WEBHOOK_TEMPLATES_DIR")

	templateIDs, err := readTemplatesIndex(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to retrieve templates: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d webhook templates\n", len(templateIDs))

	wg := sync.WaitGroup{}
	templateCh := make(chan string)

	var rc int32
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			for templateID := range templateCh {
				if err := validateTemplate(baseDir, templateID); err != nil {
					fmt.Fprintf(os.Stderr, "Webhook template %s ERROR: %v\n", templateID, err)
					atomic.StoreInt32(&rc, 1)
				} else {
					fmt.Printf("Webhook template %s OK\n", templateID)
				}
			}
			wg.Done()
		}()
	}

	go func() {
		for _, templateID := range templateIDs {
			templateCh <- templateID
		}
		close(templateCh)
	}()

	wg.Wait()
	os.Exit(int(rc))
}
