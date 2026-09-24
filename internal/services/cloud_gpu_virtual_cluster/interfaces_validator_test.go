package cloud_gpu_virtual_cluster_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/G-Core/terraform-provider-gcore/internal/customvalidator"
	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_gpu_virtual_cluster"
)

// interfacesAttribute pulls servers_settings.interfaces out of the resource schema.
func interfacesAttribute(t *testing.T) schema.ListNestedAttribute {
	t.Helper()

	resourceSchema := cloud_gpu_virtual_cluster.ResourceSchema(context.TODO())

	serversSettings, ok := resourceSchema.Attributes["servers_settings"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("servers_settings is %T, want schema.SingleNestedAttribute", resourceSchema.Attributes["servers_settings"])
	}

	interfaces, ok := serversSettings.Attributes["interfaces"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("servers_settings.interfaces is %T, want schema.ListNestedAttribute", serversSettings.Attributes["interfaces"])
	}

	return interfaces
}

// TestCloudGPUVirtualClusterInterfacesHasVariantValidator guards the wiring. The
// validator's own unit tests prove the rules; only this proves the schema uses it.
func TestCloudGPUVirtualClusterInterfacesHasVariantValidator(t *testing.T) {
	t.Parallel()

	want := reflect.TypeFor[customvalidator.GPUClusterInterfaceVariantValidator]()
	for _, v := range interfacesAttribute(t).Validators {
		if reflect.TypeOf(v) == want {
			return
		}
	}

	t.Errorf("servers_settings.interfaces is missing %s", want)
}

// TestCloudGPUVirtualClusterIPFamilyPreservesNullState guards the second half of
// GCLOUD2-28207. The built-in UseStateForUnknown leaves the plan unknown when the
// prior state is null, which is exactly the case for "subnet" interfaces, where the
// API never returns ip_family.
func TestCloudGPUVirtualClusterIPFamilyPreservesNullState(t *testing.T) {
	t.Parallel()

	ipFamilyAttr := interfacesAttribute(t).NestedObject.Attributes["ip_family"]
	ipFamily, ok := ipFamilyAttr.(schema.StringAttribute)
	if !ok {
		t.Fatalf("ip_family is %T, want schema.StringAttribute", ipFamilyAttr)
	}

	want := reflect.TypeOf(planmodifiers.StringUseStateForUnknownInclNull())
	for _, pm := range ipFamily.PlanModifiers {
		if reflect.TypeOf(pm) == want {
			return
		}
	}

	t.Errorf("servers_settings.interfaces[].ip_family is missing %s", want)
}

// TestCloudGPUVirtualClusterInterfaceNameIsReadOnly pins the spec annotation that
// makes interfaces[].name computed-only. The API never honours a requested name
// (external interfaces get the network's name, subnet interfaces store it without
// applying it), so an Optional name can never match state and servers_settings
// forces replacement on every plan.
func TestCloudGPUVirtualClusterInterfaceNameIsReadOnly(t *testing.T) {
	t.Parallel()

	nameAttr := interfacesAttribute(t).NestedObject.Attributes["name"]
	name, ok := nameAttr.(schema.StringAttribute)
	if !ok {
		t.Fatalf("name is %T, want schema.StringAttribute", nameAttr)
	}

	if !name.Computed || name.Optional || name.Required {
		t.Errorf("servers_settings.interfaces[].name must be computed-only, got computed=%t optional=%t required=%t",
			name.Computed, name.Optional, name.Required)
	}

	// Read-only is not enough on its own: without UseStateForUnknown the name goes
	// unknown on every plan, and the servers_settings replace modifier turns any
	// unknown inside the object into a replacement.
	want := reflect.TypeOf(stringplanmodifier.UseStateForUnknown())
	for _, pm := range name.PlanModifiers {
		if reflect.TypeOf(pm) == want {
			return
		}
	}

	t.Errorf("servers_settings.interfaces[].name is missing %s", want)
}
