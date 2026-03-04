// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_certificate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*CDNCertificateResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "CDN SSL certificates enable HTTPS content delivery, supporting both uploaded certificates and automated Let's Encrypt provisioning.",
		Attributes: map[string]schema.Attribute{
			"ssl_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "SSL certificate name.\n\nIt must be unique.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"automated": schema.BoolAttribute{
				Description:   "Must be **true** to issue certificate automatically.",
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"ssl_certificate": schema.StringAttribute{
				Description:   "Public part of the SSL certificate.\n\nAll chain of the SSL certificate should be added.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"ssl_private_key": schema.StringAttribute{
				Description:   "Private key of the SSL certificate.",
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"validate_root_ca": schema.BoolAttribute{
				Description:   "Defines whether to check the SSL certificate for a signature from a trusted certificate authority.\n\nPossible values:\n\n- **true** - SSL certificate must be verified to be signed by a trusted certificate authority.\n- **false** - SSL certificate will not be verified to be signed by a trusted certificate authority.",
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
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

func (r *CDNCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CDNCertificateResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
