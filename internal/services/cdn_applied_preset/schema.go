// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_applied_preset

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ resource.ResourceWithConfigValidators = (*CDNAppliedPresetResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Applied presets represent the association between a preset and a CDN resource or rule. Use them to apply a preset to an object, list the objects a preset is applied to, unapply it, and inspect which object fields a preset manages.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:   "ID of the object (CDN resource or rule, according to the preset `object_type`) to apply the preset to.",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseNonNullStateForUnknown(), int64planmodifier.RequiresReplace()},
			},
			"object_id": schema.Int64Attribute{
				Description:   "ID of the object (CDN resource or rule, according to the preset `object_type`) to apply the preset to.",
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseNonNullStateForUnknown(), int64planmodifier.RequiresReplace()},
			},
			"preset_id": schema.Int64Attribute{
				Required:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *CDNAppliedPresetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CDNAppliedPresetResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
