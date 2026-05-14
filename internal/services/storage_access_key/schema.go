// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_access_key

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*StorageAccessKeyResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Object storage access keys provide secure credentials for API access to object storage resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "Access key ID used as the username in S3 authentication. Pass this in the `AWS_ACCESS_KEY_ID` field of your S3 client.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseNonNullStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"access_key": schema.StringAttribute{
				Description:   "Access key ID used as the username in S3 authentication. Pass this in the `AWS_ACCESS_KEY_ID` field of your S3 client.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseNonNullStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"storage_id": schema.Int64Attribute{
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"created_at": schema.StringAttribute{
				Description: "ISO 8601 timestamp when the access key was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"secret_key": schema.StringAttribute{
				Description: "Secret key used as the password in S3 authentication. Save this now — it cannot be retrieved again.",
				Computed:    true,
			},
		},
	}
}

func (r *StorageAccessKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *StorageAccessKeyResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
