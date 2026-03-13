// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_binary

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.ResourceWithConfigValidators = (*FastedgeBinaryResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a FastEdge WebAssembly binary. Binaries are immutable - any change to the file content will trigger resource replacement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "Binary ID",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown(), int64planmodifier.RequiresReplace()},
			},
			"filename": schema.StringAttribute{
				Description: "Path to the WebAssembly binary file to upload. Changes to file content (detected via checksum) will trigger resource replacement.",
				Required:    true,
			},
			"api_type": schema.StringAttribute{
				Description:   "Wasm API type",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"checksum": schema.StringAttribute{
				Description:   "MD5 hash of the binary. Computed from the local file and verified against the API response.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"source": schema.Int64Attribute{
				Description:   "Source language:  \n0 - unknown  \n1 - Rust  \n2 - JavaScript",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
				Validators: []validator.Int64{
					int64validator.Between(0, 3),
				},
			},
			"status": schema.Int64Attribute{
				Description:   "Status code:  \n0 - pending  \n1 - compiled  \n2 - compilation failed (errors available)  \n3 - compilation failed (errors not available)  \n4 - resulting binary exceeded the limit  \n5 - unsupported source language",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"unref_since": schema.StringAttribute{
				Description:   "Not used since (UTC)",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *FastedgeBinaryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *FastedgeBinaryResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
