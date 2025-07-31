// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package security_profile

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type SecurityProfileDataSourceModel struct {
	ID              types.Int64                                                             `tfsdk:"id" path:"id,required"`
	IPAddress       types.String                                                            `tfsdk:"ip_address" json:"ip_address,computed"`
	Plan            types.String                                                            `tfsdk:"plan" json:"plan,computed"`
	Site            types.String                                                            `tfsdk:"site" json:"site,computed"`
	Protocols       customfield.List[customfield.Map[jsontypes.Normalized]]                 `tfsdk:"protocols" json:"protocols,computed"`
	Status          customfield.Map[jsontypes.Normalized]                                   `tfsdk:"status" json:"status,computed"`
	Fields          customfield.NestedObjectList[SecurityProfileFieldsDataSourceModel]      `tfsdk:"fields" json:"fields,computed"`
	Options         customfield.NestedObject[SecurityProfileOptionsDataSourceModel]         `tfsdk:"options" json:"options,computed"`
	ProfileTemplate customfield.NestedObject[SecurityProfileProfileTemplateDataSourceModel] `tfsdk:"profile_template" json:"profile_template,computed"`
}

type SecurityProfileFieldsDataSourceModel struct {
	ID               types.Int64                           `tfsdk:"id" json:"id,computed"`
	BaseField        types.Int64                           `tfsdk:"base_field" json:"base_field,computed"`
	Default          types.String                          `tfsdk:"default" json:"default,computed"`
	Description      types.String                          `tfsdk:"description" json:"description,computed"`
	FieldType        types.String                          `tfsdk:"field_type" json:"field_type,computed"`
	Name             types.String                          `tfsdk:"name" json:"name,computed"`
	Required         types.Bool                            `tfsdk:"required" json:"required,computed"`
	ValidationSchema customfield.Map[jsontypes.Normalized] `tfsdk:"validation_schema" json:"validation_schema,computed"`
	FieldValue       customfield.Map[jsontypes.Normalized] `tfsdk:"field_value" json:"field_value,computed"`
}

type SecurityProfileOptionsDataSourceModel struct {
	Active types.Bool   `tfsdk:"active" json:"active,computed"`
	Bgp    types.Bool   `tfsdk:"bgp" json:"bgp,computed"`
	Price  types.String `tfsdk:"price" json:"price,computed"`
}

type SecurityProfileProfileTemplateDataSourceModel struct {
	ID             types.Int64                                                                       `tfsdk:"id" json:"id,computed"`
	Created        timetypes.RFC3339                                                                 `tfsdk:"created" json:"created,computed" format:"date-time"`
	Fields         customfield.NestedObjectList[SecurityProfileProfileTemplateFieldsDataSourceModel] `tfsdk:"fields" json:"fields,computed"`
	Name           types.String                                                                      `tfsdk:"name" json:"name,computed"`
	Version        types.String                                                                      `tfsdk:"version" json:"version,computed"`
	BaseTemplate   types.Int64                                                                       `tfsdk:"base_template" json:"base_template,computed"`
	Description    types.String                                                                      `tfsdk:"description" json:"description,computed"`
	TemplateSifter types.String                                                                      `tfsdk:"template_sifter" json:"template_sifter,computed"`
}

type SecurityProfileProfileTemplateFieldsDataSourceModel struct {
	ID               types.Int64                           `tfsdk:"id" json:"id,computed"`
	Name             types.String                          `tfsdk:"name" json:"name,computed"`
	Default          types.String                          `tfsdk:"default" json:"default,computed"`
	Description      types.String                          `tfsdk:"description" json:"description,computed"`
	FieldType        types.String                          `tfsdk:"field_type" json:"field_type,computed"`
	Required         types.Bool                            `tfsdk:"required" json:"required,computed"`
	ValidationSchema customfield.Map[jsontypes.Normalized] `tfsdk:"validation_schema" json:"validation_schema,computed"`
}
