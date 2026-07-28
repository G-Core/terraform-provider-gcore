// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_cluster_pool

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudK8SClusterPoolDataSourceModel struct {
	ClusterName        types.String                   `tfsdk:"cluster_name" path:"cluster_name,required"`
	PoolName           types.String                   `tfsdk:"pool_name" path:"pool_name,required"`
	ProjectID          types.Int64                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID           types.Int64                    `tfsdk:"region_id" path:"region_id,optional"`
	AutoHealingEnabled types.Bool                     `tfsdk:"auto_healing_enabled" json:"auto_healing_enabled,computed"`
	BootVolumeSize     types.Int64                    `tfsdk:"boot_volume_size" json:"boot_volume_size,computed"`
	BootVolumeType     types.String                   `tfsdk:"boot_volume_type" json:"boot_volume_type,computed"`
	CreatedAt          types.String                   `tfsdk:"created_at" json:"created_at,computed"`
	FlavorID           types.String                   `tfsdk:"flavor_id" json:"flavor_id,computed"`
	ID                 types.String                   `tfsdk:"id" json:"id,computed"`
	IsPublicIpv4       types.Bool                     `tfsdk:"is_public_ipv4" json:"is_public_ipv4,computed"`
	MaxNodeCount       types.Int64                    `tfsdk:"max_node_count" json:"max_node_count,computed"`
	MinNodeCount       types.Int64                    `tfsdk:"min_node_count" json:"min_node_count,computed"`
	Name               types.String                   `tfsdk:"name" json:"name,computed"`
	NodeCount          types.Int64                    `tfsdk:"node_count" json:"node_count,computed"`
	ServergroupID      types.String                   `tfsdk:"servergroup_id" json:"servergroup_id,computed"`
	ServergroupName    types.String                   `tfsdk:"servergroup_name" json:"servergroup_name,computed"`
	ServergroupPolicy  types.String                   `tfsdk:"servergroup_policy" json:"servergroup_policy,computed"`
	Status             types.String                   `tfsdk:"status" json:"status,computed"`
	CrioConfig         customfield.Map[types.String]  `tfsdk:"crio_config" json:"crio_config,computed"`
	KubeletConfig      customfield.Map[types.String]  `tfsdk:"kubelet_config" json:"kubelet_config,computed"`
	Labels             customfield.Map[types.String]  `tfsdk:"labels" json:"labels,computed"`
	SecurityGroupIDs   customfield.List[types.String] `tfsdk:"security_group_ids" json:"security_group_ids,computed"`
	Taints             customfield.Map[types.String]  `tfsdk:"taints" json:"taints,computed"`
}

func (m *CloudK8SClusterPoolDataSourceModel) toReadParams(_ context.Context) (params cloud.K8SClusterPoolGetParams, diags diag.Diagnostics) {
	params = cloud.K8SClusterPoolGetParams{
		ClusterName: m.ClusterName.ValueString(),
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}
