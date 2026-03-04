// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_project

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudProjectsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.Int64Attribute{
				Description: "Client ID filter for administrators.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name to filter the results by.",
				Optional:    true,
			},
			"include_deleted": schema.BoolAttribute{
				Description: "Whether to include deleted projects in the response.",
				Computed:    true,
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
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudProjectsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Project ID, which is automatically generated upon creation.",
							Computed:    true,
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
					},
				},
			},
		},
	}
}

func (d *CloudProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudProjectsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
