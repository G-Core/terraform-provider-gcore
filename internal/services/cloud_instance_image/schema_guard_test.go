package cloud_instance_image_test

import (
	"context"
	"testing"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_instance_image"
)

// TestCloudInstanceImageResourceOmitsTagsV2 fails if the read-only tags_v2
// mirror comes back to the resource schema.
//
// It was dropped in codegen config, by renaming the response attribute to
// `tags` so it collides with the write-shape map and the generator keeps the
// request side. Nothing in this repo would otherwise notice that being
// reverted, and deleted attributes on this provider have come back on their own
// after a reseal before - which is why the load balancer carries the same kind
// of guard.
func TestCloudInstanceImageResourceOmitsTagsV2(t *testing.T) {
	t.Parallel()

	if _, found := cloud_instance_image.ResourceSchema(context.TODO()).Attributes["tags_v2"]; found {
		t.Error("tags_v2 is present in the gcore_cloud_instance_image resource schema; " +
			"the read-only tag list belongs on the data sources, as `tags`")
	}
}

// TestCloudInstanceImageResourceTagsStayWritable fails if `tags` on the resource
// stops being the user-writable map.
//
// The test above is a check on a name, so it only fires if the mirror returns
// under the name it had. Given how it was dropped - a rename onto `tags` that
// collides with the write-shape map - the likelier regression is the other one:
// the config entry reverted, or made a silent no-op by a spec-side rename of
// tags_v2 upstream, and the server-derived list arrives as `tags`. That is a
// name the resource legitimately carries, so no name check can flag it.
//
// Hence the positive assertion: `tags` must be the Optional map the caller
// writes. Between the two, both halves of the intended schema are pinned instead
// of one wrong version of it being blacklisted. Optional is part of the claim -
// a read-only computed map would satisfy the type and still be the mirror.
func TestCloudInstanceImageResourceTagsStayWritable(t *testing.T) {
	t.Parallel()

	attribute, found := cloud_instance_image.ResourceSchema(context.TODO()).Attributes["tags"]
	if !found {
		t.Fatal("tags is missing from the gcore_cloud_instance_image resource schema; " +
			"it is the map the caller writes and must stay")
	}

	tags, ok := attribute.(resourceschema.MapAttribute)
	if !ok {
		t.Fatalf("tags on the gcore_cloud_instance_image resource is a %T, want "+
			"resourceschema.MapAttribute; a nested list here is the server-derived "+
			"tag mirror back under a new name", attribute)
	}

	if !tags.Optional {
		t.Error("tags on the gcore_cloud_instance_image resource is not Optional; " +
			"the writable map has turned into a read-only one, which is the " +
			"server-derived mirror in the map's clothing")
	}
}

// TestCloudInstanceImageDataSourceExposesTags fails if the read-only tag list
// stops being exposed as `tags` on the data source.
func TestCloudInstanceImageDataSourceExposesTags(t *testing.T) {
	t.Parallel()

	attribute, found := cloud_instance_image.DataSourceSchema(context.TODO()).Attributes["tags"]
	if !found {
		t.Fatal("tags is missing from the gcore_cloud_instance_image data source schema")
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
