// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudProjectDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Projects are organizational units that group cloud resources for access control and billing.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Project ID",
				Computed:    true,
			},
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
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"include_deleted": schema.BoolAttribute{
						Description: "Whether to include deleted projects in the response.",
						Computed:    true,
						Optional:    true,
					},
					"name": schema.StringAttribute{
						Description: "Name to filter the results by.",
						Optional:    true,
					},
					"order_by": schema.StringAttribute{
						Description: "Order by field and direction.\nAvailable values: \"created_at.asc\", \"created_at.desc\", \"name.asc\", \"name.desc\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"created_at.asc",
								"created_at.desc",
								"name.asc",
								"name.desc",
							),
						},
					},
				},
			},
		},
	}
}

func (d *CloudProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudProjectDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("project_id"), path.MatchRoot("find_one_by")),
	}
}
