package shared

import "gopkg.in/yaml.v2"

type WebhookTemplateField struct {
	ID           string `yaml:"id"`
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	DefaultValue string `yaml:"default-value,omitempty"`
	Optional     bool   `yaml:"optional,omitempty"`
	Secret       bool   `yaml:"secret,omitempty"`
}

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

type WebhookTemplate struct {
	TemplateID           string                 `yaml:"template-id"`
	Name                 string                 `yaml:"name"`
	Description          string                 `yaml:"description"`
	LogoURL              string                 `yaml:"logo-url"`
	InfoURL              string                 `yaml:"info-url"`
	DocumentationURL     string                 `yaml:"documentation-url"`
	Fields               []WebhookTemplateField `yaml:"fields,omitempty"`
	Format               string                 `yaml:"format"`
	Headers              yaml.MapSlice          `yaml:"headers,omitempty"`
	CreateDownlinkAPIKey bool                   `yaml:"create-downlink-api-key"`
	BaseURL              string                 `yaml:"base-url"`
	Paths                WebhookTemplatePaths   `yaml:"paths,omitempty"`
}
