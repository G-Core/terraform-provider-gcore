// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudLoadBalancersResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudLoadBalancersItemsDataSourceModel] `json:"results,computed"`
}

type CloudLoadBalancersDataSourceModel struct {
	ProjectID        types.Int64                                                          `tfsdk:"project_id" path:"project_id,optional"`
	RegionID         types.Int64                                                          `tfsdk:"region_id" path:"region_id,optional"`
	AssignedFloating types.Bool                                                           `tfsdk:"assigned_floating" query:"assigned_floating,optional"`
	LoggingEnabled   types.Bool                                                           `tfsdk:"logging_enabled" query:"logging_enabled,optional"`
	Name             types.String                                                         `tfsdk:"name" query:"name,optional"`
	TagKeyValue      types.String                                                         `tfsdk:"tag_key_value" query:"tag_key_value,optional"`
	TagKey           *[]types.String                                                      `tfsdk:"tag_key" query:"tag_key,optional"`
	OrderBy          types.String                                                         `tfsdk:"order_by" query:"order_by,computed_optional"`
	ShowStats        types.Bool                                                           `tfsdk:"show_stats" query:"show_stats,computed_optional"`
	WithDDOS         types.Bool                                                           `tfsdk:"with_ddos" query:"with_ddos,computed_optional"`
	MaxItems         types.Int64                                                          `tfsdk:"max_items"`
	Items            customfield.NestedObjectList[CloudLoadBalancersItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudLoadBalancersDataSourceModel) toListParams(_ context.Context) (params cloud.LoadBalancerListParams, diags diag.Diagnostics) {
	mTagKey := []string{}
	if m.TagKey != nil {
		for _, item := range *m.TagKey {
			mTagKey = append(mTagKey, item.ValueString())
		}
	}

	params = cloud.LoadBalancerListParams{
		TagKey: mTagKey,
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.AssignedFloating.IsNull() {
		params.AssignedFloating = param.NewOpt(m.AssignedFloating.ValueBool())
	}
	if !m.LoggingEnabled.IsNull() {
		params.LoggingEnabled = param.NewOpt(m.LoggingEnabled.ValueBool())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = cloud.LoadBalancerListParamsOrderBy(m.OrderBy.ValueString())
	}
	if !m.ShowStats.IsNull() {
		params.ShowStats = param.NewOpt(m.ShowStats.ValueBool())
	}
	if !m.TagKeyValue.IsNull() {
		params.TagKeyValue = param.NewOpt(m.TagKeyValue.ValueString())
	}
	if !m.WithDDOS.IsNull() {
		params.WithDDOS = param.NewOpt(m.WithDDOS.ValueBool())
	}

	return
}

type CloudLoadBalancersItemsDataSourceModel struct {
	ID                    types.String                                                                  `tfsdk:"id" json:"id,computed"`
	AdminStateUp          types.Bool                                                                    `tfsdk:"admin_state_up" json:"admin_state_up,computed"`
	CreatedAt             timetypes.RFC3339                                                             `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Name                  types.String                                                                  `tfsdk:"name" json:"name,computed"`
	OperatingStatus       types.String                                                                  `tfsdk:"operating_status" json:"operating_status,computed"`
	ProjectID             types.Int64                                                                   `tfsdk:"project_id" json:"project_id,computed"`
	ProvisioningStatus    types.String                                                                  `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	Region                types.String                                                                  `tfsdk:"region" json:"region,computed"`
	RegionID              types.Int64                                                                   `tfsdk:"region_id" json:"region_id,computed"`
	TagsV2                customfield.NestedObjectList[CloudLoadBalancersTagsV2DataSourceModel]         `tfsdk:"tags_v2" json:"tags_v2,computed"`
	AdditionalVips        customfield.NestedObjectList[CloudLoadBalancersAdditionalVipsDataSourceModel] `tfsdk:"additional_vips" json:"additional_vips,computed"`
	CreatorTaskID         types.String                                                                  `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	DDOSProfile           customfield.NestedObject[CloudLoadBalancersDDOSProfileDataSourceModel]        `tfsdk:"ddos_profile" json:"ddos_profile,computed"`
	Flavor                customfield.NestedObject[CloudLoadBalancersFlavorDataSourceModel]             `tfsdk:"flavor" json:"flavor,computed"`
	FloatingIPs           customfield.NestedObjectList[CloudLoadBalancersFloatingIPsDataSourceModel]    `tfsdk:"floating_ips" json:"floating_ips,computed"`
	Listeners             customfield.NestedObjectList[CloudLoadBalancersListenersDataSourceModel]      `tfsdk:"listeners" json:"listeners,computed"`
	Logging               customfield.NestedObject[CloudLoadBalancersLoggingDataSourceModel]            `tfsdk:"logging" json:"logging,computed"`
	PreferredConnectivity types.String                                                                  `tfsdk:"preferred_connectivity" json:"preferred_connectivity,computed"`
	Stats                 customfield.NestedObject[CloudLoadBalancersStatsDataSourceModel]              `tfsdk:"stats" json:"stats,computed"`
	TaskID                types.String                                                                  `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt             timetypes.RFC3339                                                             `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	VipAddress            types.String                                                                  `tfsdk:"vip_address" json:"vip_address,computed"`
	VipIPFamily           types.String                                                                  `tfsdk:"vip_ip_family" json:"vip_ip_family,computed"`
	VipPortID             types.String                                                                  `tfsdk:"vip_port_id" json:"vip_port_id,computed"`
	VrrpIPs               customfield.NestedObjectList[CloudLoadBalancersVrrpIPsDataSourceModel]        `tfsdk:"vrrp_ips" json:"vrrp_ips,computed"`
}

type CloudLoadBalancersTagsV2DataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudLoadBalancersAdditionalVipsDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudLoadBalancersDDOSProfileDataSourceModel struct {
	ID                         types.Int64                                                                           `tfsdk:"id" json:"id,computed"`
	Fields                     customfield.NestedObjectList[CloudLoadBalancersDDOSProfileFieldsDataSourceModel]      `tfsdk:"fields" json:"fields,computed"`
	Options                    customfield.NestedObject[CloudLoadBalancersDDOSProfileOptionsDataSourceModel]         `tfsdk:"options" json:"options,computed"`
	ProfileTemplate            customfield.NestedObject[CloudLoadBalancersDDOSProfileProfileTemplateDataSourceModel] `tfsdk:"profile_template" json:"profile_template,computed"`
	ProfileTemplateDescription types.String                                                                          `tfsdk:"profile_template_description" json:"profile_template_description,computed"`
	Protocols                  customfield.NestedObjectList[CloudLoadBalancersDDOSProfileProtocolsDataSourceModel]   `tfsdk:"protocols" json:"protocols,computed"`
	Site                       types.String                                                                          `tfsdk:"site" json:"site,computed"`
	Status                     customfield.NestedObject[CloudLoadBalancersDDOSProfileStatusDataSourceModel]          `tfsdk:"status" json:"status,computed"`
}

type CloudLoadBalancersDDOSProfileFieldsDataSourceModel struct {
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

type CloudLoadBalancersDDOSProfileOptionsDataSourceModel struct {
	Active types.Bool `tfsdk:"active" json:"active,computed"`
	Bgp    types.Bool `tfsdk:"bgp" json:"bgp,computed"`
}

type CloudLoadBalancersDDOSProfileProfileTemplateDataSourceModel struct {
	ID          types.Int64                                                                                     `tfsdk:"id" json:"id,computed"`
	Description types.String                                                                                    `tfsdk:"description" json:"description,computed"`
	Fields      customfield.NestedObjectList[CloudLoadBalancersDDOSProfileProfileTemplateFieldsDataSourceModel] `tfsdk:"fields" json:"fields,computed"`
	Name        types.String                                                                                    `tfsdk:"name" json:"name,computed"`
}

type CloudLoadBalancersDDOSProfileProfileTemplateFieldsDataSourceModel struct {
	ID               types.Int64          `tfsdk:"id" json:"id,computed"`
	Default          types.String         `tfsdk:"default" json:"default,computed"`
	Description      types.String         `tfsdk:"description" json:"description,computed"`
	FieldType        types.String         `tfsdk:"field_type" json:"field_type,computed"`
	Name             types.String         `tfsdk:"name" json:"name,computed"`
	Required         types.Bool           `tfsdk:"required" json:"required,computed"`
	ValidationSchema jsontypes.Normalized `tfsdk:"validation_schema" json:"validation_schema,computed"`
}

type CloudLoadBalancersDDOSProfileProtocolsDataSourceModel struct {
	Port      types.String                   `tfsdk:"port" json:"port,computed"`
	Protocols customfield.List[types.String] `tfsdk:"protocols" json:"protocols,computed"`
}

type CloudLoadBalancersDDOSProfileStatusDataSourceModel struct {
	ErrorDescription types.String `tfsdk:"error_description" json:"error_description,computed"`
	Status           types.String `tfsdk:"status" json:"status,computed"`
}

type CloudLoadBalancersFlavorDataSourceModel struct {
	FlavorID   types.String `tfsdk:"flavor_id" json:"flavor_id,computed"`
	FlavorName types.String `tfsdk:"flavor_name" json:"flavor_name,computed"`
	Ram        types.Int64  `tfsdk:"ram" json:"ram,computed"`
	Vcpus      types.Int64  `tfsdk:"vcpus" json:"vcpus,computed"`
}

type CloudLoadBalancersFloatingIPsDataSourceModel struct {
	ID                types.String                                                                   `tfsdk:"id" json:"id,computed"`
	CreatedAt         timetypes.RFC3339                                                              `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID     types.String                                                                   `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	FixedIPAddress    types.String                                                                   `tfsdk:"fixed_ip_address" json:"fixed_ip_address,computed"`
	FloatingIPAddress types.String                                                                   `tfsdk:"floating_ip_address" json:"floating_ip_address,computed"`
	PortID            types.String                                                                   `tfsdk:"port_id" json:"port_id,computed"`
	ProjectID         types.Int64                                                                    `tfsdk:"project_id" json:"project_id,computed"`
	Region            types.String                                                                   `tfsdk:"region" json:"region,computed"`
	RegionID          types.Int64                                                                    `tfsdk:"region_id" json:"region_id,computed"`
	RouterID          types.String                                                                   `tfsdk:"router_id" json:"router_id,computed"`
	Status            types.String                                                                   `tfsdk:"status" json:"status,computed"`
	Tags              customfield.NestedObjectList[CloudLoadBalancersFloatingIPsTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	TaskID            types.String                                                                   `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt         timetypes.RFC3339                                                              `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudLoadBalancersFloatingIPsTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudLoadBalancersListenersDataSourceModel struct {
	ID types.String `tfsdk:"id" json:"id,computed"`
}

type CloudLoadBalancersLoggingDataSourceModel struct {
	DestinationRegionID types.Int64                                                                       `tfsdk:"destination_region_id" json:"destination_region_id,computed"`
	Enabled             types.Bool                                                                        `tfsdk:"enabled" json:"enabled,computed"`
	TopicName           types.String                                                                      `tfsdk:"topic_name" json:"topic_name,computed"`
	RetentionPolicy     customfield.NestedObject[CloudLoadBalancersLoggingRetentionPolicyDataSourceModel] `tfsdk:"retention_policy" json:"retention_policy,computed"`
}

type CloudLoadBalancersLoggingRetentionPolicyDataSourceModel struct {
	Period types.Int64 `tfsdk:"period" json:"period,computed"`
}

type CloudLoadBalancersStatsDataSourceModel struct {
	ActiveConnections types.Int64 `tfsdk:"active_connections" json:"active_connections,computed"`
	BytesIn           types.Int64 `tfsdk:"bytes_in" json:"bytes_in,computed"`
	BytesOut          types.Int64 `tfsdk:"bytes_out" json:"bytes_out,computed"`
	RequestErrors     types.Int64 `tfsdk:"request_errors" json:"request_errors,computed"`
	TotalConnections  types.Int64 `tfsdk:"total_connections" json:"total_connections,computed"`
}

type CloudLoadBalancersVrrpIPsDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	Role      types.String `tfsdk:"role" json:"role,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}
