// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_certificate

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNCertificateDataSourceModel struct {
	SslID               types.Int64  `tfsdk:"ssl_id" path:"ssl_id,required"`
	Automated           types.Bool   `tfsdk:"automated" json:"automated,computed"`
	CertIssuer          types.String `tfsdk:"cert_issuer" json:"cert_issuer,computed"`
	CertSubjectAlt      types.String `tfsdk:"cert_subject_alt" json:"cert_subject_alt,computed"`
	CertSubjectCn       types.String `tfsdk:"cert_subject_cn" json:"cert_subject_cn,computed"`
	Deleted             types.Bool   `tfsdk:"deleted" json:"deleted,computed"`
	HasRelatedResources types.Bool   `tfsdk:"has_related_resources" json:"hasRelatedResources,computed"`
	ID                  types.Int64  `tfsdk:"id" json:"id,computed"`
	Name                types.String `tfsdk:"name" json:"name,computed"`
	SslCertificateChain types.String `tfsdk:"ssl_certificate_chain" json:"sslCertificateChain,computed"`
	ValidityNotAfter    types.String `tfsdk:"validity_not_after" json:"validity_not_after,computed"`
	ValidityNotBefore   types.String `tfsdk:"validity_not_before" json:"validity_not_before,computed"`
}
