// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_usage_report

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudUsageReportDataSourceModel struct {
	TimeFrom      timetypes.RFC3339                                                      `tfsdk:"time_from" json:"time_from,required" format:"date-time"`
	TimeTo        timetypes.RFC3339                                                      `tfsdk:"time_to" json:"time_to,required" format:"date-time"`
	Projects      *[]types.Int64                                                         `tfsdk:"projects" json:"projects,optional"`
	Regions       *[]types.Int64                                                         `tfsdk:"regions" json:"regions,optional"`
	Types         *[]types.String                                                        `tfsdk:"types" json:"types,optional"`
	SchemaFilter  *CloudUsageReportSchemaFilterDataSourceModel                           `tfsdk:"schema_filter" json:"schema_filter,optional"`
	Tags          *CloudUsageReportTagsDataSourceModel                                   `tfsdk:"tags" json:"tags,optional"`
	EnableLastDay types.Bool                                                             `tfsdk:"enable_last_day" json:"enable_last_day,computed_optional"`
	Limit         types.Int64                                                            `tfsdk:"limit" json:"limit,computed_optional"`
	Offset        types.Int64                                                            `tfsdk:"offset" json:"offset,computed_optional"`
	Sorting       customfield.NestedObjectList[CloudUsageReportSortingDataSourceModel]   `tfsdk:"sorting" json:"sorting,computed_optional"`
	Count         types.Int64                                                            `tfsdk:"count" json:"count,computed"`
	Resources     customfield.NestedObjectList[CloudUsageReportResourcesDataSourceModel] `tfsdk:"resources" json:"resources,computed"`
	Totals        customfield.NestedObjectList[CloudUsageReportTotalsDataSourceModel]    `tfsdk:"totals" json:"totals,computed"`
}

func (m *CloudUsageReportDataSourceModel) toReadParams(_ context.Context) (params cloud.UsageReportGetParams, diags diag.Diagnostics) {
	params = cloud.UsageReportGetParams{}

	return
}

type CloudUsageReportSchemaFilterDataSourceModel struct {
	Field  types.String    `tfsdk:"field" json:"field,required"`
	Type   types.String    `tfsdk:"type" json:"type,required"`
	Values *[]types.String `tfsdk:"values" json:"values,required"`
}

type CloudUsageReportTagsDataSourceModel struct {
	Conditions    *[]*CloudUsageReportTagsConditionsDataSourceModel `tfsdk:"conditions" json:"conditions,required"`
	ConditionType types.String                                      `tfsdk:"condition_type" json:"condition_type,computed_optional"`
}

type CloudUsageReportTagsConditionsDataSourceModel struct {
	Key    types.String `tfsdk:"key" json:"key,optional"`
	Strict types.Bool   `tfsdk:"strict" json:"strict,optional"`
	Value  types.String `tfsdk:"value" json:"value,optional"`
}

type CloudUsageReportSortingDataSourceModel struct {
	BillingValue types.String `tfsdk:"billing_value" json:"billing_value,optional"`
	FirstSeen    types.String `tfsdk:"first_seen" json:"first_seen,optional"`
	LastName     types.String `tfsdk:"last_name" json:"last_name,optional"`
	LastSeen     types.String `tfsdk:"last_seen" json:"last_seen,optional"`
	Project      types.String `tfsdk:"project" json:"project,optional"`
	Region       types.String `tfsdk:"region" json:"region,optional"`
	Type         types.String `tfsdk:"type" json:"type,optional"`
}

type CloudUsageReportResourcesDataSourceModel struct {
	BillingMetricName types.String                                    `tfsdk:"billing_metric_name" json:"billing_metric_name,computed"`
	BillingValue      types.Float64                                   `tfsdk:"billing_value" json:"billing_value,computed"`
	BillingValueUnit  types.String                                    `tfsdk:"billing_value_unit" json:"billing_value_unit,computed"`
	FirstSeen         timetypes.RFC3339                               `tfsdk:"first_seen" json:"first_seen,computed" format:"date-time"`
	Flavor            types.String                                    `tfsdk:"flavor" json:"flavor,computed"`
	LastName          types.String                                    `tfsdk:"last_name" json:"last_name,computed"`
	LastSeen          timetypes.RFC3339                               `tfsdk:"last_seen" json:"last_seen,computed" format:"date-time"`
	ProjectID         types.Int64                                     `tfsdk:"project_id" json:"project_id,computed"`
	Region            types.Int64                                     `tfsdk:"region" json:"region,computed"`
	RegionID          types.Int64                                     `tfsdk:"region_id" json:"region_id,computed"`
	Tags              customfield.List[customfield.Map[types.String]] `tfsdk:"tags" json:"tags,computed"`
	Type              types.String                                    `tfsdk:"type" json:"type,computed"`
	Uuid              types.String                                    `tfsdk:"uuid" json:"uuid,computed"`
	LastSize          types.Int64                                     `tfsdk:"last_size" json:"last_size,computed"`
	ScheduleID        types.String                                    `tfsdk:"schedule_id" json:"schedule_id,computed"`
	SourceVolumeUuid  types.String                                    `tfsdk:"source_volume_uuid" json:"source_volume_uuid,computed"`
	InstanceType      types.String                                    `tfsdk:"instance_type" json:"instance_type,computed"`
	PortID            types.String                                    `tfsdk:"port_id" json:"port_id,computed"`
	SizeUnit          types.String                                    `tfsdk:"size_unit" json:"size_unit,computed"`
	VmID              types.String                                    `tfsdk:"vm_id" json:"vm_id,computed"`
	InstanceName      types.String                                    `tfsdk:"instance_name" json:"instance_name,computed"`
	AttachedToVm      types.String                                    `tfsdk:"attached_to_vm" json:"attached_to_vm,computed"`
	IPAddress         types.String                                    `tfsdk:"ip_address" json:"ip_address,computed"`
	NetworkID         types.String                                    `tfsdk:"network_id" json:"network_id,computed"`
	SubnetID          types.String                                    `tfsdk:"subnet_id" json:"subnet_id,computed"`
	FileShareType     types.String                                    `tfsdk:"file_share_type" json:"file_share_type,computed"`
	VolumeType        types.String                                    `tfsdk:"volume_type" json:"volume_type,computed"`
}

type CloudUsageReportTotalsDataSourceModel struct {
	BillingMetricName types.String  `tfsdk:"billing_metric_name" json:"billing_metric_name,computed"`
	BillingValue      types.Float64 `tfsdk:"billing_value" json:"billing_value,computed"`
	BillingValueUnit  types.String  `tfsdk:"billing_value_unit" json:"billing_value_unit,computed"`
	Flavor            types.String  `tfsdk:"flavor" json:"flavor,computed"`
	Region            types.Int64   `tfsdk:"region" json:"region,computed"`
	RegionID          types.Int64   `tfsdk:"region_id" json:"region_id,computed"`
	Type              types.String  `tfsdk:"type" json:"type,computed"`
	InstanceType      types.String  `tfsdk:"instance_type" json:"instance_type,computed"`
	FileShareType     types.String  `tfsdk:"file_share_type" json:"file_share_type,computed"`
	VolumeType        types.String  `tfsdk:"volume_type" json:"volume_type,computed"`
}
