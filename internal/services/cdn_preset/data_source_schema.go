// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_preset

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CDNPresetDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "CDN presets are predefined sets of CDN resource or rule settings that can be applied to an object in a single request, letting you configure caching, delivery, and security options consistently.",
		Attributes: map[string]schema.Attribute{
			"preset_id": schema.Int64Attribute{
				Required: true,
			},
			"id": schema.Int64Attribute{
				Description: "Preset ID.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Preset name.",
				Computed:    true,
			},
			"object_type": schema.StringAttribute{
				Description: "Type of object the preset can be applied to.\n\nPossible values:\n- **CDNResource** - Preset is applied to a CDN resource.\n- **Rule** - Preset is applied to a rule.\nAvailable values: \"CDNResource\", \"Rule\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("CDNResource", "Rule"),
				},
			},
			"service": schema.StringAttribute{
				Description: "Service the preset belongs to.\nAvailable values: \"cdn\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("cdn"),
				},
			},
			"preset_settings": schema.MapAttribute{
				Description: "CDN resource or rule settings that the preset applies to the target object, including the **options** object.\n\nThe available keys match the writable fields of the object type the preset targets (`object_type`).\nOptions included in the preset cannot be edited on the object while the preset is applied.",
				Computed:    true,
				CustomType:  customfield.NewMapType[jsontypes.Normalized](ctx),
				ElementType: jsontypes.NormalizedType{},
			},
		},
	}
}

func (d *CDNPresetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CDNPresetDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
