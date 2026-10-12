// This file provides custom Terraform plan modifiers specific to this resource.
package cloud_volume

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
)

// sizeDerivingSource is the only volume source whose size the API derives
// server-side: a restored volume always has the size of its snapshot.
const sizeDerivingSource = "snapshot"

// RequireSizeUnlessDerived reports a plan-time error when a volume is created
// from a source that cannot derive its own size and the configuration leaves
// size out.
//
// size is computed-optional so that a snapshot-restored volume converges
// without declaring a size. The other two sources still require one: the API
// rejects a create without size for `new-volume` and `image` with an opaque 400
// ("Validation Error ... {'size': ['Field required']}"). This modifier turns
// that into an actionable diagnostic before anything is sent.
//
// It only fires while creating. Once the volume exists its size lives in state,
// so an absent size in the configuration is a legitimate no-op — including
// straight after an import, where the configuration may never have declared one.
func RequireSizeUnlessDerived() planmodifier.Int64 {
	return requireSizeUnlessDerivedModifier{}
}

type requireSizeUnlessDerivedModifier struct{}

func (m requireSizeUnlessDerivedModifier) Description(_ context.Context) string {
	return fmt.Sprintf("Requires size on creation unless source is %q.", sizeDerivingSource)
}

func (m requireSizeUnlessDerivedModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m requireSizeUnlessDerivedModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	// Destroy plan: nothing to validate.
	if req.Plan.Raw.IsNull() {
		return
	}

	// Only creation. On update the size held in state stands, so leaving size
	// out of the configuration changes nothing.
	if !req.State.Raw.IsNull() {
		return
	}

	if !req.ConfigValue.IsNull() {
		return
	}

	var source types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("source"), &source)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// A source that is not resolved yet cannot be judged here; leave it to the
	// API on apply.
	if source.IsNull() || source.IsUnknown() {
		return
	}

	if strings.EqualFold(source.ValueString(), sizeDerivingSource) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Missing required argument",
		fmt.Sprintf(
			"size is required when source is %q. Only a volume restored from a snapshot (source = %q) takes its size from the source.",
			source.ValueString(), sizeDerivingSource,
		),
	)
}

// retypableVolumeTypes are the types POST /v1/volumes/{id}/retype accepts, as
// source and as target. Every other type change needs a new volume.
var retypableVolumeTypes = []string{"standard", "ssd_hiiops"}

func isRetypable(t string) bool {
	for _, r := range retypableVolumeTypes {
		if strings.EqualFold(t, r) {
			return true
		}
	}
	return false
}

// TypeNameReplaceUnlessRetypable replaces the volume when a configured
// type_name change cannot be done in place. A change between retypable types
// is applied by Update through the retype API, and a value that matches the
// volume's real type (volume_type) is only recorded. The plan action shows
// which path an apply takes.
func TypeNameReplaceUnlessRetypable() planmodifier.String {
	return typeNameReplaceUnlessRetypableModifier{}
}

type typeNameReplaceUnlessRetypableModifier struct{}

func (m typeNameReplaceUnlessRetypableModifier) Description(_ context.Context) string {
	return fmt.Sprintf("Changes between %s are applied in place; any other change replaces the volume.",
		strings.Join(retypableVolumeTypes, " and "))
}

func (m typeNameReplaceUnlessRetypableModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m typeNameReplaceUnlessRetypableModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Create or destroy.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var volumeType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("volume_type"), &volumeType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if typeNameReplaces(req.ConfigValue, req.PlanValue, req.StateValue, volumeType) {
		resp.RequiresReplace = true
		return
	}
	if req.ConfigValue.IsNull() || !retypePlanned(req.PlanValue, req.StateValue, volumeType) {
		return
	}
	// A retype is planned. The API refuses it while the volume has snapshots;
	// snapshot_ids was just refreshed, so say so now instead of at apply.
	var snapshots customfield.List[types.String]
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("snapshot_ids"), &snapshots)...)
	if resp.Diagnostics.HasError() || snapshots.IsNull() || snapshots.IsUnknown() || len(snapshots.Elements()) == 0 {
		return
	}
	resp.Diagnostics.AddAttributeError(req.Path, "Volume type cannot be changed while the volume has snapshots",
		fmt.Sprintf("Changing type_name from %q to %q is done in place, which the API refuses while the volume "+
			"has snapshots (%s). Delete the snapshots first, or keep type_name unchanged.",
			volumeType.ValueString(), req.PlanValue.ValueString(), snapshots.String()))
}

// retypePlanned reports whether an in-place type_name change must call the
// retype API: the planned value is a change from the recorded one and differs
// from the volume's real type. Update and the snapshot check share it, so an
// unchanged record over a volume_type that contradicts it (2.0.0-rc.1) never
// retypes.
func retypePlanned(plan, state, volumeType types.String) bool {
	if plan.IsNull() || plan.IsUnknown() || strings.EqualFold(plan.ValueString(), state.ValueString()) {
		return false
	}
	current := volumeType
	if current.IsNull() || current.IsUnknown() || current.ValueString() == "" {
		current = state
	}
	return !strings.EqualFold(plan.ValueString(), current.ValueString())
}

// typeNameReplaces is the type_name replacement decision, shared with
// replacementDerived. The volume's real type (volume_type) decides, not the
// recorded type_name: they differ on a 2.0.0-rc.1 adoption that was never
// applied, and restating the real type there is only recorded.
func typeNameReplaces(config, plan, state, volumeType types.String) bool {
	// Omitted from the configuration: the recorded value stands.
	if config.IsNull() || plan.Equal(state) {
		return false
	}
	if plan.IsUnknown() {
		return true
	}
	current := volumeType
	if current.IsNull() || current.IsUnknown() || current.ValueString() == "" {
		current = state
	}
	if current.IsUnknown() {
		return true
	}
	planned := plan.ValueString()
	if strings.EqualFold(planned, current.ValueString()) {
		return false
	}
	return !(isRetypable(planned) && isRetypable(current.ValueString()))
}
