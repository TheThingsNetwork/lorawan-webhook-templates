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

package schema

// WebhookTemplate defines information about a Webhook Template for The Things Stack.
type WebhookTemplate struct {
	TemplateID           string                 `yaml:"template-id"`
	Name                 string                 `yaml:"name"`
	Description          string                 `yaml:"description"`
	LogoURL              string                 `yaml:"logo-url"`
	InfoURL              string                 `yaml:"info-url"`
	DocumentationURL     string                 `yaml:"documentation-url"`
	BaseURL              string                 `yaml:"base-url"`
	Headers              map[string]string      `yaml:"headers,omitempty"`
	Format               string                 `yaml:"format"`
	Fields               []WebhookTemplateField `yaml:"fields,omitempty"`
	CreateDownlinkAPIKey bool                   `yaml:"create-downlink-api-key"`
	Paths                WebhookTemplatePaths   `yaml:"paths,omitempty"`
}

// WebhookTemplatePaths defines optional paths for each upstream message type.
type WebhookTemplatePaths struct {
	UplinkMessage            *string `yaml:"uplink-message,omitempty"`
	JoinAccept               *string `yaml:"join-accept,omitempty"`
	DownlinkAck              *string `yaml:"downlink-ack,omitempty"`
	DownlinkNack             *string `yaml:"downlink-nack,omitempty"`
	DownlinkSent             *string `yaml:"downlink-sent,omitempty"`
	DownlinkFailed           *string `yaml:"downlink-failed,omitempty"`
	DownlinkQueued           *string `yaml:"downlink-queued,omitempty"`
	DownlinkQueueInvalidated *string `yaml:"downlink-queue-invalidated,omitempty"`
	LocationSolved           *string `yaml:"location-solved,omitempty"`
	ServiceData              *string `yaml:"service-data,omitempty"`
}

// WebhookTemplateField defines an input field for webhook templates.
type WebhookTemplateField struct {
	ID           string `yaml:"id"`
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	Secret       bool   `yaml:"secret"`
	DefaultValue string `yaml:"default-value"`
	Optional     bool   `yaml:"optional"`
}
