// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_region

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudRegionDataSourceModel struct {
	RegionID             types.Int64                                                     `tfsdk:"region_id" path:"region_id,optional"`
	ShowVolumeTypes      types.Bool                                                      `tfsdk:"show_volume_types" query:"show_volume_types,computed_optional"`
	AccessLevel          types.String                                                    `tfsdk:"access_level" json:"access_level,computed"`
	Country              types.String                                                    `tfsdk:"country" json:"country,computed"`
	CreatedAt            timetypes.RFC3339                                               `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatedOn            timetypes.RFC3339                                               `tfsdk:"created_on" json:"created_on,computed" format:"date-time"`
	DDOSEndpointID       types.Int64                                                     `tfsdk:"ddos_endpoint_id" json:"ddos_endpoint_id,computed"`
	DisplayName          types.String                                                    `tfsdk:"display_name" json:"display_name,computed"`
	EndpointType         types.String                                                    `tfsdk:"endpoint_type" json:"endpoint_type,computed"`
	ExternalNetworkID    types.String                                                    `tfsdk:"external_network_id" json:"external_network_id,computed"`
	HasAI                types.Bool                                                      `tfsdk:"has_ai" json:"has_ai,computed"`
	HasAIGPU             types.Bool                                                      `tfsdk:"has_ai_gpu" json:"has_ai_gpu,computed"`
	HasBaremetal         types.Bool                                                      `tfsdk:"has_baremetal" json:"has_baremetal,computed"`
	HasBasicVm           types.Bool                                                      `tfsdk:"has_basic_vm" json:"has_basic_vm,computed"`
	HasDbaas             types.Bool                                                      `tfsdk:"has_dbaas" json:"has_dbaas,computed"`
	HasDDOS              types.Bool                                                      `tfsdk:"has_ddos" json:"has_ddos,computed"`
	HasK8S               types.Bool                                                      `tfsdk:"has_k8s" json:"has_k8s,computed"`
	HasKvm               types.Bool                                                      `tfsdk:"has_kvm" json:"has_kvm,computed"`
	HasSfs               types.Bool                                                      `tfsdk:"has_sfs" json:"has_sfs,computed"`
	ID                   types.Int64                                                     `tfsdk:"id" json:"id,computed"`
	KeystoneID           types.Int64                                                     `tfsdk:"keystone_id" json:"keystone_id,computed"`
	KeystoneName         types.String                                                    `tfsdk:"keystone_name" json:"keystone_name,computed"`
	MetricsDatabaseID    types.Int64                                                     `tfsdk:"metrics_database_id" json:"metrics_database_id,computed"`
	Slug                 types.String                                                    `tfsdk:"slug" json:"slug,computed"`
	State                types.String                                                    `tfsdk:"state" json:"state,computed"`
	TaskID               types.String                                                    `tfsdk:"task_id" json:"task_id,computed"`
	VlanPhysicalNetwork  types.String                                                    `tfsdk:"vlan_physical_network" json:"vlan_physical_network,computed"`
	Zone                 types.String                                                    `tfsdk:"zone" json:"zone,computed"`
	AvailableVolumeTypes customfield.List[types.String]                                  `tfsdk:"available_volume_types" json:"available_volume_types,computed"`
	FileShareTypes       customfield.List[types.String]                                  `tfsdk:"file_share_types" json:"file_share_types,computed"`
	Coordinates          customfield.NestedObject[CloudRegionCoordinatesDataSourceModel] `tfsdk:"coordinates" json:"coordinates,computed"`
}

func (m *CloudRegionDataSourceModel) toReadParams(_ context.Context) (params cloud.RegionGetParams, diags diag.Diagnostics) {
	params = cloud.RegionGetParams{}

	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.ShowVolumeTypes.IsNull() {
		params.ShowVolumeTypes = param.NewOpt(m.ShowVolumeTypes.ValueBool())
	}

	return
}

type CloudRegionCoordinatesDataSourceModel struct {
	Latitude  types.String `tfsdk:"latitude" json:"latitude,computed"`
	Longitude types.String `tfsdk:"longitude" json:"longitude,computed"`
}
