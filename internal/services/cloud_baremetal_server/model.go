// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_server

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apijson"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudBaremetalServerModel struct {
	ProjectID     types.Int64                             `tfsdk:"project_id" path:"project_id,optional"`
	RegionID      types.Int64                             `tfsdk:"region_id" path:"region_id,optional"`
	Flavor        types.String                            `tfsdk:"flavor" json:"flavor,required"`
	Interfaces    *[]*CloudBaremetalServerInterfacesModel `tfsdk:"interfaces" json:"interfaces,required"`
	ApptemplateID types.String                            `tfsdk:"apptemplate_id" json:"apptemplate_id,optional"`
	ImageID       types.String                            `tfsdk:"image_id" json:"image_id,optional"`
	Name          types.String                            `tfsdk:"name" json:"name,optional"`
	NameTemplate  types.String                            `tfsdk:"name_template" json:"name_template,optional"`
	Password      types.String                            `tfsdk:"password" json:"password,optional"`
	SSHKeyName    types.String                            `tfsdk:"ssh_key_name" json:"ssh_key_name,optional"`
	UserData      types.String                            `tfsdk:"user_data" json:"user_data,optional"`
	Username      types.String                            `tfsdk:"username" json:"username,optional"`
	Tags          *map[string]types.String                `tfsdk:"tags" json:"tags,optional"`
	DDOSProfile   *CloudBaremetalServerDDOSProfileModel   `tfsdk:"ddos_profile" json:"ddos_profile,optional"`
	AppConfig     jsontypes.Normalized                    `tfsdk:"app_config" json:"app_config,optional"`
	Tasks         customfield.List[types.String]          `tfsdk:"tasks" json:"tasks,computed"`
}

func (m CloudBaremetalServerModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudBaremetalServerModel) MarshalJSONForUpdate(state CloudBaremetalServerModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
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

type CloudBaremetalServerDDOSProfileModel struct {
	ProfileTemplate types.Int64                                    `tfsdk:"profile_template" json:"profile_template,required"`
	Fields          *[]*CloudBaremetalServerDDOSProfileFieldsModel `tfsdk:"fields" json:"fields,optional"`
}

type CloudBaremetalServerDDOSProfileFieldsModel struct {
	BaseField  types.Int64             `tfsdk:"base_field" json:"base_field,optional"`
	FieldName  types.String            `tfsdk:"field_name" json:"field_name,optional"`
	FieldValue *[]jsontypes.Normalized `tfsdk:"field_value" json:"field_value,optional"`
	Value      types.String            `tfsdk:"value" json:"value,optional"`
}
