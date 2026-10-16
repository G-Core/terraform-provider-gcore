package waap_domain_advanced_rule

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"

	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
)

// unmarshalRead decodes a GET response on refresh and import, then fixes the
// shapes the generic decoder gets wrong, so an imported rule plans clean:
//   - allow/captcha/handshake/monitor come back as `{}`; the decoder skips an
//     empty JSON object for a map, which would leave the configured `{}` null.
//   - an empty map the response no longer carries is kept from prior state, so a
//     kind changed out of band would leave two kinds set. The response is
//     authoritative: every kind it does not return is cleared.
//   - a rule created without a description omits the field; the schema default
//     is "" (the API only clears a description with ""), so null would plan an
//     update to "" after every import.
func unmarshalRead(raw []byte, data **WaapDomainAdvancedRuleModel) error {
	if err := apijson.Unmarshal(raw, data); err != nil {
		return err
	}
	m := *data
	if m == nil {
		return nil
	}

	if m.Description.IsNull() {
		m.Description = types.StringValue("")
	}

	if m.Action == nil {
		return nil
	}
	action := gjson.GetBytes(raw, "action")
	for key, field := range map[string]**map[string]jsontypes.Normalized{
		"allow":     &m.Action.Allow,
		"captcha":   &m.Action.Captcha,
		"handshake": &m.Action.Handshake,
		"monitor":   &m.Action.Monitor,
	} {
		switch kind := action.Get(key); {
		case !kind.IsObject():
			*field = nil
		case len(kind.Map()) == 0:
			// the decoder skips an empty object: it keeps a prior map as is, and
			// may allocate the pointer around a nil map, which the framework also
			// stores as null
			*field = &map[string]jsontypes.Normalized{}
		}
	}
	if !action.Get("block").IsObject() {
		m.Action.Block = nil
	}
	if !action.Get("tag").IsObject() {
		m.Action.Tag = nil
	}
	return nil
}
