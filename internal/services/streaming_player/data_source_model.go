// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_player

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StreamingPlayerDataSourceModel struct {
	PlayerID             types.Int64  `tfsdk:"player_id" path:"player_id,required"`
	Autoplay             types.Bool   `tfsdk:"autoplay" json:"autoplay,computed"`
	BgColor              types.String `tfsdk:"bg_color" json:"bg_color,computed"`
	ClientID             types.Int64  `tfsdk:"client_id" json:"client_id,computed"`
	CustomCss            types.String `tfsdk:"custom_css" json:"custom_css,computed"`
	Design               types.String `tfsdk:"design" json:"design,computed"`
	DisableSkin          types.Bool   `tfsdk:"disable_skin" json:"disable_skin,computed"`
	FgColor              types.String `tfsdk:"fg_color" json:"fg_color,computed"`
	Framework            types.String `tfsdk:"framework" json:"framework,computed"`
	HoverColor           types.String `tfsdk:"hover_color" json:"hover_color,computed"`
	ID                   types.Int64  `tfsdk:"id" json:"id,computed"`
	JsURL                types.String `tfsdk:"js_url" json:"js_url,computed"`
	Logo                 types.String `tfsdk:"logo" json:"logo,computed"`
	LogoPosition         types.String `tfsdk:"logo_position" json:"logo_position,computed"`
	Mute                 types.Bool   `tfsdk:"mute" json:"mute,computed"`
	Name                 types.String `tfsdk:"name" json:"name,computed"`
	SaveOptionsToCookies types.Bool   `tfsdk:"save_options_to_cookies" json:"save_options_to_cookies,computed"`
	ShowSharing          types.Bool   `tfsdk:"show_sharing" json:"show_sharing,computed"`
	SkinIsURL            types.String `tfsdk:"skin_is_url" json:"skin_is_url,computed"`
	SpeedControl         types.Bool   `tfsdk:"speed_control" json:"speed_control,computed"`
	TextColor            types.String `tfsdk:"text_color" json:"text_color,computed"`
}
