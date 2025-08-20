// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudVolumeModel struct {
	ProjectID            types.Int64                                               `tfsdk:"project_id" path:"project_id,optional"`
	RegionID             types.Int64                                               `tfsdk:"region_id" path:"region_id,optional"`
	VolumeID             types.String                                              `tfsdk:"volume_id" path:"volume_id,optional"`
	Source               types.String                                              `tfsdk:"source" json:"source,required,no_refresh"`
	AttachmentTag        types.String                                              `tfsdk:"attachment_tag" json:"attachment_tag,optional,no_refresh"`
	ImageID              types.String                                              `tfsdk:"image_id" json:"image_id,optional,no_refresh"`
	InstanceIDToAttachTo types.String                                              `tfsdk:"instance_id_to_attach_to" json:"instance_id_to_attach_to,optional,no_refresh"`
	Size                 types.Int64                                               `tfsdk:"size" json:"size,optional"`
	SnapshotID           types.String                                              `tfsdk:"snapshot_id" json:"snapshot_id,optional,no_refresh"`
	TypeName             types.String                                              `tfsdk:"type_name" json:"type_name,optional,no_refresh"`
	LifecyclePolicyIDs   *[]types.Int64                                            `tfsdk:"lifecycle_policy_ids" json:"lifecycle_policy_ids,optional,no_refresh"`
	Name                 types.String                                              `tfsdk:"name" json:"name,required"`
	Tags                 *map[string]types.String                                  `tfsdk:"tags" json:"tags,optional,no_refresh"`
	Bootable             types.Bool                                                `tfsdk:"bootable" json:"bootable,computed"`
	CreatedAt            timetypes.RFC3339                                         `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID        types.String                                              `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	ID                   types.String                                              `tfsdk:"id" json:"id,computed"`
	IsRootVolume         types.Bool                                                `tfsdk:"is_root_volume" json:"is_root_volume,computed"`
	Region               types.String                                              `tfsdk:"region" json:"region,computed"`
	Status               types.String                                              `tfsdk:"status" json:"status,computed"`
	TaskID               types.String                                              `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt            timetypes.RFC3339                                         `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	VolumeType           types.String                                              `tfsdk:"volume_type" json:"volume_type,computed"`
	SnapshotIDs          customfield.List[types.String]                            `tfsdk:"snapshot_ids" json:"snapshot_ids,computed"`
	Tasks                customfield.List[types.String]                            `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
	VolumeImageMetadata  customfield.Map[types.String]                             `tfsdk:"volume_image_metadata" json:"volume_image_metadata,computed"`
	Attachments          customfield.NestedObjectList[CloudVolumeAttachmentsModel] `tfsdk:"attachments" json:"attachments,computed"`
	LimiterStats         customfield.NestedObject[CloudVolumeLimiterStatsModel]    `tfsdk:"limiter_stats" json:"limiter_stats,computed"`
}

func (m CloudVolumeModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudVolumeModel) MarshalJSONForUpdate(state CloudVolumeModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudVolumeAttachmentsModel struct {
	AttachmentID types.String      `tfsdk:"attachment_id" json:"attachment_id,computed"`
	VolumeID     types.String      `tfsdk:"volume_id" json:"volume_id,computed"`
	AttachedAt   timetypes.RFC3339 `tfsdk:"attached_at" json:"attached_at,computed" format:"date-time"`
	Device       types.String      `tfsdk:"device" json:"device,computed"`
	FlavorID     types.String      `tfsdk:"flavor_id" json:"flavor_id,computed"`
	InstanceName types.String      `tfsdk:"instance_name" json:"instance_name,computed"`
	ServerID     types.String      `tfsdk:"server_id" json:"server_id,computed"`
}

type CloudVolumeLimiterStatsModel struct {
	IopsBaseLimit  types.Int64 `tfsdk:"iops_base_limit" json:"iops_base_limit,computed"`
	IopsBurstLimit types.Int64 `tfsdk:"iops_burst_limit" json:"iops_burst_limit,computed"`
	MBpsBaseLimit  types.Int64 `tfsdk:"m_bps_base_limit" json:"MBps_base_limit,computed"`
	MBpsBurstLimit types.Int64 `tfsdk:"m_bps_burst_limit" json:"MBps_burst_limit,computed"`
}
