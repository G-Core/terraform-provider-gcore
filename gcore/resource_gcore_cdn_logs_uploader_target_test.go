//go:build !cloud
// +build !cloud

package gcore

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/G-Core/gcorelabscdn-go/logsuploader"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestValidateCDNLogsUploaderTargetS3Gcore(t *testing.T) {
	target := func(cfg map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"name": "terraform_acctest_target",
			"config": []interface{}{
				map[string]interface{}{
					"s3_gcore": []interface{}{cfg},
				},
			},
		}
	}
	manual := func() map[string]interface{} {
		return map[string]interface{}{
			"access_key_id":     "access-key",
			"secret_access_key": "secret-access-key",
			"region":            "s-region-1",
			"endpoint":          "https://s-region-1.cloud.gcore.lu",
			"bucket_name":       "bucket",
		}
	}
	bound := func() map[string]interface{} {
		return map[string]interface{}{
			"storage_id":  123,
			"bucket_name": "bucket",
			"directory":   "logs",
		}
	}
	with := func(cfg map[string]interface{}, key string, value interface{}) map[string]interface{} {
		cfg[key] = value
		return cfg
	}
	without := func(cfg map[string]interface{}, key string) map[string]interface{} {
		delete(cfg, key)
		return cfg
	}

	tests := []struct {
		name    string
		config  map[string]interface{}
		wantErr string
	}{
		{name: "valid manual target", config: target(manual())},
		{name: "valid bound target", config: target(bound())},
		{name: "computed storage_id does not error", config: target(with(bound(), "storage_id", unknownVariableValue))},
		{name: "computed manual endpoint does not error", config: target(with(manual(), "endpoint", unknownVariableValue))},
		{
			name:    "manual target requires endpoint",
			config:  target(without(manual(), "endpoint")),
			wantErr: "config.s3_gcore: `endpoint` is required when `storage_id` is not set",
		},
		{
			name:    "manual target requires access key",
			config:  target(without(manual(), "access_key_id")),
			wantErr: "config.s3_gcore: `access_key_id` is required when `storage_id` is not set",
		},
		{
			name:    "bound target rejects access key",
			config:  target(with(bound(), "access_key_id", "access-key")),
			wantErr: "config.s3_gcore: `access_key_id` cannot be set together with `storage_id`",
		},
		{
			name:    "bound target rejects secret key",
			config:  target(with(bound(), "secret_access_key", "secret-access-key")),
			wantErr: "config.s3_gcore: `secret_access_key` cannot be set together with `storage_id`",
		},
		{
			name:    "bound target rejects region",
			config:  target(with(bound(), "region", "s-region-1")),
			wantErr: "config.s3_gcore: `region` cannot be set together with `storage_id`",
		},
		{
			name:    "bound target rejects endpoint",
			config:  target(with(bound(), "endpoint", "https://s-region-1.cloud.gcore.lu")),
			wantErr: "config.s3_gcore: `endpoint` cannot be set together with `storage_id`",
		},
		{
			name:    "bound target rejects computed endpoint",
			config:  target(with(bound(), "endpoint", unknownVariableValue)),
			wantErr: "config.s3_gcore: `endpoint` cannot be set together with `storage_id`",
		},
	}

	targetResource := resourceCDNLogsUploaderTarget()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := targetResource.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(tt.config), nil)
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

func TestBuildTargetConfig(t *testing.T) {
	tests := []struct {
		name        string
		storageType string
		config      map[string]interface{}
		want        map[string]interface{}
	}{
		{
			name:        "bound s3_gcore sends only binding fields",
			storageType: "s3_gcore",
			config: map[string]interface{}{
				"storage_id":        123,
				"bucket_name":       "bucket",
				"directory":         "",
				"access_key_id":     "",
				"secret_access_key": "",
				"region":            "",
				"endpoint":          "",
				"use_path_style":    true,
			},
			want: map[string]interface{}{
				"storage_id":  123,
				"bucket_name": "bucket",
				"directory":   nil,
			},
		},
		{
			name:        "manual s3_gcore drops storage_id",
			storageType: "s3_gcore",
			config: map[string]interface{}{
				"storage_id":        0,
				"bucket_name":       "bucket",
				"directory":         "logs",
				"access_key_id":     "access-key",
				"secret_access_key": "secret-access-key",
				"region":            "s-region-1",
				"endpoint":          "https://s-region-1.cloud.gcore.lu",
				"use_path_style":    false,
			},
			want: map[string]interface{}{
				"bucket_name":       "bucket",
				"directory":         "logs",
				"access_key_id":     "access-key",
				"secret_access_key": "secret-access-key",
				"region":            "s-region-1",
				"endpoint":          "https://s-region-1.cloud.gcore.lu",
				"use_path_style":    false,
			},
		},
		{
			name:        "other storage type is sanitized as before",
			storageType: "s3_other",
			config: map[string]interface{}{
				"bucket_name": "bucket",
				"directory":   "",
			},
			want: map[string]interface{}{
				"bucket_name": "bucket",
				"directory":   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildTargetConfig(tt.storageType, tt.config); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("buildTargetConfig() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestMergeStateConfigS3GcoreBound(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCDNLogsUploaderTarget().Schema, map[string]interface{}{})

	merged := mergeStateConfig(&logsuploader.Target{
		StorageType: logsuploader.StorageType("s3_gcore"),
		Config: map[string]interface{}{
			"storage_id":  float64(123),
			"bucket_name": "bucket",
			"directory":   "logs",
		},
	}, d)

	if merged["storage_id"] != 123 {
		t.Fatalf("storage_id = %#v, want int 123", merged["storage_id"])
	}
	if merged["use_path_style"] != true {
		t.Fatalf("use_path_style = %#v, want schema default true", merged["use_path_style"])
	}
}

func TestAccLogsUploaderTargetS3GcoreStorage(t *testing.T) {
	fullName := "gcore_cdn_logs_uploader_target.acctest_gcore"
	prefix := fmt.Sprintf("tfacc-%d", time.Now().Unix())

	boundTarget := func(bucket, directory string) string {
		return testAccGcoreStorageConfig(prefix) + fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_target" "acctest_gcore" {
  name = "%[1]s"
  config {
    s3_gcore {
      storage_id  = gcore_storage_s3.acctest.storage_id
      bucket_name = gcore_storage_s3_bucket.%[2]s.name
      directory   = "%[3]s"
    }
  }
}
`, prefix, bucket, directory)
	}

	manualTarget := testAccGcoreStorageConfig(prefix) + fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_target" "acctest_gcore" {
  name = "%[1]s"
  config {
    s3_gcore {
      access_key_id     = gcore_storage_s3.acctest.generated_access_key
      secret_access_key = gcore_storage_s3.acctest.generated_secret_key
      region            = gcore_storage_s3.acctest.location
      endpoint          = gcore_storage_s3.acctest.generated_s3_endpoint
      bucket_name       = gcore_storage_s3_bucket.second.name
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
				Config: boundTarget("first", "logs"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttrPair(fullName, "config.0.s3_gcore.0.storage_id", "gcore_storage_s3.acctest", "storage_id"),
					resource.TestCheckResourceAttr(fullName, "config.0.s3_gcore.0.bucket_name", prefix+"-first"),
					resource.TestCheckResourceAttr(fullName, "config.0.s3_gcore.0.directory", "logs"),
					resource.TestCheckResourceAttr(fullName, "config.0.s3_gcore.0.endpoint", ""),
				),
			},
			{
				ResourceName:      fullName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: boundTarget("second", "cdn-logs"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fullName, "config.0.s3_gcore.0.bucket_name", prefix+"-second"),
					resource.TestCheckResourceAttr(fullName, "config.0.s3_gcore.0.directory", "cdn-logs"),
				),
			},
			{
				Config: manualTarget,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fullName, "config.0.s3_gcore.0.storage_id", "0"),
					resource.TestCheckResourceAttrPair(fullName, "config.0.s3_gcore.0.endpoint", "gcore_storage_s3.acctest", "generated_s3_endpoint"),
				),
			},
			{
				Config: boundTarget("first", "logs"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(fullName, "config.0.s3_gcore.0.storage_id", "gcore_storage_s3.acctest", "storage_id"),
				),
			},
		},
	})
}
