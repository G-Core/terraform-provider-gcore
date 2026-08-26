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

type FastedgeSecretsSecretsListDataSourceEnvelope struct {
	Secrets customfield.NestedObjectList[FastedgeSecretsItemsDataSourceModel] `json:"secrets,computed"`
}

type FastedgeSecretsDataSourceModel struct {
	AppID      types.Int64                                                       `tfsdk:"app_id" query:"app_id,optional"`
	Search     types.String                                                      `tfsdk:"search" query:"search,optional"`
	SecretName types.String                                                      `tfsdk:"secret_name" query:"secret_name,optional"`
	MaxItems   types.Int64                                                       `tfsdk:"max_items"`
	Items      customfield.NestedObjectList[FastedgeSecretsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *FastedgeSecretsDataSourceModel) toListParams(_ context.Context) (params fastedge.SecretListParams, diags diag.Diagnostics) {
	params = fastedge.SecretListParams{}

	if !m.AppID.IsNull() {
		params.AppID = param.NewOpt(m.AppID.ValueInt64())
	}
	if !m.Search.IsNull() {
		params.Search = param.NewOpt(m.Search.ValueString())
	}
	if !m.SecretName.IsNull() {
		params.SecretName = param.NewOpt(m.SecretName.ValueString())
	}

	return
}

type FastedgeSecretsItemsDataSourceModel struct {
	Name     types.String `tfsdk:"name" json:"name,computed"`
	ID       types.Int64  `tfsdk:"id" json:"id,computed"`
	AppCount types.Int64  `tfsdk:"app_count" json:"app_count,computed"`
	Comment  types.String `tfsdk:"comment" json:"comment,computed"`
}
