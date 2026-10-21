package waap_domain_firewall_rule

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const (
	allowBody = `{"id":7,"name":"r","enabled":true,"description":"","action":{"allow":{}},` +
		`"conditions":[{"ip":{"ip_address":"192.0.2.10","negation":false}}]}`
	blockBody = `{"id":7,"name":"r","enabled":true,"description":"","action":{"block":{"status_code":403}},` +
		`"conditions":[{"ip":{"ip_address":"192.0.2.10","negation":false}}]}`
)

// readIntoState runs the Read/ImportState decode path (apijson.Unmarshal then
// syncFromBody) on prior and writes the result to a state, so the test
// sees the values Terraform would.
func readIntoState(t *testing.T, prior *WaapDomainFirewallRuleModel, body string) (allow types.Map, block types.Object) {
	t.Helper()
	ctx := context.Background()

	data := prior
	if err := apijson.Unmarshal([]byte(body), &data); err != nil {
		t.Fatal(err)
	}
	syncFromBody([]byte(body), data)

	schema := ResourceSchema(ctx)
	state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(schema.Type().TerraformType(ctx), nil)}
	if diags := state.Set(ctx, data); diags.HasError() {
		t.Fatal(diags)
	}
	if diags := state.GetAttribute(ctx, path.Root("action").AtName("allow"), &allow); diags.HasError() {
		t.Fatal(diags)
	}
	if diags := state.GetAttribute(ctx, path.Root("action").AtName("block"), &block); diags.HasError() {
		t.Fatal(diags)
	}
	return allow, block
}

func TestReadKeepsEmptyAllow(t *testing.T) {
	// ImportState starts from an empty model.
	allow, block := readIntoState(t, &WaapDomainFirewallRuleModel{ID: types.Int64Value(7), DomainID: types.Int64Value(1)}, allowBody)
	if allow.IsNull() || len(allow.Elements()) != 0 {
		t.Fatalf("import: allow = %s, want {}", allow)
	}
	if !block.IsNull() {
		t.Fatalf("import: block = %s, want null", block)
	}
}

func TestReadFollowsActionSwitch(t *testing.T) {
	blockPrior := func() *WaapDomainFirewallRuleModel {
		return &WaapDomainFirewallRuleModel{
			ID: types.Int64Value(7), DomainID: types.Int64Value(1),
			Action: &WaapDomainFirewallRuleActionModel{Block: &WaapDomainFirewallRuleActionBlockModel{
				ActionDuration: types.StringValue("5m"), StatusCode: types.Int64Value(429),
			}},
		}
	}
	allow, block := readIntoState(t, blockPrior(), allowBody)
	if allow.IsNull() || !block.IsNull() {
		t.Fatalf("block -> allow: allow = %s, block = %s, want {} and null", allow, block)
	}

	allowPrior := &WaapDomainFirewallRuleModel{
		ID: types.Int64Value(7), DomainID: types.Int64Value(1),
		Action: &WaapDomainFirewallRuleActionModel{Allow: emptyAllow()},
	}
	allow, block = readIntoState(t, allowPrior, blockBody)
	if !allow.IsNull() || block.IsNull() {
		t.Fatalf("allow -> block: allow = %s, block = %s, want null and an object", allow, block)
	}
	if got := block.Attributes()["status_code"].(basetypes.Int64Value); got.ValueInt64() != 403 {
		t.Fatalf("allow -> block: status_code = %s, want 403", got)
	}
}

func emptyAllow() *map[string]jsontypes.Normalized {
	return &map[string]jsontypes.Normalized{}
}

func TestReadMissingDescriptionIsEmpty(t *testing.T) {
	// The API omits description when the rule has none; import must match the "" default.
	body := `{"id":7,"name":"r","enabled":true,"action":{"allow":{}},` +
		`"conditions":[{"ip":{"ip_address":"192.0.2.10","negation":false}}]}`
	data := &WaapDomainFirewallRuleModel{ID: types.Int64Value(7), DomainID: types.Int64Value(1)}
	if err := apijson.Unmarshal([]byte(body), &data); err != nil {
		t.Fatal(err)
	}
	syncFromBody([]byte(body), data)
	if data.Description.IsNull() || data.Description.ValueString() != "" {
		t.Fatalf("description = %s, want \"\"", data.Description)
	}

	data = &WaapDomainFirewallRuleModel{ID: types.Int64Value(7), DomainID: types.Int64Value(1)}
	withDesc := `{"id":7,"name":"r","enabled":true,"description":"office","action":{"allow":{}},` +
		`"conditions":[{"ip":{"ip_address":"192.0.2.10","negation":false}}]}`
	if err := apijson.Unmarshal([]byte(withDesc), &data); err != nil {
		t.Fatal(err)
	}
	syncFromBody([]byte(withDesc), data)
	if data.Description.ValueString() != "office" {
		t.Fatalf("description = %s, want office", data.Description)
	}
}

func TestReadBodyWithoutActionKeepsAction(t *testing.T) {
	data := &WaapDomainFirewallRuleModel{Action: &WaapDomainFirewallRuleActionModel{Allow: emptyAllow()}}
	syncFromBody([]byte(`{"id":7}`), data)
	if data.Action.Allow == nil {
		t.Fatal("allow cleared although the body has no action key")
	}
}
