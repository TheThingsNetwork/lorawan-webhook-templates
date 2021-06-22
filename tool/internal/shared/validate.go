package shared

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/jtacoma/uritemplates"
)

var validFieldID = regexp.MustCompile("^([A-Za-z0-9_\\.]|%[0-9A-Fa-f][0-9A-Fa-f])+$")

func (f WebhookTemplateField) Validate() error {
	if !validFieldID.MatchString(f.ID) {
		return fmt.Errorf("%q is not a valid field ID", f.ID)
	}
	// TODO: Name
	// TODO: Description
	// TODO: DefaultValue
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
		if err := validateURI(*path.value, fields); err != nil {
			return fmt.Errorf("%q is not a valid path for %s: %w", *path.value, path.name, err)
		}
	}
	return nil
}

func (t WebhookTemplate) Validate() error {
	// TODO: TemplateID
	// TODO: Name
	// TODO: Description
	if t.LogoURL != "" {
		res, err := http.Get(t.LogoURL)
		if err != nil {
			return fmt.Errorf("logo-url is not valid: %w", err)
		}
		if res.StatusCode >= 400 {
			return fmt.Errorf("GETting logo-url resulted in error: %s", res.Status)
		}
	}
	if t.InfoURL != "" {
		res, err := http.Get(t.InfoURL)
		if err != nil {
			return fmt.Errorf("info-url is not valid: %w", err)
		}
		if res.StatusCode >= 400 {
			return fmt.Errorf("GETting info-url resulted in error: %s", res.Status)
		}
	}
	if t.DocumentationURL != "" {
		res, err := http.Get(t.DocumentationURL)
		if err != nil {
			return fmt.Errorf("documentation-url is not valid: %w", err)
		}
		if res.StatusCode >= 400 {
			return fmt.Errorf("GETting documentation-url resulted in error: %s", res.Status)
		}
	}
	// TODO: Fields
	// TODO: Format
	for _, header := range t.Headers {
		// TODO: validate header.Key
		if err := validateURI(header.Value.(string), t.Fields); err != nil {
			return fmt.Errorf("%q is not a valid value for header %q: %w", header.Value.(string), header.Key.(string), err)
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
