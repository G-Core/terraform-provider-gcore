package cloud_load_balancer_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_load_balancer"
)

// forbiddenResourceAttributes must never appear in the load balancer *resource*
// schema. Some of them are still exposed by the data sources, which is
// intentional.
//
//   - listeners: listeners are managed by gcore_cloud_load_balancer_listener.
//     Keeping them inline on the resource made every plan report an in-place
//     update, because the API returns listener state the resource cannot own.
//   - task_id, tasks: asynchronous job bookkeeping returned by the API. They
//     change on every write and identify a transient task, not the resource.
//   - updated_at: server metadata that changes on every write, so no plan
//     modifier can honestly stabilise it. The data sources keep it.
//   - ddos_profile: the resource's Read never asks for the DDoS profile, and the
//     with_ddos query parameter that would return it is being retired by the
//     cloud API, so the attribute is permanently null.
//
// This is a regression guard, not a schema policy: each of these was deleted
// deliberately and has come back on its own. The deletions live in custom code
// re-applied on every regeneration, and a reseal has silently reverted the
// listeners deletion twice; ddos_profile was deleted in 2.0.0-alpha.3 and a
// later regeneration reinstated it.
var forbiddenResourceAttributes = []string{
	"ddos_profile",
	"listeners",
	"task_id",
	"tasks",
	"updated_at",
}

// TestCloudLoadBalancerResourceSchemaOmitsServerOnlyAttributes fails if an
// attribute that was deliberately deleted from the resource schema comes back.
func TestCloudLoadBalancerResourceSchemaOmitsServerOnlyAttributes(t *testing.T) {
	t.Parallel()

	attributes := cloud_load_balancer.ResourceSchema(context.TODO()).Attributes

	for _, name := range forbiddenResourceAttributes {
		if _, found := attributes[name]; found {
			t.Errorf("attribute %q is present in the gcore_cloud_load_balancer resource schema; "+
				"it was removed on purpose and must stay out of the resource", name)
		}
	}
}

// TestCloudLoadBalancerResourceTagsV2IsNotPinned fails if a plan modifier is
// attached to tags_v2.
//
// tags_v2 is derived server-side: it is the union of the user's tags and
// read-only system tags, and its element order is not guaranteed. Any modifier
// that reuses the prior state value therefore claims a planned value the API is
// free to contradict, and an in-place tags update fails with "Provider produced
// inconsistent result after apply".
//
// The assertion is deliberately fail-closed: it rejects any modifier, not just
// UseStateForUnknown, and it requires tags_v2 to stay a ListNestedAttribute. A
// future change that wants either — a modifier that is genuinely safe here, or
// a migration to a set to deal with the unstable ordering — has to come edit
// this test and say so. That is the point.
//
// The reason for being this strict: UseStateForUnknown was removed from tags_v2
// once before, with a comment explaining why. A regeneration dropped the
// comment, and a later change restored the modifier because nothing recorded
// that the removal had been deliberate. Deletions are the fragile shape in
// custom code here — they can vanish in a reseal without producing a conflict —
// so the deletion is asserted in a test rather than left to a comment.
func TestCloudLoadBalancerResourceTagsV2IsNotPinned(t *testing.T) {
	t.Parallel()

	attribute, found := cloud_load_balancer.ResourceSchema(context.TODO()).Attributes["tags_v2"]
	if !found {
		t.Fatal("tags_v2 is missing from the gcore_cloud_load_balancer resource schema")
	}

	tagsV2, ok := attribute.(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("tags_v2 is a %T, want schema.ListNestedAttribute", attribute)
	}

	for _, modifier := range tagsV2.PlanModifiers {
		t.Errorf("tags_v2 carries the plan modifier %T (%s); tags_v2 must stay unknown on "+
			"updates, because the API may return tags the plan cannot predict",
			modifier, modifier.Description(context.TODO()))
	}
}
