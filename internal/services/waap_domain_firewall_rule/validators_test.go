package waap_domain_firewall_rule

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func validRule() *WaapDomainFirewallRuleModel {
	return &WaapDomainFirewallRuleModel{
		DomainID: types.Int64Value(1),
		Enabled:  types.BoolValue(true),
		Name:     types.StringValue("block office"),
		Action: &WaapDomainFirewallRuleActionModel{Block: &WaapDomainFirewallRuleActionBlockModel{
			ActionDuration: types.StringValue("12h"),
			StatusCode:     types.Int64Value(403),
		}},
		Conditions: &[]*WaapDomainFirewallRuleConditionsModel{{
			IP: &WaapDomainFirewallRuleConditionsIPModel{IPAddress: types.StringValue("192.0.2.10"), Negation: types.BoolValue(false)},
		}},
		Description: types.StringValue(""),
	}
}

// unknownAt marks the attribute at the given tftypes path as unknown.
func unknownAt(steps ...tftypes.AttributePathStep) func(tftypes.Value) tftypes.Value {
	want := tftypes.NewAttributePathWithSteps(steps)
	return func(v tftypes.Value) tftypes.Value {
		out, err := tftypes.Transform(v, func(p *tftypes.AttributePath, v tftypes.Value) (tftypes.Value, error) {
			if p.Equal(want) {
				return tftypes.NewValue(v.Type(), tftypes.UnknownValue), nil
			}
			return v, nil
		})
		if err != nil {
			panic(err)
		}
		return out
	}
}

func TestValidateConfig(t *testing.T) {
	ip := &WaapDomainFirewallRuleConditionsIPModel{IPAddress: types.StringValue("192.0.2.10"), Negation: types.BoolValue(false)}
	ipRange := &WaapDomainFirewallRuleConditionsIPRangeModel{
		LowerBound: types.StringValue("192.0.2.1"), UpperBound: types.StringValue("192.0.2.9"), Negation: types.BoolValue(false),
	}

	tests := []struct {
		name    string
		mutate  func(m *WaapDomainFirewallRuleModel)
		raw     func(tftypes.Value) tftypes.Value
		wantErr string
	}{
		{name: "valid block"},
		{name: "valid allow", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Action = &WaapDomainFirewallRuleActionModel{Allow: emptyAllow()}
		}},
		{name: "valid ip_range, bare seconds", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Action.Block.ActionDuration = types.StringValue("30")
			(*m.Conditions)[0] = &WaapDomainFirewallRuleConditionsModel{IPRange: ipRange}
		}},
		{name: "valid 100-char name and description", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue(strings.Repeat("n", 100))
			m.Description = types.StringValue(strings.Repeat("d", 100))
		}},
		{name: "both allow and block", wantErr: "exactly one of allow or block", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Action.Allow = emptyAllow()
		}},
		{name: "neither allow nor block", wantErr: "exactly one of allow or block", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Action = &WaapDomainFirewallRuleActionModel{}
		}},
		{name: "no conditions", wantErr: "exactly one entry, got 0", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Conditions = &[]*WaapDomainFirewallRuleConditionsModel{}
		}},
		{name: "two conditions", wantErr: "exactly one entry, got 2", mutate: func(m *WaapDomainFirewallRuleModel) {
			*m.Conditions = append(*m.Conditions, &WaapDomainFirewallRuleConditionsModel{IP: ip})
		}},
		{name: "condition with ip and ip_range", wantErr: "conditions[0] must set exactly one of ip or ip_range", mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IPRange = ipRange
		}},
		{name: "empty condition", wantErr: "conditions[0] must set exactly one of ip or ip_range", mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0] = &WaapDomainFirewallRuleConditionsModel{}
		}},
		{name: "duration with unit only", wantErr: "action_duration", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Action.Block.ActionDuration = types.StringValue("h")
		}},
		{name: "duration zero", wantErr: "action_duration", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Action.Block.ActionDuration = types.StringValue("0s")
		}},
		{name: "duration bad unit", wantErr: "action_duration", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Action.Block.ActionDuration = types.StringValue("5w")
		}},
		{name: "duration too long", wantErr: "action_duration", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Action.Block.ActionDuration = types.StringValue("12345678901")
		}},
		{name: "empty name", wantErr: "name must be 1 to 100", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue("")
		}},
		{name: "long name", wantErr: "name must be 1 to 100", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue(strings.Repeat("n", 101))
		}},
		{name: "name with forbidden characters", wantErr: "name may only contain", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue("block/office!")
		}},
		{name: "non-ASCII name", wantErr: "name may only contain", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue("блок")
		}},
		{name: "valid name with all allowed punctuation", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue(`Block .:'";<>?&|\_- 1`)
		}},
		{name: "valid name with a no-break space", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue("Block office")
		}},
		// Values the API would rewrite are rejected with the canonical form.
		{name: "name with surrounding spaces", wantErr: `Use "rule"`, mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue(" rule ")
		}},
		{name: "name with leading tab", wantErr: "must not start or end with whitespace", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue("\trule")
		}},
		{name: "name with trailing no-break space", wantErr: "must not start or end with whitespace", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue("rule\u00a0")
		}},
		{name: "blank name", wantErr: "at least one non-whitespace", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue("   ")
		}},
		{name: "valid name with inner double space", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue("a  b")
		}},
		{name: "valid name with inner next-line character", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue("a\u0085b")
		}},
		{name: "name with trailing next-line character", wantErr: "must not start or end with whitespace", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue("rule\u0085")
		}},
		{name: "name with information separator", wantErr: "name may only contain", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue("a\x1cb")
		}},
		{name: "valid name with inner tab", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Name = types.StringValue("a\tb")
		}},
		{name: "uppercase IPv6", wantErr: `Use "2001:db8::20"`, mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("2001:DB8::20")
		}},
		{name: "IPv6 zero run not compressed", wantErr: `Use "2001:db8::1"`, mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("2001:db8:0:0::1")
		}},
		{name: "IPv6 shorter zero run compressed", wantErr: `Use "2001:0:0:1::1"`, mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("2001:0:0:1:0:0:0:1")
		}},
		{name: "IPv6 with embedded IPv4", wantErr: `Use "2001:db8::c000:201"`, mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("2001:db8::192.0.2.1")
		}},
		{name: "valid IPv4-mapped IPv6", mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("::ffff:192.0.2.10")
		}},
		{name: "hex IPv4-mapped IPv6", wantErr: `Use "::ffff:192.0.2.10"`, mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("::ffff:c000:20a")
		}},
		{name: "IPv6 with zone", wantErr: `Use "fe80::1"`, mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("fe80::1%eth0")
		}},
		{name: "CIDR in ip_address", wantErr: "Use ip_range", mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("10.0.0.0/24")
		}},
		{name: "IPv4 with leading zero", wantErr: "single IPv4 or IPv6", mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("192.0.2.010")
		}},
		{name: "IPv4 with leading space", wantErr: "single IPv4 or IPv6", mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue(" 192.0.2.10")
		}},
		{name: "valid canonical IPv6", mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("2001:db8::20")
		}},
		{name: "valid unspecified IPv6", mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("::")
		}},
		{name: "valid IPv6 single zero group not compressed", mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("2001:db8:0:1:1:1:1:1")
		}},
		{name: "uppercase IPv6 lower_bound", wantErr: "conditions[0].ip_range.lower_bound", mutate: func(m *WaapDomainFirewallRuleModel) {
			r := *ipRange
			r.LowerBound = types.StringValue("2001:DB8::1")
			r.UpperBound = types.StringValue("2001:db8::ff")
			(*m.Conditions)[0] = &WaapDomainFirewallRuleConditionsModel{IPRange: &r}
		}},
		{name: "non-canonical IPv6 upper_bound", wantErr: `Use "2001:db8::ff"`, mutate: func(m *WaapDomainFirewallRuleModel) {
			r := *ipRange
			r.LowerBound = types.StringValue("2001:db8::1")
			r.UpperBound = types.StringValue("2001:db8:0::ff")
			(*m.Conditions)[0] = &WaapDomainFirewallRuleConditionsModel{IPRange: &r}
		}},
		{name: "long description", wantErr: "description must be 0 to 100", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Description = types.StringValue(strings.Repeat("d", 101))
		}},
		// Unknown values must never fail validation.
		{name: "unknown action", raw: unknownAt(tftypes.AttributeName("action"))},
		{name: "unknown allow next to block", raw: unknownAt(tftypes.AttributeName("action"), tftypes.AttributeName("allow"))},
		{name: "unknown action_duration", mutate: func(m *WaapDomainFirewallRuleModel) {
			m.Action.Block.ActionDuration = types.StringValue("bad")
		}, raw: unknownAt(tftypes.AttributeName("action"), tftypes.AttributeName("block"), tftypes.AttributeName("action_duration"))},
		{name: "unknown conditions", raw: unknownAt(tftypes.AttributeName("conditions"))},
		{name: "unknown ip_range next to ip", raw: unknownAt(tftypes.AttributeName("conditions"), tftypes.ElementKeyInt(0), tftypes.AttributeName("ip_range"))},
		{name: "unknown name", raw: unknownAt(tftypes.AttributeName("name"))},
		{name: "unknown ip_address", mutate: func(m *WaapDomainFirewallRuleModel) {
			(*m.Conditions)[0].IP.IPAddress = types.StringValue("2001:DB8::20")
		}, raw: unknownAt(tftypes.AttributeName("conditions"), tftypes.ElementKeyInt(0), tftypes.AttributeName("ip"), tftypes.AttributeName("ip_address"))},
	}

	ctx := context.Background()
	schema := ResourceSchema(ctx)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := validRule()
			if tt.mutate != nil {
				tt.mutate(m)
			}
			state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(schema.Type().TerraformType(ctx), nil)}
			if diags := state.Set(ctx, m); diags.HasError() {
				t.Fatal(diags)
			}
			raw := state.Raw
			if tt.raw != nil {
				raw = tt.raw(raw)
			}

			resp := &resource.ValidateConfigResponse{}
			(&WaapDomainFirewallRuleResource{}).ValidateConfig(ctx, resource.ValidateConfigRequest{
				Config: tfsdk.Config{Schema: schema, Raw: raw},
			}, resp)

			if tt.wantErr == "" {
				if resp.Diagnostics.HasError() {
					t.Fatalf("unexpected error: %v", resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.ErrorsCount() != 1 {
				t.Fatalf("want one error containing %q, got %v", tt.wantErr, resp.Diagnostics)
			}
			if d := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(d, tt.wantErr) {
				t.Fatalf("error %q does not contain %q", d, tt.wantErr)
			}
		})
	}
}
