package cloud_volume

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/G-Core/terraform-provider-gcore/internal/logging"
)

// retypeIfChanged changes the volume type in place when retypePlanned says so
// and reports whether it did. The plan modifier only lets a retypable change
// reach Update; anything else was planned as a replacement.
func (r *CloudVolumeResource) retypeIfChanged(ctx context.Context, data, state *CloudVolumeModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	if !retypePlanned(data.TypeName, state.TypeName, state.VolumeType) {
		return false, diags
	}

	params := cloud.VolumeChangeTypeParams{
		VolumeType: cloud.VolumeChangeTypeParamsVolumeType(data.TypeName.ValueString()),
	}
	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}
	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	_, err := r.client.Cloud.Volumes.ChangeTypeAndPoll(
		ctx,
		data.ID.ValueString(),
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		diags.AddError("failed to change volume type", err.Error())
		return false, diags
	}
	return true, diags
}
