// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster

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

type CloudGPUBaremetalClusterDataSourceModel struct {
	ClusterID     types.String                                                                    `tfsdk:"cluster_id" path:"cluster_id,required"`
	ProjectID     types.Int64                                                                     `tfsdk:"project_id" path:"project_id,required"`
	RegionID      types.Int64                                                                     `tfsdk:"region_id" path:"region_id,required"`
	ClusterName   types.String                                                                    `tfsdk:"cluster_name" json:"cluster_name,computed"`
	ClusterStatus types.String                                                                    `tfsdk:"cluster_status" json:"cluster_status,computed"`
	CreatedAt     types.String                                                                    `tfsdk:"created_at" json:"created_at,computed"`
	CreatorTaskID types.String                                                                    `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	Flavor        types.String                                                                    `tfsdk:"flavor" json:"flavor,computed"`
	ImageID       types.String                                                                    `tfsdk:"image_id" json:"image_id,computed"`
	ImageName     types.String                                                                    `tfsdk:"image_name" json:"image_name,computed"`
	Password      types.String                                                                    `tfsdk:"password" json:"password,computed"`
	Region        types.String                                                                    `tfsdk:"region" json:"region,computed"`
	SSHKeyName    types.String                                                                    `tfsdk:"ssh_key_name" json:"ssh_key_name,computed"`
	TaskID        types.String                                                                    `tfsdk:"task_id" json:"task_id,computed"`
	TaskStatus    types.String                                                                    `tfsdk:"task_status" json:"task_status,computed"`
	UserData      types.String                                                                    `tfsdk:"user_data" json:"user_data,computed"`
	Username      types.String                                                                    `tfsdk:"username" json:"username,computed"`
	Interfaces    customfield.NestedObjectList[CloudGPUBaremetalClusterInterfacesDataSourceModel] `tfsdk:"interfaces" json:"interfaces,computed"`
	Servers       customfield.NestedObjectList[CloudGPUBaremetalClusterServersDataSourceModel]    `tfsdk:"servers" json:"servers,computed"`
	Tags          customfield.NestedObjectList[CloudGPUBaremetalClusterTagsDataSourceModel]       `tfsdk:"tags" json:"tags,computed"`
}

func (m *CloudGPUBaremetalClusterDataSourceModel) toReadParams(_ context.Context) (params cloud.GPUBaremetalClusterGetParams, diags diag.Diagnostics) {
	params = cloud.GPUBaremetalClusterGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}
	if !m.RegionID.IsNull() {
		params.RegionID = param.NewOpt(m.RegionID.ValueInt64())
	}

	return
}

type CloudGPUBaremetalClusterInterfacesDataSourceModel struct {
	NetworkID types.String `tfsdk:"network_id" json:"network_id,computed"`
	PortID    types.String `tfsdk:"port_id" json:"port_id,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
	Type      types.String `tfsdk:"type" json:"type,computed"`
}

type CloudGPUBaremetalClusterServersDataSourceModel struct {
	ID                  types.String                                                                                           `tfsdk:"id" json:"id,computed"`
	Addresses           customfield.Map[customfield.NestedObjectList[CloudGPUBaremetalClusterServersAddressesDataSourceModel]] `tfsdk:"addresses" json:"addresses,computed"`
	BlackholePorts      customfield.NestedObjectList[CloudGPUBaremetalClusterServersBlackholePortsDataSourceModel]             `tfsdk:"blackhole_ports" json:"blackhole_ports,computed"`
	CreatedAt           timetypes.RFC3339                                                                                      `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	CreatorTaskID       types.String                                                                                           `tfsdk:"creator_task_id" json:"creator_task_id,computed"`
	DDOSProfile         customfield.NestedObject[CloudGPUBaremetalClusterServersDDOSProfileDataSourceModel]                    `tfsdk:"ddos_profile" json:"ddos_profile,computed"`
	FixedIPAssignments  customfield.NestedObjectList[CloudGPUBaremetalClusterServersFixedIPAssignmentsDataSourceModel]         `tfsdk:"fixed_ip_assignments" json:"fixed_ip_assignments,computed"`
	Flavor              customfield.NestedObject[CloudGPUBaremetalClusterServersFlavorDataSourceModel]                         `tfsdk:"flavor" json:"flavor,computed"`
	InstanceDescription types.String                                                                                           `tfsdk:"instance_description" json:"instance_description,computed"`
	InstanceIsolation   customfield.NestedObject[CloudGPUBaremetalClusterServersInstanceIsolationDataSourceModel]              `tfsdk:"instance_isolation" json:"instance_isolation,computed"`
	Name                types.String                                                                                           `tfsdk:"name" json:"name,computed"`
	ProjectID           types.Int64                                                                                            `tfsdk:"project_id" json:"project_id,computed"`
	Region              types.String                                                                                           `tfsdk:"region" json:"region,computed"`
	RegionID            types.Int64                                                                                            `tfsdk:"region_id" json:"region_id,computed"`
	SecurityGroups      customfield.NestedObjectList[CloudGPUBaremetalClusterServersSecurityGroupsDataSourceModel]             `tfsdk:"security_groups" json:"security_groups,computed"`
	SSHKeyName          types.String                                                                                           `tfsdk:"ssh_key_name" json:"ssh_key_name,computed"`
	Status              types.String                                                                                           `tfsdk:"status" json:"status,computed"`
	Tags                customfield.NestedObjectList[CloudGPUBaremetalClusterServersTagsDataSourceModel]                       `tfsdk:"tags" json:"tags,computed"`
	TaskID              types.String                                                                                           `tfsdk:"task_id" json:"task_id,computed"`
	TaskState           types.String                                                                                           `tfsdk:"task_state" json:"task_state,computed"`
	VmState             types.String                                                                                           `tfsdk:"vm_state" json:"vm_state,computed"`
}

type CloudGPUBaremetalClusterServersAddressesDataSourceModel struct {
	Addr          types.String `tfsdk:"addr" json:"addr,required"`
	Type          types.String `tfsdk:"type" json:"type,required"`
	InterfaceName types.String `tfsdk:"interface_name" json:"interface_name,optional"`
	SubnetID      types.String `tfsdk:"subnet_id" json:"subnet_id,optional"`
	SubnetName    types.String `tfsdk:"subnet_name" json:"subnet_name,optional"`
}

type CloudGPUBaremetalClusterServersBlackholePortsDataSourceModel struct {
	AlarmEnd      timetypes.RFC3339 `tfsdk:"alarm_end" json:"AlarmEnd,computed" format:"date-time"`
	AlarmStart    timetypes.RFC3339 `tfsdk:"alarm_start" json:"AlarmStart,computed" format:"date-time"`
	AlarmState    types.String      `tfsdk:"alarm_state" json:"AlarmState,computed"`
	AlertDuration types.String      `tfsdk:"alert_duration" json:"AlertDuration,computed"`
	DestinationIP types.String      `tfsdk:"destination_ip" json:"DestinationIP,computed"`
	ID            types.Int64       `tfsdk:"id" json:"ID,computed"`
}

type CloudGPUBaremetalClusterServersDDOSProfileDataSourceModel struct {
	ID                         types.Int64                                                                                        `tfsdk:"id" json:"id,computed"`
	ProfileTemplate            customfield.NestedObject[CloudGPUBaremetalClusterServersDDOSProfileProfileTemplateDataSourceModel] `tfsdk:"profile_template" json:"profile_template,computed"`
	Fields                     customfield.NestedObjectList[CloudGPUBaremetalClusterServersDDOSProfileFieldsDataSourceModel]      `tfsdk:"fields" json:"fields,computed"`
	Options                    customfield.NestedObject[CloudGPUBaremetalClusterServersDDOSProfileOptionsDataSourceModel]         `tfsdk:"options" json:"options,computed"`
	ProfileTemplateDescription types.String                                                                                       `tfsdk:"profile_template_description" json:"profile_template_description,computed"`
	Protocols                  customfield.NestedObjectList[CloudGPUBaremetalClusterServersDDOSProfileProtocolsDataSourceModel]   `tfsdk:"protocols" json:"protocols,computed"`
	Site                       types.String                                                                                       `tfsdk:"site" json:"site,computed"`
	Status                     customfield.NestedObject[CloudGPUBaremetalClusterServersDDOSProfileStatusDataSourceModel]          `tfsdk:"status" json:"status,computed"`
}

type CloudGPUBaremetalClusterServersDDOSProfileProfileTemplateDataSourceModel struct {
	ID          types.Int64                                                                                                  `tfsdk:"id" json:"id,computed"`
	Name        types.String                                                                                                 `tfsdk:"name" json:"name,computed"`
	Description types.String                                                                                                 `tfsdk:"description" json:"description,computed"`
	Fields      customfield.NestedObjectList[CloudGPUBaremetalClusterServersDDOSProfileProfileTemplateFieldsDataSourceModel] `tfsdk:"fields" json:"fields,computed"`
}

type CloudGPUBaremetalClusterServersDDOSProfileProfileTemplateFieldsDataSourceModel struct {
	ID               types.Int64          `tfsdk:"id" json:"id,computed"`
	Name             types.String         `tfsdk:"name" json:"name,computed"`
	Default          types.String         `tfsdk:"default" json:"default,computed"`
	Description      types.String         `tfsdk:"description" json:"description,computed"`
	FieldType        types.String         `tfsdk:"field_type" json:"field_type,computed"`
	Required         types.Bool           `tfsdk:"required" json:"required,computed"`
	ValidationSchema jsontypes.Normalized `tfsdk:"validation_schema" json:"validation_schema,computed"`
}

type CloudGPUBaremetalClusterServersDDOSProfileFieldsDataSourceModel struct {
	ID               types.Int64          `tfsdk:"id" json:"id,computed"`
	Default          jsontypes.Normalized `tfsdk:"default" json:"default,computed"`
	Description      types.String         `tfsdk:"description" json:"description,computed"`
	FieldValue       jsontypes.Normalized `tfsdk:"field_value" json:"field_value,computed"`
	Name             types.String         `tfsdk:"name" json:"name,computed"`
	BaseField        types.Int64          `tfsdk:"base_field" json:"base_field,computed"`
	FieldName        types.String         `tfsdk:"field_name" json:"field_name,computed"`
	FieldType        types.String         `tfsdk:"field_type" json:"field_type,computed"`
	Required         types.Bool           `tfsdk:"required" json:"required,computed"`
	ValidationSchema jsontypes.Normalized `tfsdk:"validation_schema" json:"validation_schema,computed"`
	Value            types.String         `tfsdk:"value" json:"value,computed"`
}

type CloudGPUBaremetalClusterServersDDOSProfileOptionsDataSourceModel struct {
	Active types.Bool `tfsdk:"active" json:"active,computed"`
	Bgp    types.Bool `tfsdk:"bgp" json:"bgp,computed"`
}

type CloudGPUBaremetalClusterServersDDOSProfileProtocolsDataSourceModel struct {
	Port      types.String                   `tfsdk:"port" json:"port,computed"`
	Protocols customfield.List[types.String] `tfsdk:"protocols" json:"protocols,computed"`
}

type CloudGPUBaremetalClusterServersDDOSProfileStatusDataSourceModel struct {
	ErrorDescription types.String `tfsdk:"error_description" json:"error_description,computed"`
	Status           types.String `tfsdk:"status" json:"status,computed"`
}

type CloudGPUBaremetalClusterServersFixedIPAssignmentsDataSourceModel struct {
	External  types.Bool   `tfsdk:"external" json:"external,computed"`
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudGPUBaremetalClusterServersFlavorDataSourceModel struct {
	Architecture        types.String                                                                                      `tfsdk:"architecture" json:"architecture,computed"`
	FlavorID            types.String                                                                                      `tfsdk:"flavor_id" json:"flavor_id,computed"`
	FlavorName          types.String                                                                                      `tfsdk:"flavor_name" json:"flavor_name,computed"`
	HardwareDescription customfield.NestedObject[CloudGPUBaremetalClusterServersFlavorHardwareDescriptionDataSourceModel] `tfsdk:"hardware_description" json:"hardware_description,computed"`
	OsType              types.String                                                                                      `tfsdk:"os_type" json:"os_type,computed"`
	Ram                 types.Int64                                                                                       `tfsdk:"ram" json:"ram,computed"`
	ResourceClass       types.String                                                                                      `tfsdk:"resource_class" json:"resource_class,computed"`
	Vcpus               types.Int64                                                                                       `tfsdk:"vcpus" json:"vcpus,computed"`
}

type CloudGPUBaremetalClusterServersFlavorHardwareDescriptionDataSourceModel struct {
	CPU     types.String `tfsdk:"cpu" json:"cpu,computed"`
	Disk    types.String `tfsdk:"disk" json:"disk,computed"`
	GPU     types.String `tfsdk:"gpu" json:"gpu,computed"`
	License types.String `tfsdk:"license" json:"license,computed"`
	Network types.String `tfsdk:"network" json:"network,computed"`
	Ram     types.String `tfsdk:"ram" json:"ram,computed"`
}

type CloudGPUBaremetalClusterServersInstanceIsolationDataSourceModel struct {
	Reason types.String `tfsdk:"reason" json:"reason,computed"`
}

type CloudGPUBaremetalClusterServersSecurityGroupsDataSourceModel struct {
	Name types.String `tfsdk:"name" json:"name,computed"`
}

type CloudGPUBaremetalClusterServersTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}

type CloudGPUBaremetalClusterTagsDataSourceModel struct {
	Key      types.String `tfsdk:"key" json:"key,computed"`
	ReadOnly types.Bool   `tfsdk:"read_only" json:"read_only,computed"`
	Value    types.String `tfsdk:"value" json:"value,computed"`
}
