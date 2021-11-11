package schema

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/jtacoma/uritemplates"
)

var (
	validFieldIDRegex     = regexp.MustCompile(`^([A-Za-z0-9_\-\\.]|%[0-9A-Fa-f][0-9A-Fa-f])+$`)
	validIdentifierRegex  = regexp.MustCompile(`^[a-z0-9](?:[-]?[a-z0-9]){2,}$`)
	maxIdentifierLength   = 32
	maxNameLength         = 20
	maxDescriptionLength  = 100
	maxPathLength         = 64
	maxDefaultValueLength = 100
)

func (f WebhookTemplateField) Validate() error {
	if !validFieldIDRegex.MatchString(f.ID) || len(f.ID) > maxIdentifierLength {
		return fmt.Errorf("%q is not a valid field ID", f.ID)
	}
	if len(f.Name) > maxNameLength {
		return fmt.Errorf("%q is not a valid field name: max length is %d", f.Name, maxNameLength)
	}
	if len(f.Description) > maxDescriptionLength {
		return fmt.Errorf("%q is not a valid description: max length is %d", f.Description, maxDescriptionLength)
	}
	if len(f.DefaultValue) > maxDefaultValueLength {
		return fmt.Errorf("%q is not a valid default value: max length is %d", f.DefaultValue, maxDefaultValueLength)
	}
	return nil
}

var builtinFields = []string{
	"appID",
	"applicationID",
	"appEUI",
	"joinEUI",
	"devID",
	"deviceID",
	"devEUI",
	"devAddr",
}

func validField(name string, fields []WebhookTemplateField) bool {
	for _, builtin := range builtinFields {
		if name == builtin {
			return true
		}
	}
	for _, field := range fields {
		if name == field.ID {
			return true
		}
	}
	return false
}

func testURL(url string) error {
	var err error
	for retries := 0; retries < 10; retries++ {
		_, err = http.Get(url)
		if err == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "Retry error: %s\n", err)
	}
	return err
}

func validateURI(uri string, fields []WebhookTemplateField) error {
	tmpl, err := uritemplates.Parse(uri)
	if err != nil {
		return err
	}
	for _, name := range tmpl.Names() {
		if !validField(name, fields) {
			return fmt.Errorf("undefined field {%s}", name)
		}
	}
	return nil
}

func (p WebhookTemplatePaths) Validate(fields []WebhookTemplateField) error {
	for _, path := range []struct {
		name  string
		value *string
	}{
		{"uplink-message", p.UplinkMessage},
		{"join-accept", p.JoinAccept},
		{"downlink-ack", p.DownlinkAck},
		{"downlink-nack", p.DownlinkNack},
		{"downlink-sent", p.DownlinkSent},
		{"downlink-failed", p.DownlinkFailed},
		{"downlink-queued", p.DownlinkQueued},
		{"downlink-queue-invalidated", p.DownlinkQueueInvalidated},
		{"location-solved", p.LocationSolved},
		{"service-data", p.ServiceData},
	} {
		if path.value == nil {
			continue
		}
		if len(*path.value) > maxPathLength {
			return fmt.Errorf("%s path %q is not valid: max length is %d", path.name, *path.value, maxPathLength)
		}
		if err := validateURI(*path.value, fields); err != nil {
			return fmt.Errorf("%s path %q is not valid: %w", path.name, *path.value, err)
		}
	}
	return nil
}

func (t WebhookTemplate) Validate() error {
	if !validIdentifierRegex.MatchString(t.TemplateID) || len(t.TemplateID) > maxIdentifierLength {
		return fmt.Errorf("%q is not a valid template ID", t.TemplateID)
	}
	if len(t.Name) > maxNameLength {
		return fmt.Errorf("%q is not a valid template name: max length is %d", t.Name, maxNameLength)
	}
	if len(t.Description) > maxDescriptionLength {
		return fmt.Errorf("description is not valid: max length is %d", maxDescriptionLength)
	}
	if t.LogoURL != "" {
		if err := testURL(t.LogoURL); err != nil {
			return fmt.Errorf("logo-url is not valid: %w", err)
		}
	}
	if t.InfoURL != "" {
		if err := testURL(t.InfoURL); err != nil {
			return fmt.Errorf("info-url is not valid: %w", err)
		}
	}
	if t.DocumentationURL != "" {
		if err := testURL(t.DocumentationURL); err != nil {
			return fmt.Errorf("documentation-url is not valid: %w", err)
		}
	}
	if !validIdentifierRegex.MatchString(t.Format) {
		return fmt.Errorf("%q is not a valid format", t.Format)
	}
	for headerKey, headerValue := range t.Headers {
		if err := validateURI(headerKey, t.Fields); err != nil {
			return fmt.Errorf("%q is not a valid header key: %w", headerKey, err)
		}
		if err := validateURI(headerValue, t.Fields); err != nil {
			return fmt.Errorf("%q is not a valid value for header %q: %w", headerValue, headerKey, err)
		}
	}
	if err := validateURI(t.BaseURL, t.Fields); err != nil {
		return fmt.Errorf("%q is not a valid base URL: %w", t.BaseURL, err)
	}
	if err := t.Paths.Validate(t.Fields); err != nil {
		return err
	}
	return nil
}
