package cloud_gpu_baremetal_cluster_image

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// cowFormatPlanModifier plans cow_format against the image's real storage
// format, read from the disk_format already in state (raw = true, qcow2 =
// false), rather than from the state copy of cow_format.
//
// The state copy can be stale: state written before cow_format was recovered
// holds null, or the old default false on a raw image, and `-refresh=false`
// skips the Read that would repair it. So:
//
//   - not configured: the planned value is the one disk_format implies, which
//     is what a refresh would have recorded; when disk_format does not map,
//     an unknown plan keeps the prior state value;
//   - configured: replacement is required when the value differs from what
//     disk_format implies, and when disk_format does not map, when it differs
//     from the state copy (RequiresReplaceIfConfigured semantics, so a null
//     prior still replaces).
//
// Creation and destruction are untouched.
func cowFormatPlanModifier() planmodifier.Bool {
	return cowFormatPlanModifierImpl{}
}

type cowFormatPlanModifierImpl struct{}

func (m cowFormatPlanModifierImpl) Description(_ context.Context) string {
	return "Plans cow_format from the image's storage format (disk_format); a configured value that differs from it requires replacement."
}

func (m cowFormatPlanModifierImpl) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m cowFormatPlanModifierImpl) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var diskFormat types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("disk_format"), &diskFormat)...)
	if resp.Diagnostics.HasError() {
		return
	}
	actual := cowFormatFromDiskFormat(diskFormat)

	if req.ConfigValue.IsNull() {
		switch {
		case !actual.IsNull():
			resp.PlanValue = actual
		case req.PlanValue.IsUnknown() && !req.StateValue.IsNull():
			resp.PlanValue = req.StateValue
		}
		return
	}

	if actual.IsNull() {
		actual = req.StateValue
	}

	if !req.PlanValue.Equal(actual) {
		resp.RequiresReplace = true
	}
}
