// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_advanced_rule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*WaapDomainAdvancedRuleResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "The unique identifier for the rule",
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
			"source": schema.StringAttribute{
				Description: "A CEL syntax expression that contains the rule's conditions. Allowed objects are: request, whois, session, response, tags, `user_defined_tags`, `user_agent`, `client_data`.\n\nMore info can be found here: https://gcore.com/docs/waap/waap-rules/advanced-rules",
				Required:    true,
			},
			"action": schema.SingleNestedAttribute{
				Description: "The action that the rule takes when triggered. Only one action can be set per rule.",
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
								Description: "How long a rule's block action will apply to subsequent requests. Can be specified in seconds or by using a numeral followed by 's', 'm', 'h', or 'd' to represent time format (seconds, minutes, hours, or days). Empty time intervals are not allowed.",
								Optional:    true,
							},
							"status_code": schema.Int64Attribute{
								Description: "A custom HTTP status code that the WAAP returns if a rule blocks a request\nAvailable values: 403, 405, 418, 429.",
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
					"captcha": schema.StringAttribute{
						Description: "The WAAP presented the user with a captcha",
						Optional:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"handshake": schema.StringAttribute{
						Description: "The WAAP performed automatic browser validation",
						Optional:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"monitor": schema.StringAttribute{
						Description: "The WAAP monitored the request but took no action",
						Optional:    true,
						CustomType:  jsontypes.NormalizedType{},
					},
					"tag": schema.SingleNestedAttribute{
						Description: "WAAP tag action gets a list of tags to tag the request scope with",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"tags": schema.ListAttribute{
								Description: "The list of user defined tags to tag the request with",
								Required:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
			},
			"description": schema.StringAttribute{
				Description: "The description assigned to the rule",
				Optional:    true,
			},
			"phase": schema.StringAttribute{
				Description: "The WAAP request/response phase for applying the rule. Default is \"access\".\n\n\nThe \"access\" phase is responsible for modifying the request before it is sent to the origin server. \n\nThe \"header_filter\" phase is responsible for modifying the HTTP headers of a response before they are sent back to the client.\n\nThe \"body_filter\" phase is responsible for modifying the body of a response before it is sent back to the client.\nAvailable values: \"access\", \"header_filter\", \"body_filter\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"access",
						"header_filter",
						"body_filter",
					),
				},
				Default: stringdefault.StaticString("access"),
			},
		},
	}
}

func (r *WaapDomainAdvancedRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *WaapDomainAdvancedRuleResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
