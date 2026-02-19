// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapDomainDataSourceModel struct {
	DomainID      types.Int64                                                  `tfsdk:"domain_id" path:"domain_id,required"`
	CreatedAt     timetypes.RFC3339                                            `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CustomPageSet types.Int64                                                  `tfsdk:"custom_page_set" json:"custom_page_set,computed"`
	ID            types.Int64                                                  `tfsdk:"id" json:"id,computed"`
	Name          types.String                                                 `tfsdk:"name" json:"name,computed"`
	Status        types.String                                                 `tfsdk:"status" json:"status,computed"`
	Quotas        customfield.NestedObjectMap[WaapDomainQuotasDataSourceModel] `tfsdk:"quotas" json:"quotas,computed"`
}

type WaapDomainQuotasDataSourceModel struct {
	Allowed types.Int64 `tfsdk:"allowed" json:"allowed,computed"`
	Current types.Int64 `tfsdk:"current" json:"current,computed"`
}
