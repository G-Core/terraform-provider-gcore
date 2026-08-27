// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WaapDomainDataSourceModel struct {
	ID            types.Int64                                                  `tfsdk:"id" path:"domain_id,computed"`
	DomainID      types.Int64                                                  `tfsdk:"domain_id" path:"domain_id,optional"`
	CreatedAt     timetypes.RFC3339                                            `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CustomPageSet types.Int64                                                  `tfsdk:"custom_page_set" json:"custom_page_set,computed"`
	Name          types.String                                                 `tfsdk:"name" json:"name,computed"`
	Status        types.String                                                 `tfsdk:"status" json:"status,computed"`
	Aliases       customfield.List[types.String]                               `tfsdk:"aliases" json:"aliases,computed"`
	Quotas        customfield.NestedObjectMap[WaapDomainQuotasDataSourceModel] `tfsdk:"quotas" json:"quotas,computed"`
	FindOneBy     *WaapDomainFindOneByDataSourceModel                          `tfsdk:"find_one_by"`
}

func (m *WaapDomainDataSourceModel) toListParams(_ context.Context) (params waap.DomainListParams, diags diag.Diagnostics) {
	mFindOneByIDs := []int64{}
	if m.FindOneBy.IDs != nil {
		for _, item := range *m.FindOneBy.IDs {
			mFindOneByIDs = append(mFindOneByIDs, item.ValueInt64())
		}
	}

	params = waap.DomainListParams{
		IDs: mFindOneByIDs,
	}

	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}
	if !m.FindOneBy.Ordering.IsNull() {
		params.Ordering = waap.DomainListParamsOrdering(m.FindOneBy.Ordering.ValueString())
	}
	if !m.FindOneBy.Status.IsNull() {
		params.Status = waap.DomainListParamsStatus(m.FindOneBy.Status.ValueString())
	}

	return
}

type WaapDomainQuotasDataSourceModel struct {
	Allowed types.Int64 `tfsdk:"allowed" json:"allowed,computed"`
	Current types.Int64 `tfsdk:"current" json:"current,computed"`
}

type WaapDomainFindOneByDataSourceModel struct {
	IDs      *[]types.Int64 `tfsdk:"ids" query:"ids,optional"`
	Name     types.String   `tfsdk:"name" query:"name,optional"`
	Ordering types.String   `tfsdk:"ordering" query:"ordering,optional"`
	Status   types.String   `tfsdk:"status" query:"status,optional"`
}
