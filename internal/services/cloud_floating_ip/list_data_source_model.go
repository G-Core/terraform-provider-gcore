// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_floating_ip

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudFloatingIPsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudFloatingIPsItemsDataSourceModel] `json:"results,computed"`
}

type CloudFloatingIPsDataSourceModel struct {
	ProjectID   types.Int64                                                        `tfsdk:"project_id" path:"project_id,optional"`
	RegionID    types.Int64                                                        `tfsdk:"region_id" path:"region_id,optional"`
	Status      types.String                                                       `tfsdk:"status" query:"status,optional"`
	TagKeyValue types.String                                                       `tfsdk:"tag_key_value" query:"tag_key_value,optional"`
	TagKey      *[]types.String                                                    `tfsdk:"tag_key" query:"tag_key,optional"`
	Limit       types.Int64                                                        `tfsdk:"limit" query:"limit,computed_optional"`
	MaxItems    types.Int64                                                        `tfsdk:"max_items"`
	Items       customfield.NestedObjectList[CloudFloatingIPsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudFloatingIPsDataSourceModel) toListParams(_ context.Context) (params cloud.FloatingIPListParams, diags diag.Diagnostics) {
	mTagKey := []string{}
	if m.TagKey != nil {
		for _, item := range *m.TagKey {
			mTagKey = append(mTagKey, item.ValueString())
		}
	}

	params = cloud.FloatingIPListParams{
		TagKey: mTagKey,
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.Limit.IsNull() {
		params.Limit = param.NewOpt(m.Limit.ValueInt64())
	}
	if !m.Status.IsNull() {
		params.Status = cloud.FloatingIPStatus(m.Status.ValueString())
	}
	if !m.TagKeyValue.IsNull() {
		params.TagKeyValue = param.NewOpt(m.TagKeyValue.ValueString())
	}

	return
}

type CloudFloatingIPsItemsDataSourceModel struct {
	ID                types.String                                                          `tfsdk:"id" json:"id,computed"`
	CreatedAt         timetypes.RFC3339                                                     `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID     types.String                                                          `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	FixedIPAddress    types.String                                                          `tfsdk:"fixed_ip_address" json:"fixed_ip_address,computed"`
	FloatingIPAddress types.String                                                          `tfsdk:"floating_ip_address" json:"floating_ip_address,computed"`
	Instance          customfield.NestedObject[CloudFloatingIPsInstanceDataSourceModel]     `tfsdk:"instance" json:"instance,computed"`
	Loadbalancer      customfield.NestedObject[CloudFloatingIPsLoadbalancerDataSourceModel] `tfsdk:"loadbalancer" json:"loadbalancer,computed"`
	PortID            types.String                                                          `tfsdk:"port_id" json:"port_id,computed"`
	ProjectID         types.Int64                                                           `tfsdk:"project_id" json:"project_id,computed"`
	Region            types.String                                                          `tfsdk:"region" json:"region,computed"`
	RegionID          types.Int64                                                           `tfsdk:"region_id" json:"region_id,computed"`
	RouterID          types.String                                                          `tfsdk:"router_id" json:"router_id,computed"`
	Status            types.String                                                          `tfsdk:"status" json:"status,computed"`
	Tags              customfield.NestedObjectList[CloudFloatingIPsTagsDataSourceModel]     `tfsdk:"tags" json:"tags,computed"`
	UpdatedAt         timetypes.RFC3339                                                     `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudFloatingIPsInstanceDataSourceModel struct {
	ID                  types.String                                                                                    `tfsdk:"id" json:"id,computed"`
	Addresses           customfield.Map[customfield.NestedObjectList[CloudFloatingIPsInstanceAddressesDataSourceModel]] `tfsdk:"addresses" json:"addresses,computed"`
	CreatedAt           timetypes.RFC3339                                                                               `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                                                    `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Flavor              customfield.NestedObject[CloudFloatingIPsInstanceFlavorDataSourceModel]                         `tfsdk:"flavor" json:"flavor,computed"`
	InstanceDescription types.String                                                                                    `tfsdk:"instance_description" json:"instance_description,computed"`
	Name                types.String                                                                                    `tfsdk:"name" json:"name,computed"`
	ProjectID           types.Int64                                                                                     `tfsdk:"project_id" json:"project_id,computed"`
	Region              types.String                                                                                    `tfsdk:"region" json:"region,computed"`
	RegionID            types.Int64                                                                                     `tfsdk:"region_id" json:"region_id,computed"`
	SecurityGroups      customfield.NestedObjectList[CloudFloatingIPsInstanceSecurityGroupsDataSourceModel]             `tfsdk:"security_groups" json:"security_groups,computed"`
	SSHKeyName          types.String                                                                                    `tfsdk:"ssh_key_name" json:"ssh_key_name,computed"`
	Status              types.String                                                                                    `tfsdk:"status" json:"status,computed"`
	Tags                customfield.NestedObjectList[CloudFloatingIPsInstanceTagsDataSourceModel]                       `tfsdk:"tags" json:"tags,computed"`
	TaskState           types.String                                                                                    `tfsdk:"task_state" json:"task_state,computed"`
	VmState             types.String                                                                                    `tfsdk:"vm_state" json:"vm_state,computed"`
	Volumes             customfield.NestedObjectList[CloudFloatingIPsInstanceVolumesDataSourceModel]                    `tfsdk:"volumes" json:"volumes,computed"`
}

type CloudFloatingIPsInstanceAddressesDataSourceModel struct {
	Addr          types.String `tfsdk:"addr" json:"addr,required"`
	Type          types.String `tfsdk:"type" json:"type,required"`
	InterfaceName types.String `tfsdk:"interface_name" json:"interface_name,optional"`
	SubnetID      types.String `tfsdk:"subnet_id" json:"subnet_id,optional"`
	SubnetName    types.String `tfsdk:"subnet_name" json:"subnet_name,optional"`
}

type CloudFloatingIPsInstanceFlavorDataSourceModel struct {
	FlavorID   types.String `tfsdk:"flavor_id" json:"flavor_id,computed"`
	FlavorName types.String `tfsdk:"flavor_name" json:"flavor_name,computed"`
	Ram        types.Int64  `tfsdk:"ram" json:"ram,computed"`
	Vcpus      types.Int64  `tfsdk:"vcpus" json:"vcpus,computed"`
}

type CloudFloatingIPsInstanceSecurityGroupsDataSourceModel struct {
	Name types.String `tfsdk:"name" json:"name,computed"`
}

type CloudFloatingIPsInstanceTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudFloatingIPsInstanceVolumesDataSourceModel struct {
	ID                  types.String `tfsdk:"id" json:"id,computed"`
	DeleteOnTermination types.Bool   `tfsdk:"delete_on_termination" json:"delete_on_termination,computed"`
}

type CloudFloatingIPsLoadbalancerDataSourceModel struct {
	ID                    types.String                                                                            `tfsdk:"id" json:"id,computed"`
	AdminStateUp          types.Bool                                                                              `tfsdk:"admin_state_up" json:"admin_state_up,computed"`
	CreatedAt             timetypes.RFC3339                                                                       `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Name                  types.String                                                                            `tfsdk:"name" json:"name,computed"`
	OperatingStatus       types.String                                                                            `tfsdk:"operating_status" json:"operating_status,computed"`
	ProjectID             types.Int64                                                                             `tfsdk:"project_id" json:"project_id,computed"`
	ProvisioningStatus    types.String                                                                            `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	Region                types.String                                                                            `tfsdk:"region" json:"region,computed"`
	RegionID              types.Int64                                                                             `tfsdk:"region_id" json:"region_id,computed"`
	TagsV2                customfield.NestedObjectList[CloudFloatingIPsLoadbalancerTagsV2DataSourceModel]         `tfsdk:"tags_v2" json:"tags_v2,computed"`
	AdditionalVips        customfield.NestedObjectList[CloudFloatingIPsLoadbalancerAdditionalVipsDataSourceModel] `tfsdk:"additional_vips" json:"additional_vips,computed"`
	CreatorTaskID         types.String                                                                            `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	DDOSProfile           customfield.NestedObject[CloudFloatingIPsLoadbalancerDDOSProfileDataSourceModel]        `tfsdk:"ddos_profile" json:"ddos_profile,computed"`
	Flavor                customfield.NestedObject[CloudFloatingIPsLoadbalancerFlavorDataSourceModel]             `tfsdk:"flavor" json:"flavor,computed"`
	FloatingIPs           customfield.NestedObjectList[CloudFloatingIPsLoadbalancerFloatingIPsDataSourceModel]    `tfsdk:"floating_ips" json:"floating_ips,computed"`
	Listeners             customfield.NestedObjectList[CloudFloatingIPsLoadbalancerListenersDataSourceModel]      `tfsdk:"listeners" json:"listeners,computed"`
	Logging               customfield.NestedObject[CloudFloatingIPsLoadbalancerLoggingDataSourceModel]            `tfsdk:"logging" json:"logging,computed"`
	PreferredConnectivity types.String                                                                            `tfsdk:"preferred_connectivity" json:"preferred_connectivity,computed"`
	Stats                 customfield.NestedObject[CloudFloatingIPsLoadbalancerStatsDataSourceModel]              `tfsdk:"stats" json:"stats,computed"`
	UpdatedAt             timetypes.RFC3339                                                                       `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	VipAddress            types.String                                                                            `tfsdk:"vip_address" json:"vip_address,computed"`
	VipFqdn               types.String                                                                            `tfsdk:"vip_fqdn" json:"vip_fqdn,computed"`
	VipIPFamily           types.String                                                                            `tfsdk:"vip_ip_family" json:"vip_ip_family,computed"`
	VipPortID             types.String                                                                            `tfsdk:"vip_port_id" json:"vip_port_id,computed"`
	VrrpIPs               customfield.NestedObjectList[CloudFloatingIPsLoadbalancerVrrpIPsDataSourceModel]        `tfsdk:"vrrp_ips" json:"vrrp_ips,computed"`
}

type CloudFloatingIPsLoadbalancerTagsV2DataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudFloatingIPsLoadbalancerAdditionalVipsDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudFloatingIPsLoadbalancerDDOSProfileDataSourceModel struct {
	ID                         types.Int64                                                                                     `tfsdk:"id" json:"id,computed"`
	Fields                     customfield.NestedObjectList[CloudFloatingIPsLoadbalancerDDOSProfileFieldsDataSourceModel]      `tfsdk:"fields" json:"fields,computed"`
	Options                    customfield.NestedObject[CloudFloatingIPsLoadbalancerDDOSProfileOptionsDataSourceModel]         `tfsdk:"options" json:"options,computed"`
	ProfileTemplate            customfield.NestedObject[CloudFloatingIPsLoadbalancerDDOSProfileProfileTemplateDataSourceModel] `tfsdk:"profile_template" json:"profile_template,computed"`
	ProfileTemplateDescription types.String                                                                                    `tfsdk:"profile_template_description" json:"profile_template_description,computed"`
	Protocols                  customfield.NestedObjectList[CloudFloatingIPsLoadbalancerDDOSProfileProtocolsDataSourceModel]   `tfsdk:"protocols" json:"protocols,computed"`
	Site                       types.String                                                                                    `tfsdk:"site" json:"site,computed"`
	Status                     customfield.NestedObject[CloudFloatingIPsLoadbalancerDDOSProfileStatusDataSourceModel]          `tfsdk:"status" json:"status,computed"`
}

type CloudFloatingIPsLoadbalancerDDOSProfileFieldsDataSourceModel struct {
	ID               types.Int64          `tfsdk:"id" json:"id,computed"`
	BaseField        types.Int64          `tfsdk:"base_field" json:"base_field,computed"`
	Default          types.String         `tfsdk:"default" json:"default,computed"`
	Description      types.String         `tfsdk:"description" json:"description,computed"`
	FieldType        types.String         `tfsdk:"field_type" json:"field_type,computed"`
	FieldValue       jsontypes.Normalized `tfsdk:"field_value" json:"field_value,computed"`
	Name             types.String         `tfsdk:"name" json:"name,computed"`
	Required         types.Bool           `tfsdk:"required" json:"required,computed"`
	ValidationSchema jsontypes.Normalized `tfsdk:"validation_schema" json:"validation_schema,computed"`
}

type CloudFloatingIPsLoadbalancerDDOSProfileOptionsDataSourceModel struct {
	Active types.Bool `tfsdk:"active" json:"active,computed"`
	Bgp    types.Bool `tfsdk:"bgp" json:"bgp,computed"`
}

type CloudFloatingIPsLoadbalancerDDOSProfileProfileTemplateDataSourceModel struct {
	ID          types.Int64                                                                                               `tfsdk:"id" json:"id,computed"`
	Description types.String                                                                                              `tfsdk:"description" json:"description,computed"`
	Fields      customfield.NestedObjectList[CloudFloatingIPsLoadbalancerDDOSProfileProfileTemplateFieldsDataSourceModel] `tfsdk:"fields" json:"fields,computed"`
	Name        types.String                                                                                              `tfsdk:"name" json:"name,computed"`
}

type CloudFloatingIPsLoadbalancerDDOSProfileProfileTemplateFieldsDataSourceModel struct {
	ID               types.Int64          `tfsdk:"id" json:"id,computed"`
	Default          types.String         `tfsdk:"default" json:"default,computed"`
	Description      types.String         `tfsdk:"description" json:"description,computed"`
	FieldType        types.String         `tfsdk:"field_type" json:"field_type,computed"`
	Name             types.String         `tfsdk:"name" json:"name,computed"`
	Required         types.Bool           `tfsdk:"required" json:"required,computed"`
	ValidationSchema jsontypes.Normalized `tfsdk:"validation_schema" json:"validation_schema,computed"`
}

type CloudFloatingIPsLoadbalancerDDOSProfileProtocolsDataSourceModel struct {
	Port      types.String                   `tfsdk:"port" json:"port,computed"`
	Protocols customfield.List[types.String] `tfsdk:"protocols" json:"protocols,computed"`
}

type CloudFloatingIPsLoadbalancerDDOSProfileStatusDataSourceModel struct {
	ErrorDescription types.String `tfsdk:"error_description" json:"error_description,computed"`
	Status           types.String `tfsdk:"status" json:"status,computed"`
}

type CloudFloatingIPsLoadbalancerFlavorDataSourceModel struct {
	FlavorID   types.String `tfsdk:"flavor_id" json:"flavor_id,computed"`
	FlavorName types.String `tfsdk:"flavor_name" json:"flavor_name,computed"`
	Ram        types.Int64  `tfsdk:"ram" json:"ram,computed"`
	Vcpus      types.Int64  `tfsdk:"vcpus" json:"vcpus,computed"`
}

type CloudFloatingIPsLoadbalancerFloatingIPsDataSourceModel struct {
	ID                types.String                                                                             `tfsdk:"id" json:"id,computed"`
	CreatedAt         timetypes.RFC3339                                                                        `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID     types.String                                                                             `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	FixedIPAddress    types.String                                                                             `tfsdk:"fixed_ip_address" json:"fixed_ip_address,computed"`
	FloatingIPAddress types.String                                                                             `tfsdk:"floating_ip_address" json:"floating_ip_address,computed"`
	PortID            types.String                                                                             `tfsdk:"port_id" json:"port_id,computed"`
	ProjectID         types.Int64                                                                              `tfsdk:"project_id" json:"project_id,computed"`
	Region            types.String                                                                             `tfsdk:"region" json:"region,computed"`
	RegionID          types.Int64                                                                              `tfsdk:"region_id" json:"region_id,computed"`
	RouterID          types.String                                                                             `tfsdk:"router_id" json:"router_id,computed"`
	Status            types.String                                                                             `tfsdk:"status" json:"status,computed"`
	Tags              customfield.NestedObjectList[CloudFloatingIPsLoadbalancerFloatingIPsTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	UpdatedAt         timetypes.RFC3339                                                                        `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudFloatingIPsLoadbalancerFloatingIPsTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudFloatingIPsLoadbalancerListenersDataSourceModel struct {
	ID types.String `tfsdk:"id" json:"id,computed"`
}

type CloudFloatingIPsLoadbalancerLoggingDataSourceModel struct {
	DestinationRegionID types.Int64                                                                                 `tfsdk:"destination_region_id" json:"destination_region_id,computed"`
	Enabled             types.Bool                                                                                  `tfsdk:"enabled" json:"enabled,computed"`
	TopicName           types.String                                                                                `tfsdk:"topic_name" json:"topic_name,computed"`
	RetentionPolicy     customfield.NestedObject[CloudFloatingIPsLoadbalancerLoggingRetentionPolicyDataSourceModel] `tfsdk:"retention_policy" json:"retention_policy,computed"`
}

type CloudFloatingIPsLoadbalancerLoggingRetentionPolicyDataSourceModel struct {
	Period types.Int64 `tfsdk:"period" json:"period,computed"`
}

type CloudFloatingIPsLoadbalancerStatsDataSourceModel struct {
	ActiveConnections types.Int64 `tfsdk:"active_connections" json:"active_connections,computed"`
	BytesIn           types.Int64 `tfsdk:"bytes_in" json:"bytes_in,computed"`
	BytesOut          types.Int64 `tfsdk:"bytes_out" json:"bytes_out,computed"`
	RequestErrors     types.Int64 `tfsdk:"request_errors" json:"request_errors,computed"`
	TotalConnections  types.Int64 `tfsdk:"total_connections" json:"total_connections,computed"`
}

type CloudFloatingIPsLoadbalancerVrrpIPsDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	Role      types.String `tfsdk:"role" json:"role,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudFloatingIPsTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}
