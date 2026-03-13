// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_binary

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*FastedgeBinaryDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "FastEdge binaries are immutable WebAssembly modules that implement edge application logic.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"binary_id": schema.Int64Attribute{
				Required: true,
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
				Description: "Source language:  \n0 - unknown  \n1 - Rust  \n2 - JavaScript  \n3 - Go",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 3),
				},
			},
			"status": schema.Int64Attribute{
				Description: "Status code:  \n0 - pending  \n1 - compiled  \n2 - compilation failed (errors available)  \n3 - compilation failed (errors not available)  \n4 - resulting binary exceeded the limit  \n5 - unsupported source language",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 5),
				},
			},
			"unref_since": schema.StringAttribute{
				Description: "Not used since (UTC)",
				Computed:    true,
			},
		},
	}
}

func (d *FastedgeBinaryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *FastedgeBinaryDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
