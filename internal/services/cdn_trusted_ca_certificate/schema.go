// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_trusted_ca_certificate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*CDNTrustedCaCertificateResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Trusted CA certificates verify the authenticity of CDN origin servers during HTTPS connections.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "CA certificate ID.",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown(), int64planmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "CA certificate name.\n\nIt must be unique.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"ssl_certificate": schema.StringAttribute{
				Description:   "Public part of the CA certificate.\n\nIt must be in the PEM format.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
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

func (r *CDNTrustedCaCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CDNTrustedCaCertificateResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
