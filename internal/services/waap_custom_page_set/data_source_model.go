// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_custom_page_set

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type WaapCustomPageSetDataSourceModel struct {
	SetID              types.Int64                                                                  `tfsdk:"set_id" path:"set_id,required"`
	ID                 types.Int64                                                                  `tfsdk:"id" json:"id,computed"`
	Name               types.String                                                                 `tfsdk:"name" json:"name,computed"`
	Domains            customfield.List[types.Int64]                                                `tfsdk:"domains" json:"domains,computed"`
	Block              customfield.NestedObject[WaapCustomPageSetBlockDataSourceModel]              `tfsdk:"block" json:"block,computed"`
	BlockCsrf          customfield.NestedObject[WaapCustomPageSetBlockCsrfDataSourceModel]          `tfsdk:"block_csrf" json:"block_csrf,computed"`
	Captcha            customfield.NestedObject[WaapCustomPageSetCaptchaDataSourceModel]            `tfsdk:"captcha" json:"captcha,computed"`
	CookieDisabled     customfield.NestedObject[WaapCustomPageSetCookieDisabledDataSourceModel]     `tfsdk:"cookie_disabled" json:"cookie_disabled,computed"`
	Handshake          customfield.NestedObject[WaapCustomPageSetHandshakeDataSourceModel]          `tfsdk:"handshake" json:"handshake,computed"`
	JavascriptDisabled customfield.NestedObject[WaapCustomPageSetJavascriptDisabledDataSourceModel] `tfsdk:"javascript_disabled" json:"javascript_disabled,computed"`
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
