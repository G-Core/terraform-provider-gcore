package waap_domain_firewall_rule

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

// syncFromBody fixes up the model after apijson.Unmarshal in Read and
// ImportState.
//
// description: the API leaves it out of the GET body when the rule has none,
// which reads back as null. The schema defaults it to "", so an imported rule
// without a description would plan "+ description = """. Null becomes "".
//
// action.allow / action.block: apijson.Unmarshal decodes an empty JSON object
// into a nil map, so an allow rule reads back as allow = null and an imported
// allow rule plans "+ allow = {}". It also keeps the previous state value of a
// key the body no longer has, so a rule switched between allow and block
// outside Terraform shows a phantom of the old action. Only these two keys
// are synced; fields nested inside block (e.g. action_duration) are decoded
// from the body as usual.
func syncFromBody(body []byte, data *WaapDomainFirewallRuleModel) {
	if data == nil {
		return
	}
	if data.Description.IsNull() {
		data.Description = types.StringValue("")
	}
	if data.Action == nil || !gjson.GetBytes(body, "action").IsObject() {
		return
	}
	if allow := gjson.GetBytes(body, "action.allow"); !allow.IsObject() {
		data.Action.Allow = nil
	} else if len(allow.Map()) == 0 {
		data.Action.Allow = &map[string]jsontypes.Normalized{}
	}
	if !gjson.GetBytes(body, "action.block").IsObject() {
		data.Action.Block = nil
	}
}
