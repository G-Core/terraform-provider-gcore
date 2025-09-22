// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudLoadBalancerPoolModel struct {
	ID                   types.String                                                          `tfsdk:"id" json:"id,computed"`
	ProjectID            types.Int64                                                           `tfsdk:"project_id" path:"project_id,optional"`
	RegionID             types.Int64                                                           `tfsdk:"region_id" path:"region_id,optional"`
	ListenerID           types.String                                                          `tfsdk:"listener_id" json:"listener_id,optional,no_refresh"`
	LoadbalancerID       types.String                                                          `tfsdk:"loadbalancer_id" json:"loadbalancer_id,optional,no_refresh"`
	LbAlgorithm          types.String                                                          `tfsdk:"lb_algorithm" json:"lb_algorithm,required"`
	Name                 types.String                                                          `tfsdk:"name" json:"name,required"`
	Protocol             types.String                                                          `tfsdk:"protocol" json:"protocol,required"`
	CaSecretID           types.String                                                          `tfsdk:"ca_secret_id" json:"ca_secret_id,optional"`
	CrlSecretID          types.String                                                          `tfsdk:"crl_secret_id" json:"crl_secret_id,optional"`
	SecretID             types.String                                                          `tfsdk:"secret_id" json:"secret_id,optional"`
	TimeoutClientData    types.Int64                                                           `tfsdk:"timeout_client_data" json:"timeout_client_data,optional"`
	TimeoutMemberConnect types.Int64                                                           `tfsdk:"timeout_member_connect" json:"timeout_member_connect,optional"`
	TimeoutMemberData    types.Int64                                                           `tfsdk:"timeout_member_data" json:"timeout_member_data,optional"`
	Healthmonitor        *CloudLoadBalancerPoolHealthmonitorModel                              `tfsdk:"healthmonitor" json:"healthmonitor,optional"`
	SessionPersistence   *CloudLoadBalancerPoolSessionPersistenceModel                         `tfsdk:"session_persistence" json:"session_persistence,optional"`
	Members              customfield.NestedObjectList[CloudLoadBalancerPoolMembersModel]       `tfsdk:"members" json:"members,computed_optional"`
	CreatorTaskID        types.String                                                          `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	OperatingStatus      types.String                                                          `tfsdk:"operating_status" json:"operating_status,computed"`
	ProvisioningStatus   types.String                                                          `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	TaskID               types.String                                                          `tfsdk:"task_id" json:"task_id,computed"`
	Tasks                customfield.List[types.String]                                        `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
	Listeners            customfield.NestedObjectList[CloudLoadBalancerPoolListenersModel]     `tfsdk:"listeners" json:"listeners,computed"`
	Loadbalancers        customfield.NestedObjectList[CloudLoadBalancerPoolLoadbalancersModel] `tfsdk:"loadbalancers" json:"loadbalancers,computed"`
}

func (m CloudLoadBalancerPoolModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudLoadBalancerPoolModel) MarshalJSONForUpdate(state CloudLoadBalancerPoolModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudLoadBalancerPoolHealthmonitorModel struct {
	Delay          types.Int64  `tfsdk:"delay" json:"delay,required"`
	MaxRetries     types.Int64  `tfsdk:"max_retries" json:"max_retries,required"`
	Timeout        types.Int64  `tfsdk:"timeout" json:"timeout,required"`
	Type           types.String `tfsdk:"type" json:"type,required"`
	ExpectedCodes  types.String `tfsdk:"expected_codes" json:"expected_codes,optional"`
	HTTPMethod     types.String `tfsdk:"http_method" json:"http_method,optional"`
	MaxRetriesDown types.Int64  `tfsdk:"max_retries_down" json:"max_retries_down,optional"`
	URLPath        types.String `tfsdk:"url_path" json:"url_path,optional"`
}

type CloudLoadBalancerPoolSessionPersistenceModel struct {
	Type                   types.String `tfsdk:"type" json:"type,required"`
	CookieName             types.String `tfsdk:"cookie_name" json:"cookie_name,optional"`
	PersistenceGranularity types.String `tfsdk:"persistence_granularity" json:"persistence_granularity,optional"`
	PersistenceTimeout     types.Int64  `tfsdk:"persistence_timeout" json:"persistence_timeout,optional"`
}

type CloudLoadBalancerPoolMembersModel struct {
	Address        types.String `tfsdk:"address" json:"address,required"`
	ProtocolPort   types.Int64  `tfsdk:"protocol_port" json:"protocol_port,required"`
	AdminStateUp   types.Bool   `tfsdk:"admin_state_up" json:"admin_state_up,computed_optional"`
	Backup         types.Bool   `tfsdk:"backup" json:"backup,computed_optional"`
	InstanceID     types.String `tfsdk:"instance_id" json:"instance_id,optional,no_refresh"`
	MonitorAddress types.String `tfsdk:"monitor_address" json:"monitor_address,optional"`
	MonitorPort    types.Int64  `tfsdk:"monitor_port" json:"monitor_port,optional"`
	SubnetID       types.String `tfsdk:"subnet_id" json:"subnet_id,optional"`
	Weight         types.Int64  `tfsdk:"weight" json:"weight,optional"`
}

type CloudLoadBalancerPoolListenersModel struct {
	ID types.String `tfsdk:"id" json:"id,computed"`
}

type CloudLoadBalancerPoolLoadbalancersModel struct {
	ID types.String `tfsdk:"id" json:"id,computed"`
}
