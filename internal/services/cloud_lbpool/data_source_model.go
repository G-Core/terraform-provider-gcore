// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_lbpool

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudLbpoolDataSourceModel struct {
	PoolID               types.String                                                           `tfsdk:"pool_id" path:"pool_id,required"`
	ProjectID            types.Int64                                                            `tfsdk:"project_id" path:"project_id,required"`
	RegionID             types.Int64                                                            `tfsdk:"region_id" path:"region_id,required"`
	CaSecretID           types.String                                                           `tfsdk:"ca_secret_id" json:"ca_secret_id,computed"`
	CreatorTaskID        types.String                                                           `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	CrlSecretID          types.String                                                           `tfsdk:"crl_secret_id" json:"crl_secret_id,computed"`
	ID                   types.String                                                           `tfsdk:"id" json:"id,computed"`
	LbAlgorithm          types.String                                                           `tfsdk:"lb_algorithm" json:"lb_algorithm,computed"`
	Name                 types.String                                                           `tfsdk:"name" json:"name,computed"`
	OperatingStatus      types.String                                                           `tfsdk:"operating_status" json:"operating_status,computed"`
	Protocol             types.String                                                           `tfsdk:"protocol" json:"protocol,computed"`
	ProvisioningStatus   types.String                                                           `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	SecretID             types.String                                                           `tfsdk:"secret_id" json:"secret_id,computed"`
	TaskID               types.String                                                           `tfsdk:"task_id" json:"task_id,computed"`
	TimeoutClientData    types.Int64                                                            `tfsdk:"timeout_client_data" json:"timeout_client_data,computed"`
	TimeoutMemberConnect types.Int64                                                            `tfsdk:"timeout_member_connect" json:"timeout_member_connect,computed"`
	TimeoutMemberData    types.Int64                                                            `tfsdk:"timeout_member_data" json:"timeout_member_data,computed"`
	Healthmonitor        customfield.NestedObject[CloudLbpoolHealthmonitorDataSourceModel]      `tfsdk:"healthmonitor" json:"healthmonitor,computed"`
	Listeners            customfield.NestedObjectList[CloudLbpoolListenersDataSourceModel]      `tfsdk:"listeners" json:"listeners,computed"`
	Loadbalancers        customfield.NestedObjectList[CloudLbpoolLoadbalancersDataSourceModel]  `tfsdk:"loadbalancers" json:"loadbalancers,computed"`
	Members              customfield.NestedObjectList[CloudLbpoolMembersDataSourceModel]        `tfsdk:"members" json:"members,computed"`
	SessionPersistence   customfield.NestedObject[CloudLbpoolSessionPersistenceDataSourceModel] `tfsdk:"session_persistence" json:"session_persistence,computed"`
}

func (m *CloudLbpoolDataSourceModel) toReadParams(_ context.Context) (params cloud.LoadBalancerPoolGetParams, diags diag.Diagnostics) {
	params = cloud.LoadBalancerPoolGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudLbpoolHealthmonitorDataSourceModel struct {
	ID                 types.String `tfsdk:"id" json:"id,computed"`
	AdminStateUp       types.Bool   `tfsdk:"admin_state_up" json:"admin_state_up,computed"`
	Delay              types.Int64  `tfsdk:"delay" json:"delay,computed"`
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

type CloudLbpoolListenersDataSourceModel struct {
	ID types.String `tfsdk:"id" json:"id,computed"`
}

type CloudLbpoolLoadbalancersDataSourceModel struct {
	ID types.String `tfsdk:"id" json:"id,computed"`
}

type CloudLbpoolMembersDataSourceModel struct {
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

type CloudLbpoolSessionPersistenceDataSourceModel struct {
	Type                   types.String `tfsdk:"type" json:"type,computed"`
	CookieName             types.String `tfsdk:"cookie_name" json:"cookie_name,computed"`
	PersistenceGranularity types.String `tfsdk:"persistence_granularity" json:"persistence_granularity,computed"`
	PersistenceTimeout     types.Int64  `tfsdk:"persistence_timeout" json:"persistence_timeout,computed"`
}
