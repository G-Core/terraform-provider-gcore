// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume_snapshot

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudVolumeSnapshotModel struct {
	ID            types.String                  `tfsdk:"id" json:"id,computed"`
	ProjectID     types.Int64                   `tfsdk:"project_id" path:"project_id,optional"`
	RegionID      types.Int64                   `tfsdk:"region_id" path:"region_id,optional"`
	VolumeID      types.String                  `tfsdk:"volume_id" json:"volume_id,required"`
	Description   types.String                  `tfsdk:"description" json:"description,optional"`
	Name          types.String                  `tfsdk:"name" json:"name,required"`
	Tags          customfield.Map[types.String] `tfsdk:"tags" json:"tags,computed_optional,no_refresh"`
	CreatedAt     timetypes.RFC3339             `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID types.String                  `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Region        types.String                  `tfsdk:"region" json:"region,computed"`
	Size          types.Int64                   `tfsdk:"size" json:"size,computed"`
	Status        types.String                  `tfsdk:"status" json:"status,computed"`
	UpdatedAt     timetypes.RFC3339             `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

func (m CloudVolumeSnapshotModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudVolumeSnapshotModel) MarshalJSONForUpdate(state CloudVolumeSnapshotModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}
