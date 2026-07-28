// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry_artifact

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudRegistryArtifactsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"registry_id": schema.Int64Attribute{
				Required: true,
			},
			"repository_name": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.Int64Attribute{
				Optional: true,
			},
			"region_id": schema.Int64Attribute{
				Optional: true,
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
				CustomType:  customfield.NewNestedObjectListType[CloudRegistryArtifactsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Repository ID",
							Computed:    true,
						},
						"digest": schema.StringAttribute{
							Description: "Artifact digest",
							Computed:    true,
						},
						"pulled_at": schema.StringAttribute{
							Description: "Artifact last pull date-time",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"pushed_at": schema.StringAttribute{
							Description: "Artifact push date-time",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"registry_id": schema.Int64Attribute{
							Description: "Artifact registry ID",
							Computed:    true,
						},
						"repository_id": schema.Int64Attribute{
							Description: "Artifact repository ID",
							Computed:    true,
						},
						"size": schema.Int64Attribute{
							Description: "Artifact size, bytes",
							Computed:    true,
						},
						"tags": schema.ListNestedAttribute{
							Description: "Artifact tags",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectListType[CloudRegistryArtifactsTagsDataSourceModel](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										Description: "Tag ID",
										Computed:    true,
									},
									"artifact_id": schema.Int64Attribute{
										Description: "Artifact ID",
										Computed:    true,
									},
									"name": schema.StringAttribute{
										Description: "Tag name",
										Computed:    true,
									},
									"pulled_at": schema.StringAttribute{
										Description: "Tag last pull date-time",
										Computed:    true,
										CustomType:  timetypes.RFC3339Type{},
									},
									"pushed_at": schema.StringAttribute{
										Description: "Tag push date-time",
										Computed:    true,
										CustomType:  timetypes.RFC3339Type{},
									},
									"repository_id": schema.Int64Attribute{
										Description: "Repository ID",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *CloudRegistryArtifactsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudRegistryArtifactsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
