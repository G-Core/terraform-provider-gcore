// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_quota_request

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudQuotaRequestDataSourceModel struct {
	RequestID       types.Int64                                                               `tfsdk:"request_id" path:"request_id,required"`
	ClientID        types.Int64                                                               `tfsdk:"client_id" json:"client_id,computed"`
	CreatedAt       timetypes.RFC3339                                                         `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Description     types.String                                                              `tfsdk:"description" json:"description,computed"`
	ID              types.Int64                                                               `tfsdk:"id" json:"id,computed"`
	Status          types.String                                                              `tfsdk:"status" json:"status,computed"`
	UpdatedAt       timetypes.RFC3339                                                         `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	RequestedLimits customfield.NestedObject[CloudQuotaRequestRequestedLimitsDataSourceModel] `tfsdk:"requested_limits" json:"requested_limits,computed"`
}

type CloudQuotaRequestRequestedLimitsDataSourceModel struct {
	GlobalLimits   customfield.NestedObject[CloudQuotaRequestRequestedLimitsGlobalLimitsDataSourceModel]       `tfsdk:"global_limits" json:"global_limits,computed"`
	RegionalLimits customfield.NestedObjectList[CloudQuotaRequestRequestedLimitsRegionalLimitsDataSourceModel] `tfsdk:"regional_limits" json:"regional_limits,computed"`
}

type CloudQuotaRequestRequestedLimitsGlobalLimitsDataSourceModel struct {
	InferenceCPUMillicoreCountLimit types.Int64 `tfsdk:"inference_cpu_millicore_count_limit" json:"inference_cpu_millicore_count_limit,computed"`
	InferenceGPUA100CountLimit      types.Int64 `tfsdk:"inference_gpu_a100_count_limit" json:"inference_gpu_a100_count_limit,computed"`
	InferenceGPUH100CountLimit      types.Int64 `tfsdk:"inference_gpu_h100_count_limit" json:"inference_gpu_h100_count_limit,computed"`
	InferenceGPUL40sCountLimit      types.Int64 `tfsdk:"inference_gpu_l40s_count_limit" json:"inference_gpu_l40s_count_limit,computed"`
	InferenceInstanceCountLimit     types.Int64 `tfsdk:"inference_instance_count_limit" json:"inference_instance_count_limit,computed"`
	KeypairCountLimit               types.Int64 `tfsdk:"keypair_count_limit" json:"keypair_count_limit,computed"`
	ProjectCountLimit               types.Int64 `tfsdk:"project_count_limit" json:"project_count_limit,computed"`
}

type CloudQuotaRequestRequestedLimitsRegionalLimitsDataSourceModel struct {
	BaremetalBasicCountLimit          types.Int64 `tfsdk:"baremetal_basic_count_limit" json:"baremetal_basic_count_limit,computed"`
	BaremetalGPUA100CountLimit        types.Int64 `tfsdk:"baremetal_gpu_a100_count_limit" json:"baremetal_gpu_a100_count_limit,computed"`
	BaremetalGPUCountLimit            types.Int64 `tfsdk:"baremetal_gpu_count_limit" json:"baremetal_gpu_count_limit,computed"`
	BaremetalGPUH100CountLimit        types.Int64 `tfsdk:"baremetal_gpu_h100_count_limit" json:"baremetal_gpu_h100_count_limit,computed"`
	BaremetalGPUH200CountLimit        types.Int64 `tfsdk:"baremetal_gpu_h200_count_limit" json:"baremetal_gpu_h200_count_limit,computed"`
	BaremetalGPUL40sCountLimit        types.Int64 `tfsdk:"baremetal_gpu_l40s_count_limit" json:"baremetal_gpu_l40s_count_limit,computed"`
	BaremetalHfCountLimit             types.Int64 `tfsdk:"baremetal_hf_count_limit" json:"baremetal_hf_count_limit,computed"`
	BaremetalInfrastructureCountLimit types.Int64 `tfsdk:"baremetal_infrastructure_count_limit" json:"baremetal_infrastructure_count_limit,computed"`
	BaremetalNetworkCountLimit        types.Int64 `tfsdk:"baremetal_network_count_limit" json:"baremetal_network_count_limit,computed"`
	BaremetalStorageCountLimit        types.Int64 `tfsdk:"baremetal_storage_count_limit" json:"baremetal_storage_count_limit,computed"`
	CaasContainerCountLimit           types.Int64 `tfsdk:"caas_container_count_limit" json:"caas_container_count_limit,computed"`
	CaasCPUCountLimit                 types.Int64 `tfsdk:"caas_cpu_count_limit" json:"caas_cpu_count_limit,computed"`
	CaasGPUCountLimit                 types.Int64 `tfsdk:"caas_gpu_count_limit" json:"caas_gpu_count_limit,computed"`
	CaasRamSizeLimit                  types.Int64 `tfsdk:"caas_ram_size_limit" json:"caas_ram_size_limit,computed"`
	ClusterCountLimit                 types.Int64 `tfsdk:"cluster_count_limit" json:"cluster_count_limit,computed"`
	CPUCountLimit                     types.Int64 `tfsdk:"cpu_count_limit" json:"cpu_count_limit,computed"`
	DbaasPostgresClusterCountLimit    types.Int64 `tfsdk:"dbaas_postgres_cluster_count_limit" json:"dbaas_postgres_cluster_count_limit,computed"`
	ExternalIPCountLimit              types.Int64 `tfsdk:"external_ip_count_limit" json:"external_ip_count_limit,computed"`
	FaasCPUCountLimit                 types.Int64 `tfsdk:"faas_cpu_count_limit" json:"faas_cpu_count_limit,computed"`
	FaasFunctionCountLimit            types.Int64 `tfsdk:"faas_function_count_limit" json:"faas_function_count_limit,computed"`
	FaasNamespaceCountLimit           types.Int64 `tfsdk:"faas_namespace_count_limit" json:"faas_namespace_count_limit,computed"`
	FaasRamSizeLimit                  types.Int64 `tfsdk:"faas_ram_size_limit" json:"faas_ram_size_limit,computed"`
	FirewallCountLimit                types.Int64 `tfsdk:"firewall_count_limit" json:"firewall_count_limit,computed"`
	FloatingCountLimit                types.Int64 `tfsdk:"floating_count_limit" json:"floating_count_limit,computed"`
	GPUCountLimit                     types.Int64 `tfsdk:"gpu_count_limit" json:"gpu_count_limit,computed"`
	GPUVirtualA100CountLimit          types.Int64 `tfsdk:"gpu_virtual_a100_count_limit" json:"gpu_virtual_a100_count_limit,computed"`
	GPUVirtualH100CountLimit          types.Int64 `tfsdk:"gpu_virtual_h100_count_limit" json:"gpu_virtual_h100_count_limit,computed"`
	GPUVirtualH200CountLimit          types.Int64 `tfsdk:"gpu_virtual_h200_count_limit" json:"gpu_virtual_h200_count_limit,computed"`
	GPUVirtualL40sCountLimit          types.Int64 `tfsdk:"gpu_virtual_l40s_count_limit" json:"gpu_virtual_l40s_count_limit,computed"`
	ImageCountLimit                   types.Int64 `tfsdk:"image_count_limit" json:"image_count_limit,computed"`
	ImageSizeLimit                    types.Int64 `tfsdk:"image_size_limit" json:"image_size_limit,computed"`
	IpuCountLimit                     types.Int64 `tfsdk:"ipu_count_limit" json:"ipu_count_limit,computed"`
	LaasTopicCountLimit               types.Int64 `tfsdk:"laas_topic_count_limit" json:"laas_topic_count_limit,computed"`
	LoadbalancerCountLimit            types.Int64 `tfsdk:"loadbalancer_count_limit" json:"loadbalancer_count_limit,computed"`
	NetworkCountLimit                 types.Int64 `tfsdk:"network_count_limit" json:"network_count_limit,computed"`
	RamLimit                          types.Int64 `tfsdk:"ram_limit" json:"ram_limit,computed"`
	RegionID                          types.Int64 `tfsdk:"region_id" json:"region_id,computed"`
	RegistryCountLimit                types.Int64 `tfsdk:"registry_count_limit" json:"registry_count_limit,computed"`
	RegistryStorageLimit              types.Int64 `tfsdk:"registry_storage_limit" json:"registry_storage_limit,computed"`
	RouterCountLimit                  types.Int64 `tfsdk:"router_count_limit" json:"router_count_limit,computed"`
	SecretCountLimit                  types.Int64 `tfsdk:"secret_count_limit" json:"secret_count_limit,computed"`
	ServergroupCountLimit             types.Int64 `tfsdk:"servergroup_count_limit" json:"servergroup_count_limit,computed"`
	SfsCountLimit                     types.Int64 `tfsdk:"sfs_count_limit" json:"sfs_count_limit,computed"`
	SfsSizeLimit                      types.Int64 `tfsdk:"sfs_size_limit" json:"sfs_size_limit,computed"`
	SharedVmCountLimit                types.Int64 `tfsdk:"shared_vm_count_limit" json:"shared_vm_count_limit,computed"`
	SnapshotScheduleCountLimit        types.Int64 `tfsdk:"snapshot_schedule_count_limit" json:"snapshot_schedule_count_limit,computed"`
	SubnetCountLimit                  types.Int64 `tfsdk:"subnet_count_limit" json:"subnet_count_limit,computed"`
	VmCountLimit                      types.Int64 `tfsdk:"vm_count_limit" json:"vm_count_limit,computed"`
	VolumeCountLimit                  types.Int64 `tfsdk:"volume_count_limit" json:"volume_count_limit,computed"`
	VolumeSizeLimit                   types.Int64 `tfsdk:"volume_size_limit" json:"volume_size_limit,computed"`
	VolumeSnapshotsCountLimit         types.Int64 `tfsdk:"volume_snapshots_count_limit" json:"volume_snapshots_count_limit,computed"`
	VolumeSnapshotsSizeLimit          types.Int64 `tfsdk:"volume_snapshots_size_limit" json:"volume_snapshots_size_limit,computed"`
}
