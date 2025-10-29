// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudInstanceModel struct {
	ID                  types.String                                                               `tfsdk:"id" json:"id,computed"`
	ProjectID           types.Int64                                                                `tfsdk:"project_id" path:"project_id,optional"`
	RegionID            types.Int64                                                                `tfsdk:"region_id" path:"region_id,optional"`
	Flavor              types.String                                                               `tfsdk:"flavor" json:"flavor,required,no_refresh"`
	Interfaces          *[]*CloudInstanceInterfacesModel                                           `tfsdk:"interfaces" json:"interfaces,required,no_refresh"`
	Volumes             *[]*CloudInstanceVolumesModel                                              `tfsdk:"volumes" json:"volumes,required"`
	AllowAppPorts       types.Bool                                                                 `tfsdk:"allow_app_ports" json:"allow_app_ports,optional,no_refresh"`
	NameTemplate        types.String                                                               `tfsdk:"name_template" json:"name_template,optional,no_refresh"`
	Password            types.String                                                               `tfsdk:"password" json:"password,optional,no_refresh"`
	ServergroupID       types.String                                                               `tfsdk:"servergroup_id" json:"servergroup_id,optional,no_refresh"`
	SSHKeyName          types.String                                                               `tfsdk:"ssh_key_name" json:"ssh_key_name,optional"`
	UserData            types.String                                                               `tfsdk:"user_data" json:"user_data,optional,no_refresh"`
	Username            types.String                                                               `tfsdk:"username" json:"username,optional,no_refresh"`
	SecurityGroups      *[]*CloudInstanceSecurityGroupsModel                                       `tfsdk:"security_groups" json:"security_groups,optional"`
	Configuration       jsontypes.Normalized                                                       `tfsdk:"configuration" json:"configuration,optional,no_refresh"`
	Name                types.String                                                               `tfsdk:"name" json:"name,optional"`
	Tags                *map[string]types.String                                                   `tfsdk:"tags" json:"tags,optional,no_refresh"`
	CreatedAt           timetypes.RFC3339                                                          `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                               `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	InstanceDescription types.String                                                               `tfsdk:"instance_description" json:"instance_description,computed"`
	Region              types.String                                                               `tfsdk:"region" json:"region,computed"`
	Status              types.String                                                               `tfsdk:"status" json:"status,computed"`
	TaskID              types.String                                                               `tfsdk:"task_id" json:"task_id,computed"`
	TaskState           types.String                                                               `tfsdk:"task_state" json:"task_state,computed"`
	VmState             types.String                                                               `tfsdk:"vm_state" json:"vm_state,computed"`
	Addresses           customfield.Map[customfield.NestedObjectList[CloudInstanceAddressesModel]] `tfsdk:"addresses" json:"addresses,computed"`
	Tasks               customfield.List[types.String]                                             `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
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
	IPAddress      types.String                                   `tfsdk:"ip_address" json:"ip_address,optional"`
	PortID         types.String                                   `tfsdk:"port_id" json:"port_id,optional"`
}

type CloudInstanceInterfacesSecurityGroupsModel struct {
	ID types.String `tfsdk:"id" json:"id,required"`
}

type CloudInstanceInterfacesFloatingIPModel struct {
	Source             types.String `tfsdk:"source" json:"source,required"`
	ExistingFloatingID types.String `tfsdk:"existing_floating_id" json:"existing_floating_id,optional"`
}

type CloudInstanceVolumesModel struct {
	Size                types.Int64              `tfsdk:"size" json:"size,optional,no_refresh"`
	Source              types.String             `tfsdk:"source" json:"source,required,no_refresh"`
	AttachmentTag       types.String             `tfsdk:"attachment_tag" json:"attachment_tag,optional,no_refresh"`
	DeleteOnTermination types.Bool               `tfsdk:"delete_on_termination" json:"delete_on_termination,computed_optional"`
	Name                types.String             `tfsdk:"name" json:"name,optional,no_refresh"`
	Tags                *map[string]types.String `tfsdk:"tags" json:"tags,optional,no_refresh"`
	TypeName            types.String             `tfsdk:"type_name" json:"type_name,optional,no_refresh"`
	ImageID             types.String             `tfsdk:"image_id" json:"image_id,optional,no_refresh"`
	BootIndex           types.Int64              `tfsdk:"boot_index" json:"boot_index,optional,no_refresh"`
	SnapshotID          types.String             `tfsdk:"snapshot_id" json:"snapshot_id,optional,no_refresh"`
	ApptemplateID       types.String             `tfsdk:"apptemplate_id" json:"apptemplate_id,optional,no_refresh"`
	VolumeID            types.String             `tfsdk:"volume_id" json:"volume_id,optional,no_refresh"`
}

type CloudInstanceSecurityGroupsModel struct {
	ID types.String `tfsdk:"id" json:"id,required,no_refresh"`
}

type CloudInstanceAddressesModel struct {
	Addr          types.String `tfsdk:"addr" json:"addr,required"`
	Type          types.String `tfsdk:"type" json:"type,required"`
	InterfaceName types.String `tfsdk:"interface_name" json:"interface_name,optional"`
	SubnetID      types.String `tfsdk:"subnet_id" json:"subnet_id,optional"`
	SubnetName    types.String `tfsdk:"subnet_name" json:"subnet_name,optional"`
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
	FieldName        types.String         `tfsdk:"field_name" json:"field_name,computed"`
	FieldType        types.String         `tfsdk:"field_type" json:"field_type,computed"`
	FieldValue       jsontypes.Normalized `tfsdk:"field_value" json:"field_value,computed"`
	Name             types.String         `tfsdk:"name" json:"name,computed"`
	Required         types.Bool           `tfsdk:"required" json:"required,computed"`
	ValidationSchema jsontypes.Normalized `tfsdk:"validation_schema" json:"validation_schema,computed"`
	Value            types.String         `tfsdk:"value" json:"value,computed"`
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
