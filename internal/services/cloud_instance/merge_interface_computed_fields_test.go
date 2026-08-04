package cloud_instance

import (
	"testing"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestMergeInterfaceComputedFieldsNilElement pins the nil-element guard. A null
// element in the configured interfaces list decodes to a nil pointer and is
// skipped by the interfaces validator, so it reaches this helper; dereferencing
// it panics the provider. The surviving element must still be matched.
func TestMergeInterfaceComputedFieldsNilElement(t *testing.T) {
	t.Parallel()

	ifaces := &[]*CloudInstanceInterfacesModel{
		nil,
		{
			PortID:    types.StringUnknown(),
			IPAddress: types.StringUnknown(),
		},
	}
	apiIfaces := []cloud.NetworkInterfaceUnion{
		{
			PortID:        "port-abc",
			IPAssignments: []cloud.IPAssignment{{IPAddress: "10.0.0.7"}},
		},
	}

	mergeInterfaceComputedFields(ifaces, apiIfaces, false)

	if (*ifaces)[0] != nil {
		t.Errorf("expected nil element to stay nil, got %#v", (*ifaces)[0])
	}
	got := (*ifaces)[1]
	if got.PortID.ValueString() != "port-abc" {
		t.Errorf("expected PortID merged from the API interface, got %#v", got.PortID)
	}
	if got.IPAddress.ValueString() != "10.0.0.7" {
		t.Errorf("expected IPAddress merged from the API interface, got %#v", got.IPAddress)
	}
}
