// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_config

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CDNLogsUploaderConfigDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Logs uploader configs tie a logs uploader policy to one or more targets and a set of CDN resources, controlling which access logs are uploaded and where they are delivered.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
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
			"enabled": schema.BoolAttribute{
				Description: "Enables or disables the config.",
				Computed:    true,
			},
			"for_all_resources": schema.BoolAttribute{
				Description: "If set to true, the config will be applied to all CDN resources.\nIf set to false, the config will be applied to the resources specified in the `resources` field.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the config.",
				Computed:    true,
			},
			"policy": schema.Int64Attribute{
				Description: "ID of the policy that should be assigned to given config.",
				Computed:    true,
			},
			"target": schema.Int64Attribute{
				Description: "ID of the target to which logs should be uploaded.",
				Computed:    true,
			},
			"updated": schema.StringAttribute{
				Description: "Time when the config was updated.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"resources": schema.ListAttribute{
				Description: "List of resource IDs to which the config should be applied.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.Int64](ctx),
				ElementType: types.Int64Type,
			},
			"status": schema.StringAttribute{
				Description: "Validation status of the logs uploader config.",
				Computed:    true,
				CustomType:  jsontypes.NormalizedType{},
			},
		},
	}
}

func (d *CDNLogsUploaderConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CDNLogsUploaderConfigDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
