// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_region

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudRegionsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudRegionsItemsDataSourceModel] `json:"results,computed"`
}

type CloudRegionsDataSourceModel struct {
	Product         types.String                                                   `tfsdk:"product" query:"product,optional"`
	Limit           types.Int64                                                    `tfsdk:"limit" query:"limit,computed_optional"`
	OrderBy         types.String                                                   `tfsdk:"order_by" query:"order_by,computed_optional"`
	ShowVolumeTypes types.Bool                                                     `tfsdk:"show_volume_types" query:"show_volume_types,computed_optional"`
	MaxItems        types.Int64                                                    `tfsdk:"max_items"`
	Items           customfield.NestedObjectList[CloudRegionsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudRegionsDataSourceModel) toListParams(_ context.Context) (params cloud.RegionListParams, diags diag.Diagnostics) {
	params = cloud.RegionListParams{}

	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = cloud.RegionListParamsOrderBy(m.OrderBy.ValueString())
	}
	if !m.Product.IsNull() {
		params.Product = cloud.RegionListParamsProduct(m.Product.ValueString())
	}
	if !m.ShowVolumeTypes.IsNull() {
		params.ShowVolumeTypes = param.NewOpt(m.ShowVolumeTypes.ValueBool())
	}

	return
}

type CloudRegionsItemsDataSourceModel struct {
	ID                   types.Int64                                                      `tfsdk:"id" json:"id,computed"`
	AccessLevel          types.String                                                     `tfsdk:"access_level" json:"access_level,computed"`
	AIServiceEndpointID  types.Int64                                                      `tfsdk:"ai_service_endpoint_id" json:"ai_service_endpoint_id,computed"`
	AvailableVolumeTypes customfield.List[types.String]                                   `tfsdk:"available_volume_types" json:"available_volume_types,computed"`
	Coordinates          customfield.NestedObject[CloudRegionsCoordinatesDataSourceModel] `tfsdk:"coordinates" json:"coordinates,computed"`
	Country              types.String                                                     `tfsdk:"country" json:"country,computed"`
	CreatedAt            timetypes.RFC3339                                                `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatedOn            timetypes.RFC3339                                                `tfsdk:"created_on" json:"created_on,computed" format:"date-time"`
	DisplayName          types.String                                                     `tfsdk:"display_name" json:"display_name,computed"`
	EndpointType         types.String                                                     `tfsdk:"endpoint_type" json:"endpoint_type,computed"`
	ExternalNetworkID    types.String                                                     `tfsdk:"external_network_id" json:"external_network_id,computed"`
	FileShareTypes       customfield.List[types.String]                                   `tfsdk:"file_share_types" json:"file_share_types,computed"`
	HasAI                types.Bool                                                       `tfsdk:"has_ai" json:"has_ai,computed"`
	HasAIGPU             types.Bool                                                       `tfsdk:"has_ai_gpu" json:"has_ai_gpu,computed"`
	HasBaremetal         types.Bool                                                       `tfsdk:"has_baremetal" json:"has_baremetal,computed"`
	HasBasicVm           types.Bool                                                       `tfsdk:"has_basic_vm" json:"has_basic_vm,computed"`
	HasDbaas             types.Bool                                                       `tfsdk:"has_dbaas" json:"has_dbaas,computed"`
	HasDDOS              types.Bool                                                       `tfsdk:"has_ddos" json:"has_ddos,computed"`
	HasK8s               types.Bool                                                       `tfsdk:"has_k8s" json:"has_k8s,computed"`
	HasKvm               types.Bool                                                       `tfsdk:"has_kvm" json:"has_kvm,computed"`
	HasSfs               types.Bool                                                       `tfsdk:"has_sfs" json:"has_sfs,computed"`
	KeystoneID           types.Int64                                                      `tfsdk:"keystone_id" json:"keystone_id,computed"`
	KeystoneName         types.String                                                     `tfsdk:"keystone_name" json:"keystone_name,computed"`
	MetricsDatabaseID    types.Int64                                                      `tfsdk:"metrics_database_id" json:"metrics_database_id,computed"`
	State                types.String                                                     `tfsdk:"state" json:"state,computed"`
	TaskID               types.String                                                     `tfsdk:"task_id" json:"task_id,computed"`
	VlanPhysicalNetwork  types.String                                                     `tfsdk:"vlan_physical_network" json:"vlan_physical_network,computed"`
	Zone                 types.String                                                     `tfsdk:"zone" json:"zone,computed"`
	DDOSEndpointID       types.Int64                                                      `tfsdk:"ddos_endpoint_id" json:"ddos_endpoint_id,computed"`
}

type CloudRegionsCoordinatesDataSourceModel struct {
	Latitude  types.Float64 `tfsdk:"latitude" json:"latitude,computed"`
	Longitude types.Float64 `tfsdk:"longitude" json:"longitude,computed"`
}
