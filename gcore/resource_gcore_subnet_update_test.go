package gcore

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	gc "github.com/G-Core/gcorelabscloud-go/gcore"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const subnetUpdateTestResponse = `{
	"id": "subnet-1", "name": "subnet", "cidr": "192.168.10.0/24", "network_id": "network-1",
	"enable_dhcp": true, "ip_version": 4, "gateway_ip": "192.168.10.1",
	"dns_nameservers": [], "host_routes": [], "region_id": 1, "project_id": 1,
	"created_at": "2020-01-01T00:00:00+0000", "updated_at": "2020-01-01T00:00:00+0000",
	"metadata": []
}`

func TestResourceSubnetUpdateTags(t *testing.T) {
	tests := []struct {
		name       string
		stateTags  map[string]string
		configName string
		configTags map[string]interface{}
		wantTags   map[string]interface{}
		wantNoTags bool
	}{
		{
			name:       "add, keep, change and remove keys",
			stateTags:  map[string]string{"keep": "v", "change": "old", "remove": "v"},
			configName: "subnet",
			configTags: map[string]interface{}{"keep": "v", "change": "new", "add": "v"},
			wantTags:   map[string]interface{}{"keep": "v", "change": "new", "add": "v", "remove": nil},
		},
		{
			name:       "remove metadata_map from config",
			stateTags:  map[string]string{"key1": "v1", "key2": "v2"},
			configName: "subnet",
			wantTags:   map[string]interface{}{"key1": nil, "key2": nil},
		},
		{
			name:       "rename without a metadata_map change",
			stateTags:  map[string]string{"key1": "v1"},
			configName: "renamed",
			configTags: map[string]interface{}{"key1": "v1"},
			wantNoTags: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requests []string
			var patchBody map[string]interface{}
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests = append(requests, r.Method+" "+r.URL.Path)
				if r.Method == http.MethodPatch {
					body, _ := io.ReadAll(r.Body)
					if err := json.Unmarshal(body, &patchBody); err != nil {
						t.Errorf("cannot decode the PATCH body: %v", err)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(subnetUpdateTestResponse))
			}))
			defer srv.Close()

			provider, err := gc.APITokenClient(gcorecloud.APITokenOptions{APIURL: srv.URL, APIToken: "token"})
			if err != nil {
				t.Fatal(err)
			}
			meta := &Config{Provider: provider}

			attributes := map[string]string{
				"id": "subnet-1", "project_id": "1", "region_id": "1", "name": "subnet",
				"cidr": "192.168.10.0/24", "network_id": "network-1", "connect_to_network_router": "true",
				"enable_dhcp": "true", "gateway_ip": "192.168.10.1", "dns_nameservers.#": "0", "host_routes.#": "0",
			}
			for k, v := range tt.stateTags {
				attributes["metadata_map."+k] = v
			}
			attributes["metadata_map.%"] = strconv.Itoa(len(tt.stateTags))
			state := &terraform.InstanceState{ID: "subnet-1", Attributes: attributes}

			raw := map[string]interface{}{
				"project_id": 1, "region_id": 1, "name": tt.configName,
				"cidr": "192.168.10.0/24", "network_id": "network-1",
			}
			if tt.configTags != nil {
				raw["metadata_map"] = tt.configTags
			}

			r := resourceSubnet()
			diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(raw), meta)
			if err != nil {
				t.Fatal(err)
			}
			if _, diags := r.Apply(context.Background(), state, diff, meta); diags.HasError() {
				t.Fatalf("apply failed: %v", diags)
			}

			wantRequests := []string{"PATCH /v1/subnets/1/1/subnet-1", "GET /v1/subnets/1/1/subnet-1"}
			if d := cmp.Diff(wantRequests, requests); d != "" {
				t.Errorf("requests mismatch (-want +got):\n%s", d)
			}
			for _, req := range requests {
				if strings.Contains(req, "/metadata") {
					t.Errorf("unexpected metadata request: %s", req)
				}
			}

			tags, ok := patchBody["tags"]
			if tt.wantNoTags {
				if ok {
					t.Errorf("PATCH body must not contain tags, got %v", tags)
				}
				return
			}
			if d := cmp.Diff(tt.wantTags, tags); d != "" {
				t.Errorf("tags mismatch (-want +got):\n%s", d)
			}
		})
	}
}
