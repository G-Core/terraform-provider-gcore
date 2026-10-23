package waap_domain_custom_rule

import (
	"context"
	"fmt"
	"net/netip"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// oneConditionKind matches the API rule that every condition sets exactly one
// condition type, so a wrong condition fails at plan time instead of with a
// 422 on apply.
type oneConditionKind struct{}

func (oneConditionKind) Description(_ context.Context) string {
	return "exactly one condition type must be set"
}

func (v oneConditionKind) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (oneConditionKind) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	var set []string
	for name, value := range req.ConfigValue.Attributes() {
		if value.IsUnknown() {
			// decided at apply time
			return
		}
		if !value.IsNull() {
			set = append(set, name)
		}
	}
	if len(set) != 1 {
		sort.Strings(set)
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid condition",
			fmt.Sprintf("Each condition must set exactly one condition type, got %d: %v.", len(set), set))
	}
}

// The API stores some condition values in a canonical form and returns that
// form, not the configured one. A non-canonical value would make the apply fail
// with an inconsistent result and then drift on every plan, so these validators
// reject it at plan time and name the form to use instead.

var (
	// the API lowercases JA3 and JA4 fingerprints
	ja3Fingerprint = stringvalidator.RegexMatches(
		regexp.MustCompile(`^[a-f0-9]{32}$`),
		"must be a lowercase 32-character hexadecimal JA3 fingerprint",
	)
	ja4Fingerprint = stringvalidator.RegexMatches(
		regexp.MustCompile(`^[tq][0-9]{2}[di][0-9]{2}[0-9]{2}[a-z0-9_]{2}_[a-f0-9]{12}_[a-f0-9]{12}$`),
		"must be a lowercase JA4 fingerprint",
	)
	// the API uppercases country codes
	countryCode = stringvalidator.RegexMatches(
		regexp.MustCompile(`^[A-Z]{2}$`),
		"must be an uppercase ISO 3166-1 alpha-2 country code",
	)
)

// canonicalIP accepts an IPv4 or IPv6 address only in the form the API
// returns it, for example "2001:db8::1" and not "2001:DB8:0::1".
type canonicalIP struct{}

func (canonicalIP) Description(_ context.Context) string {
	return "must be an IPv4 or IPv6 address in canonical form"
}

func (v canonicalIP) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (canonicalIP) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	addr, err := netip.ParseAddr(value)
	if err != nil || addr.Zone() != "" {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid IP address",
			fmt.Sprintf("%q is not an IPv4 or IPv6 address.", value))
		return
	}
	if canonical := addr.String(); canonical != value {
		resp.Diagnostics.AddAttributeError(req.Path, "Non-canonical IP address",
			fmt.Sprintf("The API stores %q as %q. Use %q to avoid a permanent diff.", value, canonical, canonical))
	}
}

// defaultOwnerTypes is what the API stores when an owner_types condition omits
// the list, as the spec's default declares.
var defaultOwnerTypes = types.ListValueMust(types.StringType, []attr.Value{types.StringValue("COMMERCIAL")})

// noSurroundingSpace rejects leading or trailing whitespace, which the API
// trims from the rule name and the header condition value.
type noSurroundingSpace struct{}

func (noSurroundingSpace) Description(_ context.Context) string {
	return "must not start or end with whitespace"
}

func (v noSurroundingSpace) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (noSurroundingSpace) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	if trimmed := strings.TrimSpace(value); trimmed != value {
		resp.Diagnostics.AddAttributeError(req.Path, "Leading or trailing whitespace",
			fmt.Sprintf("The API trims %q to %q. Use %q to avoid a permanent diff.", value, trimmed, trimmed))
	}
}
