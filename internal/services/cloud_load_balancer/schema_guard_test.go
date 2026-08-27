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

// TestCloudLoadBalancerResourceHasNoTagsMirror fails if a server-derived tag
// mirror reappears on the load balancer resource.
//
// The resource used to carry tags_v2, a read-only mirror of the union of the
// user's tags and read-only system tags. It was removed deliberately: the
// mirror is derived server-side, its element order is not guaranteed, and
// nothing on the resource could pin it without an in-place tags update failing
// with "Provider produced inconsistent result after apply". Tags are read
// through the data source instead.
//
// The assertion is deliberately fail-closed: any nested-list tag mirror on the
// resource, under either name, fails. A future change that wants one back has
// to come edit this test and say so. That is the point.
//
// The reason for being this strict is unchanged from the modifier guard this
// replaces: deletions are the fragile shape in custom code here - they can
// vanish in a reseal without producing a conflict - so the deletion is asserted
// in a test rather than left to a comment.
func TestCloudLoadBalancerResourceHasNoTagsMirror(t *testing.T) {
	t.Parallel()

	attributes := cloud_load_balancer.ResourceSchema(context.TODO()).Attributes

	if attribute, found := attributes["tags_v2"]; found {
		t.Errorf("the gcore_cloud_load_balancer resource carries tags_v2 (%T); the "+
			"server-derived tag mirror was removed on purpose and tags are read "+
			"through the data source", attribute)
	}

	// tags itself is legitimate and must stay: it is the user-writable map the
	// caller sets. What must not come back is a nested-list mirror under that
	// name, which is the server-derived shape in a new disguise.
	if attribute, found := attributes["tags"]; found {
		if _, isMirror := attribute.(schema.ListNestedAttribute); isMirror {
			t.Errorf("the gcore_cloud_load_balancer resource carries tags as a "+
				"ListNestedAttribute (%T); that is the server-derived mirror under a "+
				"new name. tags must stay the user-writable map", attribute)
		}
	}
}
