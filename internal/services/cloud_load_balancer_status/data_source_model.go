// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_status

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudLoadBalancerStatusDataSourceModel struct {
	LoadbalancerID     types.String                                                                  `tfsdk:"loadbalancer_id" path:"loadbalancer_id,required"`
	ProjectID          types.Int64                                                                   `tfsdk:"project_id" path:"project_id,required"`
	RegionID           types.Int64                                                                   `tfsdk:"region_id" path:"region_id,required"`
	ID                 types.String                                                                  `tfsdk:"id" json:"id,computed"`
	Name               types.String                                                                  `tfsdk:"name" json:"name,computed"`
	OperatingStatus    types.String                                                                  `tfsdk:"operating_status" json:"operating_status,computed"`
	ProvisioningStatus types.String                                                                  `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	Listeners          customfield.NestedObjectList[CloudLoadBalancerStatusListenersDataSourceModel] `tfsdk:"listeners" json:"listeners,computed"`
	Tags               customfield.NestedObjectList[CloudLoadBalancerStatusTagsDataSourceModel]      `tfsdk:"tags" json:"tags,computed"`
}

func (m *CloudLoadBalancerStatusDataSourceModel) toReadParams(_ context.Context) (params cloud.LoadBalancerStatusGetParams, diags diag.Diagnostics) {
	params = cloud.LoadBalancerStatusGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudLoadBalancerStatusListenersDataSourceModel struct {
	ID                 types.String                                                                       `tfsdk:"id" json:"id,computed"`
	Name               types.String                                                                       `tfsdk:"name" json:"name,computed"`
	OperatingStatus    types.String                                                                       `tfsdk:"operating_status" json:"operating_status,computed"`
	Pools              customfield.NestedObjectList[CloudLoadBalancerStatusListenersPoolsDataSourceModel] `tfsdk:"pools" json:"pools,computed"`
	ProvisioningStatus types.String                                                                       `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
}

type CloudLoadBalancerStatusListenersPoolsDataSourceModel struct {
	ID                 types.String                                                                                `tfsdk:"id" json:"id,computed"`
	Members            customfield.NestedObjectList[CloudLoadBalancerStatusListenersPoolsMembersDataSourceModel]   `tfsdk:"members" json:"members,computed"`
	Name               types.String                                                                                `tfsdk:"name" json:"name,computed"`
	OperatingStatus    types.String                                                                                `tfsdk:"operating_status" json:"operating_status,computed"`
	ProvisioningStatus types.String                                                                                `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	HealthMonitor      customfield.NestedObject[CloudLoadBalancerStatusListenersPoolsHealthMonitorDataSourceModel] `tfsdk:"health_monitor" json:"health_monitor,computed"`
}

type CloudLoadBalancerStatusListenersPoolsMembersDataSourceModel struct {
	ID                 types.String `tfsdk:"id" json:"id,computed"`
	Address            types.String `tfsdk:"address" json:"address,computed"`
	OperatingStatus    types.String `tfsdk:"operating_status" json:"operating_status,computed"`
	ProtocolPort       types.Int64  `tfsdk:"protocol_port" json:"protocol_port,computed"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
}

type CloudLoadBalancerStatusListenersPoolsHealthMonitorDataSourceModel struct {
	ID                 types.String `tfsdk:"id" json:"id,computed"`
	OperatingStatus    types.String `tfsdk:"operating_status" json:"operating_status,computed"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	Type               types.String `tfsdk:"type" json:"type,computed"`
}

type CloudLoadBalancerStatusTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}
