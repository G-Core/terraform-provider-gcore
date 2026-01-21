// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_custom_page_set

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*WaapCustomPageSetsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filter page sets based on their name. Supports '*' as a wildcard character",
				Optional:    true,
			},
			"ordering": schema.StringAttribute{
				Description: "Sort the response by given field.\nAvailable values: \"name\", \"-name\", \"id\", \"-id\".",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"name",
						"-name",
						"id",
						"-id",
					),
				},
			},
			"ids": schema.ListAttribute{
				Description: "Filter page sets based on their IDs",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"max_items": schema.Int64Attribute{
				Description: "Max items to fetch, default: 1000",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"items": schema.ListNestedAttribute{
				Description: "The items returned by the data source",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[WaapCustomPageSetsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The ID of the custom page set",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the custom page set",
							Computed:    true,
						},
						"block": schema.SingleNestedAttribute{
							Computed:   true,
							CustomType: customfield.NewNestedObjectType[WaapCustomPageSetsBlockDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									Description: "Indicates whether the custom custom page is active or inactive",
									Computed:    true,
								},
								"header": schema.StringAttribute{
									Description: "The text to display in the header of the custom page",
									Computed:    true,
								},
								"logo": schema.StringAttribute{
									Description: "Supported image types are JPEG, PNG and JPG, size is limited to width 450px, height 130px. This should be a base 64 encoding of the full HTML img tag compatible image, with the header included.",
									Computed:    true,
								},
								"text": schema.StringAttribute{
									Description: "The text to display in the body of the custom page",
									Computed:    true,
								},
								"title": schema.StringAttribute{
									Description: "The text to display in the title of the custom page",
									Computed:    true,
								},
							},
						},
						"block_csrf": schema.SingleNestedAttribute{
							Computed:   true,
							CustomType: customfield.NewNestedObjectType[WaapCustomPageSetsBlockCsrfDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									Description: "Indicates whether the custom custom page is active or inactive",
									Computed:    true,
								},
								"header": schema.StringAttribute{
									Description: "The text to display in the header of the custom page",
									Computed:    true,
								},
								"logo": schema.StringAttribute{
									Description: "Supported image types are JPEG, PNG and JPG, size is limited to width 450px, height 130px. This should be a base 64 encoding of the full HTML img tag compatible image, with the header included.",
									Computed:    true,
								},
								"text": schema.StringAttribute{
									Description: "The text to display in the body of the custom page",
									Computed:    true,
								},
								"title": schema.StringAttribute{
									Description: "The text to display in the title of the custom page",
									Computed:    true,
								},
							},
						},
						"captcha": schema.SingleNestedAttribute{
							Computed:   true,
							CustomType: customfield.NewNestedObjectType[WaapCustomPageSetsCaptchaDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									Description: "Indicates whether the custom custom page is active or inactive",
									Computed:    true,
								},
								"error": schema.StringAttribute{
									Description: "Error message",
									Computed:    true,
								},
								"header": schema.StringAttribute{
									Description: "The text to display in the header of the custom page",
									Computed:    true,
								},
								"logo": schema.StringAttribute{
									Description: "Supported image types are JPEG, PNG and JPG, size is limited to width 450px, height 130px. This should be a base 64 encoding of the full HTML img tag compatible image, with the header included.",
									Computed:    true,
								},
								"text": schema.StringAttribute{
									Description: "The text to display in the body of the custom page",
									Computed:    true,
								},
								"title": schema.StringAttribute{
									Description: "The text to display in the title of the custom page",
									Computed:    true,
								},
							},
						},
						"cookie_disabled": schema.SingleNestedAttribute{
							Computed:   true,
							CustomType: customfield.NewNestedObjectType[WaapCustomPageSetsCookieDisabledDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									Description: "Indicates whether the custom custom page is active or inactive",
									Computed:    true,
								},
								"header": schema.StringAttribute{
									Description: "The text to display in the header of the custom page",
									Computed:    true,
								},
								"text": schema.StringAttribute{
									Description: "The text to display in the body of the custom page",
									Computed:    true,
								},
							},
						},
						"domains": schema.ListAttribute{
							Description: "List of domain IDs that are associated with this page set",
							Computed:    true,
							CustomType:  customfield.NewListType[types.Int64](ctx),
							ElementType: types.Int64Type,
						},
						"handshake": schema.SingleNestedAttribute{
							Computed:   true,
							CustomType: customfield.NewNestedObjectType[WaapCustomPageSetsHandshakeDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									Description: "Indicates whether the custom custom page is active or inactive",
									Computed:    true,
								},
								"header": schema.StringAttribute{
									Description: "The text to display in the header of the custom page",
									Computed:    true,
								},
								"logo": schema.StringAttribute{
									Description: "Supported image types are JPEG, PNG and JPG, size is limited to width 450px, height 130px. This should be a base 64 encoding of the full HTML img tag compatible image, with the header included.",
									Computed:    true,
								},
								"title": schema.StringAttribute{
									Description: "The text to display in the title of the custom page",
									Computed:    true,
								},
							},
						},
						"javascript_disabled": schema.SingleNestedAttribute{
							Computed:   true,
							CustomType: customfield.NewNestedObjectType[WaapCustomPageSetsJavascriptDisabledDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									Description: "Indicates whether the custom custom page is active or inactive",
									Computed:    true,
								},
								"header": schema.StringAttribute{
									Description: "The text to display in the header of the custom page",
									Computed:    true,
								},
								"text": schema.StringAttribute{
									Description: "The text to display in the body of the custom page",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *WaapCustomPageSetsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *WaapCustomPageSetsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
