package cloud_gpu_baremetal_cluster_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/G-Core/terraform-provider-gcore/internal/customvalidator"
	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_gpu_baremetal_cluster"
)

// interfacesAttribute pulls servers_settings.interfaces out of the resource schema.
func interfacesAttribute(t *testing.T) schema.ListNestedAttribute {
	t.Helper()

	resourceSchema := cloud_gpu_baremetal_cluster.ResourceSchema(context.TODO())

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

// TestCloudGPUBaremetalClusterInterfacesHasVariantValidator guards the wiring. The
// validator's own unit tests prove the rules; only this proves the schema uses it.
// It carries extra weight here: this resource has no positive acceptance coverage of
// the subnet-interface path, so the wiring is otherwise unverified.
func TestCloudGPUBaremetalClusterInterfacesHasVariantValidator(t *testing.T) {
	t.Parallel()

	want := reflect.TypeFor[customvalidator.GPUClusterInterfaceVariantValidator]()
	for _, v := range interfacesAttribute(t).Validators {
		if reflect.TypeOf(v) == want {
			return
		}
	}

	t.Errorf("servers_settings.interfaces is missing %s", want)
}

// TestCloudGPUBaremetalClusterIPFamilyPreservesNullState guards the second half of
// GCLOUD2-28207. The built-in UseStateForUnknown leaves the plan unknown when the
// prior state is null, which is exactly the case for "subnet" interfaces, where the
// API never returns ip_family.
func TestCloudGPUBaremetalClusterIPFamilyPreservesNullState(t *testing.T) {
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
