package storage_object_storage_bucket

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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
// value, and removes CORS for an empty allowed_origins set. It then returns
// the block as null, so a configuration that writes one of these values never
// matches state and the apply fails with an inconsistent result. Removing the
// block is the supported way to clear either setting: the update sends the
// clear value on the user's behalf.
//
// A cors block without allowed_origins stays valid: the attribute is computed
// and keeps the bucket's current origins.
func (r *StorageObjectStorageBucketResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	// Each check collects its own diagnostics, so an error in one setting
	// does not hide an error in the other.
	resp.Diagnostics.Append(validateLifecycle(ctx, req.Config)...)
	resp.Diagnostics.Append(validateAllowedOrigins(ctx, req.Config)...)
}

func validateLifecycle(ctx context.Context, config tfsdk.Config) (diags diag.Diagnostics) {
	var lifecycle types.Object
	diags.Append(config.GetAttribute(ctx, lifecyclePath, &lifecycle)...)
	if diags.HasError() || lifecycle.IsNull() || lifecycle.IsUnknown() {
		return diags
	}

	var days types.Int64
	diags.Append(config.GetAttribute(ctx, expirationDaysPath, &days)...)
	switch {
	case diags.HasError() || days.IsUnknown():
	case days.IsNull():
		diags.AddAttributeError(expirationDaysPath,
			"Missing expiration_days",
			"expiration_days is required when the storage_object_storage_bucket_lifecycle block is set. To remove the lifecycle rule, delete the block.")
	case days.ValueInt64() < 1:
		diags.AddAttributeError(expirationDaysPath,
			"Invalid expiration_days",
			"expiration_days must be at least 1. To remove the lifecycle rule, delete the storage_object_storage_bucket_lifecycle block.")
	}
	return diags
}

func validateAllowedOrigins(ctx context.Context, config tfsdk.Config) (diags diag.Diagnostics) {
	var origins customfield.Set[types.String]
	diags.Append(config.GetAttribute(ctx, allowedOriginsPath, &origins)...)
	if !diags.HasError() && !origins.IsNull() && !origins.IsUnknown() && len(origins.Elements()) == 0 {
		diags.AddAttributeError(allowedOriginsPath,
			"Empty allowed_origins",
			"allowed_origins must not be empty. To remove CORS, delete the cors block.")
	}
	return diags
}
