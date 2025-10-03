// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudSecurityGroupsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Optional: true,
			},
			"region_id": schema.Int64Attribute{
				Optional: true,
			},
			"tag_key_value": schema.StringAttribute{
				Description: "Filter by tag key-value pairs. Must be a valid JSON string.",
				Optional:    true,
			},
			"tag_key": schema.ListAttribute{
				Description: "Filter by tag keys.",
				Optional:    true,
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
				CustomType:  customfield.NewNestedObjectListType[CloudSecurityGroupsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Security group ID",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Datetime when the security group was created",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"name": schema.StringAttribute{
							Description: "Security group name",
							Computed:    true,
						},
						"project_id": schema.Int64Attribute{
							Description: "Project ID",
							Computed:    true,
						},
						"region": schema.StringAttribute{
							Description: "Region name",
							Computed:    true,
						},
						"region_id": schema.Int64Attribute{
							Description: "Region ID",
							Computed:    true,
						},
						"revision_number": schema.Int64Attribute{
							Description: "The number of revisions",
							Computed:    true,
						},
						"tags_v2": schema.ListNestedAttribute{
							Description: "Tags for a security group",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudSecurityGroupsTagsV2DataSourceModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Description: "Tag key. The maximum size for a key is 255 bytes.",
										Computed:    true,
									},
									"read_only": schema.BoolAttribute{
										Description: "If true, the tag is read-only and cannot be modified by the user",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Tag value. The maximum size for a value is 1024 bytes.",
										Computed:    true,
									},
								},
							},
						},
						"updated_at": schema.StringAttribute{
							Description: "Datetime when the security group was last updated",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"description": schema.StringAttribute{
							Description: "Security group description",
							Computed:    true,
						},
						"security_group_rules": schema.ListNestedAttribute{
							Description: "Security group rules",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudSecurityGroupsSecurityGroupRulesDataSourceModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "The ID of the security group rule",
										Computed:    true,
									},
									"created_at": schema.StringAttribute{
										Description: "Datetime when the rule was created",
										Computed:    true,
										CustomType:  timetypes.RFC3339Type{},
									},
									"direction": schema.StringAttribute{
										Description: "Ingress or egress, which is the direction in which the security group rule is applied\nAvailable values: \"egress\", \"ingress\".",
										Computed:    true,
										Validators: []validator.String{
											stringvalidator.OneOfCaseInsensitive("egress", "ingress"),
										},
									},
									"revision_number": schema.Int64Attribute{
										Description: "The revision number of the resource",
										Computed:    true,
									},
									"security_group_id": schema.StringAttribute{
										Description: "The security group ID to associate with this security group rule",
										Computed:    true,
									},
									"updated_at": schema.StringAttribute{
										Description: "Datetime when the rule was last updated",
										Computed:    true,
										CustomType:  timetypes.RFC3339Type{},
									},
									"description": schema.StringAttribute{
										Description: "Rule description",
										Computed:    true,
									},
									"ethertype": schema.StringAttribute{
										Description: "Must be IPv4 or IPv6, and addresses represented in CIDR must match the ingress or egress rules.\nAvailable values: \"IPv4\", \"IPv6\".",
										Computed:    true,
										Validators: []validator.String{
											stringvalidator.OneOfCaseInsensitive("IPv4", "IPv6"),
										},
									},
									"port_range_max": schema.Int64Attribute{
										Description: "The maximum port number in the range that is matched by the security group rule",
										Computed:    true,
										Validators: []validator.Int64{
											int64validator.Between(0, 65535),
										},
									},
									"port_range_min": schema.Int64Attribute{
										Description: "The minimum port number in the range that is matched by the security group rule",
										Computed:    true,
										Validators: []validator.Int64{
											int64validator.Between(0, 65535),
										},
									},
									"protocol": schema.StringAttribute{
										Description: "Protocol\nAvailable values: \"ah\", \"any\", \"dccp\", \"egp\", \"esp\", \"gre\", \"icmp\", \"igmp\", \"ipencap\", \"ipip\", \"ipv6-encap\", \"ipv6-frag\", \"ipv6-icmp\", \"ipv6-nonxt\", \"ipv6-opts\", \"ipv6-route\", \"ospf\", \"pgm\", \"rsvp\", \"sctp\", \"tcp\", \"udp\", \"udplite\", \"vrrp\".",
										Computed:    true,
										Validators: []validator.String{
											stringvalidator.OneOfCaseInsensitive(
												"ah",
												"any",
												"dccp",
												"egp",
												"esp",
												"gre",
												"icmp",
												"igmp",
												"ipencap",
												"ipip",
												"ipv6-encap",
												"ipv6-frag",
												"ipv6-icmp",
												"ipv6-nonxt",
												"ipv6-opts",
												"ipv6-route",
												"ospf",
												"pgm",
												"rsvp",
												"sctp",
												"tcp",
												"udp",
												"udplite",
												"vrrp",
											),
										},
									},
									"remote_group_id": schema.StringAttribute{
										Description: "The remote group UUID to associate with this security group rule",
										Computed:    true,
									},
									"remote_ip_prefix": schema.StringAttribute{
										Description: "The remote IP prefix that is matched by this security group rule",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *CloudSecurityGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudSecurityGroupsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
