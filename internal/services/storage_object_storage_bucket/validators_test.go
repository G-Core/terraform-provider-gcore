package storage_object_storage_bucket

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TestValidateConfigReportsEveryError guards against one invalid setting
// hiding another: an invalid lifecycle and an empty allowed_origins in the
// same configuration must produce both errors.
func TestValidateConfigReportsEveryError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	objType := ResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
	lifecycleType := objType.AttributeTypes["storage_object_storage_bucket_lifecycle"].(tftypes.Object)
	raw := rawBucketWith(ctx, t, map[string]tftypes.Value{
		"storage_id": tftypes.NewValue(tftypes.Number, 1),
		"name":       tftypes.NewValue(tftypes.String, "tf-test-bucket"),
		"storage_object_storage_bucket_lifecycle": tftypes.NewValue(lifecycleType, map[string]tftypes.Value{
			"expiration_days": tftypes.NewValue(tftypes.Number, 0),
		}),
		"cors": rawCors(ctx, t, tftypes.NewValue(originsTerraformType(ctx, t), []tftypes.Value{})),
	})

	resp := &resource.ValidateConfigResponse{}
	(&StorageObjectStorageBucketResource{}).ValidateConfig(ctx, resource.ValidateConfigRequest{
		Config: tfsdk.Config{Raw: raw, Schema: ResourceSchema(ctx)},
	}, resp)

	got := map[string]bool{}
	for _, d := range resp.Diagnostics.Errors() {
		got[d.Summary()] = true
	}
	for _, want := range []string{"Invalid expiration_days", "Empty allowed_origins"} {
		if !got[want] {
			t.Errorf("missing %q error, got %v", want, resp.Diagnostics)
		}
	}
}
