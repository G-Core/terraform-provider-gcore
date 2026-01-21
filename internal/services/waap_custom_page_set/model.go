// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_custom_page_set

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type WaapCustomPageSetModel struct {
	ID                 types.Int64                               `tfsdk:"id" json:"id,computed"`
	Name               types.String                              `tfsdk:"name" json:"name,required"`
	Domains            *[]types.Int64                            `tfsdk:"domains" json:"domains,optional"`
	Block              *WaapCustomPageSetBlockModel              `tfsdk:"block" json:"block,optional"`
	BlockCsrf          *WaapCustomPageSetBlockCsrfModel          `tfsdk:"block_csrf" json:"block_csrf,optional"`
	Captcha            *WaapCustomPageSetCaptchaModel            `tfsdk:"captcha" json:"captcha,optional"`
	CookieDisabled     *WaapCustomPageSetCookieDisabledModel     `tfsdk:"cookie_disabled" json:"cookie_disabled,optional"`
	Handshake          *WaapCustomPageSetHandshakeModel          `tfsdk:"handshake" json:"handshake,optional"`
	JavascriptDisabled *WaapCustomPageSetJavascriptDisabledModel `tfsdk:"javascript_disabled" json:"javascript_disabled,optional"`
}

func (m WaapCustomPageSetModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m WaapCustomPageSetModel) MarshalJSONForUpdate(state WaapCustomPageSetModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type WaapCustomPageSetBlockModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Header  types.String `tfsdk:"header" json:"header,optional"`
	Logo    types.String `tfsdk:"logo" json:"logo,optional"`
	Text    types.String `tfsdk:"text" json:"text,optional"`
	Title   types.String `tfsdk:"title" json:"title,optional"`
}

type WaapCustomPageSetBlockCsrfModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Header  types.String `tfsdk:"header" json:"header,optional"`
	Logo    types.String `tfsdk:"logo" json:"logo,optional"`
	Text    types.String `tfsdk:"text" json:"text,optional"`
	Title   types.String `tfsdk:"title" json:"title,optional"`
}

type WaapCustomPageSetCaptchaModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Error   types.String `tfsdk:"error" json:"error,optional"`
	Header  types.String `tfsdk:"header" json:"header,optional"`
	Logo    types.String `tfsdk:"logo" json:"logo,optional"`
	Text    types.String `tfsdk:"text" json:"text,optional"`
	Title   types.String `tfsdk:"title" json:"title,optional"`
}

type WaapCustomPageSetCookieDisabledModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Header  types.String `tfsdk:"header" json:"header,optional"`
	Text    types.String `tfsdk:"text" json:"text,optional"`
}

type WaapCustomPageSetHandshakeModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Header  types.String `tfsdk:"header" json:"header,optional"`
	Logo    types.String `tfsdk:"logo" json:"logo,optional"`
	Title   types.String `tfsdk:"title" json:"title,optional"`
}

type WaapCustomPageSetJavascriptDisabledModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled,required"`
	Header  types.String `tfsdk:"header" json:"header,optional"`
	Text    types.String `tfsdk:"text" json:"text,optional"`
}
