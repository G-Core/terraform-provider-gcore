// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudVolumeDataSourceModel struct {
	ProjectID           types.Int64                                                         `tfsdk:"project_id" path:"project_id,required"`
	RegionID            types.Int64                                                         `tfsdk:"region_id" path:"region_id,required"`
	VolumeID            types.String                                                        `tfsdk:"volume_id" path:"volume_id,required"`
	Bootable            types.Bool                                                          `tfsdk:"bootable" json:"bootable,computed"`
	CreatedAt           timetypes.RFC3339                                                   `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                        `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	ID                  types.String                                                        `tfsdk:"id" json:"id,computed"`
	IsRootVolume        types.Bool                                                          `tfsdk:"is_root_volume" json:"is_root_volume,computed"`
	Name                types.String                                                        `tfsdk:"name" json:"name,computed"`
	Region              types.String                                                        `tfsdk:"region" json:"region,computed"`
	Size                types.Int64                                                         `tfsdk:"size" json:"size,computed"`
	Status              types.String                                                        `tfsdk:"status" json:"status,computed"`
	TaskID              types.String                                                        `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt           timetypes.RFC3339                                                   `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	VolumeType          types.String                                                        `tfsdk:"volume_type" json:"volume_type,computed"`
	SnapshotIDs         customfield.List[types.String]                                      `tfsdk:"snapshot_ids" json:"snapshot_ids,computed"`
	VolumeImageMetadata customfield.Map[types.String]                                       `tfsdk:"volume_image_metadata" json:"volume_image_metadata,computed"`
	Attachments         customfield.NestedObjectList[CloudVolumeAttachmentsDataSourceModel] `tfsdk:"attachments" json:"attachments,computed"`
	LimiterStats        customfield.NestedObject[CloudVolumeLimiterStatsDataSourceModel]    `tfsdk:"limiter_stats" json:"limiter_stats,computed"`
	Tags                customfield.NestedObjectList[CloudVolumeTagsDataSourceModel]        `tfsdk:"tags" json:"tags,computed"`
}

func (m *CloudVolumeDataSourceModel) toReadParams(_ context.Context) (params cloud.VolumeGetParams, diags diag.Diagnostics) {
	params = cloud.VolumeGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudVolumeAttachmentsDataSourceModel struct {
	AttachmentID types.String      `tfsdk:"attachment_id" json:"attachment_id,computed"`
	VolumeID     types.String      `tfsdk:"volume_id" json:"volume_id,computed"`
	AttachedAt   timetypes.RFC3339 `tfsdk:"attached_at" json:"attached_at,computed" format:"date-time"`
	Device       types.String      `tfsdk:"device" json:"device,computed"`
	FlavorID     types.String      `tfsdk:"flavor_id" json:"flavor_id,computed"`
	InstanceName types.String      `tfsdk:"instance_name" json:"instance_name,computed"`
	ServerID     types.String      `tfsdk:"server_id" json:"server_id,computed"`
}

type CloudVolumeLimiterStatsDataSourceModel struct {
	IopsBaseLimit  types.Int64 `tfsdk:"iops_base_limit" json:"iops_base_limit,computed"`
	IopsBurstLimit types.Int64 `tfsdk:"iops_burst_limit" json:"iops_burst_limit,computed"`
	MBpsBaseLimit  types.Int64 `tfsdk:"m_bps_base_limit" json:"MBps_base_limit,computed"`
	MBpsBurstLimit types.Int64 `tfsdk:"m_bps_burst_limit" json:"MBps_burst_limit,computed"`
}

type CloudVolumeTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}
