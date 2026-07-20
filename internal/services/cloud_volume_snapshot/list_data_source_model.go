// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume_snapshot

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudVolumeSnapshotsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudVolumeSnapshotsItemsDataSourceModel] `json:"results,computed"`
}

type CloudVolumeSnapshotsDataSourceModel struct {
	ProjectID         types.Int64                                                            `tfsdk:"project_id" path:"project_id,optional"`
	RegionID          types.Int64                                                            `tfsdk:"region_id" path:"region_id,optional"`
	InstanceID        types.String                                                           `tfsdk:"instance_id" query:"instance_id,optional"`
	LifecyclePolicyID types.Int64                                                            `tfsdk:"lifecycle_policy_id" query:"lifecycle_policy_id,optional"`
	ScheduleID        types.String                                                           `tfsdk:"schedule_id" query:"schedule_id,optional"`
	VolumeID          types.String                                                           `tfsdk:"volume_id" query:"volume_id,optional"`
	MaxItems          types.Int64                                                            `tfsdk:"max_items"`
	Items             customfield.NestedObjectList[CloudVolumeSnapshotsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudVolumeSnapshotsDataSourceModel) toListParams(_ context.Context) (params cloud.VolumeSnapshotListParams, diags diag.Diagnostics) {
	params = cloud.VolumeSnapshotListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.InstanceID.IsNull() {
		params.InstanceID = param.NewOpt(m.InstanceID.ValueString())
	}
	if !m.LifecyclePolicyID.IsNull() {
		params.LifecyclePolicyID = param.NewOpt(m.LifecyclePolicyID.ValueInt64())
	}
	if !m.ScheduleID.IsNull() {
		params.ScheduleID = param.NewOpt(m.ScheduleID.ValueString())
	}
	if !m.VolumeID.IsNull() {
		params.VolumeID = param.NewOpt(m.VolumeID.ValueString())
	}

	return
}

type CloudVolumeSnapshotsItemsDataSourceModel struct {
	ID            types.String                                                          `tfsdk:"id" json:"id,computed"`
	CreatedAt     timetypes.RFC3339                                                     `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID types.String                                                          `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Description   types.String                                                          `tfsdk:"description" json:"description,computed"`
	Name          types.String                                                          `tfsdk:"name" json:"name,computed"`
	ProjectID     types.Int64                                                           `tfsdk:"project_id" json:"project_id,computed"`
	Region        types.String                                                          `tfsdk:"region" json:"region,computed"`
	RegionID      types.Int64                                                           `tfsdk:"region_id" json:"region_id,computed"`
	Size          types.Int64                                                           `tfsdk:"size" json:"size,computed"`
	Status        types.String                                                          `tfsdk:"status" json:"status,computed"`
	Tags          customfield.NestedObjectList[CloudVolumeSnapshotsTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	TaskID        types.String                                                          `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt     timetypes.RFC3339                                                     `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	VolumeID      types.String                                                          `tfsdk:"volume_id" json:"volume_id,computed"`
}

type CloudVolumeSnapshotsTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}
