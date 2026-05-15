// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage_bucket

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*StorageObjectStorageBucketsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Buckets are containers within object storage that hold files (objects) and define their CORS, lifecycle, and access policy configuration.",
		Attributes: map[string]schema.Attribute{
			"storage_id": schema.Int64Attribute{
				Required: true,
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
				CustomType:  customfield.NewNestedObjectListType[StorageObjectStorageBucketsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cors": schema.SingleNestedAttribute{
							Computed:   true,
							CustomType: customfield.NewNestedObjectType[StorageObjectStorageBucketsCorsDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"allowed_origins": schema.ListAttribute{
									Description: `Web domains allowed to make direct browser requests to this bucket (e.g., "https://myapp.com"). Use "*" to allow any origin.`,
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
							},
						},
						"storage_object_storage_bucket_lifecycle": schema.SingleNestedAttribute{
							Computed:   true,
							CustomType: customfield.NewNestedObjectType[StorageObjectStorageBucketsStorageObjectStorageBucketLifecycleDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"expiration_days": schema.Int64Attribute{
									Description: "Days after upload before objects are automatically deleted. For example, 30 means files are removed 30 days after creation.",
									Computed:    true,
								},
							},
						},
						"name": schema.StringAttribute{
							Description: "Globally unique bucket name within the storage. Used as the path prefix when accessing objects via S3 API.",
							Computed:    true,
						},
						"policy": schema.SingleNestedAttribute{
							Computed:   true,
							CustomType: customfield.NewNestedObjectType[StorageObjectStorageBucketsPolicyDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"is_public": schema.BoolAttribute{
									Description: "When true, anyone can download objects without credentials. When false, all requests require valid S3 authentication.",
									Computed:    true,
								},
							},
						},
						"storage_id": schema.Int64Attribute{
							Description: "Parent storage this bucket belongs to. Use this ID in the URL path for bucket operations.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *StorageObjectStorageBucketsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *StorageObjectStorageBucketsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
