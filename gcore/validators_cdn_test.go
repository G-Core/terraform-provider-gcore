//go:build !cloud
// +build !cloud

package gcore

import (
	"regexp"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestRawConfigAttrSet(t *testing.T) {
	s3Gcore := func(usePathStyle cty.Value) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{
			"config": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"s3_gcore": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
					"use_path_style": usePathStyle,
					"endpoint":       cty.StringVal(""),
				})}),
			})}),
		})
	}
	path := []interface{}{"config", 0, "s3_gcore", 0}

	tests := []struct {
		name  string
		raw   cty.Value
		field string
		want  bool
	}{
		{name: "null attribute", raw: s3Gcore(cty.NullVal(cty.Bool)), field: "use_path_style", want: false},
		{name: "explicit true", raw: s3Gcore(cty.True), field: "use_path_style", want: true},
		{name: "explicit false", raw: s3Gcore(cty.False), field: "use_path_style", want: true},
		{name: "unknown value", raw: s3Gcore(cty.UnknownVal(cty.Bool)), field: "use_path_style", want: true},
		{name: "explicit empty string", raw: s3Gcore(cty.NullVal(cty.Bool)), field: "endpoint", want: true},
		{name: "missing attribute", raw: s3Gcore(cty.NullVal(cty.Bool)), field: "region", want: false},
		{name: "null raw config", raw: cty.NullVal(cty.DynamicPseudoType), field: "use_path_style", want: false},
		{
			name:  "index out of range",
			raw:   cty.ObjectVal(map[string]cty.Value{"config": cty.ListValEmpty(cty.EmptyObject)}),
			field: "use_path_style",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rawConfigAttrSet(tt.raw, append(path, tt.field)...); got != tt.want {
				t.Fatalf("rawConfigAttrSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestPlanRejectsStorageBindingOwnedFields runs the checks that rely on the raw configuration
// through a real plan, because Resource.Diff in unit tests does not provide one.
func TestPlanRejectsStorageBindingOwnedFields(t *testing.T) {
	t.Setenv("GCORE_PERMANENT_TOKEN", "dummy")
	t.Setenv("GCORE_API_ENDPOINT", "https://api.example.invalid")

	target := func(field string) string {
		return `
resource "gcore_cdn_logs_uploader_target" "t" {
  name = "t"
  config {
    s3_gcore {
      storage_id  = 123
      bucket_name = "bucket"
      ` + field + `
    }
  }
}
`
	}
	origin := func(originField, configField string) string {
		return `
resource "gcore_cdn_origingroup" "og" {
  name = "og"
  origin {
    origin_type = "s3"
    ` + originField + `
    config {
      s3_type        = "gcore"
      storage_id     = 123
      s3_bucket_name = "bucket"
      ` + configField + `
    }
  }
}
`
	}

	tests := []struct {
		name    string
		config  string
		wantErr string
	}{
		{name: "target use_path_style true", config: target(`use_path_style = true`), wantErr: "`use_path_style` cannot be set together with `storage_id`"},
		{name: "target use_path_style false", config: target(`use_path_style = false`), wantErr: "`use_path_style` cannot be set together with `storage_id`"},
		{name: "target empty access key", config: target(`access_key_id = ""`), wantErr: "`access_key_id` cannot be set together with `storage_id`"},
		{name: "target empty endpoint", config: target(`endpoint = ""`), wantErr: "`endpoint` cannot be set together with `storage_id`"},
		{name: "origin empty region", config: origin("", `s3_region = ""`), wantErr: "`s3_region` cannot be set when `s3_type` is 'gcore'"},
		{name: "origin empty secret key", config: origin("", `s3_secret_access_key = ""`), wantErr: "`s3_secret_access_key` cannot be set when `s3_type` is 'gcore'"},
		{name: "origin empty host header override", config: origin(`host_header_override = ""`, ""), wantErr: "`host_header_override` cannot be set when `s3_type` is 'gcore'"},
		{name: "valid bound target", config: target("")},
		{name: "valid bound origin", config: origin("", "")},
		{
			name: "bound origin with s3_type known only after apply",
			config: `
resource "terraform_data" "type" {
  input = "gcore"
}

resource "gcore_cdn_origingroup" "og" {
  name = "og"
  origin {
    origin_type = "s3"
    config {
      s3_type        = terraform_data.type.output
      storage_id     = 123
      s3_bucket_name = "bucket"
    }
  }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := resource.TestStep{
				Config:             tt.config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			}
			if tt.wantErr != "" {
				step.ExpectError = regexp.MustCompile(regexp.QuoteMeta(tt.wantErr))
			}
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps:             []resource.TestStep{step},
			})
		})
	}
}

// TestApplyRejectsConflictsOnceS3TypeIsKnown checks that the checks skipped while s3_type is
// unknown run in the final plan during apply, before any API request.
func TestApplyRejectsConflictsOnceS3TypeIsKnown(t *testing.T) {
	t.Setenv("GCORE_PERMANENT_TOKEN", "dummy")
	t.Setenv("GCORE_API_ENDPOINT", "https://api.example.invalid")

	config := func(s3Type, originField, configField string) string {
		return `
resource "terraform_data" "type" {
  input = "` + s3Type + `"
}

resource "gcore_cdn_origingroup" "og" {
  name = "og"
  origin {
    origin_type = "s3"
    ` + originField + `
    config {
      s3_type        = terraform_data.type.output
      s3_bucket_name = "bucket"
      storage_id     = 123
      ` + configField + `
    }
  }
}
`
	}

	tests := []struct {
		name    string
		config  string
		wantErr string
	}{
		{
			name:    "resolves to amazon with storage_id",
			config:  config("amazon", "", `s3_region = "eu-west-1"`+"\n"+`s3_access_key_id = "ak"`+"\n"+`s3_secret_access_key = "sk"`),
			wantErr: "`storage_id` is only allowed when `s3_type` is 'gcore'",
		},
		{
			name:    "resolves to gcore with region",
			config:  config("gcore", "", `s3_region = "eu-west-1"`),
			wantErr: "`s3_region` cannot be set when `s3_type` is 'gcore'",
		},
		{
			name:    "resolves to gcore with host_header_override",
			config:  config("gcore", `host_header_override = "example.com"`, ""),
			wantErr: "`host_header_override` cannot be set when `s3_type` is 'gcore'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      tt.config,
					ExpectError: regexp.MustCompile(regexp.QuoteMeta(tt.wantErr)),
				}},
			})
		})
	}
}
