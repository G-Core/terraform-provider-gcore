// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_broadcast

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*StreamingBroadcastDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"broadcast_id": schema.Int64Attribute{
				Required: true,
			},
			"ad_id": schema.Int64Attribute{
				Description: "ID of ad to be displayed in a live stream. If empty the default ad is show. If there is no default ad, no ad is shown",
				Computed:    true,
			},
			"custom_iframe_url": schema.StringAttribute{
				Description: "Custom URL of iframe for video player to be shared via sharing button in player. Auto generated iframe URL is provided by default",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Broadcast name",
				Computed:    true,
			},
			"pending_message": schema.StringAttribute{
				Description: "A custom message that is shown if broadcast status is set to pending. If empty, a default message is shown",
				Computed:    true,
			},
			"player_id": schema.Int64Attribute{
				Description: "ID of player to be used with a broadcast. If empty the default player is used",
				Computed:    true,
			},
			"poster": schema.StringAttribute{
				Description: "Uploaded poster file",
				Computed:    true,
			},
			"share_url": schema.StringAttribute{
				Description: "Custom URL or iframe displayed in the link field when a user clicks on a sharing button in player. If empty, the link field and social network sharing is disabled",
				Computed:    true,
			},
			"show_dvr_after_finish": schema.BoolAttribute{
				Description: "Regulates if a DVR record is shown once a broadcast is finished. Has two possible values:\n\n* **true** — record is shown\n* **false** — record isn't shown\n\n  \nDefault is false",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Broadcast statuses:  \n **Pending** — default “Broadcast isn’t started yet” or custom message (see `pending_message` parameter) is shown, users don't see the live stream  \n **Live** — broadcast is live, and viewers can see it  \n **Paused** — “Broadcast is paused” message is shown, users don't see the live stream  \n **Finished** — “Broadcast is finished” message is shown, users don't see the live stream  \n The users' browsers start displaying the message/stream immediately after you change the broadcast status",
				Computed:    true,
			},
			"stream_ids": schema.ListAttribute{
				Description: "IDs of streams used in a broadcast",
				Computed:    true,
				CustomType:  customfield.NewListType[types.Int64](ctx),
				ElementType: types.Int64Type,
			},
		},
	}
}

func (d *StreamingBroadcastDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *StreamingBroadcastDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
