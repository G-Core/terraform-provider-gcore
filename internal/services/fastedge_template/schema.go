// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_template

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*FastedgeTemplateResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "FastEdge templates encapsulate reusable configurations for FastEdge applications, including a WebAssembly binary reference and configurable parameters.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "Template ID",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"binary_id": schema.Int64Attribute{
				Description: "ID of the WebAssembly binary to use for this template",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Unique name for the template (used for identification and searching)",
				Required:    true,
			},
			"owned": schema.BoolAttribute{
				Description: "Is the template owned by user?",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"params": schema.ListNestedAttribute{
				Description: "Parameters",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"data_type": schema.StringAttribute{
							Description: "Parameter type determines validation and UI rendering:  \nstring - text input  \nnumber - numeric input  \ndate - date picker  \ntime - time picker  \nsecret - references a secret  \nstore - references an edge store  \nbool - boolean toggle  \njson - JSON editor or multiline text area with JSON validation  \nenum - dropdown/select with allowed values defined via parameter metadata\nAvailable values: \"string\", \"number\", \"date\", \"time\", \"secret\", \"store\", \"bool\", \"json\", \"enum\".",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"string",
									"number",
									"date",
									"time",
									"secret",
									"store",
									"bool",
									"json",
									"enum",
								),
							},
						},
						"mandatory": schema.BoolAttribute{
							Description: "If true, this parameter must be provided when instantiating the template",
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
						},
						"name": schema.StringAttribute{
							Description: "Parameter name used as a placeholder in template (e.g., `API_KEY`, `DATABASE_URL`)",
							Required:    true,
						},
						"descr": schema.StringAttribute{
							Description: "Human-readable explanation of what this parameter controls",
							Optional:    true,
						},
						"metadata": schema.StringAttribute{
							Description: "Optional JSON-encoded string for additional parameter metadata, such as allowed values for enum types",
							Optional:    true,
						},
					},
				},
			},
			"long_descr": schema.StringAttribute{
				Description: "Detailed markdown description explaining template features and usage",
				Optional:    true,
			},
			"short_descr": schema.StringAttribute{
				Description: "Brief one-line description displayed in template listings",
				Optional:    true,
			},
			"api_type": schema.StringAttribute{
				Description: "Wasm API type",
				Computed:    true,
			},
		},
	}
}

func (r *FastedgeTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *FastedgeTemplateResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
