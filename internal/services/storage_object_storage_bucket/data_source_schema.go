// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage_bucket

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*StorageObjectStorageBucketDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Buckets are containers within object storage that hold files (objects) and define their CORS, lifecycle, and access policy configuration.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"storage_id": schema.Int64Attribute{
				Required: true,
			},
			"cors": schema.SingleNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectType[StorageObjectStorageBucketCorsDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"allowed_origins": schema.ListAttribute{
						Description: `Web domains allowed to make direct browser requests to this bucket (e.g., "https://myapp.com"). Use "*" to allow any origin.`,
						Computed:    true,
						CustomType:  customfield.NewListType[types.String](ctx),
						ElementType: types.StringType,
					},
				},
			},
			"policy": schema.SingleNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectType[StorageObjectStorageBucketPolicyDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"is_public": schema.BoolAttribute{
						Description: "When true, anyone can download objects without credentials. When false, all requests require valid S3 authentication.",
						Computed:    true,
					},
				},
			},
			"storage_object_storage_bucket_lifecycle": schema.SingleNestedAttribute{
				Computed:   true,
				CustomType: customfield.NewNestedObjectType[StorageObjectStorageBucketStorageObjectStorageBucketLifecycleDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"expiration_days": schema.Int64Attribute{
						Description: "Days after upload before objects are automatically deleted. For example, 30 means files are removed 30 days after creation.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *StorageObjectStorageBucketDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *StorageObjectStorageBucketDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
