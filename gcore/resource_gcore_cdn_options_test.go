//go:build !cloud
// +build !cloud

package gcore

import (
	"testing"

	gcdn "github.com/G-Core/gcorelabscdn-go/gcore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestNetworkErrorLoggingOptionToAPI(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCDNResource().Schema, map[string]interface{}{
		"options": []interface{}{
			map[string]interface{}{
				"network_error_logging": []interface{}{
					map[string]interface{}{"value": false},
				},
			},
		},
	})

	opts := listToOptions(d.Get("options").([]interface{}))
	if opts == nil || opts.NetworkErrorLogging == nil {
		t.Fatalf("network_error_logging must be mapped: %#v", opts)
	}
	if !opts.NetworkErrorLogging.Enabled || opts.NetworkErrorLogging.Value {
		t.Fatalf("unexpected network_error_logging: %#v", opts.NetworkErrorLogging)
	}
}

func TestNetworkErrorLoggingOptionFromAPI(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCDNResource().Schema, map[string]interface{}{})

	opts := &gcdn.Options{NetworkErrorLogging: &gcdn.NetworkErrorLogging{Enabled: true, Value: true}}
	if err := d.Set("options", optionsToList(opts)); err != nil {
		t.Fatalf("set options: %s", err)
	}

	if !d.Get("options.0.network_error_logging.0.enabled").(bool) {
		t.Fatal("enabled must be true")
	}
	if !d.Get("options.0.network_error_logging.0.value").(bool) {
		t.Fatal("value must be true")
	}
}

func TestNetworkErrorLoggingOptionNotInRuleSchemas(t *testing.T) {
	for name, s := range map[string]map[string]*schema.Schema{
		"gcore_cdn_rule":          resourceCDNRule().Schema,
		"gcore_cdn_rule_template": resourceRuleTemplate().Schema,
	} {
		opts := s["options"].Elem.(*schema.Resource).Schema
		if _, ok := opts["network_error_logging"]; ok {
			t.Errorf("%s must not expose network_error_logging", name)
		}
	}
}
