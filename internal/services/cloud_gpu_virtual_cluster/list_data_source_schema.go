// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudGPUVirtualClustersDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "GPU virtual clusters provide managed virtual GPU servers with auto-scaling for parallel computation workloads.",
		Attributes: map[string]schema.Attribute{
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
				CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClustersItemsDataSourceModel](ctx),
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
							CustomType: customfield.NewNestedObjectType[CloudGPUVirtualClustersServersSettingsDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"file_shares": schema.ListNestedAttribute{
									Description: "List of file shares mounted across the cluster.",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClustersServersSettingsFileSharesDataSourceModel](ctx),
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
									CustomType: customfield.NewNestedObjectListType[CloudGPUVirtualClustersServersSettingsInterfacesDataSourceModel](ctx),
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
												CustomType:  customfield.NewNestedObjectType[CloudGPUVirtualClustersServersSettingsInterfacesFloatingIPDataSourceModel](ctx),
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
									Description: "Security groups",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClustersServersSettingsSecurityGroupsDataSourceModel](ctx),
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
								"volumes": schema.ListNestedAttribute{
									Description: "List of volumes",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClustersServersSettingsVolumesDataSourceModel](ctx),
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"boot_index": schema.Int64Attribute{
												Description: "Boot index of the volume",
												Computed:    true,
												Validators: []validator.Int64{
													int64validator.AtLeast(0),
												},
											},
											"delete_on_termination": schema.BoolAttribute{
												Description: "Flag indicating whether the volume is deleted on instance termination",
												Computed:    true,
											},
											"image_id": schema.StringAttribute{
												Description: "Image ID for the volume",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "Volume name",
												Computed:    true,
											},
											"size": schema.Int64Attribute{
												Description: "Volume size in GiB",
												Computed:    true,
											},
											"tags": schema.ListNestedAttribute{
												Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
												Computed:    true,
												CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClustersServersSettingsVolumesTagsDataSourceModel](ctx),
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
											"type": schema.StringAttribute{
												Description: "Volume type\nAvailable values: \"cold\", \"ssd_hiiops\", \"ssd_local\", \"ssd_lowlatency\", \"standard\", \"ultra\".",
												Computed:    true,
												Validators: []validator.String{
													stringvalidator.OneOfCaseInsensitive(
														"cold",
														"ssd_hiiops",
														"ssd_local",
														"ssd_lowlatency",
														"standard",
														"ultra",
													),
												},
											},
										},
									},
								},
							},
						},
						"status": schema.StringAttribute{
							Description: "Cluster status\nAvailable values: \"active\", \"creating\", \"degraded\", \"deleting\", \"error\", \"new\", \"rebooting\", \"rebuilding\", \"resizing\", \"shutoff\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"active",
									"creating",
									"degraded",
									"deleting",
									"error",
									"new",
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
							CustomType:  customfield.NewNestedObjectListType[CloudGPUVirtualClustersTagsDataSourceModel](ctx),
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

func (d *CloudGPUVirtualClustersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudGPUVirtualClustersDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
