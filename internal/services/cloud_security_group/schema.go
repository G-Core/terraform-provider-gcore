// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudSecurityGroupResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Security groups act as virtual firewalls controlling inbound and outbound traffic for instances and other resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseNonNullStateForUnknown()},
			},
			"project_id": schema.Int64Attribute{
				Description:   "Project ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
			},
			"region_id": schema.Int64Attribute{
				Description:   "Region ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
			},
			"name": schema.StringAttribute{
				Description: "Security group name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description:   "Security group description",
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{
				Description:   "Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.",
				Computed:      true,
				Optional:      true,
				CustomType:    customfield.NewMapType[types.String](ctx),
				ElementType:   types.StringType,
				PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
			},
			"rules": schema.ListNestedAttribute{
				Description: "Security group rules. Omit to apply the default template (ingress + egress); send [] to create no rules.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudSecurityGroupRulesModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"direction": schema.StringAttribute{
							Description: "Ingress or egress, which is the direction in which the security group is applied\nAvailable values: \"egress\", \"ingress\".",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("egress", "ingress"),
							},
						},
						"description": schema.StringAttribute{
							Description: "Rule description",
							Computed:    true,
							Optional:    true,
							Default:     stringdefault.StaticString(""),
						},
						"ethertype": schema.StringAttribute{
							Description: "Ether type\nAvailable values: \"IPv4\", \"IPv6\".",
							Computed:    true,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("IPv4", "IPv6"),
							},
							Default: stringdefault.StaticString("IPv4"),
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
							Description: "V2 protocol enum without 'any'. Use null for all protocols instead.\nAvailable values: \"ah\", \"dccp\", \"egp\", \"esp\", \"gre\", \"icmp\", \"igmp\", \"ipencap\", \"ipip\", \"ipv6-encap\", \"ipv6-frag\", \"ipv6-icmp\", \"ipv6-nonxt\", \"ipv6-opts\", \"ipv6-route\", \"ospf\", \"pgm\", \"rsvp\", \"sctp\", \"tcp\", \"udp\", \"udplite\", \"vrrp\".",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"ah",
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
			"created_at": schema.StringAttribute{
				Description:   "Datetime when the security group was created",
				Computed:      true,
				CustomType:    timetypes.RFC3339Type{},
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"region": schema.StringAttribute{
				Description:   "Region name",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
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
		},
	}
}

func (r *CloudSecurityGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudSecurityGroupResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
