// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_server

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

type CloudBaremetalServersResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudBaremetalServersItemsDataSourceModel] `json:"results,computed"`
}

type CloudBaremetalServersDataSourceModel struct {
	ProjectID               types.Int64                                                             `tfsdk:"project_id" path:"project_id,optional"`
	RegionID                types.Int64                                                             `tfsdk:"region_id" path:"region_id,optional"`
	ChangesBefore           timetypes.RFC3339                                                       `tfsdk:"changes_before" query:"changes-before,optional" format:"date-time"`
	ChangesSince            timetypes.RFC3339                                                       `tfsdk:"changes_since" query:"changes-since,optional" format:"date-time"`
	FlavorID                types.String                                                            `tfsdk:"flavor_id" query:"flavor_id,optional"`
	FlavorPrefix            types.String                                                            `tfsdk:"flavor_prefix" query:"flavor_prefix,optional"`
	IP                      types.String                                                            `tfsdk:"ip" query:"ip,optional"`
	Name                    types.String                                                            `tfsdk:"name" query:"name,optional"`
	OnlyWithFixedExternalIP types.Bool                                                              `tfsdk:"only_with_fixed_external_ip" query:"only_with_fixed_external_ip,optional"`
	ProfileName             types.String                                                            `tfsdk:"profile_name" query:"profile_name,optional"`
	ProtectionStatus        types.String                                                            `tfsdk:"protection_status" query:"protection_status,optional"`
	Status                  types.String                                                            `tfsdk:"status" query:"status,optional"`
	TagKeyValue             types.String                                                            `tfsdk:"tag_key_value" query:"tag_key_value,optional"`
	TypeDDOSProfile         types.String                                                            `tfsdk:"type_ddos_profile" query:"type_ddos_profile,optional"`
	Uuid                    types.String                                                            `tfsdk:"uuid" query:"uuid,optional"`
	TagValue                *[]types.String                                                         `tfsdk:"tag_value" query:"tag_value,optional"`
	IncludeK8S              types.Bool                                                              `tfsdk:"include_k8s" query:"include_k8s,computed_optional"`
	OnlyIsolated            types.Bool                                                              `tfsdk:"only_isolated" query:"only_isolated,computed_optional"`
	OrderBy                 types.String                                                            `tfsdk:"order_by" query:"order_by,computed_optional"`
	WithDDOS                types.Bool                                                              `tfsdk:"with_ddos" query:"with_ddos,computed_optional"`
	WithInterfacesName      types.Bool                                                              `tfsdk:"with_interfaces_name" query:"with_interfaces_name,computed_optional"`
	MaxItems                types.Int64                                                             `tfsdk:"max_items"`
	Items                   customfield.NestedObjectList[CloudBaremetalServersItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudBaremetalServersDataSourceModel) toListParams(_ context.Context) (params cloud.BaremetalServerListParams, diags diag.Diagnostics) {
	mTagValue := []string{}
	if m.TagValue != nil {
		for _, item := range *m.TagValue {
			mTagValue = append(mTagValue, item.ValueString())
		}
	}
	mChangesBefore, errs := m.ChangesBefore.ValueRFC3339Time()
	diags.Append(errs...)
	mChangesSince, errs := m.ChangesSince.ValueRFC3339Time()
	diags.Append(errs...)

	params = cloud.BaremetalServerListParams{
		TagValue: mTagValue,
	}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}
	if !m.ChangesBefore.IsNull() {
		params.ChangesBefore = param.NewOpt(mChangesBefore)
	}
	if !m.ChangesSince.IsNull() {
		params.ChangesSince = param.NewOpt(mChangesSince)
	}
	if !m.FlavorID.IsNull() {
		params.FlavorID = param.NewOpt(m.FlavorID.ValueString())
	}
	if !m.FlavorPrefix.IsNull() {
		params.FlavorPrefix = param.NewOpt(m.FlavorPrefix.ValueString())
	}
	if !m.IncludeK8S.IsNull() {
		params.IncludeK8S = param.NewOpt(m.IncludeK8S.ValueBool())
	}
	if !m.IP.IsNull() {
		params.IP = param.NewOpt(m.IP.ValueString())
	}
	if !m.Name.IsNull() {
		params.Name = param.NewOpt(m.Name.ValueString())
	}
	if !m.OnlyIsolated.IsNull() {
		params.OnlyIsolated = param.NewOpt(m.OnlyIsolated.ValueBool())
	}
	if !m.OnlyWithFixedExternalIP.IsNull() {
		params.OnlyWithFixedExternalIP = param.NewOpt(m.OnlyWithFixedExternalIP.ValueBool())
	}
	if !m.OrderBy.IsNull() {
		params.OrderBy = cloud.BaremetalServerListParamsOrderBy(m.OrderBy.ValueString())
	}
	if !m.ProfileName.IsNull() {
		params.ProfileName = param.NewOpt(m.ProfileName.ValueString())
	}
	if !m.ProtectionStatus.IsNull() {
		params.ProtectionStatus = cloud.BaremetalServerListParamsProtectionStatus(m.ProtectionStatus.ValueString())
	}
	if !m.Status.IsNull() {
		params.Status = cloud.BaremetalServerListParamsStatus(m.Status.ValueString())
	}
	if !m.TagKeyValue.IsNull() {
		params.TagKeyValue = param.NewOpt(m.TagKeyValue.ValueString())
	}
	if !m.TypeDDOSProfile.IsNull() {
		params.TypeDDOSProfile = cloud.BaremetalServerListParamsTypeDDOSProfile(m.TypeDDOSProfile.ValueString())
	}
	if !m.Uuid.IsNull() {
		params.Uuid = param.NewOpt(m.Uuid.ValueString())
	}
	if !m.WithDDOS.IsNull() {
		params.WithDDOS = param.NewOpt(m.WithDDOS.ValueBool())
	}
	if !m.WithInterfacesName.IsNull() {
		params.WithInterfacesName = param.NewOpt(m.WithInterfacesName.ValueBool())
	}

	return
}

type CloudBaremetalServersItemsDataSourceModel struct {
	ID                 types.String                                                                                 `tfsdk:"id" json:"id,computed"`
	Addresses          customfield.Map[customfield.NestedObjectList[CloudBaremetalServersAddressesDataSourceModel]] `tfsdk:"addresses" json:"addresses,computed"`
	BlackholePorts     customfield.NestedObjectList[CloudBaremetalServersBlackholePortsDataSourceModel]             `tfsdk:"blackhole_ports" json:"blackhole_ports,computed"`
	CreatedAt          timetypes.RFC3339                                                                            `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID      types.String                                                                                 `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	DDOSProfile        customfield.NestedObject[CloudBaremetalServersDDOSProfileDataSourceModel]                    `tfsdk:"ddos_profile" json:"ddos_profile,computed"`
	FixedIPAssignments customfield.NestedObjectList[CloudBaremetalServersFixedIPAssignmentsDataSourceModel]         `tfsdk:"fixed_ip_assignments" json:"fixed_ip_assignments,computed"`
	Flavor             customfield.NestedObject[CloudBaremetalServersFlavorDataSourceModel]                         `tfsdk:"flavor" json:"flavor,computed"`
	InstanceIsolation  customfield.NestedObject[CloudBaremetalServersInstanceIsolationDataSourceModel]              `tfsdk:"instance_isolation" json:"instance_isolation,computed"`
	Name               types.String                                                                                 `tfsdk:"name" json:"name,computed"`
	ProjectID          types.Int64                                                                                  `tfsdk:"project_id" json:"project_id,computed"`
	Region             types.String                                                                                 `tfsdk:"region" json:"region,computed"`
	RegionID           types.Int64                                                                                  `tfsdk:"region_id" json:"region_id,computed"`
	SSHKeyName         types.String                                                                                 `tfsdk:"ssh_key_name" json:"ssh_key_name,computed"`
	Status             types.String                                                                                 `tfsdk:"status" json:"status,computed"`
	Tags               customfield.NestedObjectList[CloudBaremetalServersTagsDataSourceModel]                       `tfsdk:"tags" json:"tags,computed"`
	TaskID             types.String                                                                                 `tfsdk:"task_id" json:"task_id,computed"`
	TaskState          types.String                                                                                 `tfsdk:"task_state" json:"task_state,computed"`
	VmState            types.String                                                                                 `tfsdk:"vm_state" json:"vm_state,computed"`
}

type CloudBaremetalServersAddressesDataSourceModel struct {
	Addr          types.String `tfsdk:"addr" json:"addr,required"`
	Type          types.String `tfsdk:"type" json:"type,required"`
	InterfaceName types.String `tfsdk:"interface_name" json:"interface_name,optional"`
	SubnetID      types.String `tfsdk:"subnet_id" json:"subnet_id,optional"`
	SubnetName    types.String `tfsdk:"subnet_name" json:"subnet_name,optional"`
}

type CloudBaremetalServersBlackholePortsDataSourceModel struct {
	AlarmEnd      timetypes.RFC3339 `tfsdk:"alarm_end" json:"AlarmEnd,computed" format:"date-time"`
	AlarmStart    timetypes.RFC3339 `tfsdk:"alarm_start" json:"AlarmStart,computed" format:"date-time"`
	AlarmState    types.String      `tfsdk:"alarm_state" json:"AlarmState,computed"`
	AlertDuration types.String      `tfsdk:"alert_duration" json:"AlertDuration,computed"`
	DestinationIP types.String      `tfsdk:"destination_ip" json:"DestinationIP,computed"`
	ID            types.Int64       `tfsdk:"id" json:"ID,computed"`
}

type CloudBaremetalServersDDOSProfileDataSourceModel struct {
	ID                         types.Int64                                                                              `tfsdk:"id" json:"id,computed"`
	Fields                     customfield.NestedObjectList[CloudBaremetalServersDDOSProfileFieldsDataSourceModel]      `tfsdk:"fields" json:"fields,computed"`
	Options                    customfield.NestedObject[CloudBaremetalServersDDOSProfileOptionsDataSourceModel]         `tfsdk:"options" json:"options,computed"`
	ProfileTemplate            customfield.NestedObject[CloudBaremetalServersDDOSProfileProfileTemplateDataSourceModel] `tfsdk:"profile_template" json:"profile_template,computed"`
	ProfileTemplateDescription types.String                                                                             `tfsdk:"profile_template_description" json:"profile_template_description,computed"`
	Protocols                  customfield.NestedObjectList[CloudBaremetalServersDDOSProfileProtocolsDataSourceModel]   `tfsdk:"protocols" json:"protocols,computed"`
	Site                       types.String                                                                             `tfsdk:"site" json:"site,computed"`
	Status                     customfield.NestedObject[CloudBaremetalServersDDOSProfileStatusDataSourceModel]          `tfsdk:"status" json:"status,computed"`
}

type CloudBaremetalServersDDOSProfileFieldsDataSourceModel struct {
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

type CloudBaremetalServersDDOSProfileOptionsDataSourceModel struct {
	Active types.Bool `tfsdk:"active" json:"active,computed"`
	Bgp    types.Bool `tfsdk:"bgp" json:"bgp,computed"`
}

type CloudBaremetalServersDDOSProfileProfileTemplateDataSourceModel struct {
	ID          types.Int64                                                                                        `tfsdk:"id" json:"id,computed"`
	Description types.String                                                                                       `tfsdk:"description" json:"description,computed"`
	Fields      customfield.NestedObjectList[CloudBaremetalServersDDOSProfileProfileTemplateFieldsDataSourceModel] `tfsdk:"fields" json:"fields,computed"`
	Name        types.String                                                                                       `tfsdk:"name" json:"name,computed"`
}

type CloudBaremetalServersDDOSProfileProfileTemplateFieldsDataSourceModel struct {
	ID               types.Int64          `tfsdk:"id" json:"id,computed"`
	Default          types.String         `tfsdk:"default" json:"default,computed"`
	Description      types.String         `tfsdk:"description" json:"description,computed"`
	FieldType        types.String         `tfsdk:"field_type" json:"field_type,computed"`
	Name             types.String         `tfsdk:"name" json:"name,computed"`
	Required         types.Bool           `tfsdk:"required" json:"required,computed"`
	ValidationSchema jsontypes.Normalized `tfsdk:"validation_schema" json:"validation_schema,computed"`
}

type CloudBaremetalServersDDOSProfileProtocolsDataSourceModel struct {
	Port      types.String                   `tfsdk:"port" json:"port,computed"`
	Protocols customfield.List[types.String] `tfsdk:"protocols" json:"protocols,computed"`
}

type CloudBaremetalServersDDOSProfileStatusDataSourceModel struct {
	ErrorDescription types.String `tfsdk:"error_description" json:"error_description,computed"`
	Status           types.String `tfsdk:"status" json:"status,computed"`
}

type CloudBaremetalServersFixedIPAssignmentsDataSourceModel struct {
	External  types.Bool   `tfsdk:"external" json:"external,computed"`
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudBaremetalServersFlavorDataSourceModel struct {
	Architecture        types.String                                                                            `tfsdk:"architecture" json:"architecture,computed"`
	FlavorID            types.String                                                                            `tfsdk:"flavor_id" json:"flavor_id,computed"`
	FlavorName          types.String                                                                            `tfsdk:"flavor_name" json:"flavor_name,computed"`
	HardwareDescription customfield.NestedObject[CloudBaremetalServersFlavorHardwareDescriptionDataSourceModel] `tfsdk:"hardware_description" json:"hardware_description,computed"`
	OsType              types.String                                                                            `tfsdk:"os_type" json:"os_type,computed"`
	Ram                 types.Int64                                                                             `tfsdk:"ram" json:"ram,computed"`
	ResourceClass       types.String                                                                            `tfsdk:"resource_class" json:"resource_class,computed"`
	Vcpus               types.Int64                                                                             `tfsdk:"vcpus" json:"vcpus,computed"`
}

type CloudBaremetalServersFlavorHardwareDescriptionDataSourceModel struct {
	CPU     types.String `tfsdk:"cpu" json:"cpu,computed"`
	Disk    types.String `tfsdk:"disk" json:"disk,computed"`
	License types.String `tfsdk:"license" json:"license,computed"`
	Network types.String `tfsdk:"network" json:"network,computed"`
	Ram     types.String `tfsdk:"ram" json:"ram,computed"`
}

type CloudBaremetalServersInstanceIsolationDataSourceModel struct {
	Reason types.String `tfsdk:"reason" json:"reason,computed"`
}

type CloudBaremetalServersTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}
