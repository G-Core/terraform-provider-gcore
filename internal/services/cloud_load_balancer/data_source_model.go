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

type CloudLoadBalancerDataSourceModel struct {
	ID                    types.String                                                                 `tfsdk:"id" path:"loadbalancer_id,computed"`
	LoadbalancerID        types.String                                                                 `tfsdk:"loadbalancer_id" path:"loadbalancer_id,optional"`
	ProjectID             types.Int64                                                                  `tfsdk:"project_id" path:"project_id,optional"`
	RegionID              types.Int64                                                                  `tfsdk:"region_id" path:"region_id,optional"`
	ShowStats             types.Bool                                                                   `tfsdk:"show_stats" query:"show_stats,optional"`
	WithDDOS              types.Bool                                                                   `tfsdk:"with_ddos" query:"with_ddos,optional"`
	CreatedAt             timetypes.RFC3339                                                            `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID         types.String                                                                 `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Name                  types.String                                                                 `tfsdk:"name" json:"name,computed"`
	OperatingStatus       types.String                                                                 `tfsdk:"operating_status" json:"operating_status,computed"`
	PreferredConnectivity types.String                                                                 `tfsdk:"preferred_connectivity" json:"preferred_connectivity,computed"`
	ProvisioningStatus    types.String                                                                 `tfsdk:"provisioning_status" json:"provisioning_status,computed"`
	Region                types.String                                                                 `tfsdk:"region" json:"region,computed"`
	TaskID                types.String                                                                 `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt             timetypes.RFC3339                                                            `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
	VipAddress            types.String                                                                 `tfsdk:"vip_address" json:"vip_address,computed"`
	VipIPFamily           types.String                                                                 `tfsdk:"vip_ip_family" json:"vip_ip_family,computed"`
	VipPortID             types.String                                                                 `tfsdk:"vip_port_id" json:"vip_port_id,computed"`
	AdditionalVips        customfield.NestedObjectList[CloudLoadBalancerAdditionalVipsDataSourceModel] `tfsdk:"additional_vips" json:"additional_vips,computed"`
	DDOSProfile           customfield.NestedObject[CloudLoadBalancerDDOSProfileDataSourceModel]        `tfsdk:"ddos_profile" json:"ddos_profile,computed"`
	Flavor                customfield.NestedObject[CloudLoadBalancerFlavorDataSourceModel]             `tfsdk:"flavor" json:"flavor,computed"`
	FloatingIPs           customfield.NestedObjectList[CloudLoadBalancerFloatingIPsDataSourceModel]    `tfsdk:"floating_ips" json:"floating_ips,computed"`
	Listeners             customfield.NestedObjectList[CloudLoadBalancerListenersDataSourceModel]      `tfsdk:"listeners" json:"listeners,computed"`
	Logging               customfield.NestedObject[CloudLoadBalancerLoggingDataSourceModel]            `tfsdk:"logging" json:"logging,computed"`
	Stats                 customfield.NestedObject[CloudLoadBalancerStatsDataSourceModel]              `tfsdk:"stats" json:"stats,computed"`
	TagsV2                customfield.NestedObjectList[CloudLoadBalancerTagsV2DataSourceModel]         `tfsdk:"tags_v2" json:"tags_v2,computed"`
	VrrpIPs               customfield.NestedObjectList[CloudLoadBalancerVrrpIPsDataSourceModel]        `tfsdk:"vrrp_ips" json:"vrrp_ips,computed"`
	FindOneBy             *CloudLoadBalancerFindOneByDataSourceModel                                   `tfsdk:"find_one_by"`
}

func (m *CloudLoadBalancerDataSourceModel) toReadParams(_ context.Context) (params cloud.LoadBalancerGetParams, diags diag.Diagnostics) {
	params = cloud.LoadBalancerGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

func (m *CloudLoadBalancerDataSourceModel) toListParams(_ context.Context) (params cloud.LoadBalancerListParams, diags diag.Diagnostics) {
	mFindOneByTagKey := []string{}
	if m.FindOneBy.TagKey != nil {
		for _, item := range *m.FindOneBy.TagKey {
			mFindOneByTagKey = append(mFindOneByTagKey, item.ValueString())
		}
	}

	params = cloud.LoadBalancerListParams{
		TagKey: mFindOneByTagKey,
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.FindOneBy.AssignedFloating.IsNull() {
		params.AssignedFloating = param.NewOpt(m.FindOneBy.AssignedFloating.ValueBool())
	}
	if !m.FindOneBy.LoggingEnabled.IsNull() {
		params.LoggingEnabled = param.NewOpt(m.FindOneBy.LoggingEnabled.ValueBool())
	}
	if !m.FindOneBy.Name.IsNull() {
		params.Name = param.NewOpt(m.FindOneBy.Name.ValueString())
	}
	if !m.FindOneBy.OrderBy.IsNull() {
		params.OrderBy = param.NewOpt(m.FindOneBy.OrderBy.ValueString())
	}
	if !m.FindOneBy.ShowStats.IsNull() {
		params.ShowStats = param.NewOpt(m.FindOneBy.ShowStats.ValueBool())
	}
	if !m.FindOneBy.TagKeyValue.IsNull() {
		params.TagKeyValue = param.NewOpt(m.FindOneBy.TagKeyValue.ValueString())
	}
	if !m.FindOneBy.WithDDOS.IsNull() {
		params.WithDDOS = param.NewOpt(m.FindOneBy.WithDDOS.ValueBool())
	}

	return
}

type CloudLoadBalancerAdditionalVipsDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudLoadBalancerDDOSProfileDataSourceModel struct {
	ID                         types.Int64                                                                          `tfsdk:"id" json:"id,computed"`
	Fields                     customfield.NestedObjectList[CloudLoadBalancerDDOSProfileFieldsDataSourceModel]      `tfsdk:"fields" json:"fields,computed"`
	Options                    customfield.NestedObject[CloudLoadBalancerDDOSProfileOptionsDataSourceModel]         `tfsdk:"options" json:"options,computed"`
	ProfileTemplate            customfield.NestedObject[CloudLoadBalancerDDOSProfileProfileTemplateDataSourceModel] `tfsdk:"profile_template" json:"profile_template,computed"`
	ProfileTemplateDescription types.String                                                                         `tfsdk:"profile_template_description" json:"profile_template_description,computed"`
	Protocols                  customfield.NestedObjectList[CloudLoadBalancerDDOSProfileProtocolsDataSourceModel]   `tfsdk:"protocols" json:"protocols,computed"`
	Site                       types.String                                                                         `tfsdk:"site" json:"site,computed"`
	Status                     customfield.NestedObject[CloudLoadBalancerDDOSProfileStatusDataSourceModel]          `tfsdk:"status" json:"status,computed"`
}

type CloudLoadBalancerDDOSProfileFieldsDataSourceModel struct {
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

type CloudLoadBalancerDDOSProfileOptionsDataSourceModel struct {
	Active types.Bool `tfsdk:"active" json:"active,computed"`
	Bgp    types.Bool `tfsdk:"bgp" json:"bgp,computed"`
}

type CloudLoadBalancerDDOSProfileProfileTemplateDataSourceModel struct {
	ID          types.Int64                                                                                    `tfsdk:"id" json:"id,computed"`
	Description types.String                                                                                   `tfsdk:"description" json:"description,computed"`
	Fields      customfield.NestedObjectList[CloudLoadBalancerDDOSProfileProfileTemplateFieldsDataSourceModel] `tfsdk:"fields" json:"fields,computed"`
	Name        types.String                                                                                   `tfsdk:"name" json:"name,computed"`
}

type CloudLoadBalancerDDOSProfileProfileTemplateFieldsDataSourceModel struct {
	ID               types.Int64          `tfsdk:"id" json:"id,computed"`
	Default          types.String         `tfsdk:"default" json:"default,computed"`
	Description      types.String         `tfsdk:"description" json:"description,computed"`
	FieldType        types.String         `tfsdk:"field_type" json:"field_type,computed"`
	Name             types.String         `tfsdk:"name" json:"name,computed"`
	Required         types.Bool           `tfsdk:"required" json:"required,computed"`
	ValidationSchema jsontypes.Normalized `tfsdk:"validation_schema" json:"validation_schema,computed"`
}

type CloudLoadBalancerDDOSProfileProtocolsDataSourceModel struct {
	Port      types.String                   `tfsdk:"port" json:"port,computed"`
	Protocols customfield.List[types.String] `tfsdk:"protocols" json:"protocols,computed"`
}

type CloudLoadBalancerDDOSProfileStatusDataSourceModel struct {
	ErrorDescription types.String `tfsdk:"error_description" json:"error_description,computed"`
	Status           types.String `tfsdk:"status" json:"status,computed"`
}

type CloudLoadBalancerFlavorDataSourceModel struct {
	FlavorID   types.String `tfsdk:"flavor_id" json:"flavor_id,computed"`
	FlavorName types.String `tfsdk:"flavor_name" json:"flavor_name,computed"`
	Ram        types.Int64  `tfsdk:"ram" json:"ram,computed"`
	Vcpus      types.Int64  `tfsdk:"vcpus" json:"vcpus,computed"`
}

type CloudLoadBalancerFloatingIPsDataSourceModel struct {
	ID                types.String                                                                  `tfsdk:"id" json:"id,computed"`
	CreatedAt         timetypes.RFC3339                                                             `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID     types.String                                                                  `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	FixedIPAddress    types.String                                                                  `tfsdk:"fixed_ip_address" json:"fixed_ip_address,computed"`
	FloatingIPAddress types.String                                                                  `tfsdk:"floating_ip_address" json:"floating_ip_address,computed"`
	PortID            types.String                                                                  `tfsdk:"port_id" json:"port_id,computed"`
	ProjectID         types.Int64                                                                   `tfsdk:"project_id" json:"project_id,computed"`
	Region            types.String                                                                  `tfsdk:"region" json:"region,computed"`
	RegionID          types.Int64                                                                   `tfsdk:"region_id" json:"region_id,computed"`
	RouterID          types.String                                                                  `tfsdk:"router_id" json:"router_id,computed"`
	Status            types.String                                                                  `tfsdk:"status" json:"status,computed"`
	Tags              customfield.NestedObjectList[CloudLoadBalancerFloatingIPsTagsDataSourceModel] `tfsdk:"tags" json:"tags,computed"`
	TaskID            types.String                                                                  `tfsdk:"task_id" json:"task_id,computed"`
	UpdatedAt         timetypes.RFC3339                                                             `tfsdk:"updated_at" json:"updated_at,computed" format:"date-time"`
}

type CloudLoadBalancerFloatingIPsTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudLoadBalancerListenersDataSourceModel struct {
	ID types.String `tfsdk:"id" json:"id,computed"`
}

type CloudLoadBalancerLoggingDataSourceModel struct {
	DestinationRegionID types.Int64                                                                      `tfsdk:"destination_region_id" json:"destination_region_id,computed"`
	Enabled             types.Bool                                                                       `tfsdk:"enabled" json:"enabled,computed"`
	TopicName           types.String                                                                     `tfsdk:"topic_name" json:"topic_name,computed"`
	RetentionPolicy     customfield.NestedObject[CloudLoadBalancerLoggingRetentionPolicyDataSourceModel] `tfsdk:"retention_policy" json:"retention_policy,computed"`
}

type CloudLoadBalancerLoggingRetentionPolicyDataSourceModel struct {
	Period types.Int64 `tfsdk:"period" json:"period,computed"`
}

type CloudLoadBalancerStatsDataSourceModel struct {
	ActiveConnections types.Int64 `tfsdk:"active_connections" json:"active_connections,computed"`
	BytesIn           types.Int64 `tfsdk:"bytes_in" json:"bytes_in,computed"`
	BytesOut          types.Int64 `tfsdk:"bytes_out" json:"bytes_out,computed"`
	RequestErrors     types.Int64 `tfsdk:"request_errors" json:"request_errors,computed"`
	TotalConnections  types.Int64 `tfsdk:"total_connections" json:"total_connections,computed"`
}

type CloudLoadBalancerTagsV2DataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudLoadBalancerVrrpIPsDataSourceModel struct {
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	Role      types.String `tfsdk:"role" json:"role,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudLoadBalancerFindOneByDataSourceModel struct {
	AssignedFloating types.Bool      `tfsdk:"assigned_floating" query:"assigned_floating,optional"`
	LoggingEnabled   types.Bool      `tfsdk:"logging_enabled" query:"logging_enabled,optional"`
	Name             types.String    `tfsdk:"name" query:"name,optional"`
	OrderBy          types.String    `tfsdk:"order_by" query:"order_by,optional"`
	ShowStats        types.Bool      `tfsdk:"show_stats" query:"show_stats,optional"`
	TagKey           *[]types.String `tfsdk:"tag_key" query:"tag_key,optional"`
	TagKeyValue      types.String    `tfsdk:"tag_key_value" query:"tag_key_value,optional"`
	WithDDOS         types.Bool      `tfsdk:"with_ddos" query:"with_ddos,optional"`
}
