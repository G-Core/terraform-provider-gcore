package customvalidator

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ validator.List = GPUClusterInterfaceVariantValidator{}

// GPUClusterInterfaceVariantValidator enforces the per-variant request contract for
// GPU cluster server interfaces.
//
// The API models interfaces as a discriminated union on `type`, where each variant
// accepts a different set of fields. Codegen flattens that union into a single flat
// object exposing every field on every variant, so nothing stops a user from setting
// a field the selected variant does not accept.
//
// That is not a harmless no-op. The API silently drops the extra field, Read then
// nulls it in state (apijson.Unmarshal sets absent keys to null), and the config-vs-state
// walk in planmodifiers.RequiresReplaceOnConfigChange sees "set in config, null in state"
// as a change - forcing a full cluster replacement on every plan, forever. See
// GCLOUD2-28207, where `ip_family` on a `subnet` interface did exactly this.
//
// Requiring the variant's mandatory fields is the same idea in reverse: without it,
// omitting one fails at apply with an API 400 instead of at plan.
//
// Used by both cloud_gpu_baremetal_cluster and cloud_gpu_virtual_cluster, whose
// interface contracts are identical.
type GPUClusterInterfaceVariantValidator struct{}

// interfaceVariantContract is the request contract for a single interface `type`.
type interfaceVariantContract struct {
	// allowed lists every field the API accepts for this variant. Anything set in
	// config and absent here is rejected.
	allowed map[string]struct{}
	// required lists fields the API demands for this variant.
	required []string
	// reasons holds a per-field explanation of why a prohibited field makes no
	// sense for this variant, so the diagnostic can say why and not just what.
	reasons map[string]string
}

// interfaceVariantContracts mirrors the SDK request params structs
// (GPU{Virtual,Baremetal}ClusterNewParamsServersSettingsInterface{External,Subnet,AnySubnet});
// fields tagged `api:"required"` there appear in required here.
var interfaceVariantContracts = map[string]interfaceVariantContract{
	"external": {
		allowed: fieldSet("type", "name", "ip_family", "port_security_enabled", "security_groups"),
		reasons: map[string]string{
			"network_id":  "a public interface is not attached to a specific network",
			"subnet_id":   "a public interface is not attached to a specific subnet",
			"floating_ip": "a public interface already has a public address",
		},
	},
	"subnet": {
		allowed:  fieldSet("type", "name", "port_security_enabled", "security_groups", "network_id", "subnet_id", "floating_ip"),
		required: []string{"network_id", "subnet_id"},
		reasons: map[string]string{
			"ip_family": "the IP family is determined by the referenced subnet",
		},
	},
	"any_subnet": {
		allowed:  fieldSet("type", "name", "port_security_enabled", "security_groups", "network_id", "ip_family", "floating_ip"),
		required: []string{"network_id"},
		reasons: map[string]string{
			"subnet_id": `use type = "subnet" to pin a specific subnet`,
		},
	},
}

func (v GPUClusterInterfaceVariantValidator) Description(_ context.Context) string {
	return "validates that each interface only sets fields its type accepts, and sets the ones its type requires"
}

func (v GPUClusterInterfaceVariantValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v GPUClusterInterfaceVariantValidator) ValidateList(_ context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	for i, elem := range req.ConfigValue.Elements() {
		if elem.IsNull() || elem.IsUnknown() {
			continue
		}

		obj, ok := elem.(basetypes.ObjectValue)
		if !ok {
			continue
		}
		attrs := obj.Attributes()

		typeValue, ok := attrs["type"]
		if !ok || typeValue.IsNull() || typeValue.IsUnknown() {
			continue
		}
		typeString, ok := typeValue.(basetypes.StringValue)
		if !ok {
			continue
		}

		// `type` is validated with stringvalidator.OneOfCaseInsensitive, so "SUBNET"
		// is a legal config value. Normalize before the lookup - a raw lookup would
		// miss and then report every configured field as prohibited.
		variant := strings.ToLower(typeString.ValueString())
		contract, ok := interfaceVariantContracts[variant]
		if !ok {
			// Unrecognized type. The OneOfCaseInsensitive validator on `type` already
			// reports it; don't double-report.
			continue
		}

		elemPath := req.Path.AtListIndex(i)

		// Fields this variant does not accept. Iterate in sorted order so diagnostics
		// come out deterministically rather than in Go's randomized map order.
		for _, field := range sortedKeys(attrs) {
			if _, allowed := contract.allowed[field]; allowed {
				continue
			}
			value := attrs[field]
			if value.IsNull() || value.IsUnknown() {
				continue
			}

			detail := fmt.Sprintf("'%s' is not supported when type = %q.", field, variant)
			if reason, ok := contract.reasons[field]; ok {
				detail = fmt.Sprintf("'%s' is not supported when type = %q; %s. Remove it.", field, variant, reason)
			}
			resp.Diagnostics.AddAttributeError(
				elemPath.AtName(field),
				"Invalid interface configuration",
				detail,
			)
		}

		// Fields this variant requires.
		for _, field := range contract.required {
			value, ok := attrs[field]
			if !ok {
				continue
			}
			// An unknown value is resolved at apply time - typically a reference to a
			// network or subnet that does not exist yet - so it must be accepted.
			if value.IsUnknown() || !isNullOrEmptyString(value) {
				continue
			}
			resp.Diagnostics.AddAttributeError(
				elemPath.AtName(field),
				"Missing required attribute",
				fmt.Sprintf("'%s' is required when type = %q.", field, variant),
			)
		}
	}
}

func fieldSet(fields ...string) map[string]struct{} {
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		set[f] = struct{}{}
	}
	return set
}

func sortedKeys(attrs map[string]attr.Value) []string {
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// isNullOrEmptyString reports whether a known value is absent for practical
// purposes: either null, or a string set to "".
func isNullOrEmptyString(value attr.Value) bool {
	if value.IsNull() {
		return true
	}
	if s, ok := value.(basetypes.StringValue); ok {
		return s.ValueString() == ""
	}
	return false
}
