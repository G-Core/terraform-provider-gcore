package cloud_load_balancer

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// importAdoptionPrivateKey marks a load balancer whose create-only attributes
// still have to be adopted from config into state after a `terraform import`.
//
// It is written by ImportState, read by the plan modifiers on those attributes
// and by stripAdoptedCreateOnlyFields, and cleared by Update once the adoption
// has happened.
//
// The marker exists because "prior state is null" does NOT by itself mean "just
// imported". It is equally true of a load balancer created with the attribute
// omitted - vip_network_id defaults to the external network - and a user may set
// such an attribute on day two. Treating that as an adoption would plan a clean
// in-place update, send nothing to the API, and write the configured value to
// state, leaving state describing a load balancer that does not exist. Keying on
// an explicit import marker keeps that case on the replacement path, where it
// belongs.
//
// The marker alone cannot carry the whole distinction either: it can outlive
// the import. A config that sets neither attribute produces no-op plans after
// the import, no update ever runs, and Terraform never persists a
// private-state clear made during a no-op plan - so there is no hook that can
// reliably retire the marker before a day-two edit arrives. Safety therefore
// does not rest on the marker's lifetime: every adoption is verified against
// the load balancer the API returns before anything is written to state
// (verifyAdoptedCreateOnlyFields), and the marker is cleared only after a
// successful adopting apply. A stale marker can at worst turn a
// destroy+recreate plan into a loud apply-time error, never into silent state.
const importAdoptionPrivateKey = "lb_import_adoption_pending"

var importAdoptionPrivateValue = []byte(`{"pending":true}`)

// createOnlyFields maps the JSON keys of the load balancer attributes the create
// endpoint accepts but the update endpoint does not - the PATCH body supports
// only name, logging, preferred_connectivity and tags - to a predicate reporting
// whether the prior state for that attribute was null.
//
// A single map rather than a list plus a parallel lookup table, so the two
// cannot drift apart. The keys must match the `json:` tags in model.go;
// TestStripAdoptedCreateOnlyFieldsAfterImport round-trips a real
// MarshalJSONForUpdate body to catch a rename.
var createOnlyFields = map[string]func(CloudLoadBalancerModel) bool{
	"vip_network_id": func(m CloudLoadBalancerModel) bool { return m.VipNetworkID.IsNull() },
	"vip_subnet_id":  func(m CloudLoadBalancerModel) bool { return m.VipSubnetID.IsNull() },
}

// privateStateReader and privateStateWriter are the subsets of the framework's
// private state data used here. The concrete type lives in an internal framework
// package, so it is reached through interfaces rather than named directly.
type privateStateReader interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
}

type privateStateWriter interface {
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

// markImportAdoptionPending records that the create-only attributes of a
// freshly imported load balancer still have to be adopted into state.
func markImportAdoptionPending(ctx context.Context, private privateStateWriter) diag.Diagnostics {
	if private == nil {
		return nil
	}

	return private.SetKey(ctx, importAdoptionPrivateKey, importAdoptionPrivateValue)
}

// importAdoptionPending reports whether the load balancer is in the window
// between an import and the apply that adopts its create-only attributes.
func importAdoptionPending(ctx context.Context, private privateStateReader) (bool, diag.Diagnostics) {
	if private == nil {
		return false, nil
	}

	value, diags := private.GetKey(ctx, importAdoptionPrivateKey)

	return len(value) > 0, diags
}

// clearImportAdoptionPending closes the adoption window. Setting an empty value
// deletes the key.
func clearImportAdoptionPending(ctx context.Context, private privateStateWriter) diag.Diagnostics {
	if private == nil {
		return nil
	}

	return private.SetKey(ctx, importAdoptionPrivateKey, nil)
}

// stripAdoptedCreateOnlyFields removes the create-only fields from a PATCH body
// produced by MarshalJSONForUpdate, so that adopting them into state after an
// import sends nothing the update endpoint would reject. It returns "{}" when
// nothing else is left to send, which the caller treats as "no changes".
//
// It only ever strips while adoption is pending, and only for attributes whose
// prior state was null. Anything else is left in place: if a create-only field
// reaches the request body outside the adoption window it is a genuine bug, and
// failing against the API is better than swallowing it.
func stripAdoptedCreateOnlyFields(body []byte, state CloudLoadBalancerModel, adopting bool) ([]byte, error) {
	if !adopting {
		return body, nil
	}

	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" || trimmed == "null" || trimmed == "{}" {
		return body, nil
	}

	fields := map[string]json.RawMessage{}
	if err := json.Unmarshal(body, &fields); err != nil {
		return nil, fmt.Errorf("failed to inspect update request body: %w", err)
	}

	var changed bool
	for name, priorStateIsNull := range createOnlyFields {
		if _, ok := fields[name]; !ok {
			continue
		}
		if !priorStateIsNull(state) {
			continue
		}
		delete(fields, name)
		changed = true
	}

	if !changed {
		return body, nil
	}

	return json.Marshal(fields)
}

// adoptionVerdict classifies one adopted attribute against what the API
// reports about the live load balancer.
type adoptionVerdict int

const (
	// adoptionConfirmed: the configured value matches the API's evidence.
	adoptionConfirmed adoptionVerdict = iota
	// adoptionUnverifiable: the API response carried no evidence either way.
	adoptionUnverifiable
	// adoptionMismatch: the API's evidence contradicts the configured value.
	adoptionMismatch
)

// adoptionFinding is the verdict for a single adopted attribute, with a
// human-readable detail for the diagnostic.
type adoptionFinding struct {
	attribute string
	verdict   adoptionVerdict
	detail    string
}

// adoptionEvidence is the slice of a load balancer GET response that can
// corroborate or contradict an adoption. Both adopted attributes are
// no_refresh - the API never echoes them - but the response does describe the
// subnets the balancer's addresses live in, which is enough to check every
// adopted value indirectly.
//
// Slices are left nil when the key is absent from the response, so "the API
// said nothing" (unverifiable) is distinguishable from "the API said none"
// (hard evidence).
type adoptionEvidence struct {
	VrrpIPs []struct {
		SubnetID string `json:"subnet_id"`
	} `json:"vrrp_ips"`
	AdditionalVips []struct {
		SubnetID string `json:"subnet_id"`
	} `json:"additional_vips"`
}

func parseAdoptionEvidence(body []byte) (adoptionEvidence, error) {
	var evidence adoptionEvidence
	if err := json.Unmarshal(body, &evidence); err != nil {
		return adoptionEvidence{}, fmt.Errorf("failed to parse load balancer response: %w", err)
	}

	return evidence, nil
}

// subnetIDs returns the distinct subnets the API places this load balancer in,
// in response order. An empty result means the response carried no placement
// evidence at all.
func (e adoptionEvidence) subnetIDs() []string {
	var ids []string
	for _, vrrp := range e.VrrpIPs {
		if vrrp.SubnetID != "" && !slices.Contains(ids, vrrp.SubnetID) {
			ids = append(ids, vrrp.SubnetID)
		}
	}
	for _, vip := range e.AdditionalVips {
		if vip.SubnetID != "" && !slices.Contains(ids, vip.SubnetID) {
			ids = append(ids, vip.SubnetID)
		}
	}

	return ids
}

// subnetNetworkLookup resolves a subnet ID to its network ID. It is injected
// so the verification below stays a pure function: the real implementation is
// one Networks.Subnets.Get call, tests hand in a map.
type subnetNetworkLookup func(ctx context.Context, subnetID string) (string, error)

// verifyAdoptedCreateOnlyFields checks every create-only value this apply
// would adopt into state - plan value set, prior state null - against the
// evidence in the load balancer response. Adoption writes config to state
// without sending anything to the API, so this check is the only thing
// standing between a wrong configured value and state that permanently lies
// about live infrastructure (no_refresh means no later refresh can correct
// it). It returns one finding per adopted attribute; attributes that are not
// being adopted produce no finding. A lookup failure aborts with an error so
// the caller can fail the apply and retry with the adoption window intact.
func verifyAdoptedCreateOnlyFields(ctx context.Context, plan, state CloudLoadBalancerModel, evidence adoptionEvidence, subnetNetwork subnetNetworkLookup) ([]adoptionFinding, error) {
	var findings []adoptionFinding

	evidenceSubnets := evidence.subnetIDs()

	if !plan.VipNetworkID.IsNull() && state.VipNetworkID.IsNull() {
		switch {
		case len(evidenceSubnets) == 0:
			findings = append(findings, adoptionFinding{
				attribute: "vip_network_id",
				verdict:   adoptionUnverifiable,
				detail:    "The API reported no subnets (vrrp_ips, additional_vips) for this load balancer, so the configured network cannot be checked.",
			})
		default:
			var networks []string
			for _, subnetID := range evidenceSubnets {
				networkID, err := subnetNetwork(ctx, subnetID)
				if err != nil {
					return nil, fmt.Errorf("failed to resolve subnet %s to its network: %w", subnetID, err)
				}
				if !slices.Contains(networks, networkID) {
					networks = append(networks, networkID)
				}
			}
			if slices.Contains(networks, plan.VipNetworkID.ValueString()) {
				findings = append(findings, adoptionFinding{attribute: "vip_network_id", verdict: adoptionConfirmed})
			} else {
				findings = append(findings, adoptionFinding{
					attribute: "vip_network_id",
					verdict:   adoptionMismatch,
					detail: fmt.Sprintf("The load balancer's addresses live in network(s) %s, not %q.",
						strings.Join(networks, ", "), plan.VipNetworkID.ValueString()),
				})
			}
		}
	}

	if !plan.VipSubnetID.IsNull() && state.VipSubnetID.IsNull() {
		switch {
		case len(evidenceSubnets) == 0:
			findings = append(findings, adoptionFinding{
				attribute: "vip_subnet_id",
				verdict:   adoptionUnverifiable,
				detail:    "The API reported no subnets (vrrp_ips, additional_vips) for this load balancer, so the configured subnet cannot be checked.",
			})
		case slices.Contains(evidenceSubnets, plan.VipSubnetID.ValueString()):
			findings = append(findings, adoptionFinding{attribute: "vip_subnet_id", verdict: adoptionConfirmed})
		default:
			findings = append(findings, adoptionFinding{
				attribute: "vip_subnet_id",
				verdict:   adoptionMismatch,
				detail: fmt.Sprintf("The load balancer's addresses live in subnet(s) %s, not %q.",
					strings.Join(evidenceSubnets, ", "), plan.VipSubnetID.ValueString()),
			})
		}
	}

	return findings, nil
}
