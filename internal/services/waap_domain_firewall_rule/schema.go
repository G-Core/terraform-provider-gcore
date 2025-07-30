// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_firewall_rule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*WaapDomainFirewallRuleResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "The unique identifier of the rule",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"domain_id": schema.Int64Attribute{
				Description:   "The domain ID",
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether or not the rule is enabled",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name assigned to the rule",
				Required:    true,
			},
			"action": schema.SingleNestedAttribute{
				Description: "The action that a firewall rule takes when triggered",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"allow": schema.StringAttribute{
						Description: "The WAAP allowed the request",
						Optional:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"block": schema.SingleNestedAttribute{
						Description: "WAAP block action behavior could be configured with response status code and action duration.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"action_duration": schema.StringAttribute{
								Description: "How long a rule's block action will apply to subsequent requests. Can be specified in seconds or by using a numeral followed by 's', 'm', 'h', or 'd' to represent time format (seconds, minutes, hours, or days)",
								Optional:    true,
							},
							"status_code": schema.Int64Attribute{
								Description: "Designates the HTTP status code to deliver when a request is blocked.\nAvailable values: 403, 405, 418, 429.",
								Optional:    true,
								Validators: []validator.Int64{
									int64validator.OneOf(
										403,
										405,
										418,
										429,
									),
								},
							},
						},
					},
				},
			},
			"conditions": schema.ListNestedAttribute{
				Description: "The condition required for the WAAP engine to trigger the rule.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.SingleNestedAttribute{
							Description: "Match the incoming request against a single IP address",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"ip_address": schema.StringAttribute{
									Description: "A single IPv4 or IPv6 address",
									Required:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						"ip_range": schema.SingleNestedAttribute{
							Description: "Match the incoming request against an IP range",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"lower_bound": schema.StringAttribute{
									Description: "The lower bound IPv4 or IPv6 address to match against",
									Required:    true,
								},
								"upper_bound": schema.StringAttribute{
									Description: "The upper bound IPv4 or IPv6 address to match against",
									Required:    true,
								},
								"negation": schema.BoolAttribute{
									Description: "Whether or not to apply a boolean NOT operation to the rule's condition",
									Computed:    true,
									Optional:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
					},
				},
			},
			"description": schema.StringAttribute{
				Description: "The description assigned to the rule",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
			},
		},
	}
}

func (r *WaapDomainFirewallRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *WaapDomainFirewallRuleResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
