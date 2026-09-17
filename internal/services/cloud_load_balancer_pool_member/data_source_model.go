// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_member

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudLoadBalancerPoolMemberDataSourceModel struct {
	ID                 types.String                                         `tfsdk:"id" path:"member_id,computed"`
	MemberID           types.String                                         `tfsdk:"member_id" path:"member_id,optional"`
	PoolID             types.String                                         `tfsdk:"pool_id" path:"pool_id,required"`
	ProjectID          types.Int64                                          `tfsdk:"project_id" path:"project_id,optional"`
	RegionID           types.Int64                                          `tfsdk:"region_id" path:"region_id,optional"`
	Address            types.String                                         `tfsdk:"address" json:"address,computed"`
	AdminStateUp       types.Bool                                           `tfsdk:"admin_state_up" json:"admin_state_up,computed"`
	Backup             types.Bool                                           `tfsdk:"backup" json:"backup,computed"`
	MonitorAddress     types.String                                         `tfsdk:"monitor_address" json:"monitor_address,computed"`
	MonitorPort        types.Int64                                          `tfsdk:"monitor_port" json:"monitor_port,computed"`
	OperatingStatus    types.String                                         `tfsdk:"operating_status" json:"operating_status,computed"`
	ProtocolPort       types.Int64                                          `tfsdk:"protocol_port" json:"protocol_port,computed"`
	ProvisioningStatus types.String                                         `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	SubnetID           types.String                                         `tfsdk:"subnet_id" json:"subnet_id,computed"`
	Weight             types.Int64                                          `tfsdk:"weight" json:"weight,computed"`
	FindOneBy          *CloudLoadBalancerPoolMemberFindOneByDataSourceModel `tfsdk:"find_one_by"`
}

func (m *CloudLoadBalancerPoolMemberDataSourceModel) toReadParams(_ context.Context) (params cloud.LoadBalancerPoolMemberGetParams, diags diag.Diagnostics) {
	params = cloud.LoadBalancerPoolMemberGetParams{
		PoolID: m.PoolID.ValueString(),
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

func (m *CloudLoadBalancerPoolMemberDataSourceModel) toListParams(_ context.Context) (params cloud.LoadBalancerPoolMemberListParams, diags diag.Diagnostics) {
	params = cloud.LoadBalancerPoolMemberListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.FindOneBy.OrderBy.IsNull() {
		params.OrderBy = cloud.LoadBalancerPoolMemberListParamsOrderBy(m.FindOneBy.OrderBy.ValueString())
	}

	return
}

type CloudLoadBalancerPoolMemberFindOneByDataSourceModel struct {
	OrderBy types.String `tfsdk:"order_by" query:"order_by,computed_optional"`
}
