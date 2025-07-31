// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package security_profile

import (
	"bytes"
	"mime/multipart"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/apiform"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type SecurityProfileModel struct {
	ID              types.Int64                                             `tfsdk:"id" json:"id,computed"`
	ProfileTemplate types.Int64                                             `tfsdk:"profile_template" json:"profile_template,required,no_refresh"`
	Fields          *[]*SecurityProfileFieldsModel                          `tfsdk:"fields" json:"fields,required"`
	IPAddress       types.String                                            `tfsdk:"ip_address" json:"ip_address,optional"`
	Site            types.String                                            `tfsdk:"site" json:"site,optional"`
	Plan            types.String                                            `tfsdk:"plan" json:"plan,computed"`
	Protocols       customfield.List[customfield.Map[jsontypes.Normalized]] `tfsdk:"protocols" json:"protocols,computed"`
	Status          customfield.Map[jsontypes.Normalized]                   `tfsdk:"status" json:"status,computed"`
	Options         customfield.NestedObject[SecurityProfileOptionsModel]   `tfsdk:"options" json:"options,computed"`
}

func (r SecurityProfileModel) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err != nil {
		writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}

type SecurityProfileFieldsModel struct {
	ID               types.Int64                           `tfsdk:"id" json:"id,computed"`
	BaseField        types.Int64                           `tfsdk:"base_field" json:"base_field,required"`
	Default          types.String                          `tfsdk:"default" json:"default,computed"`
	Description      types.String                          `tfsdk:"description" json:"description,computed"`
	FieldType        types.String                          `tfsdk:"field_type" json:"field_type,computed"`
	Name             types.String                          `tfsdk:"name" json:"name,computed"`
	Required         types.Bool                            `tfsdk:"required" json:"required,computed"`
	ValidationSchema customfield.Map[jsontypes.Normalized] `tfsdk:"validation_schema" json:"validation_schema,computed"`
	FieldValue       *map[string]jsontypes.Normalized      `tfsdk:"field_value" json:"field_value,optional"`
}

type SecurityProfileOptionsModel struct {
	Active types.Bool   `tfsdk:"active" json:"active,computed"`
	Bgp    types.Bool   `tfsdk:"bgp" json:"bgp,computed"`
	Price  types.String `tfsdk:"price" json:"price,computed"`
}
