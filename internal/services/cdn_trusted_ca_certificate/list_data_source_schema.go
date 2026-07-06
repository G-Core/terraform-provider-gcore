// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_trusted_ca_certificate

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CDNTrustedCaCertificatesDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Trusted CA certificates verify the authenticity of CDN origin servers during HTTPS connections.",
		Attributes: map[string]schema.Attribute{
			"automated": schema.BoolAttribute{
				Description: "How the certificate was issued.\n\nPossible values:\n- **true** – Certificate was issued automatically.\n- **false** – Certificate was added by a user.",
				Optional:    true,
			},
			"resource_id": schema.Int64Attribute{
				Description: "CDN resource ID for which the certificates are requested.",
				Optional:    true,
			},
			"validity_not_after_lte": schema.StringAttribute{
				Description: "Date and time when the certificate become untrusted (ISO 8601/RFC 3339 format, UTC.)\n\nResponse will contain certificates valid until the specified time.",
				Optional:    true,
			},
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CDNTrustedCaCertificatesItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "CA certificate ID.",
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *CDNTrustedCaCertificatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CDNTrustedCaCertificatesDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
