// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*StorageObjectStorageResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "S3-compatible object storages provide scalable cloud storage with S3 API compatibility. Each storage is provisioned in a specific location and exposes one or more access keys for authentication.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "Unique identifier for the storage instance",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseNonNullStateForUnknown(), int64planmodifier.RequiresReplace()},
			},
			"location_name": schema.StringAttribute{
				Description:   "Location code where the storage should be created",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "User-defined name for the storage instance",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"address": schema.StringAttribute{
				Description: "Full hostname/address for accessing the storage endpoint",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "ISO 8601 timestamp when the storage was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"full_name": schema.StringAttribute{
				Description: "Read-only internal full name of the storage, composed as \"{`client_id`}-{name}\".\nUsed internally by the backend. Clients should continue to identify the storage by `name`.",
				Computed:    true,
			},
			"provisioning_status": schema.StringAttribute{
				Description: "Lifecycle status of the storage. Use this to check readiness before operations.\nAvailable values: \"creating\", \"active\", \"updating\", \"deleting\", \"deleted\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"creating",
						"active",
						"updating",
						"deleting",
						"deleted",
					),
				},
			},
			"access_keys": schema.ListNestedAttribute{
				Description: "S3 access keys",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[StorageObjectStorageAccessKeysModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_key": schema.StringAttribute{
							Description: "Access key ID used as the username in S3 authentication. Pass this in the `AWS_ACCESS_KEY_ID` field of your S3 client.",
							Computed:    true,
						},
						"secret_key": schema.StringAttribute{
							Description: "Secret key used as the password in S3 authentication. Save this now — it cannot be retrieved again.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *StorageObjectStorageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *StorageObjectStorageResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
