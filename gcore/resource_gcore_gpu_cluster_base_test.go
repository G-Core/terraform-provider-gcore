//go:build !cloud
// +build !cloud

package gcore

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// gpuClusterRawConfig builds the raw config value Terraform sends the SDK for
// a cluster with a single interface: unset attributes are null, set ones carry
// the user's values.
func gpuClusterRawConfig(clusterWideSGs, ifaceSGs cty.Value) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"servers_settings": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"security_groups": clusterWideSGs,
				"interface": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"name":            cty.StringVal("iface0"),
						"type":            cty.StringVal("external"),
						"security_groups": ifaceSGs,
					}),
				}),
			}),
		}),
	})
}

func TestGPUClusterCustomizeDiffSecurityGroupConflict(t *testing.T) {
	r := resourceGPUCluster(GPUNodeTypeVirtual)

	sgID := "11111111-2222-3333-4444-555555555555"
	ifaceHash := instanceInterfaceUniqueIDByName(map[string]interface{}{"name": "iface0"})
	sgHash := schema.HashString(sgID)

	// State as Read leaves it after a refresh of a legacy cluster-wide config:
	// the API resolves security groups per interface, so the Computed
	// interface.security_groups is populated for every interface even though
	// the user never configured it.
	refreshedState := func(rawConfig cty.Value) *terraform.InstanceState {
		return &terraform.InstanceState{
			ID: "cluster-1",
			Attributes: map[string]string{
				"name":                                  "legacy",
				"servers_count":                         "1",
				"servers_settings.#":                    "1",
				"servers_settings.0.security_groups.#":  "1",
				"servers_settings.0.security_groups.0":  sgID,
				"servers_settings.0.interface.#":        "1",
				fmt.Sprintf("servers_settings.0.interface.%d.name", ifaceHash):              "iface0",
				fmt.Sprintf("servers_settings.0.interface.%d.type", ifaceHash):              "external",
				fmt.Sprintf("servers_settings.0.interface.%d.security_groups.#", ifaceHash): "1",
				fmt.Sprintf("servers_settings.0.interface.%d.security_groups.%d", ifaceHash, sgHash): sgID,
			},
			RawConfig: rawConfig,
		}
	}

	// What the SDK receives as "config" during plan is the proposed new state
	// (grpc_provider.go shims proposedNewStateVal, not the raw config), where
	// Terraform Core has already merged the prior state's Computed
	// interface.security_groups into the unset config attribute.
	legacyProposedConfig := map[string]interface{}{
		"name":          "legacy",
		"servers_count": 1,
		"servers_settings": []interface{}{
			map[string]interface{}{
				"security_groups": []interface{}{sgID},
				"interface": []interface{}{
					map[string]interface{}{
						"name":            "iface0",
						"type":            "external",
						"security_groups": []interface{}{sgID},
					},
				},
			},
		},
	}

	tests := []struct {
		name    string
		state   *terraform.InstanceState
		config  map[string]interface{}
		wantErr string
	}{
		{
			// The user only set the deprecated cluster-wide field; the
			// per-interface groups in state are API-resolved (Computed) and
			// must not be mistaken for user configuration.
			name: "legacy cluster-wide config plans after refresh",
			state: refreshedState(gpuClusterRawConfig(
				cty.ListVal([]cty.Value{cty.StringVal(sgID)}),
				cty.NullVal(cty.Set(cty.String)),
			)),
			config: legacyProposedConfig,
		},
		{
			name: "cluster-wide combined with per-interface fails",
			state: &terraform.InstanceState{
				RawConfig: gpuClusterRawConfig(
					cty.ListVal([]cty.Value{cty.StringVal(sgID)}),
					cty.SetVal([]cty.Value{cty.StringVal("sg-per-iface")}),
				),
			},
			config: map[string]interface{}{
				"name":          "conflict",
				"servers_count": 1,
				"servers_settings": []interface{}{
					map[string]interface{}{
						"security_groups": []interface{}{sgID},
						"interface": []interface{}{
							map[string]interface{}{
								"name":            "iface0",
								"type":            "external",
								"security_groups": []interface{}{"sg-per-iface"},
							},
						},
					},
				},
			},
			wantErr: "cannot be combined with",
		},
		{
			name: "per-interface only plans",
			state: &terraform.InstanceState{
				RawConfig: gpuClusterRawConfig(
					cty.NullVal(cty.List(cty.String)),
					cty.SetVal([]cty.Value{cty.StringVal("sg-per-iface")}),
				),
			},
			config: map[string]interface{}{
				"name":          "per-iface",
				"servers_count": 1,
				"servers_settings": []interface{}{
					map[string]interface{}{
						"interface": []interface{}{
							map[string]interface{}{
								"name":            "iface0",
								"type":            "external",
								"security_groups": []interface{}{"sg-per-iface"},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// SimpleDiff mirrors PlanResourceChange in grpc_provider.go.
			_, err := r.SimpleDiff(context.Background(), tt.state, terraform.NewResourceConfigRaw(tt.config), nil)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
