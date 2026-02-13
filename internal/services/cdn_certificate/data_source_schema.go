// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_certificate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigValidators = (*CDNCertificateDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ssl_id": schema.Int64Attribute{
				Required: true,
			},
			"automated": schema.BoolAttribute{
				Description: "How the SSL certificate was issued.\n\nPossible values:\n- **true** - Certificate was issued automatically.\n- **false** - Certificate was added by a use.",
				Computed:    true,
			},
			"cert_issuer": schema.StringAttribute{
				Description: "Name of the certification center issued the SSL certificate.",
				Computed:    true,
			},
			"cert_subject_alt": schema.StringAttribute{
				Description: "Alternative domain names that the SSL certificate secures.",
				Computed:    true,
			},
			"cert_subject_cn": schema.StringAttribute{
				Description: "Domain name that the SSL certificate secures.",
				Computed:    true,
			},
			"deleted": schema.BoolAttribute{
				Description: "Defines whether the certificate has been deleted. Parameter is **deprecated**.\n\nPossible values:\n- **true** - Certificate has been deleted.\n- **false** - Certificate has not been deleted.",
				Computed:    true,
			},
			"has_related_resources": schema.BoolAttribute{
				Description: "Defines whether the SSL certificate is used by a CDN resource.\n\nPossible values:\n- **true** - Certificate is used by a CDN resource.\n- **false** - Certificate is not used by a CDN resource.",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "SSL certificate ID.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "SSL certificate name.",
				Computed:    true,
			},
			"ssl_certificate_chain": schema.StringAttribute{
				Description: "Parameter is **deprecated**.",
				Computed:    true,
			},
			"validity_not_after": schema.StringAttribute{
				Description: "Date when certificate become untrusted (ISO 8601/RFC 3339 format, UTC.)",
				Computed:    true,
			},
			"validity_not_before": schema.StringAttribute{
				Description: "Date when certificate become valid (ISO 8601/RFC 3339 format, UTC.)",
				Computed:    true,
			},
		},
	}
}

func (d *CDNCertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CDNCertificateDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
