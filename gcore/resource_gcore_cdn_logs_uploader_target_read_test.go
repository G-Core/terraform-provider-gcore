//go:build !cloud
// +build !cloud

package gcore

import (
	"testing"

	"github.com/G-Core/gcorelabscdn-go/logsuploader"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestMergeStateConfigSkipsUnknownFields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCDNLogsUploaderTarget().Schema, map[string]interface{}{})

	merged := mergeStateConfig(&logsuploader.Target{
		StorageType: logsuploader.StorageType("s3_oss"),
		Config: map[string]interface{}{
			"access_key_id": "access-key",
			"bucket_name":   "bucket",
			"endpoint":      nil,
			"future_field":  "value",
		},
	}, d)

	if _, ok := merged["future_field"]; ok {
		t.Fatalf("unknown field must be skipped: %#v", merged)
	}
	if _, ok := merged["endpoint"]; !ok {
		t.Fatalf("endpoint must be kept: %#v", merged)
	}

	configData := map[string]interface{}{}
	for storageType := range resourceCDNLogsUploaderTarget().Schema["config"].Elem.(*schema.Resource).Schema {
		configData[storageType] = []interface{}{}
	}
	configData["s3_oss"] = []interface{}{merged}
	if err := d.Set("config", []interface{}{configData}); err != nil {
		t.Fatalf("merged config must be storable in state: %v", err)
	}
}

func TestMergeStateConfigKeepsStateSecrets(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCDNLogsUploaderTarget().Schema, map[string]interface{}{
		"config": []interface{}{
			map[string]interface{}{
				"s3_oss": []interface{}{
					map[string]interface{}{
						"access_key_id":     "access-key",
						"secret_access_key": "secret-access-key",
						"bucket_name":       "bucket",
					},
				},
			},
		},
	})

	merged := mergeStateConfig(&logsuploader.Target{
		StorageType: logsuploader.StorageType("s3_oss"),
		Config: map[string]interface{}{
			"access_key_id":     "access-key",
			"secret_access_key": "*****",
			"bucket_name":       "bucket",
		},
	}, d)

	if merged["secret_access_key"] != "secret-access-key" {
		t.Fatalf("secret_access_key = %#v, want the state value", merged["secret_access_key"])
	}
}

func TestMergeStateConfigUnsupportedStorageType(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCDNLogsUploaderTarget().Schema, map[string]interface{}{})

	merged := mergeStateConfig(&logsuploader.Target{
		StorageType: logsuploader.StorageType("azure_blob"),
		Config:      map[string]interface{}{"container": "logs"},
	}, d)

	if len(merged) != 0 {
		t.Fatalf("unsupported storage type must yield no fields: %#v", merged)
	}
	if _, ok := targetConfigFieldSchemas("azure_blob"); ok {
		t.Fatal("azure_blob must not be reported as supported")
	}
	if _, ok := targetConfigFieldSchemas("s3_oss"); !ok {
		t.Fatal("s3_oss must be reported as supported")
	}
}

func TestMergeStateConfigSkipsUnknownNestedHTTPFields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCDNLogsUploaderTarget().Schema, map[string]interface{}{})

	merged := mergeStateConfig(&logsuploader.Target{
		StorageType: logsuploader.StorageType("http"),
		Config: map[string]interface{}{
			"content_type": "json",
			"upload": map[string]interface{}{
				"url":          "https://logs.example.com",
				"method":       "POST",
				"headers":      map[string]interface{}{"X-Custom": "value"},
				"future_field": "value",
				"response_actions": []interface{}{
					map[string]interface{}{"action": "drop", "match_status_code": float64(400), "future_field": "value"},
				},
			},
			"auth": map[string]interface{}{
				"type":   "token",
				"config": map[string]interface{}{"header_name": "Authorization", "future_field": "value"},
			},
		},
	}, d)

	upload := merged["upload"].(map[string]interface{})
	if _, ok := upload["future_field"]; ok {
		t.Fatalf("unknown upload field must be skipped: %#v", upload)
	}
	if upload["headers"].(map[string]interface{})["X-Custom"] != "value" {
		t.Fatalf("free-form headers must be kept: %#v", upload["headers"])
	}
	action := upload["response_actions"].([]interface{})[0].(map[string]interface{})
	if _, ok := action["future_field"]; ok || action["action"] != "drop" {
		t.Fatalf("unexpected response action: %#v", action)
	}
	authConfig := merged["auth"].(map[string]interface{})["config"].(map[string]interface{})
	if _, ok := authConfig["future_field"]; ok || authConfig["header_name"] != "Authorization" {
		t.Fatalf("unexpected auth config: %#v", authConfig)
	}
}
