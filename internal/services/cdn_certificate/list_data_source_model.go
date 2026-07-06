// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_certificate

import (
	"context"

	"github.com/G-Core/gcore-go/cdn"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CDNCertificatesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CDNCertificatesItemsDataSourceModel] `json:"results,computed"`
}

type CDNCertificatesDataSourceModel struct {
	Automated           types.Bool                                                        `tfsdk:"automated" query:"automated,optional"`
	ResourceID          types.Int64                                                       `tfsdk:"resource_id" query:"resource_id,optional"`
	ValidityNotAfterLte types.String                                                      `tfsdk:"validity_not_after_lte" query:"validity_not_after_lte,optional"`
	MaxItems            types.Int64                                                       `tfsdk:"max_items"`
	Items               customfield.NestedObjectList[CDNCertificatesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CDNCertificatesDataSourceModel) toListParams(_ context.Context) (params cdn.CertificateListParams, diags diag.Diagnostics) {
	params = cdn.CertificateListParams{}

	if !m.Automated.IsNull() {
		params.Automated = param.NewOpt(m.Automated.ValueBool())
	}
	if !m.ResourceID.IsNull() {
		params.ResourceID = param.NewOpt(m.ResourceID.ValueInt64())
	}
	if !m.ValidityNotAfterLte.IsNull() {
		params.ValidityNotAfterLte = param.NewOpt(m.ValidityNotAfterLte.ValueString())
	}

	return
}

type CDNCertificatesItemsDataSourceModel struct {
	ID                  types.Int64  `tfsdk:"id" json:"id,computed"`
	Automated           types.Bool   `tfsdk:"automated" json:"automated,computed"`
	CertIssuer          types.String `tfsdk:"cert_issuer" json:"cert_issuer,computed"`
	CertSubjectAlt      types.String `tfsdk:"cert_subject_alt" json:"cert_subject_alt,computed"`
	CertSubjectCn       types.String `tfsdk:"cert_subject_cn" json:"cert_subject_cn,computed"`
	Deleted             types.Bool   `tfsdk:"deleted" json:"deleted,computed"`
	HasRelatedResources types.Bool   `tfsdk:"has_related_resources" json:"hasRelatedResources,computed"`
	Name                types.String `tfsdk:"name" json:"name,computed"`
	SslCertificateChain types.String `tfsdk:"ssl_certificate_chain" json:"sslCertificateChain,computed"`
	ValidityNotAfter    types.String `tfsdk:"validity_not_after" json:"validity_not_after,computed"`
	ValidityNotBefore   types.String `tfsdk:"validity_not_before" json:"validity_not_before,computed"`
}
