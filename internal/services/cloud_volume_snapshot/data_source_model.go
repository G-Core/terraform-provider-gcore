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

type CloudVolumeSnapshotDataSourceModel struct {
	ID            types.String                                                         `tfsdk:"id" path:"snapshot_id,computed"`
	SnapshotID    types.String                                                         `tfsdk:"snapshot_id" path:"snapshot_id,optional"`
	ProjectID     types.Int64                                                          `tfsdk:"project_id" path:"project_id,optional"`
	RegionID      types.Int64                                                          `tfsdk:"region_id" path:"region_id,optional"`
	CreatedAt     timetypes.RFC3339                                                    `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID types.String                                                         `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Description   types.String                                                         `tfsdk:"description" json:"description,computed"`
	Name          types.String                                                         `tfsdk:"name" json:"name,computed"`
	Region        types.String                                                         `tfsdk:"region" json:"region,computed"`
	Size          types.Int64                                                          `tfsdk:"size" json:"size,computed"`
	Status        types.String                                                         `tfsdk:"status" json:"status,computed"`
	TaskID        types.String                                                         `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt     timetypes.RFC3339                                                    `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	VolumeID      types.String                                                         `tfsdk:"volume_id" json:"volume_id,computed"`
	Tags          customfield.NestedObjectList[CloudVolumeSnapshotTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	FindOneBy     *CloudVolumeSnapshotFindOneByDataSourceModel                         `tfsdk:"find_one_by"`
}

func (m *CloudVolumeSnapshotDataSourceModel) toReadParams(_ context.Context) (params cloud.VolumeSnapshotGetParams, diags diag.Diagnostics) {
	params = cloud.VolumeSnapshotGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

func (m *CloudVolumeSnapshotDataSourceModel) toListParams(_ context.Context) (params cloud.VolumeSnapshotListParams, diags diag.Diagnostics) {
	params = cloud.VolumeSnapshotListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.FindOneBy.InstanceID.IsNull() {
		params.InstanceID = param.NewOpt(m.FindOneBy.InstanceID.ValueString())
	}
	if !m.FindOneBy.LifecyclePolicyID.IsNull() {
		params.LifecyclePolicyID = param.NewOpt(m.FindOneBy.LifecyclePolicyID.ValueInt64())
	}
	if !m.FindOneBy.ScheduleID.IsNull() {
		params.ScheduleID = param.NewOpt(m.FindOneBy.ScheduleID.ValueString())
	}
	if !m.FindOneBy.VolumeID.IsNull() {
		params.VolumeID = param.NewOpt(m.FindOneBy.VolumeID.ValueString())
	}

	return
}

type CloudVolumeSnapshotTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudVolumeSnapshotFindOneByDataSourceModel struct {
	InstanceID        types.String `tfsdk:"instance_id" query:"instance_id,optional"`
	LifecyclePolicyID types.Int64  `tfsdk:"lifecycle_policy_id" query:"lifecycle_policy_id,optional"`
	ScheduleID        types.String `tfsdk:"schedule_id" query:"schedule_id,optional"`
	VolumeID          types.String `tfsdk:"volume_id" query:"volume_id,optional"`
}
