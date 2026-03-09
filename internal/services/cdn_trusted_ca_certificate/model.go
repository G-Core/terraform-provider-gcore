// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_trusted_ca_certificate

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNTrustedCaCertificateModel struct {
	ID                  types.Int64  `tfsdk:"id" json:"id,computed"`
	Name                types.String `tfsdk:"name" json:"name,required"`
	SslCertificate      types.String `tfsdk:"ssl_certificate" json:"sslCertificate,required,no_refresh"`
	CertIssuer          types.String `tfsdk:"cert_issuer" json:"cert_issuer,computed"`
	CertSubjectAlt      types.String `tfsdk:"cert_subject_alt" json:"cert_subject_alt,computed"`
	CertSubjectCn       types.String `tfsdk:"cert_subject_cn" json:"cert_subject_cn,computed"`
	Deleted             types.Bool   `tfsdk:"deleted" json:"deleted,computed"`
	HasRelatedResources types.Bool   `tfsdk:"has_related_resources" json:"hasRelatedResources,computed"`
	SslCertificateChain types.String `tfsdk:"ssl_certificate_chain" json:"sslCertificateChain,computed"`
	ValidityNotAfter    types.String `tfsdk:"validity_not_after" json:"validity_not_after,computed"`
	ValidityNotBefore   types.String `tfsdk:"validity_not_before" json:"validity_not_before,computed"`
}

func (m CDNTrustedCaCertificateModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CDNTrustedCaCertificateModel) MarshalJSONForUpdate(state CDNTrustedCaCertificateModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}
