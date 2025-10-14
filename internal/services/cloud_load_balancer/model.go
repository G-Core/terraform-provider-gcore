// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudLoadBalancerModel struct {
	ID                    types.String                                                  `tfsdk:"id" json:"id,computed"`
	ProjectID             types.Int64                                                   `tfsdk:"project_id" path:"project_id,optional"`
	RegionID              types.Int64                                                   `tfsdk:"region_id" path:"region_id,optional"`
	Flavor                types.String                                                  `tfsdk:"flavor" json:"flavor,optional"`
	Name                  types.String                                                  `tfsdk:"name" json:"name,optional"`
	NameTemplate          types.String                                                  `tfsdk:"name_template" json:"name_template,optional"`
	VipNetworkID          types.String                                                  `tfsdk:"vip_network_id" json:"vip_network_id,optional"`
	VipSubnetID           types.String                                                  `tfsdk:"vip_subnet_id" json:"vip_subnet_id,optional"`
	Tags                  *map[string]types.String                                      `tfsdk:"tags" json:"tags,optional"`
	FloatingIP            *CloudLoadBalancerFloatingIPModel                             `tfsdk:"floating_ip" json:"floating_ip,optional"`
	PreferredConnectivity types.String                                                  `tfsdk:"preferred_connectivity" json:"preferred_connectivity,computed_optional"`
	VipIPFamily           types.String                                                  `tfsdk:"vip_ip_family" json:"vip_ip_family,computed_optional"`
	VipPortID             types.String                                                  `tfsdk:"vip_port_id" json:"vip_port_id,computed_optional"`
	Listeners             customfield.NestedObjectList[CloudLoadBalancerListenersModel] `tfsdk:"listeners" json:"listeners,computed_optional"`
	Logging               customfield.NestedObject[CloudLoadBalancerLoggingModel]       `tfsdk:"logging" json:"logging,computed_optional"`
	Tasks                 customfield.List[types.String]                                `tfsdk:"tasks" json:"tasks,computed"`
}

func (m CloudLoadBalancerModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudLoadBalancerModel) MarshalJSONForUpdate(state CloudLoadBalancerModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

type CloudLoadBalancerFloatingIPModel struct {
	Source             types.String `tfsdk:"source" json:"source,required"`
	ExistingFloatingID types.String `tfsdk:"existing_floating_id" json:"existing_floating_id,optional"`
}

type CloudLoadBalancerListenersModel struct {
	Name                 types.String                                                       `tfsdk:"name" json:"name,required"`
	Protocol             types.String                                                       `tfsdk:"protocol" json:"protocol,required"`
	ProtocolPort         types.Int64                                                        `tfsdk:"protocol_port" json:"protocol_port,required"`
	AllowedCidrs         *[]types.String                                                    `tfsdk:"allowed_cidrs" json:"allowed_cidrs,optional"`
	ConnectionLimit      types.Int64                                                        `tfsdk:"connection_limit" json:"connection_limit,computed_optional"`
	InsertXForwarded     types.Bool                                                         `tfsdk:"insert_x_forwarded" json:"insert_x_forwarded,optional"`
	Pools                customfield.NestedObjectList[CloudLoadBalancerListenersPoolsModel] `tfsdk:"pools" json:"pools,computed_optional"`
	SecretID             types.String                                                       `tfsdk:"secret_id" json:"secret_id,optional"`
	SniSecretID          *[]types.String                                                    `tfsdk:"sni_secret_id" json:"sni_secret_id,optional"`
	TimeoutClientData    types.Int64                                                        `tfsdk:"timeout_client_data" json:"timeout_client_data,optional"`
	TimeoutMemberConnect types.Int64                                                        `tfsdk:"timeout_member_connect" json:"timeout_member_connect,optional"`
	TimeoutMemberData    types.Int64                                                        `tfsdk:"timeout_member_data" json:"timeout_member_data,optional"`
	UserList             *[]*CloudLoadBalancerListenersUserListModel                        `tfsdk:"user_list" json:"user_list,optional"`
}

type CloudLoadBalancerListenersPoolsModel struct {
	LbAlgorithm          types.String                                                              `tfsdk:"lb_algorithm" json:"lb_algorithm,required"`
	Name                 types.String                                                              `tfsdk:"name" json:"name,required"`
	Protocol             types.String                                                              `tfsdk:"protocol" json:"protocol,required"`
	CaSecretID           types.String                                                              `tfsdk:"ca_secret_id" json:"ca_secret_id,optional"`
	CrlSecretID          types.String                                                              `tfsdk:"crl_secret_id" json:"crl_secret_id,optional"`
	Healthmonitor        *CloudLoadBalancerListenersPoolsHealthmonitorModel                        `tfsdk:"healthmonitor" json:"healthmonitor,optional"`
	ListenerID           types.String                                                              `tfsdk:"listener_id" json:"listener_id,optional"`
	LoadbalancerID       types.String                                                              `tfsdk:"loadbalancer_id" json:"loadbalancer_id,optional"`
	Members              customfield.NestedObjectList[CloudLoadBalancerListenersPoolsMembersModel] `tfsdk:"members" json:"members,computed_optional"`
	SecretID             types.String                                                              `tfsdk:"secret_id" json:"secret_id,optional"`
	SessionPersistence   *CloudLoadBalancerListenersPoolsSessionPersistenceModel                   `tfsdk:"session_persistence" json:"session_persistence,optional"`
	TimeoutClientData    types.Int64                                                               `tfsdk:"timeout_client_data" json:"timeout_client_data,optional"`
	TimeoutMemberConnect types.Int64                                                               `tfsdk:"timeout_member_connect" json:"timeout_member_connect,optional"`
	TimeoutMemberData    types.Int64                                                               `tfsdk:"timeout_member_data" json:"timeout_member_data,optional"`
}

type CloudLoadBalancerListenersPoolsHealthmonitorModel struct {
	Delay          types.Int64  `tfsdk:"delay" json:"delay,required"`
	MaxRetries     types.Int64  `tfsdk:"max_retries" json:"max_retries,required"`
	Timeout        types.Int64  `tfsdk:"timeout" json:"timeout,required"`
	Type           types.String `tfsdk:"type" json:"type,required"`
	ExpectedCodes  types.String `tfsdk:"expected_codes" json:"expected_codes,optional"`
	HTTPMethod     types.String `tfsdk:"http_method" json:"http_method,optional"`
	MaxRetriesDown types.Int64  `tfsdk:"max_retries_down" json:"max_retries_down,optional"`
	URLPath        types.String `tfsdk:"url_path" json:"url_path,optional"`
}

type CloudLoadBalancerListenersPoolsMembersModel struct {
	Address        types.String `tfsdk:"address" json:"address,required"`
	ProtocolPort   types.Int64  `tfsdk:"protocol_port" json:"protocol_port,required"`
	AdminStateUp   types.Bool   `tfsdk:"admin_state_up" json:"admin_state_up,computed_optional"`
	Backup         types.Bool   `tfsdk:"backup" json:"backup,computed_optional"`
	InstanceID     types.String `tfsdk:"instance_id" json:"instance_id,optional"`
	MonitorAddress types.String `tfsdk:"monitor_address" json:"monitor_address,optional"`
	MonitorPort    types.Int64  `tfsdk:"monitor_port" json:"monitor_port,optional"`
	SubnetID       types.String `tfsdk:"subnet_id" json:"subnet_id,optional"`
	Weight         types.Int64  `tfsdk:"weight" json:"weight,optional"`
}

type CloudLoadBalancerListenersPoolsSessionPersistenceModel struct {
	Type                   types.String `tfsdk:"type" json:"type,required"`
	CookieName             types.String `tfsdk:"cookie_name" json:"cookie_name,optional"`
	PersistenceGranularity types.String `tfsdk:"persistence_granularity" json:"persistence_granularity,optional"`
	PersistenceTimeout     types.Int64  `tfsdk:"persistence_timeout" json:"persistence_timeout,optional"`
}

type CloudLoadBalancerListenersUserListModel struct {
	EncryptedPassword types.String `tfsdk:"encrypted_password" json:"encrypted_password,required"`
	Username          types.String `tfsdk:"username" json:"username,required"`
}

type CloudLoadBalancerLoggingModel struct {
	DestinationRegionID types.Int64                                   `tfsdk:"destination_region_id" json:"destination_region_id,optional"`
	Enabled             types.Bool                                    `tfsdk:"enabled" json:"enabled,computed_optional"`
	RetentionPolicy     *CloudLoadBalancerLoggingRetentionPolicyModel `tfsdk:"retention_policy" json:"retention_policy,optional"`
	TopicName           types.String                                  `tfsdk:"topic_name" json:"topic_name,optional"`
}

type CloudLoadBalancerLoggingRetentionPolicyModel struct {
	Period types.Int64 `tfsdk:"period" json:"period,required"`
}
