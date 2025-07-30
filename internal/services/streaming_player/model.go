// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_player

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type StreamingPlayerModel struct {
	PlayerID             types.Int64                 `tfsdk:"player_id" path:"player_id,optional"`
	Player               *StreamingPlayerPlayerModel `tfsdk:"player" json:"player,optional,no_refresh"`
	Autoplay             types.Bool                  `tfsdk:"autoplay" json:"autoplay,computed"`
	BgColor              types.String                `tfsdk:"bg_color" json:"bg_color,computed"`
	ClientID             types.Int64                 `tfsdk:"client_id" json:"client_id,computed"`
	CustomCss            types.String                `tfsdk:"custom_css" json:"custom_css,computed"`
	Design               types.String                `tfsdk:"design" json:"design,computed"`
	DisableSkin          types.Bool                  `tfsdk:"disable_skin" json:"disable_skin,computed"`
	FgColor              types.String                `tfsdk:"fg_color" json:"fg_color,computed"`
	Framework            types.String                `tfsdk:"framework" json:"framework,computed"`
	HoverColor           types.String                `tfsdk:"hover_color" json:"hover_color,computed"`
	ID                   types.Int64                 `tfsdk:"id" json:"id,computed"`
	JsURL                types.String                `tfsdk:"js_url" json:"js_url,computed"`
	Logo                 types.String                `tfsdk:"logo" json:"logo,computed"`
	LogoPosition         types.String                `tfsdk:"logo_position" json:"logo_position,computed"`
	Mute                 types.Bool                  `tfsdk:"mute" json:"mute,computed"`
	Name                 types.String                `tfsdk:"name" json:"name,computed"`
	SaveOptionsToCookies types.Bool                  `tfsdk:"save_options_to_cookies" json:"save_options_to_cookies,computed"`
	ShowSharing          types.Bool                  `tfsdk:"show_sharing" json:"show_sharing,computed"`
	SkinIsURL            types.String                `tfsdk:"skin_is_url" json:"skin_is_url,computed"`
	SpeedControl         types.Bool                  `tfsdk:"speed_control" json:"speed_control,computed"`
	TextColor            types.String                `tfsdk:"text_color" json:"text_color,computed"`
}

func (m StreamingPlayerModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m StreamingPlayerModel) MarshalJSONForUpdate(state StreamingPlayerModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type StreamingPlayerPlayerModel struct {
	Name                 types.String `tfsdk:"name" json:"name,required"`
	ID                   types.Int64  `tfsdk:"id" json:"id,optional"`
	Autoplay             types.Bool   `tfsdk:"autoplay" json:"autoplay,optional"`
	BgColor              types.String `tfsdk:"bg_color" json:"bg_color,optional"`
	ClientID             types.Int64  `tfsdk:"client_id" json:"client_id,optional"`
	CustomCss            types.String `tfsdk:"custom_css" json:"custom_css,optional"`
	Design               types.String `tfsdk:"design" json:"design,optional"`
	DisableSkin          types.Bool   `tfsdk:"disable_skin" json:"disable_skin,optional"`
	FgColor              types.String `tfsdk:"fg_color" json:"fg_color,optional"`
	Framework            types.String `tfsdk:"framework" json:"framework,optional"`
	HoverColor           types.String `tfsdk:"hover_color" json:"hover_color,optional"`
	JsURL                types.String `tfsdk:"js_url" json:"js_url,optional"`
	Logo                 types.String `tfsdk:"logo" json:"logo,optional"`
	LogoPosition         types.String `tfsdk:"logo_position" json:"logo_position,optional"`
	Mute                 types.Bool   `tfsdk:"mute" json:"mute,optional"`
	SaveOptionsToCookies types.Bool   `tfsdk:"save_options_to_cookies" json:"save_options_to_cookies,optional"`
	ShowSharing          types.Bool   `tfsdk:"show_sharing" json:"show_sharing,optional"`
	SkinIsURL            types.String `tfsdk:"skin_is_url" json:"skin_is_url,optional"`
	SpeedControl         types.Bool   `tfsdk:"speed_control" json:"speed_control,optional"`
	TextColor            types.String `tfsdk:"text_color" json:"text_color,optional"`
}
