// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_trusted_ca_certificate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*CDNTrustedCaCertificateDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"cert_issuer": schema.StringAttribute{
				Description: "Name of the certification center that issued the CA certificate.",
				Computed:    true,
			},
			"cert_subject_alt": schema.StringAttribute{
				Description: "Alternative domain names that the CA certificate secures.",
				Computed:    true,
			},
			"cert_subject_cn": schema.StringAttribute{
				Description: "Domain name that the CA certificate secures.",
				Computed:    true,
			},
			"deleted": schema.BoolAttribute{
				Description: "Defines whether the certificate has been deleted. Parameter is **deprecated**.\n\nPossible values:\n- **true** - Certificate has been deleted.\n- **false** - Certificate has not been deleted.",
				Computed:    true,
			},
			"has_related_resources": schema.BoolAttribute{
				Description: "Defines whether the CA certificate is used by a CDN resource.\n\nPossible values:\n- **true** - Certificate is used by a CDN resource.\n- **false** - Certificate is not used by a CDN resource.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "CA certificate name.",
				Computed:    true,
			},
			"ssl_certificate_chain": schema.StringAttribute{
				Description: "Parameter is **deprecated**.",
				Computed:    true,
			},
			"validity_not_after": schema.StringAttribute{
				Description: "Date when the CA certificate become untrusted (ISO 8601/RFC 3339 format, UTC.)",
				Computed:    true,
			},
			"validity_not_before": schema.StringAttribute{
				Description: "Date when the CA certificate become valid (ISO 8601/RFC 3339 format, UTC.)",
				Computed:    true,
			},
		},
	}
}

func (d *CDNTrustedCaCertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CDNTrustedCaCertificateDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
