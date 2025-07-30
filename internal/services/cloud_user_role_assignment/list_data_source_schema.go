// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_user_role_assignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudUserRoleAssignmentsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "User ID for filtering",
				Optional:    true,
			},
			"limit": schema.Int64Attribute{
				Description: "Limit the number of returned items. Falls back to default of 1000 if not specified. Limited by max limit value of 1000",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtMost(1000),
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
				CustomType:  customfield.NewNestedObjectListType[CloudUserRoleAssignmentsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Assignment ID",
							Computed:    true,
						},
						"assigned_by": schema.Int64Attribute{
							Computed: true,
						},
						"client_id": schema.Int64Attribute{
							Description: "Client ID",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Created timestamp",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"project_id": schema.Int64Attribute{
							Description: "Project ID",
							Computed:    true,
						},
						"role": schema.StringAttribute{
							Description: "User role",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Updated timestamp",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"user_id": schema.Int64Attribute{
							Description: "User ID",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudUserRoleAssignmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudUserRoleAssignmentsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
