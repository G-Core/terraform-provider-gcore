// This file provides custom Terraform plan modifiers specific to this resource.
package cloud_volume

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
