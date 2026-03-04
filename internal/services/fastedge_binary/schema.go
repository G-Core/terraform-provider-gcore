// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_binary

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*FastedgeBinaryResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "FastEdge binaries are immutable WebAssembly modules that implement edge application logic.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "Binary ID",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown(), int64planmodifier.RequiresReplace()},
			},
			"body": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"api_type": schema.StringAttribute{
				Description: "Wasm API type",
				Computed:    true,
			},
			"checksum": schema.StringAttribute{
				Description: "MD5 hash of the binary",
				Computed:    true,
			},
			"source": schema.Int64Attribute{
				Description: "Source language:  \n0 - unknown  \n1 - Rust  \n2 - JavaScript",
				Computed:    true,
			},
			"status": schema.Int64Attribute{
				Description: "Status code:  \n0 - pending  \n1 - compiled  \n2 - compilation failed (errors available)  \n3 - compilation failed (errors not available)  \n4 - resulting binary exceeded the limit  \n5 - unsupported source language",
				Computed:    true,
			},
			"unref_since": schema.StringAttribute{
				Description: "Not used since (UTC)",
				Computed:    true,
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
