package cloud_secret

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPayloadUnpaired(t *testing.T) {
	t.Parallel()

	nullObj := types.ObjectNull(payloadAttrTypes)
	setObj := types.ObjectValueMust(payloadAttrTypes, map[string]attr.Value{
		"certificate_wo":       types.StringNull(),
		"certificate_chain_wo": types.StringNull(),
		"private_key_wo":       types.StringNull(),
	})
	unknownObj := types.ObjectUnknown(payloadAttrTypes)

	nullVer := types.Int64Null()
	setVer := types.Int64Value(1)
	unknownVer := types.Int64Unknown()

	cases := []struct {
		name    string
		payload types.Object
		version types.Int64
		want    bool
	}{
		{"both null", nullObj, nullVer, false},
		{"both set", setObj, setVer, false},
		{"payload set, version null", setObj, nullVer, true},
		{"payload null, version set", nullObj, setVer, true},
		{"payload unknown, version null", unknownObj, nullVer, false},
		{"payload unknown, version set", unknownObj, setVer, false},
		{"payload null, version unknown", nullObj, unknownVer, false},
		{"payload set, version unknown", setObj, unknownVer, false},
		{"both unknown", unknownObj, unknownVer, false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := payloadUnpaired(tc.payload, tc.version); got != tc.want {
				t.Errorf("payloadUnpaired(%s) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestPayloadMissing(t *testing.T) {
	t.Parallel()

	nullObj := types.ObjectNull(payloadAttrTypes)
	setObj := types.ObjectValueMust(payloadAttrTypes, map[string]attr.Value{
		"certificate_wo":       types.StringNull(),
		"certificate_chain_wo": types.StringNull(),
		"private_key_wo":       types.StringNull(),
	})
	unknownObj := types.ObjectUnknown(payloadAttrTypes)

	nullVer := types.Int64Null()
	setVer := types.Int64Value(1)
	unknownVer := types.Int64Unknown()

	cases := []struct {
		name    string
		payload types.Object
		version types.Int64
		want    bool
	}{
		{"both null", nullObj, nullVer, true},
		{"both set", setObj, setVer, false},
		{"payload set, version null", setObj, nullVer, true},
		{"payload null, version set", nullObj, setVer, true},
		{"payload unknown, version null", unknownObj, nullVer, true},
		{"payload unknown, version set", unknownObj, setVer, false},
		{"payload null, version unknown", nullObj, unknownVer, true},
		{"payload set, version unknown", setObj, unknownVer, false},
		{"both unknown", unknownObj, unknownVer, false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := payloadMissing(tc.payload, tc.version); got != tc.want {
				t.Errorf("payloadMissing(%s) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}
