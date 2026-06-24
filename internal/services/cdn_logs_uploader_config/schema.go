// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CDNLogsUploaderConfigResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Logs uploader configs tie a logs uploader policy to one or more targets and a set of CDN resources, controlling which access logs are uploaded and where they are delivered.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseNonNullStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description: "Name of the config.",
				Required:    true,
			},
			"policy": schema.Int64Attribute{
				Description: "ID of the policy that should be assigned to given config.",
				Required:    true,
			},
			"target": schema.Int64Attribute{
				Description: "ID of the target to which logs should be uploaded.",
				Required:    true,
			},
			"resources": schema.ListAttribute{
				Description: "List of resource IDs to which the config should be applied.",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"enabled": schema.BoolAttribute{
				Description: "Enables or disables the config.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"for_all_resources": schema.BoolAttribute{
				Description: "If set to true, the config will be applied to all CDN resources.\nIf set to false, the config will be applied to the resources specified in the `resources` field.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"client_id": schema.Int64Attribute{
				Description: "Client that owns the config.",
				Computed:    true,
			},
			"created": schema.StringAttribute{
				Description: "Time when the config was created.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"updated": schema.StringAttribute{
				Description: "Time when the config was updated.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"status": schema.StringAttribute{
				Description: "Validation status of the logs uploader config.",
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
			},
		},
	}
}

func (r *CDNLogsUploaderConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CDNLogsUploaderConfigResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
