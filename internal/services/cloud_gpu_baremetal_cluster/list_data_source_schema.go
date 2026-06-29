// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudGPUBaremetalClustersDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"region_id": schema.Int64Attribute{
				Description: "Region ID",
				Optional:    true,
			},
			"ids": schema.ListAttribute{
				Description: "Return only clusters with these IDs, e.g. `ids=<id1>&ids=<id2>`.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"image_ids": schema.ListAttribute{
				Description: "Return only clusters built from these source image IDs, e.g. `image_ids=<id1>&image_ids=<id2>`.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"tags": schema.MapAttribute{
				Description: "Filter by exact tag key-value pairs, e.g. `tags[env]=prod&tags[team]=core`. Pairs are ANDed; values match case-insensitively.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"created_at": schema.SingleNestedAttribute{
				Description: "Filter by creation time (UTC), e.g. `created_at[gte]=2026-01-01T00:00:00Z`.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"gt": schema.StringAttribute{
						Description: "Strictly after this timestamp, e.g. `[gt]=2026-01-01T00:00:00Z`.",
						Optional:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"gte": schema.StringAttribute{
						Description: "At or after this timestamp (inclusive), e.g. `[gte]=2026-01-01T00:00:00Z`.",
						Optional:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"lt": schema.StringAttribute{
						Description: "Strictly before this timestamp, e.g. `[lt]=2026-02-01T00:00:00Z`.",
						Optional:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"lte": schema.StringAttribute{
						Description: "At or before this timestamp (inclusive), e.g. `[lte]=2026-02-01T00:00:00Z`.",
						Optional:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
				},
			},
			"flavor": schema.SingleNestedAttribute{
				Description: "Filter by flavor (case-insensitive), e.g. `flavor[prefix]=bm3-`, `flavor[exact]=bm3-ai-1xlarge-h100-80-8`.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"contains": schema.ListAttribute{
						Description: "Case-insensitive substring, e.g. `[contains]=web`. Repeat the key to match any substring.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"exact": schema.ListAttribute{
						Description: "Case-insensitive exact match, e.g. `[exact]=web-1`. Repeat the key to match any of several.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"prefix": schema.ListAttribute{
						Description: "Case-insensitive starts-with, e.g. `[prefix]=prod-`. Repeat the key to match any prefix.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"suffix": schema.ListAttribute{
						Description: "Case-insensitive ends-with, e.g. `[suffix]=-db`. Repeat the key to match any suffix.",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"name": schema.SingleNestedAttribute{
				Description: "Filter by name (case-insensitive), e.g. `name[contains]=gpu`, `name[prefix]=prod-`.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"contains": schema.ListAttribute{
						Description: "Case-insensitive substring, e.g. `[contains]=web`. Repeat the key to match any substring.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"exact": schema.ListAttribute{
						Description: "Case-insensitive exact match, e.g. `[exact]=web-1`. Repeat the key to match any of several.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"prefix": schema.ListAttribute{
						Description: "Case-insensitive starts-with, e.g. `[prefix]=prod-`. Repeat the key to match any prefix.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"suffix": schema.ListAttribute{
						Description: "Case-insensitive ends-with, e.g. `[suffix]=-db`. Repeat the key to match any suffix.",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"servers_count": schema.SingleNestedAttribute{
				Description: "Filter by node count, e.g. `servers_count[gte]=2`, `servers_count[gte]=2&servers_count[lt]=8`.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"gt": schema.Int64Attribute{
						Description: "Strictly greater than, e.g. `[gt]=1`.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
					"gte": schema.Int64Attribute{
						Description: "Greater than or equal, e.g. `[gte]=2`.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
					"lt": schema.Int64Attribute{
						Description: "Strictly less than, e.g. `[lt]=8`.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
					"lte": schema.Int64Attribute{
						Description: "Less than or equal, e.g. `[lte]=4`.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
				},
			},
			"tag_key": schema.SingleNestedAttribute{
				Description: "Filter by tag key regardless of value, e.g. `tag_key[contains]=team`.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"contains": schema.ListAttribute{
						Description: "Case-insensitive substring, e.g. `[contains]=web`. Repeat the key to match any substring.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"exact": schema.ListAttribute{
						Description: "Case-insensitive exact match, e.g. `[exact]=web-1`. Repeat the key to match any of several.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"prefix": schema.ListAttribute{
						Description: "Case-insensitive starts-with, e.g. `[prefix]=prod-`. Repeat the key to match any prefix.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"suffix": schema.ListAttribute{
						Description: "Case-insensitive ends-with, e.g. `[suffix]=-db`. Repeat the key to match any suffix.",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"tag_value": schema.SingleNestedAttribute{
				Description: "Filter by tag value regardless of key, e.g. `tag_value[prefix]=prod`.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"contains": schema.ListAttribute{
						Description: "Case-insensitive substring, e.g. `[contains]=web`. Repeat the key to match any substring.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"exact": schema.ListAttribute{
						Description: "Case-insensitive exact match, e.g. `[exact]=web-1`. Repeat the key to match any of several.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"prefix": schema.ListAttribute{
						Description: "Case-insensitive starts-with, e.g. `[prefix]=prod-`. Repeat the key to match any prefix.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"suffix": schema.ListAttribute{
						Description: "Case-insensitive ends-with, e.g. `[suffix]=-db`. Repeat the key to match any suffix.",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"updated_at": schema.SingleNestedAttribute{
				Description: "Filter by last-change time (UTC), e.g. `updated_at[gte]=2026-06-01T00:00:00Z`.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"gt": schema.StringAttribute{
						Description: "Strictly after this timestamp, e.g. `[gt]=2026-01-01T00:00:00Z`.",
						Optional:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"gte": schema.StringAttribute{
						Description: "At or after this timestamp (inclusive), e.g. `[gte]=2026-01-01T00:00:00Z`.",
						Optional:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"lt": schema.StringAttribute{
						Description: "Strictly before this timestamp, e.g. `[lt]=2026-02-01T00:00:00Z`.",
						Optional:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
					"lte": schema.StringAttribute{
						Description: "At or before this timestamp (inclusive), e.g. `[lte]=2026-02-01T00:00:00Z`.",
						Optional:    true,
						CustomType:  timetypes.RFC3339Type{},
					},
				},
			},
			"managed_by": schema.ListAttribute{
				Description: "Filter by who manages the cluster: `user` (default) or `k8s` (Managed Kubernetes). Pass both to include all, e.g. `managed_by=user&managed_by=k8s`.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOfCaseInsensitive("k8s", "user"),
					),
				},
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
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
				CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClustersItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Cluster unique identifier",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Cluster creation date time",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"flavor": schema.StringAttribute{
							Description: "Cluster flavor name",
							Computed:    true,
						},
						"has_pending_changes": schema.BoolAttribute{
							Description: "True if any server in the cluster has pending (not yet applied) settings changes",
							Computed:    true,
						},
						"image_id": schema.StringAttribute{
							Description: "Image ID",
							Computed:    true,
						},
						"managed_by": schema.StringAttribute{
							Description: "User type managing the resource\nAvailable values: \"k8s\", \"user\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("k8s", "user"),
							},
						},
						"name": schema.StringAttribute{
							Description: "Cluster name",
							Computed:    true,
						},
						"servers_count": schema.Int64Attribute{
							Description: "Cluster servers count",
							Computed:    true,
						},
						"servers_ids": schema.ListAttribute{
							Description: "List of cluster nodes",
							Computed:    true,
							CustomType:  customfield.NewListType[types.String](ctx),
							ElementType: types.StringType,
						},
						"servers_settings": schema.SingleNestedAttribute{
							Computed:   true,
							CustomType: customfield.NewNestedObjectType[CloudGPUBaremetalClustersServersSettingsDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"file_shares": schema.ListNestedAttribute{
									Description: "List of file shares mounted across the cluster.",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClustersServersSettingsFileSharesDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description: "Unique identifier of the file share in UUID format.",
												Computed:    true,
											},
											"mount_path": schema.StringAttribute{
												Description: "Absolute mount path inside the system where the file share will be mounted.",
												Computed:    true,
											},
										},
									},
								},
								"interfaces": schema.ListNestedAttribute{
									Computed:   true,
									CustomType: customfield.NewNestedObjectListType[CloudGPUBaremetalClustersServersSettingsInterfacesDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"ip_family": schema.StringAttribute{
												Description: "Which subnets should be selected: IPv4, IPv6, or use dual stack.\nAvailable values: \"dual\", \"ipv4\", \"ipv6\".",
												Computed:    true,
												Validators: []validator.String{
													stringvalidator.OneOfCaseInsensitive(
														"dual",
														"ipv4",
														"ipv6",
													),
												},
											},
											"name": schema.StringAttribute{
												Description: "Interface name",
												Computed:    true,
											},
											"security_groups": schema.ListNestedAttribute{
												Description: "Resolved security groups applied to this interface.",
												Computed:    true,
												CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClustersServersSettingsInterfacesSecurityGroupsDataSourceModel](ctx),
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															Description: "Security group ID",
															Computed:    true,
														},
														"name": schema.StringAttribute{
															Description: "Security group name",
															Computed:    true,
														},
													},
												},
											},
											"type": schema.StringAttribute{
												Description: `Available values: "external", "subnet", "any_subnet".`,
												Computed:    true,
												Validators: []validator.String{
													stringvalidator.OneOfCaseInsensitive(
														"external",
														"subnet",
														"any_subnet",
													),
												},
											},
											"floating_ip": schema.SingleNestedAttribute{
												Description: "Floating IP config for this subnet attachment",
												Computed:    true,
												CustomType:  customfield.NewNestedObjectType[CloudGPUBaremetalClustersServersSettingsInterfacesFloatingIPDataSourceModel](ctx),
												Attributes: map[string]schema.Attribute{
													"source": schema.StringAttribute{
														Description: `Available values: "new".`,
														Computed:    true,
														Validators: []validator.String{
															stringvalidator.OneOfCaseInsensitive("new"),
														},
													},
												},
											},
											"network_id": schema.StringAttribute{
												Description: "Network ID the subnet belongs to. Port will be plugged in this network",
												Computed:    true,
											},
											"subnet_id": schema.StringAttribute{
												Description: "Port is assigned an IP address from this subnet",
												Computed:    true,
											},
											"ip_address": schema.StringAttribute{
												Description: "Fixed IP address",
												Computed:    true,
											},
										},
									},
								},
								"security_groups": schema.ListNestedAttribute{
									Description:        "Deprecated. Deduplicated union of security groups across all interfaces; the actual assignment may differ per interface. Use `interfaces[].security_groups` for the authoritative per-interface list.",
									Computed:           true,
									DeprecationMessage: "This attribute is deprecated.",
									CustomType:         customfield.NewNestedObjectListType[CloudGPUBaremetalClustersServersSettingsSecurityGroupsDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description: "Security group ID",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "Security group name",
												Computed:    true,
											},
										},
									},
								},
								"ssh_key_name": schema.StringAttribute{
									Description: "SSH key name",
									Computed:    true,
								},
								"user_data": schema.StringAttribute{
									Description: "Optional custom user data",
									Computed:    true,
								},
							},
						},
						"status": schema.StringAttribute{
							Description: "Cluster status\nAvailable values: \"active\", \"creating\", \"degraded\", \"deleting\", \"error\", \"rebooting\", \"rebuilding\", \"resizing\", \"shutoff\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"active",
									"creating",
									"degraded",
									"deleting",
									"error",
									"rebooting",
									"rebuilding",
									"resizing",
									"shutoff",
								),
							},
						},
						"tags": schema.ListNestedAttribute{
							Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudGPUBaremetalClustersTagsDataSourceModel](ctx),
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
						"updated_at": schema.StringAttribute{
							Description: "Cluster update date time",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
					},
				},
			},
		},
	}
}

func (d *CloudGPUBaremetalClustersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudGPUBaremetalClustersDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
