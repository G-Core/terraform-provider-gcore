//go:build !cloud
// +build !cloud

package gcore

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/G-Core/gcorelabscdn-go/origingroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOriginGroup(t *testing.T) {
	fullName := "gcore_cdn_origingroup.acctest"

	type Params struct {
		Source  string
		Enabled string
	}

	create := Params{"google.com", "true"}
	update := Params{"tut.by", "false"}

	template := func(params *Params) string {
		return fmt.Sprintf(`
            resource "gcore_cdn_origingroup" "acctest" {
			  name = "terraform_acctest_group"
			  use_next = true

			  origin {
			    source = "%s"
				enabled = %s
			  }

			  origin {
			    source = "yandex.ru"
			    enabled = true
			  }
			}
		`, params.Source, params.Enabled)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_CDN_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: template(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "terraform_acctest_group"),
					resource.TestCheckResourceAttr(fullName, "origin.0.source", create.Source),
					resource.TestCheckResourceAttr(fullName, "origin.0.enabled", create.Enabled),
					resource.TestCheckResourceAttr(fullName, "origin.1.source", "yandex.ru"),
					resource.TestCheckResourceAttr(fullName, "origin.1.enabled", "true"),
				),
			},
			{
				Config: template(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "terraform_acctest_group"),
					resource.TestCheckResourceAttr(fullName, "origin.0.source", update.Source),
					resource.TestCheckResourceAttr(fullName, "origin.0.enabled", update.Enabled),
					resource.TestCheckResourceAttr(fullName, "origin.1.source", "yandex.ru"),
					resource.TestCheckResourceAttr(fullName, "origin.1.enabled", "true"),
				),
			},
		},
	})
}

func TestAccOriginGroupMixed(t *testing.T) {
	fullName := "gcore_cdn_origingroup.acctest_mixed"

	createConfig := `
			resource "gcore_cdn_origingroup" "acctest_mixed" {
			  name     = "terraform_acctest_mixed_group"
			  use_next = true

			  origin {
			    source               = "cdn.example.com"
			    enabled              = true
			    host_header_override = "origin.example.com"
			  }

			  origin {
			    origin_type = "s3"
			    enabled     = true
		    backup      = true
		    config {
		      s3_type              = "amazon"
		      s3_bucket_name       = "test-bucket"
		      s3_region            = "eu-west-1"
		      s3_access_key_id     = "dummy-access-key"
		      s3_secret_access_key = "dummy-secret-key"
			    }
			  }
			}
		`

	updateConfig := `
			resource "gcore_cdn_origingroup" "acctest_mixed" {
			  name     = "terraform_acctest_mixed_group"
			  use_next = true

			  origin {
			    source               = "cdn.example.com"
			    enabled              = true
			    host_header_override = "static.example.com"
			  }

			  origin {
			    origin_type          = "s3"
			    enabled              = true
			    backup               = false
			    host_header_override = "storage.example.com"
			    config {
			      s3_type              = "amazon"
			      s3_bucket_name       = "test-bucket"
			      s3_region            = "us-east-1"
			      s3_access_key_id     = "dummy-access-key"
			      s3_secret_access_key = "dummy-secret-key"
			    }
			  }
			}
		`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_CDN_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "terraform_acctest_mixed_group"),
					resource.TestCheckResourceAttr(fullName, "origin.#", "2"),
					resource.TestCheckResourceAttr(fullName, "origin.0.host_header_override", "origin.example.com"),
				),
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "terraform_acctest_mixed_group"),
					resource.TestCheckResourceAttr(fullName, "origin.#", "2"),
					resource.TestCheckResourceAttr(fullName, "origin.0.host_header_override", "static.example.com"),
					resource.TestCheckResourceAttr(fullName, "origin.1.host_header_override", "storage.example.com"),
				),
			},
		},
	})
}

func TestAccOriginGroupS3Origin(t *testing.T) {
	fullName := "gcore_cdn_origingroup.acctest_s3"

	createConfig := `
			resource "gcore_cdn_origingroup" "acctest_s3" {
			  name     = "terraform_acctest_s3_group"
			  use_next = true

			  origin {
			    origin_type          = "s3"
			    enabled              = true
			    host_header_override = "bucket-origin.example.com"
			    config {
			      s3_type              = "amazon"
			      s3_bucket_name       = "test-bucket"
			      s3_region            = "eu-west-1"
			      s3_access_key_id     = "dummy-access-key"
		      s3_secret_access_key = "dummy-secret-key"
		    }
			  }
			}
		`

	updateConfig := `
			resource "gcore_cdn_origingroup" "acctest_s3" {
			  name     = "terraform_acctest_s3_group"
			  use_next = true

			  origin {
			    origin_type          = "s3"
			    enabled              = true
			    host_header_override = "storage.example.com"
			    config {
			      s3_type              = "other"
			      s3_storage_hostname  = "s3.example.com"
			      s3_bucket_name       = "test-bucket"
			      s3_access_key_id     = "dummy-access-key"
			      s3_secret_access_key = "dummy-secret-key"
			    }
			  }
			}
		`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_CDN_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "terraform_acctest_s3_group"),
					resource.TestCheckResourceAttr(fullName, "origin.#", "1"),
					resource.TestCheckResourceAttr(fullName, "origin.0.host_header_override", "bucket-origin.example.com"),
				),
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "terraform_acctest_s3_group"),
					resource.TestCheckResourceAttr(fullName, "origin.#", "1"),
					resource.TestCheckResourceAttr(fullName, "origin.0.host_header_override", "storage.example.com"),
				),
			},
		},
	})
}

func TestValidateCDNOriginGroupConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]interface{}
		wantErr string
	}{
		{
			name: "valid host origin",
			config: map[string]interface{}{
				"name": "terraform_acctest_group",
				"origin": []interface{}{
					map[string]interface{}{
						"source": "example.com",
					},
				},
			},
		},
		{
			name: "host origin rejects config",
			config: map[string]interface{}{
				"name": "terraform_acctest_group",
				"origin": []interface{}{
					map[string]interface{}{
						"source": "example.com",
						"config": []interface{}{
							testOriginGroupS3Config("bucket", "eu-west-1", "dummy-access-key", "dummy-secret-key"),
						},
					},
				},
			},
			wantErr: "origin.0: `config` cannot be specified for host origins",
		},
		{
			name: "valid s3 origin",
			config: map[string]interface{}{
				"name": "terraform_acctest_group",
				"origin": []interface{}{
					map[string]interface{}{
						"origin_type": "s3",
						"config": []interface{}{
							testOriginGroupS3Config("bucket", "eu-west-1", "dummy-access-key", "dummy-secret-key"),
						},
					},
				},
			},
		},
		{
			name: "valid s3 other origin",
			config: map[string]interface{}{
				"name": "terraform_acctest_group",
				"origin": []interface{}{
					map[string]interface{}{
						"origin_type": "s3",
						"config": []interface{}{
							testOriginGroupS3OtherConfig("bucket", "s3.example.com", "dummy-access-key", "dummy-secret-key"),
						},
					},
				},
			},
		},
		{
			name: "s3 origin requires config",
			config: map[string]interface{}{
				"name": "terraform_acctest_group",
				"origin": []interface{}{
					map[string]interface{}{
						"origin_type": "s3",
					},
				},
			},
			wantErr: "origin.0: `config` block is required for s3 origins",
		},
		{
			name: "s3 origin rejects source",
			config: map[string]interface{}{
				"name": "terraform_acctest_group",
				"origin": []interface{}{
					map[string]interface{}{
						"origin_type": "s3",
						"source":      "example.com",
						"config": []interface{}{
							testOriginGroupS3Config("bucket", "eu-west-1", "dummy-access-key", "dummy-secret-key"),
						},
					},
				},
			},
			wantErr: "origin.0: `source` cannot be specified for s3 origins",
		},
		{
			name: "amazon s3 origin requires region",
			config: map[string]interface{}{
				"name": "terraform_acctest_group",
				"origin": []interface{}{
					map[string]interface{}{
						"origin_type": "s3",
						"config": []interface{}{
							map[string]interface{}{
								"s3_type":              "amazon",
								"s3_bucket_name":       "bucket",
								"s3_access_key_id":     "dummy-access-key",
								"s3_secret_access_key": "dummy-secret-key",
							},
						},
					},
				},
			},
			wantErr: "origin.0.config: `s3_region` is required when `s3_type` is 'amazon'",
		},
		{
			name: "other s3 origin requires hostname",
			config: map[string]interface{}{
				"name": "terraform_acctest_group",
				"origin": []interface{}{
					map[string]interface{}{
						"origin_type": "s3",
						"config": []interface{}{
							map[string]interface{}{
								"s3_type":              "other",
								"s3_bucket_name":       "bucket",
								"s3_access_key_id":     "dummy-access-key",
								"s3_secret_access_key": "dummy-secret-key",
							},
						},
					},
				},
			},
			wantErr: "origin.0.config: `s3_storage_hostname` is required when `s3_type` is 'other'",
		},
	}

	originGroupResource := resourceCDNOriginGroup()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := originGroupResource.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(tt.config), nil)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error containing %q", tt.wantErr)
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// unknownVariableValue mirrors the SDK's hcl2shim.UnknownVariableValue sentinel (which
// lives in an internal package). Setting an attribute to this value makes the diff treat it
// as a computed ("known after apply") value, as if it referenced another resource.
const unknownVariableValue = "74D93920-ED26-11E3-AC10-0800200C9A66"

// TestValidateCDNOriginGroupConfigComputed verifies that a required attribute which is
// "known after apply" does not fail validation at plan time. The cases drive the real
// Resource.Diff path, where the diff collapses such a value to an empty string but
// NewValueKnown still reports it as unknown.
func TestValidateCDNOriginGroupConfigComputed(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]interface{}
		wantErr string
	}{
		{
			name: "computed host source does not error",
			config: map[string]interface{}{
				"name": "terraform_acctest_group",
				"origin": []interface{}{
					map[string]interface{}{
						"origin_type": "host",
						"source":      unknownVariableValue,
					},
				},
			},
		},
		{
			name: "empty host source still errors",
			config: map[string]interface{}{
				"name": "terraform_acctest_group",
				"origin": []interface{}{
					map[string]interface{}{
						"origin_type": "host",
						"source":      "",
					},
				},
			},
			wantErr: "origin.0: `source` is required for host origins",
		},
		{
			name: "computed amazon s3_region does not error",
			config: map[string]interface{}{
				"name": "terraform_acctest_group",
				"origin": []interface{}{
					map[string]interface{}{
						"origin_type": "s3",
						"config": []interface{}{
							map[string]interface{}{
								"s3_type":              "amazon",
								"s3_bucket_name":       "bucket",
								"s3_access_key_id":     "dummy-access-key",
								"s3_secret_access_key": "dummy-secret-key",
								"s3_region":            unknownVariableValue,
							},
						},
					},
				},
			},
		},
		{
			name: "computed other s3_storage_hostname does not error",
			config: map[string]interface{}{
				"name": "terraform_acctest_group",
				"origin": []interface{}{
					map[string]interface{}{
						"origin_type": "s3",
						"config": []interface{}{
							map[string]interface{}{
								"s3_type":              "other",
								"s3_bucket_name":       "bucket",
								"s3_access_key_id":     "dummy-access-key",
								"s3_secret_access_key": "dummy-secret-key",
								"s3_storage_hostname":  unknownVariableValue,
							},
						},
					},
				},
			},
		},
	}

	originGroupResource := resourceCDNOriginGroup()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := originGroupResource.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(tt.config), nil)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestPreserveRestoreS3OriginCredentialsByIndex(t *testing.T) {
	originGroupResource := resourceCDNOriginGroup()
	state := schema.TestResourceDataRaw(t, originGroupResource.Schema, map[string]interface{}{
		"name": "terraform_acctest_group",
		"origin": []interface{}{
			map[string]interface{}{
				"origin_type": "s3",
				"config": []interface{}{
					testOriginGroupS3Config("shared-bucket", "eu-west-1", "first-access-key", "first-secret-key"),
				},
			},
			map[string]interface{}{
				"source": "example.com",
			},
			map[string]interface{}{
				"origin_type": "s3",
				"config": []interface{}{
					testOriginGroupS3Config("shared-bucket", "us-east-1", "second-access-key", "second-secret-key"),
				},
			},
		},
	})

	creds := preserveS3OriginCredentials(state)
	if len(creds) != 2 {
		t.Fatalf("expected 2 preserved credential sets, got %d", len(creds))
	}

	origins := []interface{}{
		map[string]interface{}{
			"origin_type": "s3",
			"config": []interface{}{
				map[string]interface{}{
					"s3_type":        "amazon",
					"s3_bucket_name": "shared-bucket",
					"s3_region":      "eu-west-1",
				},
			},
		},
		map[string]interface{}{
			"source": "example.com",
		},
		map[string]interface{}{
			"origin_type": "s3",
			"config": []interface{}{
				map[string]interface{}{
					"s3_type":        "amazon",
					"s3_bucket_name": "shared-bucket",
					"s3_region":      "us-east-1",
				},
			},
		},
	}

	restoreS3OriginCredentials(origins, creds)

	first := origins[0].(map[string]interface{})["config"].([]interface{})[0].(map[string]interface{})
	if got := first["s3_access_key_id"]; got != "first-access-key" {
		t.Fatalf("first origin access key = %v, want %q", got, "first-access-key")
	}
	if got := first["s3_secret_access_key"]; got != "first-secret-key" {
		t.Fatalf("first origin secret key = %v, want %q", got, "first-secret-key")
	}

	second := origins[2].(map[string]interface{})["config"].([]interface{})[0].(map[string]interface{})
	if got := second["s3_access_key_id"]; got != "second-access-key" {
		t.Fatalf("second origin access key = %v, want %q", got, "second-access-key")
	}
	if got := second["s3_secret_access_key"]; got != "second-secret-key" {
		t.Fatalf("second origin secret key = %v, want %q", got, "second-secret-key")
	}
}

func TestOriginGroupHostHeaderOverrideMapping(t *testing.T) {
	hostHeader := "origin.example.com"

	reqs := listToSourceRequests([]interface{}{
		map[string]interface{}{
			"source":               "example.com",
			"enabled":              true,
			"backup":               false,
			"origin_type":          "host",
			"host_header_override": hostHeader,
			"config":               []interface{}{},
		},
	})

	if len(reqs) != 1 {
		t.Fatalf("expected 1 request, got %d", len(reqs))
	}
	if reqs[0].HostHeaderOverride == nil || *reqs[0].HostHeaderOverride != hostHeader {
		t.Fatalf("unexpected host header override in request: %#v", reqs[0].HostHeaderOverride)
	}

	sources := sourcesToList([]origingroups.Source{
		{
			Source:             "example.com",
			Enabled:            true,
			HostHeaderOverride: &hostHeader,
		},
	})

	if len(sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(sources))
	}

	fields := sources[0].(map[string]interface{})
	if got := fields["host_header_override"]; got != hostHeader {
		t.Fatalf("host_header_override = %v, want %q", got, hostHeader)
	}
}

func testOriginGroupS3Config(bucketName, region, accessKeyID, secretAccessKey string) map[string]interface{} {
	return map[string]interface{}{
		"s3_type":              "amazon",
		"s3_bucket_name":       bucketName,
		"s3_region":            region,
		"s3_access_key_id":     accessKeyID,
		"s3_secret_access_key": secretAccessKey,
	}
}

func testOriginGroupS3OtherConfig(bucketName, storageHostname, accessKeyID, secretAccessKey string) map[string]interface{} {
	return map[string]interface{}{
		"s3_type":              "other",
		"s3_bucket_name":       bucketName,
		"s3_storage_hostname":  storageHostname,
		"s3_access_key_id":     accessKeyID,
		"s3_secret_access_key": secretAccessKey,
	}
}

func TestValidateCDNOriginGroupConfigS3Gcore(t *testing.T) {
	s3Origin := func(cfg map[string]interface{}, extra map[string]interface{}) map[string]interface{} {
		origin := map[string]interface{}{
			"origin_type": "s3",
			"config":      []interface{}{cfg},
		}
		for k, v := range extra {
			origin[k] = v
		}
		return map[string]interface{}{
			"name":   "terraform_acctest_group",
			"origin": []interface{}{origin},
		}
	}
	withField := func(cfg map[string]interface{}, key string, value interface{}) map[string]interface{} {
		cfg[key] = value
		return cfg
	}

	tests := []struct {
		name    string
		config  map[string]interface{}
		wantErr string
	}{
		{
			name:   "valid gcore origin",
			config: s3Origin(testOriginGroupS3GcoreConfig("bucket", 123), nil),
		},
		{
			name:   "computed storage_id does not error",
			config: s3Origin(withField(testOriginGroupS3GcoreConfig("bucket", 0), "storage_id", unknownVariableValue), nil),
		},
		{
			name:    "gcore origin requires storage_id",
			config:  s3Origin(testOriginGroupS3GcoreConfig("bucket", 0), nil),
			wantErr: "origin.0.config: `storage_id` is required when `s3_type` is 'gcore'",
		},
		{
			name:    "gcore origin rejects access key",
			config:  s3Origin(withField(testOriginGroupS3GcoreConfig("bucket", 123), "s3_access_key_id", "ak"), nil),
			wantErr: "origin.0.config: `s3_access_key_id` cannot be set when `s3_type` is 'gcore'",
		},
		{
			name:    "gcore origin rejects secret key",
			config:  s3Origin(withField(testOriginGroupS3GcoreConfig("bucket", 123), "s3_secret_access_key", "sk"), nil),
			wantErr: "origin.0.config: `s3_secret_access_key` cannot be set when `s3_type` is 'gcore'",
		},
		{
			name:    "gcore origin rejects region",
			config:  s3Origin(withField(testOriginGroupS3GcoreConfig("bucket", 123), "s3_region", "eu-west-1"), nil),
			wantErr: "origin.0.config: `s3_region` cannot be set when `s3_type` is 'gcore'",
		},
		{
			name:    "gcore origin rejects storage hostname",
			config:  s3Origin(withField(testOriginGroupS3GcoreConfig("bucket", 123), "s3_storage_hostname", "s3.example.com"), nil),
			wantErr: "origin.0.config: `s3_storage_hostname` cannot be set when `s3_type` is 'gcore'",
		},
		{
			name:    "gcore origin rejects computed region",
			config:  s3Origin(withField(testOriginGroupS3GcoreConfig("bucket", 123), "s3_region", unknownVariableValue), nil),
			wantErr: "origin.0.config: `s3_region` cannot be set when `s3_type` is 'gcore'",
		},
		{
			name:    "gcore origin rejects host_header_override",
			config:  s3Origin(testOriginGroupS3GcoreConfig("bucket", 123), map[string]interface{}{"host_header_override": "example.com"}),
			wantErr: "origin.0: `host_header_override` cannot be set when `s3_type` is 'gcore'",
		},
		{
			name:    "amazon origin rejects storage_id",
			config:  s3Origin(withField(testOriginGroupS3Config("bucket", "eu-west-1", "ak", "sk"), "storage_id", 123), nil),
			wantErr: "origin.0.config: `storage_id` is only allowed when `s3_type` is 'gcore'",
		},
		{
			name: "amazon origin requires access key",
			config: s3Origin(map[string]interface{}{
				"s3_type":              "amazon",
				"s3_bucket_name":       "bucket",
				"s3_region":            "eu-west-1",
				"s3_secret_access_key": "sk",
			}, nil),
			wantErr: "origin.0.config: `s3_access_key_id` is required when `s3_type` is 'amazon'",
		},
		{
			name: "other origin requires secret key",
			config: s3Origin(map[string]interface{}{
				"s3_type":             "other",
				"s3_bucket_name":      "bucket",
				"s3_storage_hostname": "s3.example.com",
				"s3_access_key_id":    "ak",
			}, nil),
			wantErr: "origin.0.config: `s3_secret_access_key` is required when `s3_type` is 'other'",
		},
		{
			name:   "computed amazon keys do not error",
			config: s3Origin(testOriginGroupS3Config("bucket", "eu-west-1", unknownVariableValue, unknownVariableValue), nil),
		},
	}

	originGroupResource := resourceCDNOriginGroup()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := originGroupResource.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(tt.config), nil)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestOriginGroupS3GcoreMapping(t *testing.T) {
	reqs := listToSourceRequests([]interface{}{
		map[string]interface{}{
			"source":               "",
			"enabled":              true,
			"backup":               false,
			"origin_type":          "s3",
			"host_header_override": "",
			"config": []interface{}{
				withS3AuthType(testOriginGroupS3GcoreConfig("bucket", 123)),
			},
		},
	})

	if len(reqs) != 1 {
		t.Fatalf("expected 1 request, got %d", len(reqs))
	}

	body, err := json.Marshal(reqs[0])
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if _, ok := got["host_header_override"]; ok {
		t.Fatalf("host_header_override must be omitted: %s", body)
	}
	cfg := got["config"].(map[string]interface{})
	for _, key := range []string{"s3_access_key_id", "s3_secret_access_key", "s3_region", "s3_storage_hostname"} {
		if _, ok := cfg[key]; ok {
			t.Fatalf("%s must be omitted: %s", key, body)
		}
	}
	if cfg["storage_id"] != float64(123) || cfg["s3_type"] != "gcore" || cfg["s3_bucket_name"] != "bucket" {
		t.Fatalf("unexpected config: %s", body)
	}

	sources := sourcesToList([]origingroups.Source{
		{
			Enabled:    true,
			OriginType: "s3",
			Config: &origingroups.S3Config{
				S3Type:       "gcore",
				S3BucketName: "bucket",
				S3AuthType:   "awsSignatureV4",
				StorageID:    123,
			},
		},
	})
	fields := sources[0].(map[string]interface{})
	stateCfg := fields["config"].([]interface{})[0].(map[string]interface{})
	if stateCfg["storage_id"] != 123 || stateCfg["s3_type"] != "gcore" {
		t.Fatalf("unexpected state config: %#v", stateCfg)
	}
	if fields["host_header_override"] != "" {
		t.Fatalf("host_header_override = %v, want empty", fields["host_header_override"])
	}
}

func TestOriginGroupS3ManualMappingOmitsStorageID(t *testing.T) {
	cfg := withS3AuthType(testOriginGroupS3Config("bucket", "eu-west-1", "ak", "sk"))
	cfg["storage_id"] = 0
	reqs := listToSourceRequests([]interface{}{
		map[string]interface{}{
			"enabled":     true,
			"backup":      false,
			"origin_type": "s3",
			"config":      []interface{}{cfg},
		},
	})

	body, err := json.Marshal(reqs[0].Config)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if _, ok := got["storage_id"]; ok {
		t.Fatalf("storage_id must be omitted: %s", body)
	}
	if got["s3_access_key_id"] != "ak" || got["s3_secret_access_key"] != "sk" || got["s3_region"] != "eu-west-1" {
		t.Fatalf("unexpected config: %s", body)
	}
}

func TestRestoreS3OriginCredentialsSkipsGcore(t *testing.T) {
	origins := []interface{}{
		map[string]interface{}{
			"origin_type": "s3",
			"config": []interface{}{
				map[string]interface{}{"s3_type": "gcore", "s3_access_key_id": "", "s3_secret_access_key": ""},
			},
		},
	}

	restoreS3OriginCredentials(origins, map[int]s3OriginCredentials{0: {accessKeyID: "old", secretAccessKey: "old"}})

	cfg := origins[0].(map[string]interface{})["config"].([]interface{})[0].(map[string]interface{})
	if cfg["s3_access_key_id"] != "" || cfg["s3_secret_access_key"] != "" {
		t.Fatalf("credentials must not be restored into a gcore origin: %#v", cfg)
	}
}

func testOriginGroupS3GcoreConfig(bucketName string, storageID int) map[string]interface{} {
	cfg := map[string]interface{}{
		"s3_type":        "gcore",
		"s3_bucket_name": bucketName,
	}
	if storageID != 0 {
		cfg["storage_id"] = storageID
	}
	return cfg
}

func withS3AuthType(cfg map[string]interface{}) map[string]interface{} {
	cfg["s3_auth_type"] = "awsSignatureV4"
	return cfg
}

// testAccGcoreStorageConfig declares a Gcore Object Storage with two buckets for storage-binding tests.
func testAccGcoreStorageConfig(prefix string) string {
	location := os.Getenv("GCORE_STORAGE_S3_LOCATION")
	if location == "" {
		location = "s-region-1"
	}
	return fmt.Sprintf(`
resource "gcore_storage_s3" "acctest" {
  name     = "%[1]s"
  location = "%[2]s"
}

resource "gcore_storage_s3_bucket" "first" {
  storage_id = gcore_storage_s3.acctest.storage_id
  name       = "%[1]s-first"
}

resource "gcore_storage_s3_bucket" "second" {
  storage_id = gcore_storage_s3.acctest.storage_id
  name       = "%[1]s-second"
}
`, prefix, location)
}

func TestAccOriginGroupS3GcoreStorage(t *testing.T) {
	fullName := "gcore_cdn_origingroup.acctest_gcore"
	prefix := fmt.Sprintf("tfacc-%d", time.Now().Unix())

	gcoreOrigin := func(bucket string) string {
		return testAccGcoreStorageConfig(prefix) + fmt.Sprintf(`
resource "gcore_cdn_origingroup" "acctest_gcore" {
  name     = "%[1]s"
  use_next = true

  origin {
    origin_type = "s3"
    config {
      s3_type        = "gcore"
      storage_id     = gcore_storage_s3.acctest.storage_id
      s3_bucket_name = gcore_storage_s3_bucket.%[2]s.name
    }
  }
}
`, prefix, bucket)
	}

	otherOrigin := testAccGcoreStorageConfig(prefix) + fmt.Sprintf(`
resource "gcore_cdn_origingroup" "acctest_gcore" {
  name     = "%[1]s"
  use_next = true

  origin {
    origin_type = "s3"
    config {
      s3_type              = "other"
      s3_storage_hostname  = trimprefix(gcore_storage_s3.acctest.generated_s3_endpoint, "https://")
      s3_bucket_name       = gcore_storage_s3_bucket.second.name
      s3_access_key_id     = gcore_storage_s3.acctest.generated_access_key
      s3_secret_access_key = gcore_storage_s3.acctest.generated_secret_key
    }
  }
}
`, prefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_CDN_URL_VAR, GCORE_STORAGE_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: gcoreOrigin("first"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "origin.#", "1"),
					resource.TestCheckResourceAttr(fullName, "origin.0.config.0.s3_type", "gcore"),
					resource.TestCheckResourceAttrPair(fullName, "origin.0.config.0.storage_id", "gcore_storage_s3.acctest", "storage_id"),
					resource.TestCheckResourceAttr(fullName, "origin.0.config.0.s3_bucket_name", prefix+"-first"),
					resource.TestCheckResourceAttr(fullName, "origin.0.config.0.s3_access_key_id", ""),
					resource.TestCheckResourceAttr(fullName, "origin.0.host_header_override", ""),
				),
			},
			{
				ResourceName:      fullName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: gcoreOrigin("second"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fullName, "origin.0.config.0.s3_bucket_name", prefix+"-second"),
				),
			},
			{
				Config: otherOrigin,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fullName, "origin.0.config.0.s3_type", "other"),
					resource.TestCheckResourceAttr(fullName, "origin.0.config.0.storage_id", "0"),
				),
			},
			{
				Config: gcoreOrigin("first"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fullName, "origin.0.config.0.s3_type", "gcore"),
					resource.TestCheckResourceAttr(fullName, "origin.0.config.0.s3_storage_hostname", ""),
				),
			},
		},
	})
}
