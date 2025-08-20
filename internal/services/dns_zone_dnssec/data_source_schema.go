// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_dnssec

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*DNSZoneDnssecDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"algorithm": schema.StringAttribute{
				Description: "Specifies the algorithm used for the key.",
				Computed:    true,
			},
			"digest": schema.StringAttribute{
				Description: "Represents the hashed value of the DS record.",
				Computed:    true,
			},
			"digest_algorithm": schema.StringAttribute{
				Description: "Specifies the algorithm used to generate the digest.",
				Computed:    true,
			},
			"digest_type": schema.StringAttribute{
				Description: "Specifies the type of the digest algorithm used.",
				Computed:    true,
			},
			"ds": schema.StringAttribute{
				Description: "Represents the complete DS record.",
				Computed:    true,
			},
			"flags": schema.Int64Attribute{
				Description: "Represents the flag for DNSSEC record.",
				Computed:    true,
			},
			"key_tag": schema.Int64Attribute{
				Description: "Represents the identifier of the DNSKEY record.",
				Computed:    true,
			},
			"key_type": schema.StringAttribute{
				Description: "Specifies the type of the key used in the algorithm.",
				Computed:    true,
			},
			"public_key": schema.StringAttribute{
				Description: "Represents the public key used in the DS record.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *DNSZoneDnssecDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *DNSZoneDnssecDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
