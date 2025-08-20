// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_video_subtitle

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ resource.ResourceWithConfigValidators = (*StreamingVideoSubtitleResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"video_id": schema.Int64Attribute{
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"body": schema.StringAttribute{
				Required:   true,
				CustomType: jsontypes.NormalizedType{},
			},
			"language": schema.StringAttribute{
				Description: "3-letter language code according to ISO-639-2 (bibliographic code)",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of subtitle file",
				Optional:    true,
			},
			"vtt": schema.StringAttribute{
				Description: `Full text of subtitles/captions, with escaped "\n" ("\r") symbol of new line`,
				Optional:    true,
			},
		},
	}
}

func (r *StreamingVideoSubtitleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *StreamingVideoSubtitleResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
