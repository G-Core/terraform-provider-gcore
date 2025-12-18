// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_listener

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudLoadBalancerListenerModel struct {
	ID                   types.String                                                         `tfsdk:"id" json:"id,computed"`
	ProjectID            types.Int64                                                          `tfsdk:"project_id" path:"project_id,optional"`
	RegionID             types.Int64                                                          `tfsdk:"region_id" path:"region_id,optional"`
	LoadBalancerID       types.String                                                         `tfsdk:"load_balancer_id" json:"load_balancer_id,required"`
	Protocol             types.String                                                         `tfsdk:"protocol" json:"protocol,required"`
	ProtocolPort         types.Int64                                                          `tfsdk:"protocol_port" json:"protocol_port,required"`
	DefaultPoolID        types.String                                                         `tfsdk:"default_pool_id" json:"default_pool_id,optional,no_refresh"`
	InsertXForwarded     types.Bool                                                           `tfsdk:"insert_x_forwarded" json:"insert_x_forwarded,optional,no_refresh"`
	Name                 types.String                                                         `tfsdk:"name" json:"name,required"`
	SecretID             types.String                                                         `tfsdk:"secret_id" json:"secret_id,optional"`
	AllowedCidrs         *[]types.String                                                      `tfsdk:"allowed_cidrs" json:"allowed_cidrs,optional"`
	ConnectionLimit      types.Int64                                                          `tfsdk:"connection_limit" json:"connection_limit,computed_optional"`
	TimeoutClientData    types.Int64                                                          `tfsdk:"timeout_client_data" json:"timeout_client_data,computed_optional"`
	TimeoutMemberConnect types.Int64                                                          `tfsdk:"timeout_member_connect" json:"timeout_member_connect,computed_optional"`
	TimeoutMemberData    types.Int64                                                          `tfsdk:"timeout_member_data" json:"timeout_member_data,computed_optional"`
	SniSecretID          customfield.List[types.String]                                       `tfsdk:"sni_secret_id" json:"sni_secret_id,computed_optional"`
	UserList             customfield.NestedObjectList[CloudLoadBalancerListenerUserListModel] `tfsdk:"user_list" json:"user_list,computed_optional"`
	CreatorTaskID        types.String                                                         `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	OperatingStatus      types.String                                                         `tfsdk:"operating_status" json:"operating_status,computed"`
	PoolCount            types.Int64                                                          `tfsdk:"pool_count" json:"pool_count,computed"`
	ProvisioningStatus   types.String                                                         `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	TaskID               types.String                                                         `tfsdk:"task_id" json:"task_id,computed"`
	InsertHeaders        customfield.Map[jsontypes.Normalized]                                `tfsdk:"insert_headers" json:"insert_headers,computed"`
	Tasks                customfield.List[types.String]                                       `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
	Stats                customfield.NestedObject[CloudLoadBalancerListenerStatsModel]        `tfsdk:"stats" json:"stats,computed"`
}

func (m CloudLoadBalancerListenerModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudLoadBalancerListenerModel) MarshalJSONForUpdate(state CloudLoadBalancerListenerModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudLoadBalancerListenerUserListModel struct {
	EncryptedPassword types.String `tfsdk:"encrypted_password" json:"encrypted_password,required"`
	Username          types.String `tfsdk:"username" json:"username,required"`
}

type CloudLoadBalancerListenerStatsModel struct {
	ActiveConnections types.Int64 `tfsdk:"active_connections" json:"active_connections,computed"`
	BytesIn           types.Int64 `tfsdk:"bytes_in" json:"bytes_in,computed"`
	BytesOut          types.Int64 `tfsdk:"bytes_out" json:"bytes_out,computed"`
	RequestErrors     types.Int64 `tfsdk:"request_errors" json:"request_errors,computed"`
	TotalConnections  types.Int64 `tfsdk:"total_connections" json:"total_connections,computed"`
}
