// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudProjectDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Projects are organizational units that group cloud resources for access control and billing.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"client_id": schema.Int64Attribute{
				Description: "ID associated with the client.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime of creation, which is automatically generated.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"deleted_at": schema.StringAttribute{
				Description: "Datetime of deletion, which is automatically generated if the project is deleted.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"description": schema.StringAttribute{
				Description: "Description of the project.",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "Project ID, which is automatically generated upon creation.",
				Computed:    true,
			},
			"is_default": schema.BoolAttribute{
				Description: "Indicates if the project is the default one. Each client always has one default project.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Unique project name for a client.",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "The state of the project.",
				Computed:    true,
			},
		},
	}
}

func (d *CloudProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudProjectDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
