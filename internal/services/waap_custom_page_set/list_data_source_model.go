// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_custom_page_set

import (
	"context"

	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapCustomPageSetsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[WaapCustomPageSetsItemsDataSourceModel] `json:"results,computed"`
}

type WaapCustomPageSetsDataSourceModel struct {
	Name     types.String                                                         `tfsdk:"name" query:"name,optional"`
	Ordering types.String                                                         `tfsdk:"ordering" query:"ordering,optional"`
	IDs      *[]types.Int64                                                       `tfsdk:"ids" query:"ids,optional"`
	Limit    types.Int64                                                          `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems types.Int64                                                          `tfsdk:"max_items"`
	Items    customfield.NestedObjectList[WaapCustomPageSetsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *WaapCustomPageSetsDataSourceModel) toListParams(_ context.Context) (params waap.CustomPageSetListParams, diags diag.Diagnostics) {
	mIDs := []int64{}
	for _, item := range *m.IDs {
		mIDs = append(mIDs, item.ValueInt64())
	}

	params = waap.CustomPageSetListParams{
		IDs: mIDs,
	}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.Ordering.IsNull() {
		params.Ordering = waap.CustomPageSetListParamsOrdering(m.Ordering.ValueString())
	}

	return
}

type WaapCustomPageSetsItemsDataSourceModel struct {
	ID                 types.Int64                                                                   `tfsdk:"id" json:"id,computed"`
	Name               types.String                                                                  `tfsdk:"name" json:"name,computed"`
	Block              customfield.NestedObject[WaapCustomPageSetsBlockDataSourceModel]              `tfsdk:"block" json:"block,computed"`
	BlockCsrf          customfield.NestedObject[WaapCustomPageSetsBlockCsrfDataSourceModel]          `tfsdk:"block_csrf" json:"block_csrf,computed"`
	Captcha            customfield.NestedObject[WaapCustomPageSetsCaptchaDataSourceModel]            `tfsdk:"captcha" json:"captcha,computed"`
	CookieDisabled     customfield.NestedObject[WaapCustomPageSetsCookieDisabledDataSourceModel]     `tfsdk:"cookie_disabled" json:"cookie_disabled,computed"`
	Domains            customfield.List[types.Int64]                                                 `tfsdk:"domains" json:"domains,computed"`
	Handshake          customfield.NestedObject[WaapCustomPageSetsHandshakeDataSourceModel]          `tfsdk:"handshake" json:"handshake,computed"`
	JavascriptDisabled customfield.NestedObject[WaapCustomPageSetsJavascriptDisabledDataSourceModel] `tfsdk:"javascript_disabled" json:"javascript_disabled,computed"`
}

type WaapCustomPageSetsBlockDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Header  types.String `tfsdk:"header" json:"header,computed"`
	Logo    types.String `tfsdk:"logo" json:"logo,computed"`
	Text    types.String `tfsdk:"text" json:"text,computed"`
	Title   types.String `tfsdk:"title" json:"title,computed"`
}

type WaapCustomPageSetsBlockCsrfDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Header  types.String `tfsdk:"header" json:"header,computed"`
	Logo    types.String `tfsdk:"logo" json:"logo,computed"`
	Text    types.String `tfsdk:"text" json:"text,computed"`
	Title   types.String `tfsdk:"title" json:"title,computed"`
}

type WaapCustomPageSetsCaptchaDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Error   types.String `tfsdk:"error" json:"error,computed"`
	Header  types.String `tfsdk:"header" json:"header,computed"`
	Logo    types.String `tfsdk:"logo" json:"logo,computed"`
	Text    types.String `tfsdk:"text" json:"text,computed"`
	Title   types.String `tfsdk:"title" json:"title,computed"`
}

type WaapCustomPageSetsCookieDisabledDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Header  types.String `tfsdk:"header" json:"header,computed"`
	Text    types.String `tfsdk:"text" json:"text,computed"`
}

type WaapCustomPageSetsHandshakeDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Header  types.String `tfsdk:"header" json:"header,computed"`
	Logo    types.String `tfsdk:"logo" json:"logo,computed"`
	Title   types.String `tfsdk:"title" json:"title,computed"`
}

type WaapCustomPageSetsJavascriptDisabledDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Header  types.String `tfsdk:"header" json:"header,computed"`
	Text    types.String `tfsdk:"text" json:"text,computed"`
}
