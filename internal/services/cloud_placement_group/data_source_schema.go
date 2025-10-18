// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_placement_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudPlacementGroupDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Placement Groups allow you to specific a policy that determines whether Virtual Machines will be hosted on the same physical server or on different ones.",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.Int64Attribute{
				Optional: true,
			},
			"region_id": schema.Int64Attribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the server group.",
				Computed:    true,
			},
			"policy": schema.StringAttribute{
				Description: "The server group policy. Options are: anti-affinity, affinity, or soft-anti-affinity.",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"servergroup_id": schema.StringAttribute{
				Description: "The ID of the server group.",
				Computed:    true,
			},
			"instances": schema.ListNestedAttribute{
				Description: "The list of instances in this server group.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudPlacementGroupInstancesDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_id": schema.StringAttribute{
							Description: "The ID of the instance, corresponding to the attribute 'id'.",
							Computed:    true,
						},
						"instance_name": schema.StringAttribute{
							Description: "The name of the instance, corresponding to the attribute 'name'.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudPlacementGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudPlacementGroupDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
