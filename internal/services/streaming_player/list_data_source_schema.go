// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_player

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*StreamingPlayersDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[StreamingPlayersItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Player name",
							Computed:    true,
						},
						"id": schema.Int64Attribute{
							Description: "Player ID",
							Computed:    true,
						},
						"autoplay": schema.BoolAttribute{
							Description: "Enables video playback right after player load:\n\n* **true** — video starts playing right after player loads\n* **false** — video isn’t played automatically. A user must click play to start\n\nDefault is false",
							Computed:    true,
						},
						"bg_color": schema.StringAttribute{
							Description: "Color of skin background in format #AAAAAA",
							Computed:    true,
						},
						"client_id": schema.Int64Attribute{
							Description: "Client ID",
							Computed:    true,
						},
						"custom_css": schema.StringAttribute{
							Description: "Custom CSS to be added to player iframe",
							Computed:    true,
						},
						"design": schema.StringAttribute{
							Description: "String to be rendered as JS parameters to player",
							Computed:    true,
						},
						"disable_skin": schema.BoolAttribute{
							Description: "Enables/Disables player skin:\n\n* **true** — player skin is disabled\n* **false** — player skin is enabled\n\nDefault is false",
							Computed:    true,
						},
						"fg_color": schema.StringAttribute{
							Description: "Color of skin foreground (elements) in format #AAAAAA",
							Computed:    true,
						},
						"framework": schema.StringAttribute{
							Description: "Player framework type",
							Computed:    true,
						},
						"hover_color": schema.StringAttribute{
							Description: "Color of foreground elements when mouse is over in format #AAAAAA",
							Computed:    true,
						},
						"js_url": schema.StringAttribute{
							Description: "Player main JS file URL. Leave empty to use JS URL from the default player",
							Computed:    true,
						},
						"logo": schema.StringAttribute{
							Description: "URL to logo image",
							Computed:    true,
						},
						"logo_position": schema.StringAttribute{
							Description: "Logotype position.   \n Has four possible values:\n\n* **tl** — top left\n* **tr** — top right\n* **bl** — bottom left\n* **br** — bottom right\n\nDefault is null",
							Computed:    true,
						},
						"mute": schema.BoolAttribute{
							Description: "Regulates the sound volume:\n\n* **true** — video starts with volume off\n* **false** — video starts with volume on\n\nDefault is false",
							Computed:    true,
						},
						"save_options_to_cookies": schema.BoolAttribute{
							Description: "Enables/Disables saving volume and other options in cookies:\n\n* **true** — user settings will be saved\n* **false** — user settings will not be saved\n\nDefault is true",
							Computed:    true,
						},
						"show_sharing": schema.BoolAttribute{
							Description: "Enables/Disables sharing button display:\n\n* **true** — sharing button is displayed\n* **false** — no sharing button is displayed\n\nDefault is true",
							Computed:    true,
						},
						"skin_is_url": schema.StringAttribute{
							Description: "URL to custom skin JS file",
							Computed:    true,
						},
						"speed_control": schema.BoolAttribute{
							Description: "Enables/Disables speed control button display:\n\n* **true** — sharing button is displayed\n* **false** — no sharing button is displayed\n\nDefault is false",
							Computed:    true,
						},
						"text_color": schema.StringAttribute{
							Description: "Color of skin text elements in format #AAAAAA",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *StreamingPlayersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *StreamingPlayersDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
