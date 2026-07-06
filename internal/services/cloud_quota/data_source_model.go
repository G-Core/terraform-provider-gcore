// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_quota

import (
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudQuotaDataSourceModel struct {
	GlobalQuotas   customfield.NestedObject[CloudQuotaGlobalQuotasDataSourceModel]       `tfsdk:"global_quotas" json:"global_quotas,computed"`
	RegionalQuotas customfield.NestedObjectList[CloudQuotaRegionalQuotasDataSourceModel] `tfsdk:"regional_quotas" json:"regional_quotas,computed"`
}

type CloudQuotaGlobalQuotasDataSourceModel struct {
	InferenceCPUMillicoreCountLimit      types.Int64 `tfsdk:"inference_cpu_millicore_count_limit" json:"inference_cpu_millicore_count_limit,computed"`
	InferenceCPUMillicoreCountUsage      types.Int64 `tfsdk:"inference_cpu_millicore_count_usage" json:"inference_cpu_millicore_count_usage,computed"`
	InferenceGPUA100CountLimit           types.Int64 `tfsdk:"inference_gpu_a100_count_limit" json:"inference_gpu_a100_count_limit,computed"`
	InferenceGPUA100CountUsage           types.Int64 `tfsdk:"inference_gpu_a100_count_usage" json:"inference_gpu_a100_count_usage,computed"`
	InferenceGPUH100CountLimit           types.Int64 `tfsdk:"inference_gpu_h100_count_limit" json:"inference_gpu_h100_count_limit,computed"`
	InferenceGPUH100CountUsage           types.Int64 `tfsdk:"inference_gpu_h100_count_usage" json:"inference_gpu_h100_count_usage,computed"`
	InferenceGPUL40sCountLimit           types.Int64 `tfsdk:"inference_gpu_l40s_count_limit" json:"inference_gpu_l40s_count_limit,computed"`
	InferenceGPUL40sCountUsage           types.Int64 `tfsdk:"inference_gpu_l40s_count_usage" json:"inference_gpu_l40s_count_usage,computed"`
	InferenceInstanceCountLimit          types.Int64 `tfsdk:"inference_instance_count_limit" json:"inference_instance_count_limit,computed"`
	InferenceInstanceCountUsage          types.Int64 `tfsdk:"inference_instance_count_usage" json:"inference_instance_count_usage,computed"`
	InferencePublicModelAPIKeyCountLimit types.Int64 `tfsdk:"inference_public_model_api_key_count_limit" json:"inference_public_model_api_key_count_limit,computed"`
	InferencePublicModelAPIKeyCountUsage types.Int64 `tfsdk:"inference_public_model_api_key_count_usage" json:"inference_public_model_api_key_count_usage,computed"`
	KeypairCountLimit                    types.Int64 `tfsdk:"keypair_count_limit" json:"keypair_count_limit,computed"`
	KeypairCountUsage                    types.Int64 `tfsdk:"keypair_count_usage" json:"keypair_count_usage,computed"`
	ProjectCountLimit                    types.Int64 `tfsdk:"project_count_limit" json:"project_count_limit,computed"`
	ProjectCountUsage                    types.Int64 `tfsdk:"project_count_usage" json:"project_count_usage,computed"`
}

type CloudQuotaRegionalQuotasDataSourceModel struct {
	BaremetalBasicCountLimit          types.Int64 `tfsdk:"baremetal_basic_count_limit" json:"baremetal_basic_count_limit,computed"`
	BaremetalBasicCountUsage          types.Int64 `tfsdk:"baremetal_basic_count_usage" json:"baremetal_basic_count_usage,computed"`
	BaremetalGPUA100CountLimit        types.Int64 `tfsdk:"baremetal_gpu_a100_count_limit" json:"baremetal_gpu_a100_count_limit,computed"`
	BaremetalGPUA100CountUsage        types.Int64 `tfsdk:"baremetal_gpu_a100_count_usage" json:"baremetal_gpu_a100_count_usage,computed"`
	BaremetalGPUCountLimit            types.Int64 `tfsdk:"baremetal_gpu_count_limit" json:"baremetal_gpu_count_limit,computed"`
	BaremetalGPUCountUsage            types.Int64 `tfsdk:"baremetal_gpu_count_usage" json:"baremetal_gpu_count_usage,computed"`
	BaremetalGPUH100CountLimit        types.Int64 `tfsdk:"baremetal_gpu_h100_count_limit" json:"baremetal_gpu_h100_count_limit,computed"`
	BaremetalGPUH100CountUsage        types.Int64 `tfsdk:"baremetal_gpu_h100_count_usage" json:"baremetal_gpu_h100_count_usage,computed"`
	BaremetalGPUH200CountLimit        types.Int64 `tfsdk:"baremetal_gpu_h200_count_limit" json:"baremetal_gpu_h200_count_limit,computed"`
	BaremetalGPUH200CountUsage        types.Int64 `tfsdk:"baremetal_gpu_h200_count_usage" json:"baremetal_gpu_h200_count_usage,computed"`
	BaremetalGPUL40sCountLimit        types.Int64 `tfsdk:"baremetal_gpu_l40s_count_limit" json:"baremetal_gpu_l40s_count_limit,computed"`
	BaremetalGPUL40sCountUsage        types.Int64 `tfsdk:"baremetal_gpu_l40s_count_usage" json:"baremetal_gpu_l40s_count_usage,computed"`
	BaremetalHfCountLimit             types.Int64 `tfsdk:"baremetal_hf_count_limit" json:"baremetal_hf_count_limit,computed"`
	BaremetalHfCountUsage             types.Int64 `tfsdk:"baremetal_hf_count_usage" json:"baremetal_hf_count_usage,computed"`
	BaremetalInfrastructureCountLimit types.Int64 `tfsdk:"baremetal_infrastructure_count_limit" json:"baremetal_infrastructure_count_limit,computed"`
	BaremetalInfrastructureCountUsage types.Int64 `tfsdk:"baremetal_infrastructure_count_usage" json:"baremetal_infrastructure_count_usage,computed"`
	BaremetalNetworkCountLimit        types.Int64 `tfsdk:"baremetal_network_count_limit" json:"baremetal_network_count_limit,computed"`
	BaremetalNetworkCountUsage        types.Int64 `tfsdk:"baremetal_network_count_usage" json:"baremetal_network_count_usage,computed"`
	BaremetalStorageCountLimit        types.Int64 `tfsdk:"baremetal_storage_count_limit" json:"baremetal_storage_count_limit,computed"`
	BaremetalStorageCountUsage        types.Int64 `tfsdk:"baremetal_storage_count_usage" json:"baremetal_storage_count_usage,computed"`
	CaasContainerCountLimit           types.Int64 `tfsdk:"caas_container_count_limit" json:"caas_container_count_limit,computed"`
	CaasContainerCountUsage           types.Int64 `tfsdk:"caas_container_count_usage" json:"caas_container_count_usage,computed"`
	CaasCPUCountLimit                 types.Int64 `tfsdk:"caas_cpu_count_limit" json:"caas_cpu_count_limit,computed"`
	CaasCPUCountUsage                 types.Int64 `tfsdk:"caas_cpu_count_usage" json:"caas_cpu_count_usage,computed"`
	CaasGPUCountLimit                 types.Int64 `tfsdk:"caas_gpu_count_limit" json:"caas_gpu_count_limit,computed"`
	CaasGPUCountUsage                 types.Int64 `tfsdk:"caas_gpu_count_usage" json:"caas_gpu_count_usage,computed"`
	CaasRamSizeLimit                  types.Int64 `tfsdk:"caas_ram_size_limit" json:"caas_ram_size_limit,computed"`
	CaasRamSizeUsage                  types.Int64 `tfsdk:"caas_ram_size_usage" json:"caas_ram_size_usage,computed"`
	ClusterCountLimit                 types.Int64 `tfsdk:"cluster_count_limit" json:"cluster_count_limit,computed"`
	ClusterCountUsage                 types.Int64 `tfsdk:"cluster_count_usage" json:"cluster_count_usage,computed"`
	CPUCountLimit                     types.Int64 `tfsdk:"cpu_count_limit" json:"cpu_count_limit,computed"`
	CPUCountUsage                     types.Int64 `tfsdk:"cpu_count_usage" json:"cpu_count_usage,computed"`
	DbaasPostgresClusterCountLimit    types.Int64 `tfsdk:"dbaas_postgres_cluster_count_limit" json:"dbaas_postgres_cluster_count_limit,computed"`
	DbaasPostgresClusterCountUsage    types.Int64 `tfsdk:"dbaas_postgres_cluster_count_usage" json:"dbaas_postgres_cluster_count_usage,computed"`
	ExternalIPCountLimit              types.Int64 `tfsdk:"external_ip_count_limit" json:"external_ip_count_limit,computed"`
	ExternalIPCountUsage              types.Int64 `tfsdk:"external_ip_count_usage" json:"external_ip_count_usage,computed"`
	FaasCPUCountLimit                 types.Int64 `tfsdk:"faas_cpu_count_limit" json:"faas_cpu_count_limit,computed"`
	FaasCPUCountUsage                 types.Int64 `tfsdk:"faas_cpu_count_usage" json:"faas_cpu_count_usage,computed"`
	FaasFunctionCountLimit            types.Int64 `tfsdk:"faas_function_count_limit" json:"faas_function_count_limit,computed"`
	FaasFunctionCountUsage            types.Int64 `tfsdk:"faas_function_count_usage" json:"faas_function_count_usage,computed"`
	FaasNamespaceCountLimit           types.Int64 `tfsdk:"faas_namespace_count_limit" json:"faas_namespace_count_limit,computed"`
	FaasNamespaceCountUsage           types.Int64 `tfsdk:"faas_namespace_count_usage" json:"faas_namespace_count_usage,computed"`
	FaasRamSizeLimit                  types.Int64 `tfsdk:"faas_ram_size_limit" json:"faas_ram_size_limit,computed"`
	FaasRamSizeUsage                  types.Int64 `tfsdk:"faas_ram_size_usage" json:"faas_ram_size_usage,computed"`
	FirewallCountLimit                types.Int64 `tfsdk:"firewall_count_limit" json:"firewall_count_limit,computed"`
	FirewallCountUsage                types.Int64 `tfsdk:"firewall_count_usage" json:"firewall_count_usage,computed"`
	FloatingCountLimit                types.Int64 `tfsdk:"floating_count_limit" json:"floating_count_limit,computed"`
	FloatingCountUsage                types.Int64 `tfsdk:"floating_count_usage" json:"floating_count_usage,computed"`
	GPUCountLimit                     types.Int64 `tfsdk:"gpu_count_limit" json:"gpu_count_limit,computed"`
	GPUCountUsage                     types.Int64 `tfsdk:"gpu_count_usage" json:"gpu_count_usage,computed"`
	GPUVirtualA100CountLimit          types.Int64 `tfsdk:"gpu_virtual_a100_count_limit" json:"gpu_virtual_a100_count_limit,computed"`
	GPUVirtualA100CountUsage          types.Int64 `tfsdk:"gpu_virtual_a100_count_usage" json:"gpu_virtual_a100_count_usage,computed"`
	GPUVirtualH100CountLimit          types.Int64 `tfsdk:"gpu_virtual_h100_count_limit" json:"gpu_virtual_h100_count_limit,computed"`
	GPUVirtualH100CountUsage          types.Int64 `tfsdk:"gpu_virtual_h100_count_usage" json:"gpu_virtual_h100_count_usage,computed"`
	GPUVirtualH200CountLimit          types.Int64 `tfsdk:"gpu_virtual_h200_count_limit" json:"gpu_virtual_h200_count_limit,computed"`
	GPUVirtualH200CountUsage          types.Int64 `tfsdk:"gpu_virtual_h200_count_usage" json:"gpu_virtual_h200_count_usage,computed"`
	GPUVirtualL40sCountLimit          types.Int64 `tfsdk:"gpu_virtual_l40s_count_limit" json:"gpu_virtual_l40s_count_limit,computed"`
	GPUVirtualL40sCountUsage          types.Int64 `tfsdk:"gpu_virtual_l40s_count_usage" json:"gpu_virtual_l40s_count_usage,computed"`
	ImageCountLimit                   types.Int64 `tfsdk:"image_count_limit" json:"image_count_limit,computed"`
	ImageCountUsage                   types.Int64 `tfsdk:"image_count_usage" json:"image_count_usage,computed"`
	ImageSizeLimit                    types.Int64 `tfsdk:"image_size_limit" json:"image_size_limit,computed"`
	ImageSizeUsage                    types.Int64 `tfsdk:"image_size_usage" json:"image_size_usage,computed"`
	IpuCountLimit                     types.Int64 `tfsdk:"ipu_count_limit" json:"ipu_count_limit,computed"`
	IpuCountUsage                     types.Int64 `tfsdk:"ipu_count_usage" json:"ipu_count_usage,computed"`
	LaasTopicCountLimit               types.Int64 `tfsdk:"laas_topic_count_limit" json:"laas_topic_count_limit,computed"`
	LaasTopicCountUsage               types.Int64 `tfsdk:"laas_topic_count_usage" json:"laas_topic_count_usage,computed"`
	LoadbalancerCountLimit            types.Int64 `tfsdk:"loadbalancer_count_limit" json:"loadbalancer_count_limit,computed"`
	LoadbalancerCountUsage            types.Int64 `tfsdk:"loadbalancer_count_usage" json:"loadbalancer_count_usage,computed"`
	NetworkCountLimit                 types.Int64 `tfsdk:"network_count_limit" json:"network_count_limit,computed"`
	NetworkCountUsage                 types.Int64 `tfsdk:"network_count_usage" json:"network_count_usage,computed"`
	RamLimit                          types.Int64 `tfsdk:"ram_limit" json:"ram_limit,computed"`
	RamUsage                          types.Int64 `tfsdk:"ram_usage" json:"ram_usage,computed"`
	RegionID                          types.Int64 `tfsdk:"region_id" json:"region_id,computed"`
	RegistryCountLimit                types.Int64 `tfsdk:"registry_count_limit" json:"registry_count_limit,computed"`
	RegistryCountUsage                types.Int64 `tfsdk:"registry_count_usage" json:"registry_count_usage,computed"`
	RegistryStorageLimit              types.Int64 `tfsdk:"registry_storage_limit" json:"registry_storage_limit,computed"`
	RegistryStorageUsage              types.Int64 `tfsdk:"registry_storage_usage" json:"registry_storage_usage,computed"`
	RouterCountLimit                  types.Int64 `tfsdk:"router_count_limit" json:"router_count_limit,computed"`
	RouterCountUsage                  types.Int64 `tfsdk:"router_count_usage" json:"router_count_usage,computed"`
	SecretCountLimit                  types.Int64 `tfsdk:"secret_count_limit" json:"secret_count_limit,computed"`
	SecretCountUsage                  types.Int64 `tfsdk:"secret_count_usage" json:"secret_count_usage,computed"`
	ServergroupCountLimit             types.Int64 `tfsdk:"servergroup_count_limit" json:"servergroup_count_limit,computed"`
	ServergroupCountUsage             types.Int64 `tfsdk:"servergroup_count_usage" json:"servergroup_count_usage,computed"`
	SfsCountLimit                     types.Int64 `tfsdk:"sfs_count_limit" json:"sfs_count_limit,computed"`
	SfsCountUsage                     types.Int64 `tfsdk:"sfs_count_usage" json:"sfs_count_usage,computed"`
	SfsSizeLimit                      types.Int64 `tfsdk:"sfs_size_limit" json:"sfs_size_limit,computed"`
	SfsSizeUsage                      types.Int64 `tfsdk:"sfs_size_usage" json:"sfs_size_usage,computed"`
	SharedVmCountLimit                types.Int64 `tfsdk:"shared_vm_count_limit" json:"shared_vm_count_limit,computed"`
	SharedVmCountUsage                types.Int64 `tfsdk:"shared_vm_count_usage" json:"shared_vm_count_usage,computed"`
	SnapshotScheduleCountLimit        types.Int64 `tfsdk:"snapshot_schedule_count_limit" json:"snapshot_schedule_count_limit,computed"`
	SnapshotScheduleCountUsage        types.Int64 `tfsdk:"snapshot_schedule_count_usage" json:"snapshot_schedule_count_usage,computed"`
	SubnetCountLimit                  types.Int64 `tfsdk:"subnet_count_limit" json:"subnet_count_limit,computed"`
	SubnetCountUsage                  types.Int64 `tfsdk:"subnet_count_usage" json:"subnet_count_usage,computed"`
	VmCountLimit                      types.Int64 `tfsdk:"vm_count_limit" json:"vm_count_limit,computed"`
	VmCountUsage                      types.Int64 `tfsdk:"vm_count_usage" json:"vm_count_usage,computed"`
	VolumeCountLimit                  types.Int64 `tfsdk:"volume_count_limit" json:"volume_count_limit,computed"`
	VolumeCountUsage                  types.Int64 `tfsdk:"volume_count_usage" json:"volume_count_usage,computed"`
	VolumeSizeLimit                   types.Int64 `tfsdk:"volume_size_limit" json:"volume_size_limit,computed"`
	VolumeSizeUsage                   types.Int64 `tfsdk:"volume_size_usage" json:"volume_size_usage,computed"`
	VolumeSnapshotsCountLimit         types.Int64 `tfsdk:"volume_snapshots_count_limit" json:"volume_snapshots_count_limit,computed"`
	VolumeSnapshotsCountUsage         types.Int64 `tfsdk:"volume_snapshots_count_usage" json:"volume_snapshots_count_usage,computed"`
	VolumeSnapshotsSizeLimit          types.Int64 `tfsdk:"volume_snapshots_size_limit" json:"volume_snapshots_size_limit,computed"`
	VolumeSnapshotsSizeUsage          types.Int64 `tfsdk:"volume_snapshots_size_usage" json:"volume_snapshots_size_usage,computed"`
}
