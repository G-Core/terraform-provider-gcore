// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WaapDomainModel struct {
	DomainID      types.Int64                                        `tfsdk:"domain_id" path:"domain_id,required"`
	Status        types.String                                       `tfsdk:"status" json:"status,required"`
	CreatedAt     timetypes.RFC3339                                  `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CustomPageSet types.Int64                                        `tfsdk:"custom_page_set" json:"custom_page_set,computed"`
	ID            types.Int64                                        `tfsdk:"id" json:"id,computed"`
	Name          types.String                                       `tfsdk:"name" json:"name,computed"`
	Quotas        customfield.NestedObjectMap[WaapDomainQuotasModel] `tfsdk:"quotas" json:"quotas,computed"`
}

func (m WaapDomainModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m WaapDomainModel) MarshalJSONForUpdate(state WaapDomainModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type WaapDomainQuotasModel struct {
	Allowed types.Int64 `tfsdk:"allowed" json:"allowed,computed"`
	Current types.Int64 `tfsdk:"current" json:"current,computed"`
}
