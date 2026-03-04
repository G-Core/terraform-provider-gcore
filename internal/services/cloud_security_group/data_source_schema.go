// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudSecurityGroupDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Group ID",
				Computed:    true,
			},
			"group_id": schema.StringAttribute{
				Description: "Group ID",
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
				Description: "Datetime when the security group was created",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"description": schema.StringAttribute{
				Description: "Security group description",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Security group name",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region name",
				Computed:    true,
			},
			"revision_number": schema.Int64Attribute{
				Description: "The number of revisions",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Datetime when the security group was last updated",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"security_group_rules": schema.ListNestedAttribute{
				Description: "Security group rules",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudSecurityGroupSecurityGroupRulesDataSourceModel](ctx),
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
						"description": schema.StringAttribute{
							Description: "Rule description",
							Computed:    true,
						},
						"direction": schema.StringAttribute{
							Description: "Ingress or egress, which is the direction in which the security group rule is applied\nAvailable values: \"egress\", \"ingress\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("egress", "ingress"),
							},
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
					},
				},
			},
			"tags_v2": schema.ListNestedAttribute{
				Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudSecurityGroupTagsV2DataSourceModel](ctx),
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
					"name": schema.StringAttribute{
						Description: "Optional. Filter by name. Must be specified a full name of the security group.",
						Optional:    true,
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

func (d *CloudSecurityGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudSecurityGroupDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("group_id"), path.MatchRoot("find_one_by")),
	}
}
