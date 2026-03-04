// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudLoadBalancerModel struct {
	ID                    types.String                                                       `tfsdk:"id" json:"id,computed"`
	ProjectID             types.Int64                                                        `tfsdk:"project_id" path:"project_id,optional"`
	RegionID              types.Int64                                                        `tfsdk:"region_id" path:"region_id,optional"`
	Flavor                types.String                                                       `tfsdk:"flavor" json:"flavor,optional,no_refresh"`
	NameTemplate          types.String                                                       `tfsdk:"name_template" json:"name_template,optional,no_refresh"`
	VipNetworkID          types.String                                                       `tfsdk:"vip_network_id" json:"vip_network_id,optional,no_refresh"`
	VipSubnetID           types.String                                                       `tfsdk:"vip_subnet_id" json:"vip_subnet_id,optional,no_refresh"`
	FloatingIP            *CloudLoadBalancerFloatingIPModel                                  `tfsdk:"floating_ip" json:"floating_ip,optional,no_refresh"`
	VipIPFamily           types.String                                                       `tfsdk:"vip_ip_family" json:"vip_ip_family,computed_optional"`
	VipPortID             types.String                                                       `tfsdk:"vip_port_id" json:"vip_port_id,computed_optional"`
	Name                  types.String                                                       `tfsdk:"name" json:"name,optional"`
	Tags                  *map[string]types.String                                           `tfsdk:"tags" json:"tags,optional,no_refresh"`
	PreferredConnectivity types.String                                                       `tfsdk:"preferred_connectivity" json:"preferred_connectivity,computed_optional"`
	Logging               customfield.NestedObject[CloudLoadBalancerLoggingModel]            `tfsdk:"logging" json:"logging,computed_optional"`
	AdminStateUp          types.Bool                                                         `tfsdk:"admin_state_up" json:"admin_state_up,computed"`
	CreatedAt             timetypes.RFC3339                                                  `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID         types.String                                                       `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	OperatingStatus       types.String                                                       `tfsdk:"operating_status" json:"operating_status,computed"`
	ProvisioningStatus    types.String                                                       `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	Region                types.String                                                       `tfsdk:"region" json:"region,computed"`
	UpdatedAt             timetypes.RFC3339                                                  `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	VipAddress            types.String                                                       `tfsdk:"vip_address" json:"vip_address,computed"`
	VipFqdn               types.String                                                       `tfsdk:"vip_fqdn" json:"vip_fqdn,computed"`
	AdditionalVips        customfield.NestedObjectList[CloudLoadBalancerAdditionalVipsModel] `tfsdk:"additional_vips" json:"additional_vips,computed"`
	DDOSProfile           customfield.NestedObject[CloudLoadBalancerDDOSProfileModel]        `tfsdk:"ddos_profile" json:"ddos_profile,computed"`
	FloatingIPs           customfield.NestedObjectList[CloudLoadBalancerFloatingIPsModel]    `tfsdk:"floating_ips" json:"floating_ips,computed"`
	Stats                 customfield.NestedObject[CloudLoadBalancerStatsModel]              `tfsdk:"stats" json:"stats,computed"`
	TagsV2                customfield.NestedObjectList[CloudLoadBalancerTagsV2Model]         `tfsdk:"tags_v2" json:"tags_v2,computed"`
	VrrpIPs               customfield.NestedObjectList[CloudLoadBalancerVrrpIPsModel]        `tfsdk:"vrrp_ips" json:"vrrp_ips,computed"`
}

func (m CloudLoadBalancerModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudLoadBalancerModel) MarshalJSONForUpdate(state CloudLoadBalancerModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudLoadBalancerFloatingIPModel struct {
	Source             types.String `tfsdk:"source" json:"source,required"`
	ExistingFloatingID types.String `tfsdk:"existing_floating_id" json:"existing_floating_id,optional"`
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

type CloudLoadBalancerAdditionalVipsModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudLoadBalancerDDOSProfileModel struct {
	ID                         types.Int64                                                                `tfsdk:"id" json:"id,computed"`
	Fields                     customfield.NestedObjectList[CloudLoadBalancerDDOSProfileFieldsModel]      `tfsdk:"fields" json:"fields,computed"`
	Options                    customfield.NestedObject[CloudLoadBalancerDDOSProfileOptionsModel]         `tfsdk:"options" json:"options,computed"`
	ProfileTemplate            customfield.NestedObject[CloudLoadBalancerDDOSProfileProfileTemplateModel] `tfsdk:"profile_template" json:"profile_template,computed"`
	ProfileTemplateDescription types.String                                                               `tfsdk:"profile_template_description" json:"profile_template_description,computed"`
	Protocols                  customfield.NestedObjectList[CloudLoadBalancerDDOSProfileProtocolsModel]   `tfsdk:"protocols" json:"protocols,computed"`
	Site                       types.String                                                               `tfsdk:"site" json:"site,computed"`
	Status                     customfield.NestedObject[CloudLoadBalancerDDOSProfileStatusModel]          `tfsdk:"status" json:"status,computed"`
}

type CloudLoadBalancerDDOSProfileFieldsModel struct {
	ID               types.Int64          `tfsdk:"id" json:"id,computed"`
	BaseField        types.Int64          `tfsdk:"base_field" json:"base_field,computed"`
	Default          types.String         `tfsdk:"default" json:"default,computed"`
	Description      types.String         `tfsdk:"description" json:"description,computed"`
	FieldName        types.String         `tfsdk:"field_name" json:"field_name,computed"`
	FieldType        types.String         `tfsdk:"field_type" json:"field_type,computed"`
	FieldValue       jsontypes.Normalized `tfsdk:"field_value" json:"field_value,computed"`
	Name             types.String         `tfsdk:"name" json:"name,computed"`
	Required         types.Bool           `tfsdk:"required" json:"required,computed"`
	ValidationSchema jsontypes.Normalized `tfsdk:"validation_schema" json:"validation_schema,computed"`
	Value            types.String         `tfsdk:"value" json:"value,computed"`
}

type CloudLoadBalancerDDOSProfileOptionsModel struct {
	Active types.Bool `tfsdk:"active" json:"active,computed"`
	Bgp    types.Bool `tfsdk:"bgp" json:"bgp,computed"`
}

type CloudLoadBalancerDDOSProfileProfileTemplateModel struct {
	ID          types.Int64                                                                          `tfsdk:"id" json:"id,computed"`
	Description types.String                                                                         `tfsdk:"description" json:"description,computed"`
	Fields      customfield.NestedObjectList[CloudLoadBalancerDDOSProfileProfileTemplateFieldsModel] `tfsdk:"fields" json:"fields,computed"`
	Name        types.String                                                                         `tfsdk:"name" json:"name,computed"`
}

type CloudLoadBalancerDDOSProfileProfileTemplateFieldsModel struct {
	ID               types.Int64          `tfsdk:"id" json:"id,computed"`
	Default          types.String         `tfsdk:"default" json:"default,computed"`
	Description      types.String         `tfsdk:"description" json:"description,computed"`
	FieldType        types.String         `tfsdk:"field_type" json:"field_type,computed"`
	Name             types.String         `tfsdk:"name" json:"name,computed"`
	Required         types.Bool           `tfsdk:"required" json:"required,computed"`
	ValidationSchema jsontypes.Normalized `tfsdk:"validation_schema" json:"validation_schema,computed"`
}

type CloudLoadBalancerDDOSProfileProtocolsModel struct {
	Port      types.String                   `tfsdk:"port" json:"port,computed"`
	Protocols customfield.List[types.String] `tfsdk:"protocols" json:"protocols,computed"`
}

type CloudLoadBalancerDDOSProfileStatusModel struct {
	ErrorDescription types.String `tfsdk:"error_description" json:"error_description,computed"`
	Status           types.String `tfsdk:"status" json:"status,computed"`
}

type CloudLoadBalancerFloatingIPsModel struct {
	ID                types.String                                                        `tfsdk:"id" json:"id,computed"`
	CreatedAt         timetypes.RFC3339                                                   `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID     types.String                                                        `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	FixedIPAddress    types.String                                                        `tfsdk:"fixed_ip_address" json:"fixed_ip_address,computed"`
	FloatingIPAddress types.String                                                        `tfsdk:"floating_ip_address" json:"floating_ip_address,computed"`
	PortID            types.String                                                        `tfsdk:"port_id" json:"port_id,computed"`
	ProjectID         types.Int64                                                         `tfsdk:"project_id" json:"project_id,computed"`
	Region            types.String                                                        `tfsdk:"region" json:"region,computed"`
	RegionID          types.Int64                                                         `tfsdk:"region_id" json:"region_id,computed"`
	RouterID          types.String                                                        `tfsdk:"router_id" json:"router_id,computed"`
	Status            types.String                                                        `tfsdk:"status" json:"status,computed"`
	Tags              customfield.NestedObjectList[CloudLoadBalancerFloatingIPsTagsModel] `tfsdk:"tags" json:"tags,computed"`
	UpdatedAt         timetypes.RFC3339                                                   `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudLoadBalancerFloatingIPsTagsModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudLoadBalancerStatsModel struct {
	ActiveConnections types.Int64 `tfsdk:"active_connections" json:"active_connections,computed"`
	BytesIn           types.Int64 `tfsdk:"bytes_in" json:"bytes_in,computed"`
	BytesOut          types.Int64 `tfsdk:"bytes_out" json:"bytes_out,computed"`
	RequestErrors     types.Int64 `tfsdk:"request_errors" json:"request_errors,computed"`
	TotalConnections  types.Int64 `tfsdk:"total_connections" json:"total_connections,computed"`
}

type CloudLoadBalancerTagsV2Model struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudLoadBalancerVrrpIPsModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	Role      types.String `tfsdk:"role" json:"role,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}
