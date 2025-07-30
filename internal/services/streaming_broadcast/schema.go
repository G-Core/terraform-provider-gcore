// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_broadcast

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ resource.ResourceWithConfigValidators = (*StreamingBroadcastResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"broadcast_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"broadcast": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Broadcast name",
						Required:    true,
					},
					"ad_id": schema.Int64Attribute{
						Description: "ID of ad to be displayed in a live stream. If empty the default ad is show. If there is no default ad, no ad is shown",
						Optional:    true,
					},
					"custom_iframe_url": schema.StringAttribute{
						Description: "Custom URL of iframe for video player to be shared via sharing button in player. Auto generated iframe URL is provided by default",
						Optional:    true,
					},
					"pending_message": schema.StringAttribute{
						Description: "A custom message that is shown if broadcast status is set to pending. If empty, a default message is shown",
						Optional:    true,
					},
					"player_id": schema.Int64Attribute{
						Description: "ID of player to be used with a broadcast. If empty the default player is used",
						Optional:    true,
					},
					"poster": schema.StringAttribute{
						Description: "Uploaded poster file",
						Optional:    true,
					},
					"share_url": schema.StringAttribute{
						Description: "Custom URL or iframe displayed in the link field when a user clicks on a sharing button in player. If empty, the link field and social network sharing is disabled",
						Optional:    true,
					},
					"show_dvr_after_finish": schema.BoolAttribute{
						Description: "Regulates if a DVR record is shown once a broadcast is finished. Has two possible values:\n\n* **true** — record is shown\n* **false** — record isn't shown\n\n  \nDefault is false",
						Optional:    true,
					},
					"status": schema.StringAttribute{
						Description: "Broadcast statuses:  \n **Pending** — default “Broadcast isn’t started yet” or custom message (see `pending_message` parameter) is shown, users don't see the live stream  \n **Live** — broadcast is live, and viewers can see it  \n **Paused** — “Broadcast is paused” message is shown, users don't see the live stream  \n **Finished** — “Broadcast is finished” message is shown, users don't see the live stream  \n The users' browsers start displaying the message/stream immediately after you change the broadcast status",
						Optional:    true,
					},
					"stream_ids": schema.ListAttribute{
						Description: "IDs of streams used in a broadcast",
						Optional:    true,
						ElementType: types.Int64Type,
					},
				},
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

func (r *StreamingBroadcastResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *StreamingBroadcastResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
