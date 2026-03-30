// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_server

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudBaremetalServerModel struct {
	ID                 types.String                                                                      `tfsdk:"id" json:"id,computed"`
	ProjectID          types.Int64                                                                       `tfsdk:"project_id" path:"project_id,optional"`
	RegionID           types.Int64                                                                       `tfsdk:"region_id" path:"region_id,optional"`
	Flavor             types.String                                                                      `tfsdk:"flavor" json:"flavor,required,no_refresh"`
	Interfaces         *[]*CloudBaremetalServerInterfacesModel                                           `tfsdk:"interfaces" json:"interfaces,required,no_refresh"`
	ApptemplateID      types.String                                                                      `tfsdk:"apptemplate_id" json:"apptemplate_id,optional,no_refresh"`
	ImageID            types.String                                                                      `tfsdk:"image_id" json:"image_id,optional,no_refresh"`
	NameTemplate       types.String                                                                      `tfsdk:"name_template" json:"name_template,optional,no_refresh"`
	Password           types.String                                                                      `tfsdk:"password_wo" json:"password,optional,no_refresh"`
	PasswordWoVersion  types.Int64                                                                       `tfsdk:"password_wo_version"`
	SSHKeyName         types.String                                                                      `tfsdk:"ssh_key_name" json:"ssh_key_name,optional,no_refresh"`
	UserData           types.String                                                                      `tfsdk:"user_data" json:"user_data,optional,no_refresh"`
	Username           types.String                                                                      `tfsdk:"username" json:"username,optional,no_refresh"`
	AppConfig          *map[string]jsontypes.Normalized                                                  `tfsdk:"app_config" json:"app_config,optional,no_refresh"`
	Name               types.String                                                                      `tfsdk:"name" json:"name,computed_optional"`
	Tags               customfield.Map[types.String]                                                     `tfsdk:"tags" json:"tags,computed_optional,no_refresh"`
	CreatedAt          timetypes.RFC3339                                                                 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
	Region             types.String                                                                      `tfsdk:"region" json:"region,computed"`
	Status             types.String                                                                      `tfsdk:"status" json:"status,computed"`
	VmState            types.String                                                                      `tfsdk:"vm_state" json:"vm_state,computed"`
	Addresses          customfield.Map[customfield.NestedObjectList[CloudBaremetalServerAddressesModel]] `tfsdk:"addresses" json:"addresses,computed"`
	BlackholePorts     customfield.NestedObjectList[CloudBaremetalServerBlackholePortsModel]             `tfsdk:"blackhole_ports" json:"blackhole_ports,computed"`
	FixedIPAssignments customfield.NestedObjectList[CloudBaremetalServerFixedIPAssignmentsModel]         `tfsdk:"fixed_ip_assignments" json:"fixed_ip_assignments,computed"`
	InstanceIsolation  customfield.NestedObject[CloudBaremetalServerInstanceIsolationModel]              `tfsdk:"instance_isolation" json:"instance_isolation,computed"`
}

func (m CloudBaremetalServerModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudBaremetalServerModel) MarshalJSONForUpdate(state CloudBaremetalServerModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudBaremetalServerInterfacesModel struct {
	Type          types.String                                   `tfsdk:"type" json:"type,required"`
	InterfaceName types.String                                   `tfsdk:"interface_name" json:"interface_name,optional"`
	IPFamily      types.String                                   `tfsdk:"ip_family" json:"ip_family,optional"`
	PortGroup     types.Int64                                    `tfsdk:"port_group" json:"port_group,computed_optional"`
	NetworkID     types.String                                   `tfsdk:"network_id" json:"network_id,optional"`
	SubnetID      types.String                                   `tfsdk:"subnet_id" json:"subnet_id,optional"`
	FloatingIP    *CloudBaremetalServerInterfacesFloatingIPModel `tfsdk:"floating_ip" json:"floating_ip,optional"`
	IPAddress     types.String                                   `tfsdk:"ip_address" json:"ip_address,optional"`
	PortID        types.String                                   `tfsdk:"port_id" json:"port_id,optional"`
}

type CloudBaremetalServerInterfacesFloatingIPModel struct {
	Source             types.String `tfsdk:"source" json:"source,required"`
	ExistingFloatingID types.String `tfsdk:"existing_floating_id" json:"existing_floating_id,optional"`
}

type CloudBaremetalServerAddressesModel struct {
	Addr          types.String `tfsdk:"addr" json:"addr,computed"`
	Type          types.String `tfsdk:"type" json:"type,computed"`
	InterfaceName types.String `tfsdk:"interface_name" json:"interface_name,computed"`
	SubnetID      types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
	SubnetName    types.String `tfsdk:"subnet_name" json:"subnet_name,computed"`
}

type CloudBaremetalServerBlackholePortsModel struct {
	AlarmEnd      timetypes.RFC3339 `tfsdk:"alarm_end" json:"AlarmEnd,computed" format:"date-time"`
	AlarmStart    timetypes.RFC3339 `tfsdk:"alarm_start" json:"AlarmStart,computed" format:"date-time"`
	AlarmState    types.String      `tfsdk:"alarm_state" json:"AlarmState,computed"`
	AlertDuration types.String      `tfsdk:"alert_duration" json:"AlertDuration,computed"`
	DestinationIP types.String      `tfsdk:"destination_ip" json:"DestinationIP,computed"`
	ID            types.Int64       `tfsdk:"id" json:"ID,computed"`
}

type CloudBaremetalServerFixedIPAssignmentsModel struct {
	External  types.Bool   `tfsdk:"external" json:"external,computed"`
	IPAddress types.String `tfsdk:"ip_address" json:"ip_address,computed"`
	SubnetID  types.String `tfsdk:"subnet_id" json:"subnet_id,computed"`
}

type CloudBaremetalServerInstanceIsolationModel struct {
	Reason types.String `tfsdk:"reason" json:"reason,computed"`
}
