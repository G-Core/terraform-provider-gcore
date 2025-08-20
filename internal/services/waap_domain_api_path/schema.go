// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_path

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*WaapDomainAPIPathResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "The path ID",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"domain_id": schema.Int64Attribute{
				Description:   "The domain ID",
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"http_scheme": schema.StringAttribute{
				Description: "The different HTTP schemes an API path can have\nAvailable values: \"HTTP\", \"HTTPS\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("HTTP", "HTTPS"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"method": schema.StringAttribute{
				Description: "The different methods an API path can have\nAvailable values: \"GET\", \"POST\", \"PUT\", \"PATCH\", \"DELETE\", \"TRACE\", \"HEAD\", \"OPTIONS\".",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"GET",
						"POST",
						"PUT",
						"PATCH",
						"DELETE",
						"TRACE",
						"HEAD",
						"OPTIONS",
					),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"api_version": schema.StringAttribute{
				Description:   "The API version",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"path": schema.StringAttribute{
				Description: "The API path, locations that are saved for resource IDs will be put in curly brackets",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Description: "The different statuses an API path can have\nAvailable values: \"CONFIRMED_API\", \"POTENTIAL_API\", \"NOT_API\", \"DELISTED_API\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"CONFIRMED_API",
						"POTENTIAL_API",
						"NOT_API",
						"DELISTED_API",
					),
				},
			},
			"api_groups": schema.ListAttribute{
				Description: "An array of api groups associated with the API path",
				Optional:    true,
				ElementType: types.StringType,
			},
			"tags": schema.ListAttribute{
				Description: "An array of tags associated with the API path",
				Optional:    true,
				ElementType: types.StringType,
			},
			"first_detected": schema.StringAttribute{
				Description: "The date and time in ISO 8601 format the API path was first detected.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"last_detected": schema.StringAttribute{
				Description: "The date and time in ISO 8601 format the API path was last detected.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"request_count": schema.Int64Attribute{
				Description: "The number of requests for this path in the last 24 hours",
				Computed:    true,
			},
			"source": schema.StringAttribute{
				Description: "The different sources an API path can have\nAvailable values: \"API_DESCRIPTION_FILE\", \"TRAFFIC_SCAN\", \"USER_DEFINED\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"API_DESCRIPTION_FILE",
						"TRAFFIC_SCAN",
						"USER_DEFINED",
					),
				},
			},
		},
	}
}

func (r *WaapDomainAPIPathResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *WaapDomainAPIPathResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
