// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudNetworkDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Networks provide software-defined networking infrastructure for connecting instances and other cloud resources within a region.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Network ID",
				Computed:    true,
			},
			"network_id": schema.StringAttribute{
				Description: "Network ID",
				Optional:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Datetime when the network was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"creator_task_id": schema.StringAttribute{
				Description: "Task that created this entity",
				Computed:    true,
			},
			"default": schema.BoolAttribute{
				Description: "True if network has `is_default` attribute",
				Computed:    true,
			},
			"external": schema.BoolAttribute{
				Description: "True if the network `router:external` attribute",
				Computed:    true,
			},
			"mtu": schema.Int64Attribute{
				Description: "MTU (maximum transmission unit). Default value is 1450",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Network name",
				Computed:    true,
			},
			"port_security_enabled": schema.BoolAttribute{
				Description: "Indicates `port_security_enabled` status of all newly created in the network ports.",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"segmentation_id": schema.Int64Attribute{
				Description: "Id of network segment",
				Computed:    true,
			},
			"shared": schema.BoolAttribute{
				Description: "True when the network is shared with your project by external owner",
				Computed:    true,
			},
			"task_id": schema.StringAttribute{
				Description: "The UUID of the active task that currently holds a lock on the resource. This lock prevents concurrent modifications to ensure consistency. If `null`, the resource is not locked.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Network type (vlan, vxlan)",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the network was last updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"subnets": schema.ListAttribute{
				Description: "List of subnetworks",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"tags": schema.ListNestedAttribute{
				Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudNetworkTagsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "Tag key. Maximum 255 characters. Cannot contain spaces, tabs, newlines, empty string or '=' character.",
							Computed:    true,
						},
						"read_only": schema.BoolAttribute{
							Description: "If true, the tag is read-only and cannot be modified by the user",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "Tag value. Maximum 255 characters. Cannot contain spaces, tabs, newlines, empty string or '=' character.",
							Computed:    true,
						},
					},
				},
			},
			"find_one_by": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"external": schema.BoolAttribute{
						Description: "Filter by external network status",
						Optional:    true,
					},
					"name": schema.StringAttribute{
						Description: "Filter networks by name",
						Optional:    true,
					},
					"network_type": schema.StringAttribute{
						Description: "Filter by network type (vlan or vxlan)\nAvailable values: \"vlan\", \"vxlan\".",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("vlan", "vxlan"),
						},
					},
					"order_by": schema.StringAttribute{
						Description: "Ordering networks list result by `name`, `created_at` or `priority` fields and directions (e.g. `created_at.desc`). Default is `created_at.desc`. Use `priority.desc` to sort by shared network priority (relevant when `owned_by=any`).\nAvailable values: \"created_at.asc\", \"created_at.desc\", \"name.asc\", \"name.desc\", \"priority.desc\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"created_at.asc",
								"created_at.desc",
								"name.asc",
								"name.desc",
								"priority.desc",
							),
						},
					},
					"owned_by": schema.StringAttribute{
						Description: "Controls which networks are returned. 'project' (default) returns only networks owned by the project. 'any' returns all networks that the project can use, including shared networks from other projects.\nAvailable values: \"any\", \"project\".",
						Computed:    true,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("any", "project"),
						},
					},
					"tag_key": schema.ListAttribute{
						Description: "Optional. Filter by tag keys. ?`tag_key`=key1&`tag_key`=key2",
						Optional:    true,
						ElementType: types.StringType,
					},
					"tag_key_value": schema.StringAttribute{
						Description: "Optional. Filter by tag key-value pairs.",
						Optional:    true,
					},
				},
			},
		},
	}
}

func (d *CloudNetworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudNetworkDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("network_id"), path.MatchRoot("find_one_by")),
	}
}
