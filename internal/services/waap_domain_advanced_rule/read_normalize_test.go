package waap_domain_advanced_rule

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Import starts from an empty model: an empty-object action kind and a missing
// description must decode to the values a matching config plans.
func TestUnmarshalReadImportedRule(t *testing.T) {
	t.Parallel()

	data := new(WaapDomainAdvancedRuleModel)
	raw := []byte(`{"name":"r","enabled":false,"action":{"monitor":{}},"source":"s","phase":"access","id":7}`)
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
}

// A refresh after the action kind changed out of band must drop the old kind.
func TestUnmarshalReadActionKindChangedOutOfBand(t *testing.T) {
	t.Parallel()

	data := &WaapDomainAdvancedRuleModel{
		Description: types.StringValue("d"),
		Action: &WaapDomainAdvancedRuleActionModel{
			Monitor: &map[string]jsontypes.Normalized{},
			Block:   &WaapDomainAdvancedRuleActionBlockModel{StatusCode: types.Int64Value(403)},
		},
	}
	raw := []byte(`{"name":"r","description":"d","enabled":true,"action":{"allow":{}},"source":"s","phase":"access","id":7}`)
	if err := unmarshalRead(raw, &data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Action.Monitor != nil || data.Action.Block != nil {
		t.Errorf("expected monitor and block dropped, got %+v", data.Action)
	}
	if data.Action.Allow == nil || *data.Action.Allow == nil {
		t.Errorf("expected allow = {}, got %+v", data.Action)
	}
	if data.Description.ValueString() != "d" {
		t.Errorf("expected description kept, got %v", data.Description)
	}
}

// A kind the response returns as {} replaces a prior map that had entries.
func TestUnmarshalReadEmptyKindReplacesPriorEntries(t *testing.T) {
	t.Parallel()

	data := &WaapDomainAdvancedRuleModel{
		Action: &WaapDomainAdvancedRuleActionModel{
			Allow: &map[string]jsontypes.Normalized{"foo": jsontypes.NewNormalizedValue(`"bar"`)},
		},
	}
	raw := []byte(`{"name":"r","enabled":true,"action":{"allow":{}},"source":"s","phase":"access","id":7}`)
	if err := unmarshalRead(raw, &data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Action.Allow == nil || *data.Action.Allow == nil || len(*data.Action.Allow) != 0 {
		t.Errorf("expected allow = {}, got %+v", data.Action.Allow)
	}
}
