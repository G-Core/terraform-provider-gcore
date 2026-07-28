package cloud_instance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestResolveUnknownInterfaceComputedFields verifies that unknown computed interface
// fields are resolved to null before state is persisted, while known and null values
// are left untouched, and that nil/empty inputs are handled safely.
func TestResolveUnknownInterfaceComputedFields(t *testing.T) {
	t.Parallel()

	t.Run("unknown resolves to null", func(t *testing.T) {
		t.Parallel()
		ifaces := &[]*CloudInstanceInterfacesModel{
			{
				PortID:    types.StringUnknown(),
				IPAddress: types.StringUnknown(),
			},
		}
		resolveUnknownInterfaceComputedFields(ifaces)

		got := (*ifaces)[0]
		if !got.PortID.IsNull() {
			t.Errorf("expected PortID to be null, got %#v", got.PortID)
		}
		if !got.IPAddress.IsNull() {
			t.Errorf("expected IPAddress to be null, got %#v", got.IPAddress)
		}
	})

	t.Run("known values preserved", func(t *testing.T) {
		t.Parallel()
		ifaces := &[]*CloudInstanceInterfacesModel{
			{
				PortID:    types.StringValue("port-123"),
				IPAddress: types.StringValue("10.0.0.5"),
			},
		}
		resolveUnknownInterfaceComputedFields(ifaces)

		got := (*ifaces)[0]
		if got.PortID.ValueString() != "port-123" {
			t.Errorf("expected PortID preserved, got %#v", got.PortID)
		}
		if got.IPAddress.ValueString() != "10.0.0.5" {
			t.Errorf("expected IPAddress preserved, got %#v", got.IPAddress)
		}
	})

	t.Run("null values preserved", func(t *testing.T) {
		t.Parallel()
		ifaces := &[]*CloudInstanceInterfacesModel{
			{
				PortID:    types.StringNull(),
				IPAddress: types.StringNull(),
			},
		}
		resolveUnknownInterfaceComputedFields(ifaces)

		got := (*ifaces)[0]
		if !got.PortID.IsNull() {
			t.Errorf("expected PortID to remain null, got %#v", got.PortID)
		}
		if !got.IPAddress.IsNull() {
			t.Errorf("expected IPAddress to remain null, got %#v", got.IPAddress)
		}
	})

	t.Run("mixed elements", func(t *testing.T) {
		t.Parallel()
		ifaces := &[]*CloudInstanceInterfacesModel{
			{
				PortID:    types.StringValue("port-known"),
				IPAddress: types.StringUnknown(),
			},
			{
				PortID:    types.StringUnknown(),
				IPAddress: types.StringValue("10.0.0.9"),
			},
		}
		resolveUnknownInterfaceComputedFields(ifaces)

		first := (*ifaces)[0]
		if first.PortID.ValueString() != "port-known" {
			t.Errorf("expected first PortID preserved, got %#v", first.PortID)
		}
		if !first.IPAddress.IsNull() {
			t.Errorf("expected first IPAddress resolved to null, got %#v", first.IPAddress)
		}
		second := (*ifaces)[1]
		if !second.PortID.IsNull() {
			t.Errorf("expected second PortID resolved to null, got %#v", second.PortID)
		}
		if second.IPAddress.ValueString() != "10.0.0.9" {
			t.Errorf("expected second IPAddress preserved, got %#v", second.IPAddress)
		}
	})

	t.Run("nil pointer is safe", func(t *testing.T) {
		t.Parallel()
		resolveUnknownInterfaceComputedFields(nil)
	})

	t.Run("empty slice is safe", func(t *testing.T) {
		t.Parallel()
		empty := &[]*CloudInstanceInterfacesModel{}
		resolveUnknownInterfaceComputedFields(empty)
	})

	t.Run("nil element is safe", func(t *testing.T) {
		t.Parallel()
		ifaces := &[]*CloudInstanceInterfacesModel{
			nil,
			{
				PortID:    types.StringUnknown(),
				IPAddress: types.StringUnknown(),
			},
		}
		resolveUnknownInterfaceComputedFields(ifaces)

		got := (*ifaces)[1]
		if !got.PortID.IsNull() {
			t.Errorf("expected PortID resolved to null, got %#v", got.PortID)
		}
		if !got.IPAddress.IsNull() {
			t.Errorf("expected IPAddress resolved to null, got %#v", got.IPAddress)
		}
	})
}
