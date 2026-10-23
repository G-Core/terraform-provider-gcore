package waap_domain_custom_rule

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func customRuleModel(action *WaapDomainCustomRuleActionModel) WaapDomainCustomRuleModel {
	return WaapDomainCustomRuleModel{
		ID:          types.Int64Value(1),
		DomainID:    types.Int64Value(2),
		Enabled:     types.BoolValue(true),
		Name:        types.StringValue("rule"),
		Description: types.StringValue(""),
		Action:      action,
		Conditions: &[]*WaapDomainCustomRuleConditionsModel{{
			IP: &WaapDomainCustomRuleConditionsIPModel{
				IPAddress: types.StringValue("192.0.2.10"),
				Negation:  types.BoolValue(false),
			},
		}},
	}
}

func marshalUpdate(t *testing.T, plan, state WaapDomainCustomRuleModel) map[string]any {
	t.Helper()
	body, err := plan.MarshalJSONForUpdate(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("invalid json %s: %v", body, err)
	}
	return got
}

// The API replaces the action as a whole and rejects null for a kind that is
// not set, so switching block -> monitor must send only monitor, never
// "block": null.
func TestMarshalJSONForUpdateSendsOnlyNewActionKind(t *testing.T) {
	t.Parallel()

	state := customRuleModel(&WaapDomainCustomRuleActionModel{
		Block: &WaapDomainCustomRuleActionBlockModel{
			StatusCode:     types.Int64Value(403),
			ActionDuration: types.StringValue("5m"),
		},
	})
	plan := customRuleModel(&WaapDomainCustomRuleActionModel{
		Monitor: &map[string]jsontypes.Normalized{},
	})

	got := marshalUpdate(t, plan, state)
	action, ok := got["action"].(map[string]any)
	if !ok || len(action) != 1 || action["monitor"] == nil {
		t.Errorf("expected only a non-null monitor action, got %v", got)
	}
}

// Dropping action_duration inside block sends the block without a null key.
func TestMarshalJSONForUpdateOmitsRemovedBlockField(t *testing.T) {
	t.Parallel()

	state := customRuleModel(&WaapDomainCustomRuleActionModel{
		Block: &WaapDomainCustomRuleActionBlockModel{
			StatusCode:     types.Int64Value(403),
			ActionDuration: types.StringValue("5m"),
		},
	})
	plan := customRuleModel(&WaapDomainCustomRuleActionModel{
		Block: &WaapDomainCustomRuleActionBlockModel{
			StatusCode:     types.Int64Value(429),
			ActionDuration: types.StringNull(),
		},
	})

	got := marshalUpdate(t, plan, state)
	block, ok := got["action"].(map[string]any)["block"].(map[string]any)
	if !ok {
		t.Fatalf("expected a block action, got %v", got)
	}
	if _, present := block["action_duration"]; present {
		t.Errorf("expected action_duration to be omitted, got %v", got)
	}
	if block["status_code"] != float64(429) {
		t.Errorf("expected status_code 429, got %v", got)
	}
}
