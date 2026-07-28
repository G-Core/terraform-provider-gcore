// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_cluster_pool

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudK8SClusterPoolsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_name": schema.StringAttribute{
				Description: "Cluster name",
				Required:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Optional:    true,
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
				CustomType:  customfield.NewNestedObjectListType[CloudK8SClusterPoolsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "UUID of the cluster pool",
							Computed:    true,
						},
						"auto_healing_enabled": schema.BoolAttribute{
							Description: "Indicates the status of auto healing",
							Computed:    true,
						},
						"boot_volume_size": schema.Int64Attribute{
							Description: "Size of the boot volume",
							Computed:    true,
						},
						"boot_volume_type": schema.StringAttribute{
							Description: "Type of the boot volume",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Date of function creation",
							Computed:    true,
						},
						"crio_config": schema.MapAttribute{
							Description: "Crio configuration for pool nodes",
							Computed:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"flavor_id": schema.StringAttribute{
							Description: "ID of the cluster pool flavor",
							Computed:    true,
						},
						"is_public_ipv4": schema.BoolAttribute{
							Description: "Indicates if the pool is public",
							Computed:    true,
						},
						"kubelet_config": schema.MapAttribute{
							Description: "Kubelet configuration for pool nodes",
							Computed:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"labels": schema.MapAttribute{
							Description: "Labels applied to the cluster pool",
							Computed:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"max_node_count": schema.Int64Attribute{
							Description: "Maximum node count in the cluster pool",
							Computed:    true,
						},
						"min_node_count": schema.Int64Attribute{
							Description: "Minimum node count in the cluster pool",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the cluster pool",
							Computed:    true,
						},
						"node_count": schema.Int64Attribute{
							Description: "Node count in the cluster pool",
							Computed:    true,
						},
						"security_group_ids": schema.ListAttribute{
							Description: "Security group IDs applied to the cluster pool nodes",
							Computed:    true,
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
						},
						"status": schema.StringAttribute{
							Description: "Status of the cluster pool",
							Computed:    true,
						},
						"taints": schema.MapAttribute{
							Description: "Taints applied to the cluster pool",
							Computed:    true,
							CustomType:  customfield.NewMapType[types.String](ctx),
							ElementType: types.StringType,
						},
						"servergroup_id": schema.StringAttribute{
							Description: "Server group ID",
							Computed:    true,
						},
						"servergroup_name": schema.StringAttribute{
							Description: "Server group name",
							Computed:    true,
						},
						"servergroup_policy": schema.StringAttribute{
							Description: "Anti-affinity, affinity or soft-anti-affinity server group policy",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudK8SClusterPoolsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudK8SClusterPoolsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
