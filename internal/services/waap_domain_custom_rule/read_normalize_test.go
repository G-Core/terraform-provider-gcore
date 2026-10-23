package waap_domain_custom_rule

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const conditionsJSON = `"conditions":[{"ip":{"negation":false,"ip_address":"192.0.2.10"}}]`

// Import starts from an empty model: an empty-object action kind and a missing
// description must decode to the values a matching config plans.
func TestUnmarshalReadImportedRule(t *testing.T) {
	t.Parallel()

	data := new(WaapDomainCustomRuleModel)
	raw := []byte(`{"name":"r","enabled":false,"action":{"monitor":{}},"id":7,` + conditionsJSON + `}`)
	if err := unmarshalRead(raw, &data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Description.ValueString() != "" || data.Description.IsNull() {
		t.Errorf("expected description \"\", got %v", data.Description)
	}
	if data.Action == nil || data.Action.Monitor == nil || *data.Action.Monitor == nil || len(*data.Action.Monitor) != 0 {
		t.Fatalf("expected monitor = {}, got %+v", data.Action)
	}
	if data.Action.Allow != nil || data.Action.Block != nil {
		t.Errorf("expected only monitor set, got %+v", data.Action)
	}
	if data.Conditions == nil || len(*data.Conditions) != 1 || (*data.Conditions)[0].IP == nil {
		t.Errorf("expected one ip condition, got %+v", data.Conditions)
	}
}

// A refresh after the action kind changed out of band must drop the old kind.
func TestUnmarshalReadActionKindChangedOutOfBand(t *testing.T) {
	t.Parallel()

	data := &WaapDomainCustomRuleModel{
		Description: types.StringValue("d"),
		Action: &WaapDomainCustomRuleActionModel{
			Monitor: &map[string]jsontypes.Normalized{},
			Block:   &WaapDomainCustomRuleActionBlockModel{StatusCode: types.Int64Value(403)},
		},
	}
	raw := []byte(`{"name":"r","description":"d","enabled":true,"action":{"captcha":{}},"id":7,` + conditionsJSON + `}`)
	if err := unmarshalRead(raw, &data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Action.Monitor != nil || data.Action.Block != nil {
		t.Errorf("expected monitor and block dropped, got %+v", data.Action)
	}
	if data.Action.Captcha == nil || *data.Action.Captcha == nil {
		t.Errorf("expected captcha = {}, got %+v", data.Action)
	}
	if data.Description.ValueString() != "d" {
		t.Errorf("expected description kept, got %v", data.Description)
	}
}

// A legacy rule with a response-only condition imports it as an empty
// condition and warns at its index; other conditions are untouched.
func TestImportResponseOnlyCondition(t *testing.T) {
	t.Parallel()

	raw := []byte(`{"name":"r","enabled":true,"action":{"monitor":{}},"id":7,"conditions":[` +
		`{"ip":{"negation":false,"ip_address":"192.0.2.10"}},` +
		`{"response_header":{"negation":false,"header":"X-Foo","value":"bar","match_type":"Contains"}}]}`)

	data := new(WaapDomainCustomRuleModel)
	if err := unmarshalRead(raw, &data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Conditions == nil || len(*data.Conditions) != 2 {
		t.Fatalf("expected two conditions, got %+v", data.Conditions)
	}
	if (*data.Conditions)[0].IP == nil {
		t.Errorf("expected the ip condition kept, got %+v", (*data.Conditions)[0])
	}
	if c := (*data.Conditions)[1]; c != nil && *c != (WaapDomainCustomRuleConditionsModel{}) {
		t.Errorf("expected the response_header condition to decode empty, got %+v", c)
	}

	diags := responseOnlyConditionWarnings(raw)
	if diags.WarningsCount() != 1 || diags.HasError() {
		t.Fatalf("expected one warning, got %v", diags)
	}
	if got := diags[0].(diag.DiagnosticWithPath).Path(); !got.Equal(path.Root("conditions").AtListIndex(1)) {
		t.Errorf("expected the warning on conditions[1], got %s", got)
	}

	if w := responseOnlyConditionWarnings([]byte(`{"conditions":[{"ip":{"ip_address":"192.0.2.10"}}]}`)); len(w) != 0 {
		t.Errorf("expected no warning without response-only conditions, got %v", w)
	}
}
