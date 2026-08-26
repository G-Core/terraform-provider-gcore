// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_secret

import (
	"context"

	"github.com/G-Core/gcore-go/fastedge"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FastedgeSecretDataSourceModel struct {
	ID          types.Int64                                                           `tfsdk:"id" path:"secret_id,computed"`
	SecretID    types.Int64                                                           `tfsdk:"secret_id" path:"secret_id,optional"`
	AppCount    types.Int64                                                           `tfsdk:"app_count" json:"app_count,computed"`
	Comment     types.String                                                          `tfsdk:"comment" json:"comment,computed"`
	Name        types.String                                                          `tfsdk:"name" json:"name,computed"`
	SecretSlots customfield.NestedObjectSet[FastedgeSecretSecretSlotsDataSourceModel] `tfsdk:"secret_slots" json:"secret_slots,computed"`
	FindOneBy   *FastedgeSecretFindOneByDataSourceModel                               `tfsdk:"find_one_by"`
}

func (m *FastedgeSecretDataSourceModel) toListParams(_ context.Context) (params fastedge.SecretListParams, diags diag.Diagnostics) {
	params = fastedge.SecretListParams{}

	if !m.FindOneBy.AppID.IsNull() {
		params.AppID = param.NewOpt(m.FindOneBy.AppID.ValueInt64())
	}
	if !m.FindOneBy.Search.IsNull() {
		params.Search = param.NewOpt(m.FindOneBy.Search.ValueString())
	}
	if !m.FindOneBy.SecretName.IsNull() {
		params.SecretName = param.NewOpt(m.FindOneBy.SecretName.ValueString())
	}

	return
}

type FastedgeSecretSecretSlotsDataSourceModel struct {
	Slot     types.Int64  `tfsdk:"slot" json:"slot,computed"`
	Checksum types.String `tfsdk:"checksum" json:"checksum,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type FastedgeSecretFindOneByDataSourceModel struct {
	AppID      types.Int64  `tfsdk:"app_id" query:"app_id,optional"`
	Search     types.String `tfsdk:"search" query:"search,optional"`
	SecretName types.String `tfsdk:"secret_name" query:"secret_name,optional"`
}
