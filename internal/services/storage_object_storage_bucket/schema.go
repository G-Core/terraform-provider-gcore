// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage_bucket

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*StorageObjectStorageBucketResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Buckets are containers within object storage that hold files (objects) and define their CORS, lifecycle, and access policy configuration.",
		Attributes: map[string]schema.Attribute{
			"storage_id": schema.Int64Attribute{
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "Name of the bucket to create",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"cors": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					// Computed+Optional with no Default, so the framework plans this
					// attribute as unknown whenever the cors block is present but
					// allowed_origins is omitted. The model field must therefore be
					// able to carry an unknown value; a plain *[]types.String cannot,
					// and Plan.Get fails before the request is ever sent. customfield.List
					// is what every other computed_optional list in this provider uses.
					// UseStateForUnknown keeps the post-create plan empty by reusing the
					// prior value; on create the state is null, so the attribute stays
					// unknown, is omitted from the request body, and is resolved to null
					// when the response is decoded.
					"allowed_origins": schema.ListAttribute{
						Description: "Web domains allowed to make direct browser requests. Send an empty array to remove CORS configuration.",
						Optional:    true,
						Computed:    true,
						CustomType:  customfield.NewListType[types.String](ctx),
						ElementType: types.StringType,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"policy": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"is_public": schema.BoolAttribute{
						Description: "Set to true to allow unauthenticated object downloads, false to require valid S3 credentials.",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"storage_object_storage_bucket_lifecycle": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"expiration_days": schema.Int64Attribute{
						Description: "Days before objects are automatically deleted. Set to a positive number to enable, or null/0 to remove the rule.",
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *StorageObjectStorageBucketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *StorageObjectStorageBucketResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
