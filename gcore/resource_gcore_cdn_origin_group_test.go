//go:build !cloud
// +build !cloud

package gcore

import (
	"context"
	"fmt"
	"strings"
	"testing"

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
