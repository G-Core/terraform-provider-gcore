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

type WaapCustomPageSetDataSourceModel struct {
	ID                 types.Int64                                                                  `tfsdk:"id" path:"set_id,computed"`
	SetID              types.Int64                                                                  `tfsdk:"set_id" path:"set_id,optional"`
	Name               types.String                                                                 `tfsdk:"name" json:"name,computed"`
	Domains            customfield.List[types.Int64]                                                `tfsdk:"domains" json:"domains,computed"`
	Block              customfield.NestedObject[WaapCustomPageSetBlockDataSourceModel]              `tfsdk:"block" json:"block,computed"`
	BlockCsrf          customfield.NestedObject[WaapCustomPageSetBlockCsrfDataSourceModel]          `tfsdk:"block_csrf" json:"block_csrf,computed"`
	Captcha            customfield.NestedObject[WaapCustomPageSetCaptchaDataSourceModel]            `tfsdk:"captcha" json:"captcha,computed"`
	CookieDisabled     customfield.NestedObject[WaapCustomPageSetCookieDisabledDataSourceModel]     `tfsdk:"cookie_disabled" json:"cookie_disabled,computed"`
	Handshake          customfield.NestedObject[WaapCustomPageSetHandshakeDataSourceModel]          `tfsdk:"handshake" json:"handshake,computed"`
	JavascriptDisabled customfield.NestedObject[WaapCustomPageSetJavascriptDisabledDataSourceModel] `tfsdk:"javascript_disabled" json:"javascript_disabled,computed"`
	FindOneBy          *WaapCustomPageSetFindOneByDataSourceModel                                   `tfsdk:"find_one_by"`
}

func (m *WaapCustomPageSetDataSourceModel) toListParams(_ context.Context) (params waap.CustomPageSetListParams, diags diag.Diagnostics) {
	mFindOneByIDs := []int64{}
	if m.FindOneBy.IDs != nil {
		for _, item := range *m.FindOneBy.IDs {
			mFindOneByIDs = append(mFindOneByIDs, item.ValueInt64())
		}
	}

	params = waap.CustomPageSetListParams{
		IDs: mFindOneByIDs,
	}

	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}
	if !m.FindOneBy.Ordering.IsNull() {
		params.Ordering = waap.CustomPageSetListParamsOrdering(m.FindOneBy.Ordering.ValueString())
	}

	return
}

type WaapCustomPageSetBlockDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Header  types.String `tfsdk:"header" json:"header,computed"`
	Logo    types.String `tfsdk:"logo" json:"logo,computed"`
	Text    types.String `tfsdk:"text" json:"text,computed"`
	Title   types.String `tfsdk:"title" json:"title,computed"`
}

type WaapCustomPageSetBlockCsrfDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Header  types.String `tfsdk:"header" json:"header,computed"`
	Logo    types.String `tfsdk:"logo" json:"logo,computed"`
	Text    types.String `tfsdk:"text" json:"text,computed"`
	Title   types.String `tfsdk:"title" json:"title,computed"`
}

type WaapCustomPageSetCaptchaDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Error   types.String `tfsdk:"error" json:"error,computed"`
	Header  types.String `tfsdk:"header" json:"header,computed"`
	Logo    types.String `tfsdk:"logo" json:"logo,computed"`
	Text    types.String `tfsdk:"text" json:"text,computed"`
	Title   types.String `tfsdk:"title" json:"title,computed"`
}

type WaapCustomPageSetCookieDisabledDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Header  types.String `tfsdk:"header" json:"header,computed"`
	Text    types.String `tfsdk:"text" json:"text,computed"`
}

type WaapCustomPageSetHandshakeDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Header  types.String `tfsdk:"header" json:"header,computed"`
	Logo    types.String `tfsdk:"logo" json:"logo,computed"`
	Title   types.String `tfsdk:"title" json:"title,computed"`
}

type WaapCustomPageSetJavascriptDisabledDataSourceModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,computed"`
	Header  types.String `tfsdk:"header" json:"header,computed"`
	Text    types.String `tfsdk:"text" json:"text,computed"`
}

type WaapCustomPageSetFindOneByDataSourceModel struct {
	IDs      *[]types.Int64 `tfsdk:"ids" query:"ids,optional"`
	Name     types.String   `tfsdk:"name" query:"name,optional"`
	Ordering types.String   `tfsdk:"ordering" query:"ordering,optional"`
}
