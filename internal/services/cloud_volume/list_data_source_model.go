// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_volume

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudVolumesResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudVolumesItemsDataSourceModel] `json:"results,computed"`
}

type CloudVolumesDataSourceModel struct {
	ProjectID      types.Int64                                                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID       types.Int64                                                    `tfsdk:"region_id" path:"region_id,optional"`
	Bootable       types.Bool                                                     `tfsdk:"bootable" query:"bootable,optional"`
	ClusterID      types.String                                                   `tfsdk:"cluster_id" query:"cluster_id,optional"`
	HasAttachments types.Bool                                                     `tfsdk:"has_attachments" query:"has_attachments,optional"`
	IDPart         types.String                                                   `tfsdk:"id_part" query:"id_part,optional"`
	InstanceID     types.String                                                   `tfsdk:"instance_id" query:"instance_id,optional"`
	NamePart       types.String                                                   `tfsdk:"name_part" query:"name_part,optional"`
	TagKeyValue    types.String                                                   `tfsdk:"tag_key_value" query:"tag_key_value,optional"`
	TagKey         *[]types.String                                                `tfsdk:"tag_key" query:"tag_key,optional"`
	Limit          types.Int64                                                    `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems       types.Int64                                                    `tfsdk:"max_items"`
	Items          customfield.NestedObjectList[CloudVolumesItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudVolumesDataSourceModel) toListParams(_ context.Context) (params cloud.VolumeListParams, diags diag.Diagnostics) {
	mTagKey := []string{}
	if m.TagKey != nil {
		for _, item := range *m.TagKey {
			mTagKey = append(mTagKey, item.ValueString())
		}
	}

	params = cloud.VolumeListParams{
		TagKey: mTagKey,
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.Bootable.IsNull() {
		params.Bootable = param.NewOpt(m.Bootable.ValueBool())
	}
	if !m.ClusterID.IsNull() {
		params.ClusterID = param.NewOpt(m.ClusterID.ValueString())
	}
	if !m.HasAttachments.IsNull() {
		params.HasAttachments = param.NewOpt(m.HasAttachments.ValueBool())
	}
	if !m.IDPart.IsNull() {
		params.IDPart = param.NewOpt(m.IDPart.ValueString())
	}
	if !m.InstanceID.IsNull() {
		params.InstanceID = param.NewOpt(m.InstanceID.ValueString())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.NamePart.IsNull() {
		params.NamePart = param.NewOpt(m.NamePart.ValueString())
	}
	if !m.TagKeyValue.IsNull() {
		params.TagKeyValue = param.NewOpt(m.TagKeyValue.ValueString())
	}

	return
}

type CloudVolumesItemsDataSourceModel struct {
	ID                  types.String                                                         `tfsdk:"id" json:"id,computed"`
	Bootable            types.Bool                                                           `tfsdk:"bootable" json:"bootable,computed"`
	CreatedAt           timetypes.RFC3339                                                    `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	IsRootVolume        types.Bool                                                           `tfsdk:"is_root_volume" json:"is_root_volume,computed"`
	Name                types.String                                                         `tfsdk:"name" json:"name,computed"`
	ProjectID           types.Int64                                                          `tfsdk:"project_id" json:"project_id,computed"`
	Region              types.String                                                         `tfsdk:"region" json:"region,computed"`
	RegionID            types.Int64                                                          `tfsdk:"region_id" json:"region_id,computed"`
	Size                types.Int64                                                          `tfsdk:"size" json:"size,computed"`
	Status              types.String                                                         `tfsdk:"status" json:"status,computed"`
	Tags                customfield.NestedObjectList[CloudVolumesTagsDataSourceModel]        `tfsdk:"tags" json:"tags,computed"`
	VolumeType          types.String                                                         `tfsdk:"volume_type" json:"volume_type,computed"`
	Attachments         customfield.NestedObjectList[CloudVolumesAttachmentsDataSourceModel] `tfsdk:"attachments" json:"attachments,computed"`
	CreatorTaskID       types.String                                                         `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	LimiterStats        customfield.NestedObject[CloudVolumesLimiterStatsDataSourceModel]    `tfsdk:"limiter_stats" json:"limiter_stats,computed"`
	SnapshotIDs         customfield.List[types.String]                                       `tfsdk:"snapshot_ids" json:"snapshot_ids,computed"`
	TaskID              types.String                                                         `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt           timetypes.RFC3339                                                    `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	VolumeImageMetadata customfield.Map[types.String]                                        `tfsdk:"volume_image_metadata" json:"volume_image_metadata,computed"`
}

type CloudVolumesTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudVolumesAttachmentsDataSourceModel struct {
	AttachmentID types.String      `tfsdk:"attachment_id" json:"attachment_id,computed"`
	VolumeID     types.String      `tfsdk:"volume_id" json:"volume_id,computed"`
	AttachedAt   timetypes.RFC3339 `tfsdk:"attached_at" json:"attached_at,computed" format:"date-time"`
	Device       types.String      `tfsdk:"device" json:"device,computed"`
	FlavorID     types.String      `tfsdk:"flavor_id" json:"flavor_id,computed"`
	InstanceName types.String      `tfsdk:"instance_name" json:"instance_name,computed"`
	ServerID     types.String      `tfsdk:"server_id" json:"server_id,computed"`
}

type CloudVolumesLimiterStatsDataSourceModel struct {
	IopsBaseLimit  types.Int64 `tfsdk:"iops_base_limit" json:"iops_base_limit,computed"`
	IopsBurstLimit types.Int64 `tfsdk:"iops_burst_limit" json:"iops_burst_limit,computed"`
	MBpsBaseLimit  types.Int64 `tfsdk:"m_bps_base_limit" json:"MBps_base_limit,computed"`
	MBpsBurstLimit types.Int64 `tfsdk:"m_bps_burst_limit" json:"MBps_burst_limit,computed"`
}
