package customvalidator_test

import (
	"context"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/G-Core/terraform-provider-gcore/internal/customvalidator"
)

var (
	securityGroupType = types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.StringType}}
	floatingIPType    = types.ObjectType{AttrTypes: map[string]attr.Type{"source": types.StringType}}

	// interfaceAttrTypes mirrors the flattened interface object in both GPU cluster
	// resource schemas. The validator only calls IsNull/IsUnknown on the non-`type`
	// fields, so plain framework types stand in for the schema's custom types.
	interfaceAttrTypes = map[string]attr.Type{
		"type":                  types.StringType,
		"ip_family":             types.StringType,
		"name":                  types.StringType,
		"port_security_enabled": types.BoolType,
		"security_groups":       types.ListType{ElemType: securityGroupType},
		"network_id":            types.StringType,
		"subnet_id":             types.StringType,
		"floating_ip":           floatingIPType,
	}

	interfaceObjectType = types.ObjectType{AttrTypes: interfaceAttrTypes}
)

func nullInterfaceAttrs() map[string]attr.Value {
	return map[string]attr.Value{
		"type":                  types.StringNull(),
		"ip_family":             types.StringNull(),
		"name":                  types.StringNull(),
		"port_security_enabled": types.BoolNull(),
		"security_groups":       types.ListNull(securityGroupType),
		"network_id":            types.StringNull(),
		"subnet_id":             types.StringNull(),
		"floating_ip":           types.ObjectNull(floatingIPType.AttrTypes),
	}
}

// iface builds one interface object, defaulting every unnamed attribute to null.
func iface(t *testing.T, set map[string]attr.Value) attr.Value {
	t.Helper()

	attrs := nullInterfaceAttrs()
	for name, value := range set {
		if _, ok := attrs[name]; !ok {
			t.Fatalf("unknown interface attribute %q in test case", name)
		}
		attrs[name] = value
	}

	obj, diags := types.ObjectValue(interfaceAttrTypes, attrs)
	if diags.HasError() {
		t.Fatalf("building interface object: %v", diags)
	}
	return obj
}

func ifaceList(t *testing.T, elems ...attr.Value) types.List {
	t.Helper()

	list, diags := types.ListValue(interfaceObjectType, elems)
	if diags.HasError() {
		t.Fatalf("building interfaces list: %v", diags)
	}
	return list
}

// errorPaths runs the validator and returns the sorted attribute paths it flagged.
func errorPaths(t *testing.T, configValue types.List) []string {
	t.Helper()

	req := validator.ListRequest{
		Path:        path.Root("servers_settings").AtName("interfaces"),
		ConfigValue: configValue,
	}
	resp := &validator.ListResponse{}
	customvalidator.GPUClusterInterfaceVariantValidator{}.ValidateList(context.Background(), req, resp)

	paths := make([]string, 0, len(resp.Diagnostics))
	for _, d := range resp.Diagnostics {
		withPath, ok := d.(interface{ Path() path.Path })
		if !ok {
			t.Fatalf("diagnostic %q has no path", d.Summary())
		}
		paths = append(paths, withPath.Path().String())
	}
	sort.Strings(paths)
	return paths
}

func sgList() attr.Value {
	return types.ListValueMust(securityGroupType, []attr.Value{
		types.ObjectValueMust(securityGroupType.AttrTypes, map[string]attr.Value{
			"id": types.StringValue("11111111-1111-1111-1111-111111111111"),
		}),
	})
}

func floatingIP() attr.Value {
	return types.ObjectValueMust(floatingIPType.AttrTypes, map[string]attr.Value{
		"source": types.StringValue("new"),
	})
}

func TestGPUClusterInterfaceVariantValidator(t *testing.T) {
	t.Parallel()

	const (
		netID    = "22222222-2222-2222-2222-222222222222"
		subnetID = "33333333-3333-3333-3333-333333333333"
	)

	testCases := map[string]struct {
		ifaces []map[string]attr.Value
		want   []string
	}{
		// --- external ---
		"external with only allowed fields": {
			ifaces: []map[string]attr.Value{{
				"type":                  types.StringValue("external"),
				"ip_family":             types.StringValue("ipv4"),
				"name":                  types.StringValue("eth0"),
				"port_security_enabled": types.BoolValue(true),
				"security_groups":       sgList(),
			}},
		},
		"external rejects network_id": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("external"),
				"network_id": types.StringValue(netID),
			}},
			want: []string{"servers_settings.interfaces[0].network_id"},
		},
		"external rejects subnet_id": {
			ifaces: []map[string]attr.Value{{
				"type":      types.StringValue("external"),
				"subnet_id": types.StringValue(subnetID),
			}},
			want: []string{"servers_settings.interfaces[0].subnet_id"},
		},
		"external rejects floating_ip": {
			ifaces: []map[string]attr.Value{{
				"type":        types.StringValue("external"),
				"floating_ip": floatingIP(),
			}},
			want: []string{"servers_settings.interfaces[0].floating_ip"},
		},
		"external reports every prohibited field": {
			ifaces: []map[string]attr.Value{{
				"type":        types.StringValue("external"),
				"network_id":  types.StringValue(netID),
				"subnet_id":   types.StringValue(subnetID),
				"floating_ip": floatingIP(),
			}},
			want: []string{
				"servers_settings.interfaces[0].floating_ip",
				"servers_settings.interfaces[0].network_id",
				"servers_settings.interfaces[0].subnet_id",
			},
		},

		// --- subnet ---
		"subnet with only allowed fields": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("subnet"),
				"network_id": types.StringValue(netID),
				"subnet_id":  types.StringValue(subnetID),
			}},
		},
		"subnet rejects ip_family": { // GCLOUD2-28207
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("subnet"),
				"network_id": types.StringValue(netID),
				"subnet_id":  types.StringValue(subnetID),
				"ip_family":  types.StringValue("ipv4"),
			}},
			want: []string{"servers_settings.interfaces[0].ip_family"},
		},
		"subnet accepts unknown ip_family": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("subnet"),
				"network_id": types.StringValue(netID),
				"subnet_id":  types.StringValue(subnetID),
				"ip_family":  types.StringUnknown(),
			}},
		},
		"subnet requires network_id": {
			ifaces: []map[string]attr.Value{{
				"type":      types.StringValue("subnet"),
				"subnet_id": types.StringValue(subnetID),
			}},
			want: []string{"servers_settings.interfaces[0].network_id"},
		},
		"subnet requires subnet_id": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("subnet"),
				"network_id": types.StringValue(netID),
			}},
			want: []string{"servers_settings.interfaces[0].subnet_id"},
		},
		"subnet rejects empty network_id": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("subnet"),
				"network_id": types.StringValue(""),
				"subnet_id":  types.StringValue(subnetID),
			}},
			want: []string{"servers_settings.interfaces[0].network_id"},
		},
		"subnet accepts unknown network_id": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("subnet"),
				"network_id": types.StringUnknown(),
				"subnet_id":  types.StringValue(subnetID),
			}},
		},

		// --- any_subnet ---
		"any_subnet with only allowed fields": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("any_subnet"),
				"network_id": types.StringValue(netID),
				"ip_family":  types.StringValue("dual"),
			}},
		},
		"any_subnet rejects subnet_id": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("any_subnet"),
				"network_id": types.StringValue(netID),
				"subnet_id":  types.StringValue(subnetID),
			}},
			want: []string{"servers_settings.interfaces[0].subnet_id"},
		},
		"any_subnet requires network_id": {
			ifaces: []map[string]attr.Value{{
				"type": types.StringValue("any_subnet"),
			}},
			want: []string{"servers_settings.interfaces[0].network_id"},
		},

		// --- discriminator normalization (schema validates type case-insensitively) ---
		"uppercase SUBNET rejects ip_family": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("SUBNET"),
				"network_id": types.StringValue(netID),
				"subnet_id":  types.StringValue(subnetID),
				"ip_family":  types.StringValue("ipv4"),
			}},
			want: []string{"servers_settings.interfaces[0].ip_family"},
		},
		"uppercase EXTERNAL rejects network_id": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("EXTERNAL"),
				"network_id": types.StringValue(netID),
			}},
			want: []string{"servers_settings.interfaces[0].network_id"},
		},
		"mixed case Any_Subnet requires network_id": {
			ifaces: []map[string]attr.Value{{
				"type": types.StringValue("Any_Subnet"),
			}},
			want: []string{"servers_settings.interfaces[0].network_id"},
		},
		"uppercase SUBNET accepts a valid config": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("Subnet"),
				"network_id": types.StringValue(netID),
				"subnet_id":  types.StringValue(subnetID),
			}},
		},

		// --- values the validator must not judge ---
		"null type is skipped": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringNull(),
				"network_id": types.StringValue(netID),
			}},
		},
		"unknown type is skipped": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringUnknown(),
				"network_id": types.StringValue(netID),
			}},
		},
		"unrecognized type is left to the OneOf validator": {
			ifaces: []map[string]attr.Value{{
				"type":       types.StringValue("bogus"),
				"network_id": types.StringValue(netID),
				"subnet_id":  types.StringValue(subnetID),
			}},
		},

		// --- multiple elements ---
		"errors carry the correct element index": {
			ifaces: []map[string]attr.Value{
				{
					"type":      types.StringValue("external"),
					"ip_family": types.StringValue("ipv4"),
				},
				{
					"type":       types.StringValue("subnet"),
					"network_id": types.StringValue(netID),
					"subnet_id":  types.StringValue(subnetID),
					"ip_family":  types.StringValue("ipv4"),
				},
				{
					"type":       types.StringValue("any_subnet"),
					"network_id": types.StringValue(netID),
					"subnet_id":  types.StringValue(subnetID),
				},
			},
			want: []string{
				"servers_settings.interfaces[1].ip_family",
				"servers_settings.interfaces[2].subnet_id",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			elems := make([]attr.Value, 0, len(tc.ifaces))
			for _, set := range tc.ifaces {
				elems = append(elems, iface(t, set))
			}

			got := errorPaths(t, ifaceList(t, elems...))
			want := tc.want
			if want == nil {
				want = []string{}
			}

			if len(got) != len(want) {
				t.Fatalf("got %d error(s) %v, want %d %v", len(got), got, len(want), want)
			}
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("error %d: got path %q, want %q", i, got[i], want[i])
				}
			}
		})
	}
}

// TestGPUClusterInterfaceVariantValidatorDiagnostics pins the wording users see,
// including the "why" clause that makes the error actionable.
func TestGPUClusterInterfaceVariantValidatorDiagnostics(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		iface       map[string]attr.Value
		wantSummary string
		wantDetail  string
	}{
		"ip_family on subnet": {
			iface: map[string]attr.Value{
				"type":       types.StringValue("subnet"),
				"network_id": types.StringValue("net"),
				"subnet_id":  types.StringValue("sub"),
				"ip_family":  types.StringValue("ipv4"),
			},
			wantSummary: "Invalid interface configuration",
			wantDetail:  `'ip_family' is not supported when type = "subnet"; the IP family is determined by the referenced subnet. Remove it.`,
		},
		"subnet_id on any_subnet": {
			iface: map[string]attr.Value{
				"type":       types.StringValue("any_subnet"),
				"network_id": types.StringValue("net"),
				"subnet_id":  types.StringValue("sub"),
			},
			wantSummary: "Invalid interface configuration",
			wantDetail:  `'subnet_id' is not supported when type = "any_subnet"; use type = "subnet" to pin a specific subnet. Remove it.`,
		},
		"missing network_id on subnet": {
			iface: map[string]attr.Value{
				"type":      types.StringValue("subnet"),
				"subnet_id": types.StringValue("sub"),
			},
			wantSummary: "Missing required attribute",
			wantDetail:  `'network_id' is required when type = "subnet".`,
		},
		// The message reports the normalized variant, not the raw casing.
		"uppercase type normalizes in the message": {
			iface: map[string]attr.Value{
				"type":       types.StringValue("SUBNET"),
				"network_id": types.StringValue("net"),
				"subnet_id":  types.StringValue("sub"),
				"ip_family":  types.StringValue("ipv4"),
			},
			wantSummary: "Invalid interface configuration",
			wantDetail:  `'ip_family' is not supported when type = "subnet"; the IP family is determined by the referenced subnet. Remove it.`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.ListRequest{
				Path:        path.Root("servers_settings").AtName("interfaces"),
				ConfigValue: ifaceList(t, iface(t, tc.iface)),
			}
			resp := &validator.ListResponse{}
			customvalidator.GPUClusterInterfaceVariantValidator{}.ValidateList(context.Background(), req, resp)

			if len(resp.Diagnostics) != 1 {
				t.Fatalf("got %d diagnostics, want 1: %v", len(resp.Diagnostics), resp.Diagnostics)
			}
			d := resp.Diagnostics[0]
			if d.Summary() != tc.wantSummary {
				t.Errorf("summary: got %q, want %q", d.Summary(), tc.wantSummary)
			}
			if d.Detail() != tc.wantDetail {
				t.Errorf("detail:\n got %q\nwant %q", d.Detail(), tc.wantDetail)
			}
		})
	}
}

func TestGPUClusterInterfaceVariantValidatorSkipsAbsentList(t *testing.T) {
	t.Parallel()

	for name, configValue := range map[string]types.List{
		"null list":    types.ListNull(interfaceObjectType),
		"unknown list": types.ListUnknown(interfaceObjectType),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := errorPaths(t, configValue); len(got) != 0 {
				t.Errorf("got errors %v, want none", got)
			}
		})
	}
}

func TestGPUClusterInterfaceVariantValidatorSkipsAbsentElements(t *testing.T) {
	t.Parallel()

	list := ifaceList(t,
		types.ObjectNull(interfaceAttrTypes),
		types.ObjectUnknown(interfaceAttrTypes),
	)

	if got := errorPaths(t, list); len(got) != 0 {
		t.Errorf("got errors %v, want none", got)
	}
}
