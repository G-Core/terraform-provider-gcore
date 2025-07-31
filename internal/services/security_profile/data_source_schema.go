// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package security_profile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*SecurityProfileDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"ip_address": schema.StringAttribute{
				Computed: true,
			},
			"plan": schema.StringAttribute{
				Computed: true,
			},
			"site": schema.StringAttribute{
				Computed: true,
			},
			"protocols": schema.ListAttribute{
				Computed:   true,
				CustomType: customfield.NewListType[customfield.Map[jsontypes.Normalized]](ctx),
				ElementType: types.MapType{
					ElemType: jsontypes.NormalizedType{},
				},
			},
			"status": schema.MapAttribute{
				Computed:    true,
				CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
				ElementType: jsontypes.NormalizedType{},
			},
			"fields": schema.ListNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectListType[SecurityProfileFieldsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"base_field": schema.Int64Attribute{
							Computed: true,
						},
						"default": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"field_type": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"required": schema.BoolAttribute{
							Computed: true,
						},
						"validation_schema": schema.MapAttribute{
							Computed:    true,
							CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
							ElementType: jsontypes.NormalizedType{},
						},
						"field_value": schema.MapAttribute{
							Computed:    true,
							CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
							ElementType: jsontypes.NormalizedType{},
						},
					},
				},
			},
			"options": schema.SingleNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectType[SecurityProfileOptionsDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"active": schema.BoolAttribute{
						Computed: true,
					},
					"bgp": schema.BoolAttribute{
						Computed: true,
					},
					"price": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"profile_template": schema.SingleNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectType[SecurityProfileProfileTemplateDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed: true,
					},
					"created": schema.StringAttribute{
						Computed:   true,
						CustomType: timetypes.RFC3339Type{},
					},
					"fields": schema.ListNestedAttribute{
						Computed:   true,
						CustomType: customfield.NewNestedObjectListType[SecurityProfileProfileTemplateFieldsDataSourceModel](ctx),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
								"default": schema.StringAttribute{
									Computed: true,
								},
								"description": schema.StringAttribute{
									Computed: true,
								},
								"field_type": schema.StringAttribute{
									Description: `Available values: "int", "bool", "str".`,
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive(
											"int",
											"bool",
											"str",
										),
									},
								},
								"required": schema.BoolAttribute{
									Computed: true,
								},
								"validation_schema": schema.MapAttribute{
									Computed:    true,
									CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
									ElementType: jsontypes.NormalizedType{},
								},
							},
						},
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
					"version": schema.StringAttribute{
						Computed: true,
					},
					"base_template": schema.Int64Attribute{
						Computed: true,
					},
					"description": schema.StringAttribute{
						Computed: true,
					},
					"template_sifter": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

func (d *SecurityProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *SecurityProfileDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
