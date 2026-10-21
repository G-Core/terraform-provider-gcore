package waap_domain_firewall_rule

import (
	"context"
	"fmt"
	"net/netip"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithValidateConfig = (*WaapDomainFirewallRuleResource)(nil)

// actionDurationRe is the API pattern for block.action_duration: a positive
// number, optionally followed by s, m, h or d.
var actionDurationRe = regexp.MustCompile(`^[1-9][0-9]*[smhd]?$`)

// nameRe is the API pattern for name. Go's \s is narrower than the API's
// (Unicode White_Space), so \v, U+0085 and Unicode spaces are added
// explicitly. U+001C-U+001F are rejected on purpose: only Python re counts
// them as whitespace, and the API may not.
var nameRe = regexp.MustCompile(`^[A-Za-z0-9\s\v\x{85}\p{Z}.:'";<>?&|\\_-]+$`)

// ValidateConfig checks at plan time what the API would otherwise reject on
// apply, and values the API would rewrite (which would show a change on every
// plan). Unknown values are skipped.
func (r *WaapDomainFirewallRuleResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var name, description types.String
	var action types.Object
	var conditions types.List
	for p, v := range map[string]any{"name": &name, "description": &description, "action": &action, "conditions": &conditions} {
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(p), v)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	validateLength(resp, path.Root("name"), name, 1, 100)
	// The API strips leading and trailing whitespace (Python str.strip).
	// strings.TrimSpace removes the same characters for every character nameRe
	// allows.
	if isSet(name) && name.ValueString() != "" {
		s := name.ValueString()
		if !nameRe.MatchString(s) {
			resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid name",
				fmt.Sprintf("name may only contain ASCII letters, digits, whitespace and . : ' \" ; < > ? & | \\ _ -; got %q.", s))
		} else if t := strings.TrimSpace(s); t == "" {
			resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid name",
				"name must contain at least one non-whitespace character; the API trims leading and trailing whitespace.")
		} else if t != s {
			resp.Diagnostics.AddAttributeError(path.Root("name"), "Invalid name",
				fmt.Sprintf("name must not start or end with whitespace; the API trims it, so %q would read back as %q and every plan would show a change. Use %q.", s, t, t))
		}
	}
	validateLength(resp, path.Root("description"), description, 0, 100)

	if isSet(action) {
		attrs := action.Attributes()
		exactlyOne(resp, path.Root("action"), "action", attrs, "allow", "block")
		if block, ok := attrs["block"].(types.Object); ok && isSet(block) {
			p := path.Root("action").AtName("block").AtName("action_duration")
			if d, ok := block.Attributes()["action_duration"].(types.String); ok && isSet(d) {
				if v := d.ValueString(); len(v) > 10 || !actionDurationRe.MatchString(v) {
					resp.Diagnostics.AddAttributeError(p, "Invalid action_duration",
						fmt.Sprintf("action_duration must be a positive number of seconds, optionally followed by s, m, h or d (for example 30, 60s, 12h), at most 10 characters; got %q.", v))
				}
			}
		}
	}

	if isSet(conditions) {
		elems := conditions.Elements()
		if len(elems) != 1 {
			resp.Diagnostics.AddAttributeError(path.Root("conditions"), "Invalid conditions",
				fmt.Sprintf("conditions must have exactly one entry, got %d.", len(elems)))
		}
		for i, e := range elems {
			if c, ok := e.(types.Object); ok && isSet(c) {
				cp := path.Root("conditions").AtListIndex(i)
				exactlyOne(resp, cp, fmt.Sprintf("conditions[%d]", i), c.Attributes(), "ip", "ip_range")
				if ip, ok := c.Attributes()["ip"].(types.Object); ok && isSet(ip) {
					validateIP(resp, cp.AtName("ip").AtName("ip_address"), ip.Attributes()["ip_address"])
				}
				if r, ok := c.Attributes()["ip_range"].(types.Object); ok && isSet(r) {
					for _, k := range []string{"lower_bound", "upper_bound"} {
						validateIP(resp, cp.AtName("ip_range").AtName(k), r.Attributes()[k])
					}
				}
			}
		}
	}
}

func isSet(v attr.Value) bool {
	return v != nil && !v.IsNull() && !v.IsUnknown()
}

// exactlyOne reports an error on p unless exactly one of the two keys is set.
// It stays silent while either is unknown.
func exactlyOne(resp *resource.ValidateConfigResponse, p path.Path, label string, attrs map[string]attr.Value, a, b string) {
	va, vb := attrs[a], attrs[b]
	if (va != nil && va.IsUnknown()) || (vb != nil && vb.IsUnknown()) {
		return
	}
	if isSet(va) == isSet(vb) {
		resp.Diagnostics.AddAttributeError(p, "Invalid "+label,
			fmt.Sprintf("%s must set exactly one of %s or %s.", label, a, b))
	}
}

func validateLength(resp *resource.ValidateConfigResponse, p path.Path, v types.String, min, max int) {
	if !isSet(v) {
		return
	}
	if n := utf8.RuneCountInString(v.ValueString()); n < min || n > max {
		resp.Diagnostics.AddAttributeError(p, "Invalid "+p.String(),
			fmt.Sprintf("%s must be %d to %d characters long, got %d.", p, min, max, n))
	}
}

// validateIP accepts one address, written exactly as the API stores it. The
// API (Python ipaddress) rewrites IPv6 to RFC 5952 form, IPv4-mapped
// addresses included (::ffff:c000:20a is stored as ::ffff:192.0.2.10);
// netip's String() gives the same form. netip parses only canonical IPv4, so
// IPv4 needs no further check. Unknown values are skipped.
func validateIP(resp *resource.ValidateConfigResponse, p path.Path, v attr.Value) {
	s, ok := v.(types.String)
	if !ok || !isSet(s) {
		return
	}
	in := s.ValueString()
	a, err := netip.ParseAddr(in)
	switch {
	case err != nil:
		resp.Diagnostics.AddAttributeError(p, "Invalid IP address",
			fmt.Sprintf("%s must be a single IPv4 or IPv6 address with no CIDR prefix or spaces; got %q. Use ip_range for a range.", p, in))
	case a.Zone() != "":
		resp.Diagnostics.AddAttributeError(p, "Invalid IP address",
			fmt.Sprintf("%s: IPv6 zone identifiers are not supported; got %q. Use %q.", p, in, a.WithZone("").String()))
	case a.String() != in:
		resp.Diagnostics.AddAttributeError(p, "Non-canonical IPv6 address",
			fmt.Sprintf("%s: the API stores IPv6 addresses in RFC 5952 form (lowercase, longest zero run compressed), so %q would read back as %q and every plan would show a change. Use %q.", p, in, a.String(), a.String()))
	}
}
