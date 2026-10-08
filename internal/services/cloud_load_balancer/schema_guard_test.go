package cloud_load_balancer_test

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

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
//   - tags_v2: the read-only server view of tags. Nothing could ever be set
//     through it, and the refresh path reads the raw tags_v2 JSON key rather
//     than the attribute, so the resource loses nothing by dropping it. The
//     read-only view now lives on the data sources, where the same list is
//     exposed as `tags`. A plan modifier on it used to break in-place tag
//     updates with "Provider produced inconsistent result after apply"; the
//     attribute being gone is the stronger version of that guard.
//   - vip_network_id, vip_subnet_id: create-only inputs the API never returns,
//     so they could not be refreshed or checked after import. A private VIP is
//     placed with vip_port_id and a gcore_cloud_reserved_fixed_ip instead. They
//     are skipped for terraform in codegen config, not deleted here, so this
//     catches the config entry being dropped or turning into a no-op.
//
// This is a regression guard, not a schema policy: each of these was deleted
// deliberately and has come back on its own. The deletions live in custom code
// re-applied on every regeneration, and a reseal has silently reverted the
// listeners deletion twice; ddos_profile was deleted in 2.0.0-alpha.3 and a
// later regeneration reinstated it.
//
// Note what this list can and cannot catch: it is a check on names, so it only
// fires when a deleted attribute comes back under the name it had. The shapes
// that would slip past it are pinned positively by the tests below.
var forbiddenResourceAttributes = []string{
	"ddos_profile",
	"listeners",
	"tags_v2",
	"task_id",
	"tasks",
	"updated_at",
	"vip_network_id",
	"vip_subnet_id",
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

// TestCloudLoadBalancerResourceDescriptionsOmitRemovedAttributes fails if an
// attribute description still points users at a removed attribute. The
// descriptions come from the OpenAPI spec, which keeps vip_network_id for the
// Go and Python SDKs, so a regeneration can bring the stale text back.
func TestCloudLoadBalancerResourceDescriptionsOmitRemovedAttributes(t *testing.T) {
	t.Parallel()

	for attrName, attribute := range cloud_load_balancer.ResourceSchema(context.TODO()).Attributes {
		for _, removed := range forbiddenResourceAttributes {
			if strings.Contains(attribute.GetDescription(), "`"+removed+"`") {
				t.Errorf("description of %q mentions the removed attribute %q: %q",
					attrName, removed, attribute.GetDescription())
			}
		}
	}
}

// TestCloudLoadBalancerResourceTagsStayWritable fails if `tags` on the resource
// stops being the user-writable map.
//
// This is the half forbiddenResourceAttributes cannot cover. Banning the name
// tags_v2 catches the server-derived list returning under its old name and
// nothing else, and the way that list was dropped makes the other case the
// likely one: codegen config renames the response attribute to `tags` so it
// collides with the write-shape map and the generator keeps the request side.
// Reverting that entry brings the mirror back under its old name, which the
// list above catches. What it cannot catch is a spec-side rename of tags_v2
// upstream turning the entry into a silent no-op: then the mirror arrives as
// `tags` - a name the resource legitimately carries, so no name check can flag
// it.
//
// Pinning the shape closes that. The assertion this restores only said `tags`
// was not a ListNestedAttribute, which left every other way the attribute could
// go wrong passing quietly; saying what it must be instead states both halves of
// the intended schema rather than blacklisting one wrong version of it. Optional
// is asserted too, because "user-writable" is the actual claim - a read-only
// computed map would satisfy the type and still be the mirror.
func TestCloudLoadBalancerResourceTagsStayWritable(t *testing.T) {
	t.Parallel()

	attribute, found := cloud_load_balancer.ResourceSchema(context.TODO()).Attributes["tags"]
	if !found {
		t.Fatal("tags is missing from the gcore_cloud_load_balancer resource schema; " +
			"it is the map the caller writes and must stay")
	}

	tags, ok := attribute.(resourceschema.MapAttribute)
	if !ok {
		t.Fatalf("tags on the gcore_cloud_load_balancer resource is a %T, want "+
			"resourceschema.MapAttribute; a nested list here is the server-derived "+
			"tag mirror back under a new name", attribute)
	}

	if !tags.Optional {
		t.Error("tags on the gcore_cloud_load_balancer resource is not Optional; " +
			"the writable map has turned into a read-only one, which is the " +
			"server-derived mirror in the map's clothing")
	}
}

// TestCloudLoadBalancerDataSourceExposesTags fails if the read-only tag list
// stops being exposed as `tags` on the data source.
//
// The rename from tags_v2 lives in codegen config, not in this repo, so nothing
// here would otherwise notice it being reverted - and a spec-side rename of
// tags_v2 upstream would turn the config entry into a silent no-op rather than
// an error. This is the assertion that catches that.
func TestCloudLoadBalancerDataSourceExposesTags(t *testing.T) {
	t.Parallel()

	attribute, found := cloud_load_balancer.DataSourceSchema(context.TODO()).Attributes["tags"]
	if !found {
		t.Fatal("tags is missing from the gcore_cloud_load_balancer data source schema")
	}

	tags, ok := attribute.(datasourceschema.ListNestedAttribute)
	if !ok {
		t.Fatalf("tags is a %T, want datasourceschema.ListNestedAttribute", attribute)
	}

	for _, name := range []string{"key", "read_only", "value"} {
		if _, found := tags.NestedObject.Attributes[name]; !found {
			t.Errorf("tags is missing the %q attribute", name)
		}
	}
}

// TestCloudLoadBalancerVipPortIDStaysReplaceOnly pins the plan modifiers on
// vip_port_id, the only way left to place the VIP.
//
// RequiresReplaceIfConfigured: the API cannot move a load balancer to another
// port, so a changed vip_port_id must destroy and recreate it.
// UseStateForUnknown: vip_port_id is computed when omitted, so without it every
// plan of a load balancer created without vip_port_id would show the attribute
// as "known after apply".
func TestCloudLoadBalancerVipPortIDStaysReplaceOnly(t *testing.T) {
	t.Parallel()

	attribute, found := cloud_load_balancer.ResourceSchema(context.TODO()).Attributes["vip_port_id"]
	if !found {
		t.Fatal("vip_port_id is missing from the gcore_cloud_load_balancer resource schema")
	}

	port, ok := attribute.(resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("vip_port_id is a %T, want resourceschema.StringAttribute", attribute)
	}

	if !port.Optional || !port.Computed {
		t.Errorf("vip_port_id must be Optional and Computed, got Optional=%v Computed=%v", port.Optional, port.Computed)
	}

	ctx := context.TODO()
	var got []string
	for _, m := range port.PlanModifiers {
		got = append(got, fmt.Sprintf("%T (%s)", m, m.Description(ctx)))
	}
	for _, want := range []planmodifier.String{
		stringplanmodifier.RequiresReplaceIfConfigured(),
		stringplanmodifier.UseStateForUnknown(),
	} {
		// RequiresReplace and RequiresReplaceIfConfigured share one type, so
		// the description is what tells them apart.
		found := slices.ContainsFunc(port.PlanModifiers, func(m planmodifier.String) bool {
			return fmt.Sprintf("%T", m) == fmt.Sprintf("%T", want) && m.Description(ctx) == want.Description(ctx)
		})
		if !found {
			t.Errorf("vip_port_id is missing the %q plan modifier; plan modifiers are %v", want.Description(ctx), got)
		}
	}
}
