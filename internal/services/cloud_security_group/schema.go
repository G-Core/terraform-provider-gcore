// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*CloudSecurityGroupResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "Security group ID",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"region_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"instances": schema.ListAttribute{
				Description:   "List of instances",
				Optional:      true,
				ElementType:   types.StringType,
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"security_group": schema.SingleNestedAttribute{
				Description: "Security group",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudSecurityGroupSecurityGroupModel](ctx),
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Security group name",
						Required:    true,
					},
					"description": schema.StringAttribute{
						Description: "Security group description",
						Optional:    true,
					},
					"security_group_rules": schema.ListNestedAttribute{
						Description: "Security group rules",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"description": schema.StringAttribute{
									Description: "Rule description",
									Optional:    true,
								},
								"direction": schema.StringAttribute{
									Description: "Ingress or egress, which is the direction in which the security group is applied\nAvailable values: \"egress\", \"ingress\".",
									Optional:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("egress", "ingress"),
									},
								},
								"ethertype": schema.StringAttribute{
									Description: "Ether type\nAvailable values: \"IPv4\", \"IPv6\".",
									Optional:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("IPv4", "IPv6"),
									},
								},
								"port_range_max": schema.Int64Attribute{
									Description: "The maximum port number in the range that is matched by the security group rule",
									Optional:    true,
									Validators: []validator.Int64{
										int64validator.Between(0, 65535),
									},
								},
								"port_range_min": schema.Int64Attribute{
									Description: "The minimum port number in the range that is matched by the security group rule",
									Optional:    true,
									Validators: []validator.Int64{
										int64validator.Between(0, 65535),
									},
								},
								"protocol": schema.StringAttribute{
									Description: "Protocol\nAvailable values: \"ah\", \"any\", \"dccp\", \"egp\", \"esp\", \"gre\", \"icmp\", \"igmp\", \"ipencap\", \"ipip\", \"ipv6-encap\", \"ipv6-frag\", \"ipv6-icmp\", \"ipv6-nonxt\", \"ipv6-opts\", \"ipv6-route\", \"ospf\", \"pgm\", \"rsvp\", \"sctp\", \"tcp\", \"udp\", \"udplite\", \"vrrp\".",
									Optional:    true,
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
									Description: "The remote group UUID to associate with this security group",
									Optional:    true,
								},
								"remote_ip_prefix": schema.StringAttribute{
									Description: "The remote IP prefix that is matched by this security group rule",
									Optional:    true,
								},
							},
						},
					},
					"tags": schema.MapAttribute{
						Description: "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
						Optional:    true,
						ElementType: jsontypes.NormalizedType{},
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplaceIfConfigured()},
			},
			"name": schema.StringAttribute{
				Description: "Name",
				Optional:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Update key-value tags using JSON Merge Patch semantics (RFC 7386). Provide key-value pairs to add or update tags. Set tag values to `null` to remove tags. Unspecified tags remain unchanged. Read-only tags are always preserved and cannot be modified.\n\n**Examples:**\n\n\\* **Add/update tags:** `{'tags': {'environment': 'production', 'team': 'backend'}}` adds new tags or updates existing ones.\n\n\\* **Delete tags:** `{'tags': {'`old_tag`': null}}` removes specific tags.\n\n\\* **Remove all tags:** `{'tags': null}` removes all user-managed tags (read-only tags are preserved).\n\n\\* **Partial update:** `{'tags': {'environment': 'staging'}}` only updates specified tags.\n\n\\* **Mixed operations:** `{'tags': {'environment': 'production', '`cost_center`': 'engineering', '`deprecated_tag`': null}}` adds/updates 'environment' and '`cost_center`' while removing '`deprecated_tag`', preserving other existing tags.\n\n\\* **Replace all:** first delete existing tags with null values, then add new ones in the same request.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"changed_rules": schema.ListNestedAttribute{
				Description: "List of rules to create or delete",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Description: "Action for a rule\nAvailable values: \"create\", \"delete\".",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("create", "delete"),
							},
						},
						"description": schema.StringAttribute{
							Description: "Security grpup rule description",
							Optional:    true,
						},
						"direction": schema.StringAttribute{
							Description: "Ingress or egress, which is the direction in which the security group rule is applied\nAvailable values: \"egress\", \"ingress\".",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("egress", "ingress"),
							},
						},
						"ethertype": schema.StringAttribute{
							Description: "Must be IPv4 or IPv6, and addresses represented in CIDR must match the ingress or egress rules.\nAvailable values: \"IPv4\", \"IPv6\".",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("IPv4", "IPv6"),
							},
						},
						"port_range_max": schema.Int64Attribute{
							Description: "The maximum port number in the range that is matched by the security group rule",
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.Between(0, 65535),
							},
						},
						"port_range_min": schema.Int64Attribute{
							Description: "The minimum port number in the range that is matched by the security group rule",
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.Between(0, 65535),
							},
						},
						"protocol": schema.StringAttribute{
							Description: "Protocol\nAvailable values: \"ah\", \"any\", \"dccp\", \"egp\", \"esp\", \"gre\", \"icmp\", \"igmp\", \"ipencap\", \"ipip\", \"ipv6-encap\", \"ipv6-frag\", \"ipv6-icmp\", \"ipv6-nonxt\", \"ipv6-opts\", \"ipv6-route\", \"ospf\", \"pgm\", \"rsvp\", \"sctp\", \"tcp\", \"udp\", \"udplite\", \"vrrp\".",
							Optional:    true,
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
							Optional:    true,
						},
						"remote_ip_prefix": schema.StringAttribute{
							Description: "The remote IP prefix that is matched by this security group rule",
							Optional:    true,
						},
						"security_group_rule_id": schema.StringAttribute{
							Description: "UUID of rule to be deleted. Required for action 'delete' only",
							Optional:    true,
						},
					},
				},
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
				CustomType:  customfield.NewNestedObjectListType[CloudSecurityGroupSecurityGroupRulesModel](ctx),
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
							Default:     stringdefault.StaticString(""),
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
			"tags_v2": schema.ListNestedAttribute{
				Description: "List of key-value tags associated with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudSecurityGroupTagsV2Model](ctx),
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
		},
	}
}

func (r *CloudSecurityGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudSecurityGroupResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
