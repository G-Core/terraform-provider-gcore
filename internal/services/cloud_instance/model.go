// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance

import (
	"encoding/json"
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudInstanceModel struct {
	ID                  types.String                                                               `tfsdk:"id" json:"id,computed"`
	ProjectID           types.Int64                                                                `tfsdk:"project_id" path:"project_id,optional"`
	RegionID            types.Int64                                                                `tfsdk:"region_id" path:"region_id,optional"`
	Flavor              types.String                                                               `tfsdk:"flavor" json:"flavor,required"`
	Interfaces          *[]*CloudInstanceInterfacesModel                                           `tfsdk:"interfaces" json:"interfaces,required"`
	Volumes             *[]*CloudInstanceVolumesModel                                              `tfsdk:"volumes" json:"volumes,required"`
	AllowAppPorts       types.Bool                                                                 `tfsdk:"allow_app_ports" json:"allow_app_ports,optional,no_refresh"`
	NameTemplate        types.String                                                               `tfsdk:"name_template" json:"name_template,optional,no_refresh"`
	Password            types.String                                                               `tfsdk:"password_wo" json:"password,optional,no_refresh"`
	PasswordWoVersion   types.Int64                                                                `tfsdk:"password_wo_version"`
	ServergroupID       types.String                                                               `tfsdk:"servergroup_id" json:"servergroup_id,optional,no_refresh"`
	SSHKeyName          types.String                                                               `tfsdk:"ssh_key_name" json:"ssh_key_name,optional"`
	UserData            types.String                                                               `tfsdk:"user_data" json:"user_data,optional,no_refresh"`
	Username            types.String                                                               `tfsdk:"username" json:"username,optional,no_refresh"`
	Configuration       *map[string]jsontypes.Normalized                                           `tfsdk:"configuration" json:"configuration,optional,no_refresh"`
	SecurityGroups      *[]*CloudInstanceSecurityGroupsModel                                       `tfsdk:"security_groups" json:"security_groups,optional"`
	Name                types.String                                                               `tfsdk:"name" json:"name,optional"`
	Tags                customfield.Map[types.String]                                              `tfsdk:"tags" json:"tags,computed_optional,no_refresh"`
	CreatedAt           timetypes.RFC3339                                                          `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                               `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	InstanceDescription types.String                                                               `tfsdk:"instance_description" json:"instance_description,computed"`
	Region              types.String                                                               `tfsdk:"region" json:"region,computed"`
	Status              types.String                                                               `tfsdk:"status" json:"status,computed"`
	TaskState           types.String                                                               `tfsdk:"task_state" json:"task_state,computed"`
	VmState             types.String                                                               `tfsdk:"vm_state" json:"vm_state,computed_optional"`
	Addresses           customfield.Map[customfield.NestedObjectList[CloudInstanceAddressesModel]] `tfsdk:"addresses" json:"addresses,computed"`
	BlackholePorts      customfield.NestedObjectList[CloudInstanceBlackholePortsModel]             `tfsdk:"blackhole_ports" json:"blackhole_ports,computed"`
	DDOSProfile         customfield.NestedObject[CloudInstanceDDOSProfileModel]                    `tfsdk:"ddos_profile" json:"ddos_profile,computed"`
	FixedIPAssignments  customfield.NestedObjectList[CloudInstanceFixedIPAssignmentsModel]         `tfsdk:"fixed_ip_assignments" json:"fixed_ip_assignments,computed"`
	InstanceIsolation   customfield.NestedObject[CloudInstanceInstanceIsolationModel]              `tfsdk:"instance_isolation" json:"instance_isolation,computed"`
}

func (m CloudInstanceModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudInstanceModel) MarshalJSONForUpdate(state CloudInstanceModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudInstanceInterfacesModel struct {
	Type           types.String                                   `tfsdk:"type" json:"type,required"`
	InterfaceName  types.String                                   `tfsdk:"interface_name" json:"interface_name,optional"`
	IPFamily       types.String                                   `tfsdk:"ip_family" json:"ip_family,optional"`
	SecurityGroups *[]*CloudInstanceInterfacesSecurityGroupsModel `tfsdk:"security_groups" json:"security_groups,optional"`
	NetworkID      types.String                                   `tfsdk:"network_id" json:"network_id,optional"`
	SubnetID       types.String                                   `tfsdk:"subnet_id" json:"subnet_id,optional"`
	FloatingIP     *CloudInstanceInterfacesFloatingIPModel        `tfsdk:"floating_ip" json:"floating_ip,optional"`
	IPAddress      types.String                                   `tfsdk:"ip_address" json:"ip_address,computed_optional"`
	PortID         types.String                                   `tfsdk:"port_id" json:"port_id,computed_optional"`
}

type CloudInstanceInterfacesSecurityGroupsModel struct {
	ID types.String `tfsdk:"id" json:"id,required"`
}

type CloudInstanceInterfacesFloatingIPModel struct {
	Source             types.String `tfsdk:"source" json:"source,required"`
	ExistingFloatingID types.String `tfsdk:"existing_floating_id" json:"existing_floating_id,optional"`
}

// CloudInstanceVolumesModel represents an existing volume to attach to an instance.
// Only existing volumes are supported - users must create volumes separately using gcore_cloud_volume.
type CloudInstanceVolumesModel struct {
	VolumeID      types.String `tfsdk:"volume_id" json:"volume_id,required"`
	BootIndex     types.Int64  `tfsdk:"boot_index" json:"boot_index,optional"`
	AttachmentTag types.String `tfsdk:"attachment_tag" json:"attachment_tag,optional"`
}

// MarshalJSONWithState implements CustomMarshaler interface for apijson.
// This adds the hardcoded fields required by the API: source="existing-volume" and delete_on_termination=false.
func (m CloudInstanceVolumesModel) MarshalJSONWithState(plan any, state any) ([]byte, error) {
	// Build the volume payload with required fields
	payload := map[string]interface{}{
		"source":                "existing-volume",
		"volume_id":             m.VolumeID.ValueString(),
		"boot_index":            m.BootIndex.ValueInt64(),
		"delete_on_termination": false,
	}

	// Only include attachment_tag if it's set
	if !m.AttachmentTag.IsNull() && !m.AttachmentTag.IsUnknown() && m.AttachmentTag.ValueString() != "" {
		payload["attachment_tag"] = m.AttachmentTag.ValueString()
	}

	return json.Marshal(payload)
}

type CloudInstanceSecurityGroupsModel struct {
	ID types.String `tfsdk:"id" json:"id,required,no_refresh"`
}

type CloudInstanceAddressesModel struct {
	Addr          types.String `tfsdk:"addr" json:"addr,computed"`
	Type          types.String `tfsdk:"type" json:"type,computed"`
	InterfaceName types.String `tfsdk:"interface_name" json:"interface_name,computed"`
	SubnetID      types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
	SubnetName    types.String `tfsdk:"subnet_name" json:"subnet_name,computed"`
}

type CloudInstanceBlackholePortsModel struct {
	AlarmEnd      timetypes.RFC3339 `tfsdk:"alarm_end" json:"AlarmEnd,computed" format:"date-time"`
	AlarmStart    timetypes.RFC3339 `tfsdk:"alarm_start" json:"AlarmStart,computed" format:"date-time"`
	AlarmState    types.String      `tfsdk:"alarm_state" json:"AlarmState,computed"`
	AlertDuration types.String      `tfsdk:"alert_duration" json:"AlertDuration,computed"`
	DestinationIP types.String      `tfsdk:"destination_ip" json:"DestinationIP,computed"`
	ID            types.Int64       `tfsdk:"id" json:"ID,computed"`
}

type CloudInstanceDDOSProfileModel struct {
	ID                         types.Int64                                                            `tfsdk:"id" json:"id,computed"`
	Fields                     customfield.NestedObjectList[CloudInstanceDDOSProfileFieldsModel]      `tfsdk:"fields" json:"fields,computed"`
	Options                    customfield.NestedObject[CloudInstanceDDOSProfileOptionsModel]         `tfsdk:"options" json:"options,computed"`
	ProfileTemplate            customfield.NestedObject[CloudInstanceDDOSProfileProfileTemplateModel] `tfsdk:"profile_template" json:"profile_template,computed"`
	ProfileTemplateDescription types.String                                                           `tfsdk:"profile_template_description" json:"profile_template_description,computed"`
	Protocols                  customfield.NestedObjectList[CloudInstanceDDOSProfileProtocolsModel]   `tfsdk:"protocols" json:"protocols,computed"`
	Site                       types.String                                                           `tfsdk:"site" json:"site,computed"`
	Status                     customfield.NestedObject[CloudInstanceDDOSProfileStatusModel]          `tfsdk:"status" json:"status,computed"`
}

type CloudInstanceDDOSProfileFieldsModel struct {
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

type CloudInstanceDDOSProfileOptionsModel struct {
	Active types.Bool `tfsdk:"active" json:"active,computed"`
	Bgp    types.Bool `tfsdk:"bgp" json:"bgp,computed"`
}

type CloudInstanceDDOSProfileProfileTemplateModel struct {
	ID          types.Int64                                                                      `tfsdk:"id" json:"id,computed"`
	Description types.String                                                                     `tfsdk:"description" json:"description,computed"`
	Fields      customfield.NestedObjectList[CloudInstanceDDOSProfileProfileTemplateFieldsModel] `tfsdk:"fields" json:"fields,computed"`
	Name        types.String                                                                     `tfsdk:"name" json:"name,computed"`
}

type CloudInstanceDDOSProfileProfileTemplateFieldsModel struct {
	ID               types.Int64          `tfsdk:"id" json:"id,computed"`
	Default          types.String         `tfsdk:"default" json:"default,computed"`
	Description      types.String         `tfsdk:"description" json:"description,computed"`
	FieldType        types.String         `tfsdk:"field_type" json:"field_type,computed"`
	Name             types.String         `tfsdk:"name" json:"name,computed"`
	Required         types.Bool           `tfsdk:"required" json:"required,computed"`
	ValidationSchema jsontypes.Normalized `tfsdk:"validation_schema" json:"validation_schema,computed"`
}

type CloudInstanceDDOSProfileProtocolsModel struct {
	Port      types.String                   `tfsdk:"port" json:"port,computed"`
	Protocols customfield.List[types.String] `tfsdk:"protocols" json:"protocols,computed"`
}

type CloudInstanceDDOSProfileStatusModel struct {
	ErrorDescription types.String `tfsdk:"error_description" json:"error_description,computed"`
	Status           types.String `tfsdk:"status" json:"status,computed"`
}

type CloudInstanceFixedIPAssignmentsModel struct {
	External  types.Bool   `tfsdk:"external" json:"external,computed"`
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudInstanceInstanceIsolationModel struct {
	Reason types.String `tfsdk:"reason" json:"reason,computed"`
}
