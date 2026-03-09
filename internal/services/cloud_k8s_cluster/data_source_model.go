// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_cluster

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudK8SClusterDataSourceModel struct {
	ID               types.String                                                           `tfsdk:"id" path:"cluster_name,computed"`
	ClusterName      types.String                                                           `tfsdk:"cluster_name" path:"cluster_name,required"`
	ProjectID        types.Int64                                                            `tfsdk:"project_id" path:"project_id,optional"`
	RegionID         types.Int64                                                            `tfsdk:"region_id" path:"region_id,optional"`
	CreatedAt        types.String                                                           `tfsdk:"created_at" json:"created_at,computed"`
	CreatorTaskID    types.String                                                           `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	FixedNetwork     types.String                                                           `tfsdk:"fixed_network" json:"fixed_network,computed"`
	FixedSubnet      types.String                                                           `tfsdk:"fixed_subnet" json:"fixed_subnet,computed"`
	IsIpv6           types.Bool                                                             `tfsdk:"is_ipv6" json:"is_ipv6,computed"`
	IsPublic         types.Bool                                                             `tfsdk:"is_public" json:"is_public,computed"`
	Keypair          types.String                                                           `tfsdk:"keypair" json:"keypair,computed"`
	Name             types.String                                                           `tfsdk:"name" json:"name,computed"`
	PodsIPPool       types.String                                                           `tfsdk:"pods_ip_pool" json:"pods_ip_pool,computed"`
	PodsIpv6Pool     types.String                                                           `tfsdk:"pods_ipv6_pool" json:"pods_ipv6_pool,computed"`
	ServicesIPPool   types.String                                                           `tfsdk:"services_ip_pool" json:"services_ip_pool,computed"`
	ServicesIpv6Pool types.String                                                           `tfsdk:"services_ipv6_pool" json:"services_ipv6_pool,computed"`
	Status           types.String                                                           `tfsdk:"status" json:"status,computed"`
	TaskID           types.String                                                           `tfsdk:"task_id" json:"task_id,computed"`
	Version          types.String                                                           `tfsdk:"version" json:"version,computed"`
	AutoscalerConfig customfield.Map[types.String]                                          `tfsdk:"autoscaler_config" json:"autoscaler_config,computed"`
	AddOns           customfield.NestedObject[CloudK8SClusterAddOnsDataSourceModel]         `tfsdk:"add_ons" json:"add_ons,computed"`
	Authentication   customfield.NestedObject[CloudK8SClusterAuthenticationDataSourceModel] `tfsdk:"authentication" json:"authentication,computed"`
	Cni              customfield.NestedObject[CloudK8SClusterCniDataSourceModel]            `tfsdk:"cni" json:"cni,computed"`
	Csi              customfield.NestedObject[CloudK8SClusterCsiDataSourceModel]            `tfsdk:"csi" json:"csi,computed"`
	Logging          customfield.NestedObject[CloudK8SClusterLoggingDataSourceModel]        `tfsdk:"logging" json:"logging,computed"`
	Pools            customfield.NestedObjectList[CloudK8SClusterPoolsDataSourceModel]      `tfsdk:"pools" json:"pools,computed"`
}

func (m *CloudK8SClusterDataSourceModel) toReadParams(_ context.Context) (params cloud.K8SClusterGetParams, diags diag.Diagnostics) {
	params = cloud.K8SClusterGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudK8SClusterAddOnsDataSourceModel struct {
	Slurm customfield.NestedObject[CloudK8SClusterAddOnsSlurmDataSourceModel] `tfsdk:"slurm" json:"slurm,computed"`
}

type CloudK8SClusterAddOnsSlurmDataSourceModel struct {
	Enabled     types.Bool                     `tfsdk:"enabled" json:"enabled,computed"`
	FileShareID types.String                   `tfsdk:"file_share_id" json:"file_share_id,computed"`
	SSHKeyIDs   customfield.List[types.String] `tfsdk:"ssh_key_ids" json:"ssh_key_ids,computed"`
	WorkerCount types.Int64                    `tfsdk:"worker_count" json:"worker_count,computed"`
}

type CloudK8SClusterAuthenticationDataSourceModel struct {
	KubeconfigCreatedAt timetypes.RFC3339                                                          `tfsdk:"kubeconfig_created_at" json:"kubeconfig_created_at,computed" format:"date-time"`
	KubeconfigExpiresAt timetypes.RFC3339                                                          `tfsdk:"kubeconfig_expires_at" json:"kubeconfig_expires_at,computed" format:"date-time"`
	Oidc                customfield.NestedObject[CloudK8SClusterAuthenticationOidcDataSourceModel] `tfsdk:"oidc" json:"oidc,computed"`
}

type CloudK8SClusterAuthenticationOidcDataSourceModel struct {
	ClientID       types.String                   `tfsdk:"client_id" json:"client_id,computed"`
	GroupsClaim    types.String                   `tfsdk:"groups_claim" json:"groups_claim,computed"`
	GroupsPrefix   types.String                   `tfsdk:"groups_prefix" json:"groups_prefix,computed"`
	IssuerURL      types.String                   `tfsdk:"issuer_url" json:"issuer_url,computed"`
	RequiredClaims customfield.Map[types.String]  `tfsdk:"required_claims" json:"required_claims,computed"`
	SigningAlgs    customfield.List[types.String] `tfsdk:"signing_algs" json:"signing_algs,computed"`
	UsernameClaim  types.String                   `tfsdk:"username_claim" json:"username_claim,computed"`
	UsernamePrefix types.String                   `tfsdk:"username_prefix" json:"username_prefix,computed"`
}

type CloudK8SClusterCniDataSourceModel struct {
	Cilium                  customfield.NestedObject[CloudK8SClusterCniCiliumDataSourceModel] `tfsdk:"cilium" json:"cilium,computed"`
	CloudK8SClusterProvider types.String                                                      `tfsdk:"cloud_k8s_cluster_provider" json:"provider,computed"`
}

type CloudK8SClusterCniCiliumDataSourceModel struct {
	Encryption     types.Bool   `tfsdk:"encryption" json:"encryption,computed"`
	HubbleRelay    types.Bool   `tfsdk:"hubble_relay" json:"hubble_relay,computed"`
	HubbleUi       types.Bool   `tfsdk:"hubble_ui" json:"hubble_ui,computed"`
	LbAcceleration types.Bool   `tfsdk:"lb_acceleration" json:"lb_acceleration,computed"`
	LbMode         types.String `tfsdk:"lb_mode" json:"lb_mode,computed"`
	MaskSize       types.Int64  `tfsdk:"mask_size" json:"mask_size,computed"`
	MaskSizeV6     types.Int64  `tfsdk:"mask_size_v6" json:"mask_size_v6,computed"`
	RoutingMode    types.String `tfsdk:"routing_mode" json:"routing_mode,computed"`
	Tunnel         types.String `tfsdk:"tunnel" json:"tunnel,computed"`
}

type CloudK8SClusterCsiDataSourceModel struct {
	Nfs customfield.NestedObject[CloudK8SClusterCsiNfsDataSourceModel] `tfsdk:"nfs" json:"nfs,computed"`
}

type CloudK8SClusterCsiNfsDataSourceModel struct {
	VastEnabled types.Bool `tfsdk:"vast_enabled" json:"vast_enabled,computed"`
}

type CloudK8SClusterLoggingDataSourceModel struct {
	DestinationRegionID types.Int64                                                                    `tfsdk:"destination_region_id" json:"destination_region_id,computed"`
	Enabled             types.Bool                                                                     `tfsdk:"enabled" json:"enabled,computed"`
	TopicName           types.String                                                                   `tfsdk:"topic_name" json:"topic_name,computed"`
	RetentionPolicy     customfield.NestedObject[CloudK8SClusterLoggingRetentionPolicyDataSourceModel] `tfsdk:"retention_policy" json:"retention_policy,computed"`
}

type CloudK8SClusterLoggingRetentionPolicyDataSourceModel struct {
	Period types.Int64 `tfsdk:"period" json:"period,computed"`
}

type CloudK8SClusterPoolsDataSourceModel struct {
	ID                 types.String                  `tfsdk:"id" json:"id,computed"`
	AutoHealingEnabled types.Bool                    `tfsdk:"auto_healing_enabled" json:"auto_healing_enabled,computed"`
	BootVolumeSize     types.Int64                   `tfsdk:"boot_volume_size" json:"boot_volume_size,computed"`
	BootVolumeType     types.String                  `tfsdk:"boot_volume_type" json:"boot_volume_type,computed"`
	CreatedAt          types.String                  `tfsdk:"created_at" json:"created_at,computed"`
	CrioConfig         customfield.Map[types.String] `tfsdk:"crio_config" json:"crio_config,computed"`
	FlavorID           types.String                  `tfsdk:"flavor_id" json:"flavor_id,computed"`
	IsPublicIpv4       types.Bool                    `tfsdk:"is_public_ipv4" json:"is_public_ipv4,computed"`
	KubeletConfig      customfield.Map[types.String] `tfsdk:"kubelet_config" json:"kubelet_config,computed"`
	Labels             customfield.Map[types.String] `tfsdk:"labels" json:"labels,computed"`
	MaxNodeCount       types.Int64                   `tfsdk:"max_node_count" json:"max_node_count,computed"`
	MinNodeCount       types.Int64                   `tfsdk:"min_node_count" json:"min_node_count,computed"`
	Name               types.String                  `tfsdk:"name" json:"name,computed"`
	NodeCount          types.Int64                   `tfsdk:"node_count" json:"node_count,computed"`
	Status             types.String                  `tfsdk:"status" json:"status,computed"`
	Taints             customfield.Map[types.String] `tfsdk:"taints" json:"taints,computed"`
	ServergroupID      types.String                  `tfsdk:"servergroup_id" json:"servergroup_id,computed"`
	ServergroupName    types.String                  `tfsdk:"servergroup_name" json:"servergroup_name,computed"`
	ServergroupPolicy  types.String                  `tfsdk:"servergroup_policy" json:"servergroup_policy,computed"`
}
