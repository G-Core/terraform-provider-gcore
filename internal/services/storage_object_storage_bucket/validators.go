package storage_object_storage_bucket

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.ResourceWithValidateConfig = (*StorageObjectStorageBucketResource)(nil)

	lifecyclePath      = path.Root("storage_object_storage_bucket_lifecycle")
	expirationDaysPath = lifecyclePath.AtName("expiration_days")
	allowedOriginsPath = path.Root("cors").AtName("allowed_origins")
)

// ValidateConfig rejects the API's "clear" values in configuration.
//
// The API removes a lifecycle rule for expiration_days 0, null or an omitted
// value, and removes CORS for an empty allowed_origins list. It then returns
// the block as null, so a configuration that writes one of these values never
// matches state and the apply fails with an inconsistent result. Removing the
// block is the supported way to clear either setting: the update sends the
// clear value on the user's behalf.
//
// A cors block without allowed_origins stays valid: the attribute is computed
// and keeps the bucket's current origins.
func (r *StorageObjectStorageBucketResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var lifecycle types.Object
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, lifecyclePath, &lifecycle)...)
	if !resp.Diagnostics.HasError() && !lifecycle.IsNull() && !lifecycle.IsUnknown() {
		var days types.Int64
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, expirationDaysPath, &days)...)
		switch {
		case resp.Diagnostics.HasError() || days.IsUnknown():
		case days.IsNull():
			resp.Diagnostics.AddAttributeError(expirationDaysPath,
				"Missing expiration_days",
				"expiration_days is required when the storage_object_storage_bucket_lifecycle block is set. To remove the lifecycle rule, delete the block.")
		case days.ValueInt64() < 1:
			resp.Diagnostics.AddAttributeError(expirationDaysPath,
				"Invalid expiration_days",
				"expiration_days must be at least 1. To remove the lifecycle rule, delete the storage_object_storage_bucket_lifecycle block.")
		}
	}

	var origins customfield.List[types.String]
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, allowedOriginsPath, &origins)...)
	if !resp.Diagnostics.HasError() && !origins.IsNull() && !origins.IsUnknown() && len(origins.Elements()) == 0 {
		resp.Diagnostics.AddAttributeError(allowedOriginsPath,
			"Empty allowed_origins",
			"allowed_origins must not be empty. To remove CORS, delete the cors block.")
	}
}
