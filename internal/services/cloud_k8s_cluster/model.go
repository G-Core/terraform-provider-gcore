// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_cluster

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudK8SClusterModel struct {
	ID               types.String                                              `tfsdk:"id" json:"-,computed"`
	Name             types.String                                              `tfsdk:"name" json:"name,required"`
	ProjectID        types.Int64                                               `tfsdk:"project_id" path:"project_id,optional"`
	RegionID         types.Int64                                               `tfsdk:"region_id" path:"region_id,optional"`
	Keypair          types.String                                              `tfsdk:"keypair" json:"keypair,required"`
	Version          types.String                                              `tfsdk:"version" json:"version,required"`
	Pools            *[]*CloudK8SClusterPoolsModel                             `tfsdk:"pools" json:"pools,required"`
	PodsIPPool       types.String                                              `tfsdk:"pods_ip_pool" json:"pods_ip_pool,optional"`
	PodsIpv6Pool     types.String                                              `tfsdk:"pods_ipv6_pool" json:"pods_ipv6_pool,optional"`
	ServicesIPPool   types.String                                              `tfsdk:"services_ip_pool" json:"services_ip_pool,optional"`
	ServicesIpv6Pool types.String                                              `tfsdk:"services_ipv6_pool" json:"services_ipv6_pool,optional"`
	FixedNetwork     types.String                                              `tfsdk:"fixed_network" json:"fixed_network,computed_optional"`
	FixedSubnet      types.String                                              `tfsdk:"fixed_subnet" json:"fixed_subnet,computed_optional"`
	IsIpv6           types.Bool                                                `tfsdk:"is_ipv6" json:"is_ipv6,computed_optional"`
	Csi              customfield.NestedObject[CloudK8SClusterCsiModel]         `tfsdk:"csi" json:"csi,computed_optional"`
	AutoscalerConfig *map[string]types.String                                  `tfsdk:"autoscaler_config" json:"autoscaler_config,optional"`
	AddOns           *CloudK8SClusterAddOnsModel                               `tfsdk:"add_ons" json:"add_ons,optional"`
	Authentication   *CloudK8SClusterAuthenticationModel                       `tfsdk:"authentication" json:"authentication,optional"`
	Cni              customfield.NestedObject[CloudK8SClusterCniModel]         `tfsdk:"cni" json:"cni,computed_optional"`
	DDOSProfile      customfield.NestedObject[CloudK8SClusterDDOSProfileModel] `tfsdk:"ddos_profile" json:"ddos_profile,computed_optional"`
	Logging          customfield.NestedObject[CloudK8SClusterLoggingModel]     `tfsdk:"logging" json:"logging,computed_optional"`
	CreatedAt        types.String                                              `tfsdk:"created_at" json:"created_at,computed"`
	CreatorTaskID    types.String                                              `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	IsPublic         types.Bool                                                `tfsdk:"is_public" json:"is_public,computed"`
	Status           types.String                                              `tfsdk:"status" json:"status,computed"`
	TaskID           types.String                                              `tfsdk:"task_id" json:"task_id,computed"`
	Tasks            customfield.List[types.String]                            `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
}

func (m CloudK8SClusterModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudK8SClusterModel) MarshalJSONForUpdate(state CloudK8SClusterModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudK8SClusterPoolsModel struct {
	FlavorID           types.String                  `tfsdk:"flavor_id" json:"flavor_id,required"`
	MinNodeCount       types.Int64                   `tfsdk:"min_node_count" json:"min_node_count,required"`
	Name               types.String                  `tfsdk:"name" json:"name,required"`
	AutoHealingEnabled types.Bool                    `tfsdk:"auto_healing_enabled" json:"auto_healing_enabled,computed_optional"`
	BootVolumeSize     types.Int64                   `tfsdk:"boot_volume_size" json:"boot_volume_size,optional"`
	BootVolumeType     types.String                  `tfsdk:"boot_volume_type" json:"boot_volume_type,optional"`
	CrioConfig         *map[string]types.String      `tfsdk:"crio_config" json:"crio_config,optional"`
	IsPublicIpv4       types.Bool                    `tfsdk:"is_public_ipv4" json:"is_public_ipv4,computed_optional"`
	KubeletConfig      *map[string]types.String      `tfsdk:"kubelet_config" json:"kubelet_config,optional"`
	Labels             customfield.Map[types.String] `tfsdk:"labels" json:"labels,computed_optional"`
	MaxNodeCount       types.Int64                   `tfsdk:"max_node_count" json:"max_node_count,optional"`
	ServergroupPolicy  types.String                  `tfsdk:"servergroup_policy" json:"servergroup_policy,optional"`
	Taints             customfield.Map[types.String] `tfsdk:"taints" json:"taints,computed_optional"`
}

type CloudK8SClusterCsiModel struct {
	Nfs customfield.NestedObject[CloudK8SClusterCsiNfsModel] `tfsdk:"nfs" json:"nfs,computed_optional"`
}

type CloudK8SClusterCsiNfsModel struct {
	VastEnabled types.Bool `tfsdk:"vast_enabled" json:"vast_enabled,computed_optional"`
}

type CloudK8SClusterAddOnsModel struct {
	Slurm *CloudK8SClusterAddOnsSlurmModel `tfsdk:"slurm" json:"slurm,optional"`
}

type CloudK8SClusterAddOnsSlurmModel struct {
	Enabled     types.Bool      `tfsdk:"enabled" json:"enabled,required"`
	FileShareID types.String    `tfsdk:"file_share_id" json:"file_share_id,required"`
	SSHKeyIDs   *[]types.String `tfsdk:"ssh_key_ids" json:"ssh_key_ids,required"`
	WorkerCount types.Int64     `tfsdk:"worker_count" json:"worker_count,required"`
}

type CloudK8SClusterAuthenticationModel struct {
	Oidc *CloudK8SClusterAuthenticationOidcModel `tfsdk:"oidc" json:"oidc,optional"`
}

type CloudK8SClusterAuthenticationOidcModel struct {
	ClientID       types.String             `tfsdk:"client_id" json:"client_id,optional"`
	GroupsClaim    types.String             `tfsdk:"groups_claim" json:"groups_claim,optional"`
	GroupsPrefix   types.String             `tfsdk:"groups_prefix" json:"groups_prefix,optional"`
	IssuerURL      types.String             `tfsdk:"issuer_url" json:"issuer_url,optional"`
	RequiredClaims *map[string]types.String `tfsdk:"required_claims" json:"required_claims,optional"`
	SigningAlgs    *[]types.String          `tfsdk:"signing_algs" json:"signing_algs,optional"`
	UsernameClaim  types.String             `tfsdk:"username_claim" json:"username_claim,optional"`
	UsernamePrefix types.String             `tfsdk:"username_prefix" json:"username_prefix,optional"`
}

type CloudK8SClusterCniModel struct {
	Cilium                  customfield.NestedObject[CloudK8SClusterCniCiliumModel] `tfsdk:"cilium" json:"cilium,computed_optional"`
	CloudK8SClusterProvider types.String                                            `tfsdk:"cloud_k8s_cluster_provider" json:"provider,computed_optional"`
}

type CloudK8SClusterCniCiliumModel struct {
	Encryption     types.Bool   `tfsdk:"encryption" json:"encryption,computed_optional"`
	HubbleRelay    types.Bool   `tfsdk:"hubble_relay" json:"hubble_relay,computed_optional"`
	HubbleUi       types.Bool   `tfsdk:"hubble_ui" json:"hubble_ui,computed_optional"`
	LbAcceleration types.Bool   `tfsdk:"lb_acceleration" json:"lb_acceleration,computed_optional"`
	LbMode         types.String `tfsdk:"lb_mode" json:"lb_mode,computed_optional"`
	MaskSize       types.Int64  `tfsdk:"mask_size" json:"mask_size,computed_optional"`
	MaskSizeV6     types.Int64  `tfsdk:"mask_size_v6" json:"mask_size_v6,computed_optional"`
	RoutingMode    types.String `tfsdk:"routing_mode" json:"routing_mode,computed_optional"`
	Tunnel         types.String `tfsdk:"tunnel" json:"tunnel,computed_optional"`
}

type CloudK8SClusterDDOSProfileModel struct {
	Enabled             types.Bool                                                          `tfsdk:"enabled" json:"enabled,required"`
	Fields              customfield.NestedObjectList[CloudK8SClusterDDOSProfileFieldsModel] `tfsdk:"fields" json:"fields,computed_optional"`
	ProfileTemplate     types.Int64                                                         `tfsdk:"profile_template" json:"profile_template,optional"`
	ProfileTemplateName types.String                                                        `tfsdk:"profile_template_name" json:"profile_template_name,optional"`
}

type CloudK8SClusterDDOSProfileFieldsModel struct {
	BaseField  types.Int64          `tfsdk:"base_field" json:"base_field,required"`
	FieldValue jsontypes.Normalized `tfsdk:"field_value" json:"field_value,optional"`
	Value      types.String         `tfsdk:"value" json:"value,optional"`
}

type CloudK8SClusterLoggingModel struct {
	DestinationRegionID types.Int64                                 `tfsdk:"destination_region_id" json:"destination_region_id,optional"`
	Enabled             types.Bool                                  `tfsdk:"enabled" json:"enabled,computed_optional"`
	RetentionPolicy     *CloudK8SClusterLoggingRetentionPolicyModel `tfsdk:"retention_policy" json:"retention_policy,optional"`
	TopicName           types.String                                `tfsdk:"topic_name" json:"topic_name,optional"`
}

type CloudK8SClusterLoggingRetentionPolicyModel struct {
	Period types.Int64 `tfsdk:"period" json:"period,required"`
}
