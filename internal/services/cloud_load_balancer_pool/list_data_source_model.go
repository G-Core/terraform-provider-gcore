// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudLoadBalancerPoolsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudLoadBalancerPoolsItemsDataSourceModel] `json:"results,computed"`
}

type CloudLoadBalancerPoolsDataSourceModel struct {
	ProjectID      types.Int64                                                              `tfsdk:"project_id" path:"project_id,optional"`
	RegionID       types.Int64                                                              `tfsdk:"region_id" path:"region_id,optional"`
	ListenerID     types.String                                                             `tfsdk:"listener_id" query:"listener_id,optional"`
	LoadBalancerID types.String                                                             `tfsdk:"load_balancer_id" query:"load_balancer_id,optional"`
	Name           types.String                                                             `tfsdk:"name" query:"name,optional"`
	Details        types.Bool                                                               `tfsdk:"details" query:"details,computed_optional"`
	MaxItems       types.Int64                                                              `tfsdk:"max_items"`
	Items          customfield.NestedObjectList[CloudLoadBalancerPoolsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudLoadBalancerPoolsDataSourceModel) toListParams(_ context.Context) (params cloud.LoadBalancerPoolListParams, diags diag.Diagnostics) {
	params = cloud.LoadBalancerPoolListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.Details.IsNull() {
		params.Details = param.NewOpt(m.Details.ValueBool())
	}
	if !m.ListenerID.IsNull() {
		params.ListenerID = param.NewOpt(m.ListenerID.ValueString())
	}
	if !m.LoadBalancerID.IsNull() {
		params.LoadBalancerID = param.NewOpt(m.LoadBalancerID.ValueString())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}

	return
}

type CloudLoadBalancerPoolsItemsDataSourceModel struct {
	ID                   types.String                                                                      `tfsdk:"id" json:"id,computed"`
	AdminStateUp         types.Bool                                                                        `tfsdk:"admin_state_up" json:"admin_state_up,computed"`
	CaSecretID           types.String                                                                      `tfsdk:"ca_secret_id" json:"ca_secret_id,computed"`
	CreatorTaskID        types.String                                                                      `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	CrlSecretID          types.String                                                                      `tfsdk:"crl_secret_id" json:"crl_secret_id,computed"`
	Healthmonitor        customfield.NestedObject[CloudLoadBalancerPoolsHealthmonitorDataSourceModel]      `tfsdk:"healthmonitor" json:"healthmonitor,computed"`
	LbAlgorithm          types.String                                                                      `tfsdk:"lb_algorithm" json:"lb_algorithm,computed"`
	Listeners            customfield.NestedObjectList[CloudLoadBalancerPoolsListenersDataSourceModel]      `tfsdk:"listeners" json:"listeners,computed"`
	Loadbalancers        customfield.NestedObjectList[CloudLoadBalancerPoolsLoadbalancersDataSourceModel]  `tfsdk:"loadbalancers" json:"loadbalancers,computed"`
	Members              customfield.NestedObjectList[CloudLoadBalancerPoolsMembersDataSourceModel]        `tfsdk:"members" json:"members,computed"`
	Name                 types.String                                                                      `tfsdk:"name" json:"name,computed"`
	OperatingStatus      types.String                                                                      `tfsdk:"operating_status" json:"operating_status,computed"`
	Protocol             types.String                                                                      `tfsdk:"protocol" json:"protocol,computed"`
	ProvisioningStatus   types.String                                                                      `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	SecretID             types.String                                                                      `tfsdk:"secret_id" json:"secret_id,computed"`
	SessionPersistence   customfield.NestedObject[CloudLoadBalancerPoolsSessionPersistenceDataSourceModel] `tfsdk:"session_persistence" json:"session_persistence,computed"`
	TaskID               types.String                                                                      `tfsdk:"task_id" json:"task_id,computed"`
	TimeoutClientData    types.Int64                                                                       `tfsdk:"timeout_client_data" json:"timeout_client_data,computed"`
	TimeoutMemberConnect types.Int64                                                                       `tfsdk:"timeout_member_connect" json:"timeout_member_connect,computed"`
	TimeoutMemberData    types.Int64                                                                       `tfsdk:"timeout_member_data" json:"timeout_member_data,computed"`
}

type CloudLoadBalancerPoolsHealthmonitorDataSourceModel struct {
	ID                 types.String `tfsdk:"id" json:"id,computed"`
	AdminStateUp       types.Bool   `tfsdk:"admin_state_up" json:"admin_state_up,computed"`
	Delay              types.Int64  `tfsdk:"delay" json:"delay,computed"`
	DomainName         types.String `tfsdk:"domain_name" json:"domain_name,computed"`
	HTTPVersion        types.String `tfsdk:"http_version" json:"http_version,computed"`
	MaxRetries         types.Int64  `tfsdk:"max_retries" json:"max_retries,computed"`
	MaxRetriesDown     types.Int64  `tfsdk:"max_retries_down" json:"max_retries_down,computed"`
	OperatingStatus    types.String `tfsdk:"operating_status" json:"operating_status,computed"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	Timeout            types.Int64  `tfsdk:"timeout" json:"timeout,computed"`
	Type               types.String `tfsdk:"type" json:"type,computed"`
	ExpectedCodes      types.String `tfsdk:"expected_codes" json:"expected_codes,computed"`
	HTTPMethod         types.String `tfsdk:"http_method" json:"http_method,computed"`
	URLPath            types.String `tfsdk:"url_path" json:"url_path,computed"`
}

type CloudLoadBalancerPoolsListenersDataSourceModel struct {
	ID types.String `tfsdk:"id" json:"id,computed"`
}

type CloudLoadBalancerPoolsLoadbalancersDataSourceModel struct {
	ID types.String `tfsdk:"id" json:"id,computed"`
}

type CloudLoadBalancerPoolsMembersDataSourceModel struct {
	ID                 types.String `tfsdk:"id" json:"id,computed"`
	Address            types.String `tfsdk:"address" json:"address,computed"`
	AdminStateUp       types.Bool   `tfsdk:"admin_state_up" json:"admin_state_up,computed"`
	Backup             types.Bool   `tfsdk:"backup" json:"backup,computed"`
	OperatingStatus    types.String `tfsdk:"operating_status" json:"operating_status,computed"`
	ProtocolPort       types.Int64  `tfsdk:"protocol_port" json:"protocol_port,computed"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	SubnetID           types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
	Weight             types.Int64  `tfsdk:"weight" json:"weight,computed"`
	MonitorAddress     types.String `tfsdk:"monitor_address" json:"monitor_address,computed"`
	MonitorPort        types.Int64  `tfsdk:"monitor_port" json:"monitor_port,computed"`
}

type CloudLoadBalancerPoolsSessionPersistenceDataSourceModel struct {
	Type                   types.String `tfsdk:"type" json:"type,computed"`
	CookieName             types.String `tfsdk:"cookie_name" json:"cookie_name,computed"`
	PersistenceGranularity types.String `tfsdk:"persistence_granularity" json:"persistence_granularity,computed"`
	PersistenceTimeout     types.Int64  `tfsdk:"persistence_timeout" json:"persistence_timeout,computed"`
}
