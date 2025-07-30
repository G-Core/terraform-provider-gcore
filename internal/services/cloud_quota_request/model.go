// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_quota_request

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
)

type CloudQuotaRequestModel struct {
	RequestID       types.String                           `tfsdk:"request_id" path:"request_id,optional"`
	Description     types.String                           `tfsdk:"description" json:"description,required"`
	RequestedLimits *CloudQuotaRequestRequestedLimitsModel `tfsdk:"requested_limits" json:"requested_limits,required"`
	ClientID        types.Int64                            `tfsdk:"client_id" json:"client_id,optional"`
	CreatedAt       timetypes.RFC3339                      `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	ID              types.Int64                            `tfsdk:"id" json:"id,computed"`
	Status          types.String                           `tfsdk:"status" json:"status,computed"`
	UpdatedAt       timetypes.RFC3339                      `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

func (m CloudQuotaRequestModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudQuotaRequestModel) MarshalJSONForUpdate(state CloudQuotaRequestModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type CloudQuotaRequestRequestedLimitsModel struct {
	GlobalLimits   *CloudQuotaRequestRequestedLimitsGlobalLimitsModel      `tfsdk:"global_limits" json:"global_limits,optional"`
	RegionalLimits *[]*CloudQuotaRequestRequestedLimitsRegionalLimitsModel `tfsdk:"regional_limits" json:"regional_limits,optional"`
}

type CloudQuotaRequestRequestedLimitsGlobalLimitsModel struct {
	InferenceCPUMillicoreCountLimit types.Int64 `tfsdk:"inference_cpu_millicore_count_limit" json:"inference_cpu_millicore_count_limit,optional"`
	InferenceGPUA100CountLimit      types.Int64 `tfsdk:"inference_gpu_a100_count_limit" json:"inference_gpu_a100_count_limit,optional"`
	InferenceGPUH100CountLimit      types.Int64 `tfsdk:"inference_gpu_h100_count_limit" json:"inference_gpu_h100_count_limit,optional"`
	InferenceGPUL40sCountLimit      types.Int64 `tfsdk:"inference_gpu_l40s_count_limit" json:"inference_gpu_l40s_count_limit,optional"`
	InferenceInstanceCountLimit     types.Int64 `tfsdk:"inference_instance_count_limit" json:"inference_instance_count_limit,optional"`
	KeypairCountLimit               types.Int64 `tfsdk:"keypair_count_limit" json:"keypair_count_limit,optional"`
	ProjectCountLimit               types.Int64 `tfsdk:"project_count_limit" json:"project_count_limit,optional"`
}

type CloudQuotaRequestRequestedLimitsRegionalLimitsModel struct {
	BaremetalBasicCountLimit          types.Int64 `tfsdk:"baremetal_basic_count_limit" json:"baremetal_basic_count_limit,optional"`
	BaremetalGPUA100CountLimit        types.Int64 `tfsdk:"baremetal_gpu_a100_count_limit" json:"baremetal_gpu_a100_count_limit,optional"`
	BaremetalGPUCountLimit            types.Int64 `tfsdk:"baremetal_gpu_count_limit" json:"baremetal_gpu_count_limit,optional"`
	BaremetalGPUH100CountLimit        types.Int64 `tfsdk:"baremetal_gpu_h100_count_limit" json:"baremetal_gpu_h100_count_limit,optional"`
	BaremetalGPUH200CountLimit        types.Int64 `tfsdk:"baremetal_gpu_h200_count_limit" json:"baremetal_gpu_h200_count_limit,optional"`
	BaremetalGPUL40sCountLimit        types.Int64 `tfsdk:"baremetal_gpu_l40s_count_limit" json:"baremetal_gpu_l40s_count_limit,optional"`
	BaremetalHfCountLimit             types.Int64 `tfsdk:"baremetal_hf_count_limit" json:"baremetal_hf_count_limit,optional"`
	BaremetalInfrastructureCountLimit types.Int64 `tfsdk:"baremetal_infrastructure_count_limit" json:"baremetal_infrastructure_count_limit,optional"`
	BaremetalNetworkCountLimit        types.Int64 `tfsdk:"baremetal_network_count_limit" json:"baremetal_network_count_limit,optional"`
	BaremetalStorageCountLimit        types.Int64 `tfsdk:"baremetal_storage_count_limit" json:"baremetal_storage_count_limit,optional"`
	CaasContainerCountLimit           types.Int64 `tfsdk:"caas_container_count_limit" json:"caas_container_count_limit,optional"`
	CaasCPUCountLimit                 types.Int64 `tfsdk:"caas_cpu_count_limit" json:"caas_cpu_count_limit,optional"`
	CaasGPUCountLimit                 types.Int64 `tfsdk:"caas_gpu_count_limit" json:"caas_gpu_count_limit,optional"`
	CaasRamSizeLimit                  types.Int64 `tfsdk:"caas_ram_size_limit" json:"caas_ram_size_limit,optional"`
	ClusterCountLimit                 types.Int64 `tfsdk:"cluster_count_limit" json:"cluster_count_limit,optional"`
	CPUCountLimit                     types.Int64 `tfsdk:"cpu_count_limit" json:"cpu_count_limit,optional"`
	DbaasPostgresClusterCountLimit    types.Int64 `tfsdk:"dbaas_postgres_cluster_count_limit" json:"dbaas_postgres_cluster_count_limit,optional"`
	ExternalIPCountLimit              types.Int64 `tfsdk:"external_ip_count_limit" json:"external_ip_count_limit,optional"`
	FaasCPUCountLimit                 types.Int64 `tfsdk:"faas_cpu_count_limit" json:"faas_cpu_count_limit,optional"`
	FaasFunctionCountLimit            types.Int64 `tfsdk:"faas_function_count_limit" json:"faas_function_count_limit,optional"`
	FaasNamespaceCountLimit           types.Int64 `tfsdk:"faas_namespace_count_limit" json:"faas_namespace_count_limit,optional"`
	FaasRamSizeLimit                  types.Int64 `tfsdk:"faas_ram_size_limit" json:"faas_ram_size_limit,optional"`
	FirewallCountLimit                types.Int64 `tfsdk:"firewall_count_limit" json:"firewall_count_limit,optional"`
	FloatingCountLimit                types.Int64 `tfsdk:"floating_count_limit" json:"floating_count_limit,optional"`
	GPUCountLimit                     types.Int64 `tfsdk:"gpu_count_limit" json:"gpu_count_limit,optional"`
	GPUVirtualA100CountLimit          types.Int64 `tfsdk:"gpu_virtual_a100_count_limit" json:"gpu_virtual_a100_count_limit,optional"`
	GPUVirtualH100CountLimit          types.Int64 `tfsdk:"gpu_virtual_h100_count_limit" json:"gpu_virtual_h100_count_limit,optional"`
	GPUVirtualH200CountLimit          types.Int64 `tfsdk:"gpu_virtual_h200_count_limit" json:"gpu_virtual_h200_count_limit,optional"`
	GPUVirtualL40sCountLimit          types.Int64 `tfsdk:"gpu_virtual_l40s_count_limit" json:"gpu_virtual_l40s_count_limit,optional"`
	ImageCountLimit                   types.Int64 `tfsdk:"image_count_limit" json:"image_count_limit,optional"`
	ImageSizeLimit                    types.Int64 `tfsdk:"image_size_limit" json:"image_size_limit,optional"`
	IpuCountLimit                     types.Int64 `tfsdk:"ipu_count_limit" json:"ipu_count_limit,optional"`
	LaasTopicCountLimit               types.Int64 `tfsdk:"laas_topic_count_limit" json:"laas_topic_count_limit,optional"`
	LoadbalancerCountLimit            types.Int64 `tfsdk:"loadbalancer_count_limit" json:"loadbalancer_count_limit,optional"`
	NetworkCountLimit                 types.Int64 `tfsdk:"network_count_limit" json:"network_count_limit,optional"`
	RamLimit                          types.Int64 `tfsdk:"ram_limit" json:"ram_limit,optional"`
	RegionID                          types.Int64 `tfsdk:"region_id" json:"region_id,optional"`
	RegistryCountLimit                types.Int64 `tfsdk:"registry_count_limit" json:"registry_count_limit,optional"`
	RegistryStorageLimit              types.Int64 `tfsdk:"registry_storage_limit" json:"registry_storage_limit,optional"`
	RouterCountLimit                  types.Int64 `tfsdk:"router_count_limit" json:"router_count_limit,optional"`
	SecretCountLimit                  types.Int64 `tfsdk:"secret_count_limit" json:"secret_count_limit,optional"`
	ServergroupCountLimit             types.Int64 `tfsdk:"servergroup_count_limit" json:"servergroup_count_limit,optional"`
	SfsCountLimit                     types.Int64 `tfsdk:"sfs_count_limit" json:"sfs_count_limit,optional"`
	SfsSizeLimit                      types.Int64 `tfsdk:"sfs_size_limit" json:"sfs_size_limit,optional"`
	SharedVmCountLimit                types.Int64 `tfsdk:"shared_vm_count_limit" json:"shared_vm_count_limit,optional"`
	SnapshotScheduleCountLimit        types.Int64 `tfsdk:"snapshot_schedule_count_limit" json:"snapshot_schedule_count_limit,optional"`
	SubnetCountLimit                  types.Int64 `tfsdk:"subnet_count_limit" json:"subnet_count_limit,optional"`
	VmCountLimit                      types.Int64 `tfsdk:"vm_count_limit" json:"vm_count_limit,optional"`
	VolumeCountLimit                  types.Int64 `tfsdk:"volume_count_limit" json:"volume_count_limit,optional"`
	VolumeSizeLimit                   types.Int64 `tfsdk:"volume_size_limit" json:"volume_size_limit,optional"`
	VolumeSnapshotsCountLimit         types.Int64 `tfsdk:"volume_snapshots_count_limit" json:"volume_snapshots_count_limit,optional"`
	VolumeSnapshotsSizeLimit          types.Int64 `tfsdk:"volume_snapshots_size_limit" json:"volume_snapshots_size_limit,optional"`
}
