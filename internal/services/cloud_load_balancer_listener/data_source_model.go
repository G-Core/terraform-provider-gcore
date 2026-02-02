// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_listener

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudLoadBalancerListenerDataSourceModel struct {
	ID                   types.String                                                                   `tfsdk:"id" path:"listener_id,computed"`
	ListenerID           types.String                                                                   `tfsdk:"listener_id" path:"listener_id,required"`
	ProjectID            types.Int64                                                                    `tfsdk:"project_id" path:"project_id,optional"`
	RegionID             types.Int64                                                                    `tfsdk:"region_id" path:"region_id,optional"`
	ShowStats            types.Bool                                                                     `tfsdk:"show_stats" query:"show_stats,computed_optional"`
	AdminStateUp         types.Bool                                                                     `tfsdk:"admin_state_up" json:"admin_state_up,computed"`
	ConnectionLimit      types.Int64                                                                    `tfsdk:"connection_limit" json:"connection_limit,computed"`
	CreatorTaskID        types.String                                                                   `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	LoadBalancerID       types.String                                                                   `tfsdk:"load_balancer_id" json:"load_balancer_id,computed"`
	Name                 types.String                                                                   `tfsdk:"name" json:"name,computed"`
	OperatingStatus      types.String                                                                   `tfsdk:"operating_status" json:"operating_status,computed"`
	PoolCount            types.Int64                                                                    `tfsdk:"pool_count" json:"pool_count,computed"`
	Protocol             types.String                                                                   `tfsdk:"protocol" json:"protocol,computed"`
	ProtocolPort         types.Int64                                                                    `tfsdk:"protocol_port" json:"protocol_port,computed"`
	ProvisioningStatus   types.String                                                                   `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	SecretID             types.String                                                                   `tfsdk:"secret_id" json:"secret_id,computed"`
	TaskID               types.String                                                                   `tfsdk:"task_id" json:"task_id,computed"`
	TimeoutClientData    types.Int64                                                                    `tfsdk:"timeout_client_data" json:"timeout_client_data,computed"`
	TimeoutMemberConnect types.Int64                                                                    `tfsdk:"timeout_member_connect" json:"timeout_member_connect,computed"`
	TimeoutMemberData    types.Int64                                                                    `tfsdk:"timeout_member_data" json:"timeout_member_data,computed"`
	AllowedCidrs         customfield.List[types.String]                                                 `tfsdk:"allowed_cidrs" json:"allowed_cidrs,computed"`
	InsertHeaders        customfield.Map[jsontypes.Normalized]                                          `tfsdk:"insert_headers" json:"insert_headers,computed"`
	SniSecretID          customfield.List[types.String]                                                 `tfsdk:"sni_secret_id" json:"sni_secret_id,computed"`
	Stats                customfield.NestedObject[CloudLoadBalancerListenerStatsDataSourceModel]        `tfsdk:"stats" json:"stats,computed"`
	UserList             customfield.NestedObjectList[CloudLoadBalancerListenerUserListDataSourceModel] `tfsdk:"user_list" json:"user_list,computed"`
}

func (m *CloudLoadBalancerListenerDataSourceModel) toReadParams(_ context.Context) (params cloud.LoadBalancerListenerGetParams, diags diag.Diagnostics) {
	params = cloud.LoadBalancerListenerGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.ShowStats.IsNull() {
		params.ShowStats = param.NewOpt(m.ShowStats.ValueBool())
	}

	return
}

type CloudLoadBalancerListenerStatsDataSourceModel struct {
	ActiveConnections types.Int64 `tfsdk:"active_connections" json:"active_connections,computed"`
	BytesIn           types.Int64 `tfsdk:"bytes_in" json:"bytes_in,computed"`
	BytesOut          types.Int64 `tfsdk:"bytes_out" json:"bytes_out,computed"`
	RequestErrors     types.Int64 `tfsdk:"request_errors" json:"request_errors,computed"`
	TotalConnections  types.Int64 `tfsdk:"total_connections" json:"total_connections,computed"`
}

type CloudLoadBalancerListenerUserListDataSourceModel struct {
	EncryptedPassword types.String `tfsdk:"encrypted_password" json:"encrypted_password,computed"`
	Username          types.String `tfsdk:"username" json:"username,computed"`
}
