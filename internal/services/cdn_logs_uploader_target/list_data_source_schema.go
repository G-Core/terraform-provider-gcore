// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_target

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CDNLogsUploaderTargetsDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Logs uploader targets define destinations for CDN log delivery, such as S3 buckets or SFTP servers, with associated authentication and configuration settings.",
		Attributes: map[string]schema.Attribute{
			"search": schema.StringAttribute{
				Description: "Search by target name or id.",
				Optional:    true,
			},
			"config_ids": schema.ListAttribute{
				Description: "Filter by ids of related logs uploader configs that use given target.",
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
				CustomType:  customfield.NewNestedObjectListType[CDNLogsUploaderTargetsItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"client_id": schema.Int64Attribute{
							Description: "Client that owns the target.",
							Computed:    true,
						},
						"config": schema.SingleNestedAttribute{
							Description: "Config for specific storage type.",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CDNLogsUploaderTargetsConfigDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"access_key_id": schema.StringAttribute{
									Computed: true,
								},
								"bucket_name": schema.StringAttribute{
									Computed: true,
								},
								"directory": schema.StringAttribute{
									Computed: true,
								},
								"endpoint": schema.StringAttribute{
									Computed: true,
								},
								"region": schema.StringAttribute{
									Computed: true,
								},
								"use_path_style": schema.BoolAttribute{
									Computed: true,
								},
								"hostname": schema.StringAttribute{
									Computed: true,
								},
								"timeout_seconds": schema.Int64Attribute{
									Computed: true,
									Validators: []validator.Int64{
										int64validator.Between(0, 300),
									},
								},
								"user": schema.StringAttribute{
									Computed: true,
								},
								"key_passphrase": schema.StringAttribute{
									Computed: true,
								},
								"password": schema.StringAttribute{
									Computed: true,
								},
								"private_key": schema.StringAttribute{
									Computed: true,
								},
								"append": schema.SingleNestedAttribute{
									Computed:   true,
									CustomType: customfield.NewNestedObjectType[CDNLogsUploaderTargetsConfigAppendDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"headers": schema.MapAttribute{
											Computed:    true,
											CustomType:  customfield.NewMapType[types.String](ctx),
											ElementType: types.StringType,
										},
										"method": schema.StringAttribute{
											Description: `Available values: "POST", "PUT".`,
											Computed:    true,
											Validators: []validator.String{
												stringvalidator.OneOfCaseInsensitive("POST", "PUT"),
											},
										},
										"response_actions": schema.ListNestedAttribute{
											Computed:   true,
											CustomType: customfield.NewNestedObjectListType[CDNLogsUploaderTargetsConfigAppendResponseActionsDataSourceModel](ctx),
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"action": schema.StringAttribute{
														Description: `Available values: "drop", "retry", "append".`,
														Computed:    true,
														Validators: []validator.String{
															stringvalidator.OneOfCaseInsensitive(
																"drop",
																"retry",
																"append",
															),
														},
													},
													"description": schema.StringAttribute{
														Computed: true,
													},
													"match_payload": schema.StringAttribute{
														Computed: true,
													},
													"match_status_code": schema.Int64Attribute{
														Computed: true,
													},
												},
											},
										},
										"timeout_seconds": schema.Int64Attribute{
											Computed: true,
										},
										"url": schema.StringAttribute{
											Computed: true,
										},
										"use_compression": schema.BoolAttribute{
											Computed: true,
										},
									},
								},
								"auth": schema.SingleNestedAttribute{
									Computed:   true,
									CustomType: customfield.NewNestedObjectType[CDNLogsUploaderTargetsConfigAuthDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"config": schema.SingleNestedAttribute{
											Computed:   true,
											CustomType: customfield.NewNestedObjectType[CDNLogsUploaderTargetsConfigAuthConfigDataSourceModel](ctx),
											Attributes: map[string]schema.Attribute{
												"token": schema.StringAttribute{
													Computed: true,
												},
												"header_name": schema.StringAttribute{
													Computed: true,
												},
												"account_key": schema.StringAttribute{
													Description: "Masked secret value.",
													Computed:    true,
												},
												"access_key_id": schema.StringAttribute{
													Description: "Alibaba access key ID.",
													Computed:    true,
												},
												"secret_access_key": schema.StringAttribute{
													Description: "Masked secret value.",
													Computed:    true,
												},
											},
										},
										"type": schema.StringAttribute{
											Description: `Available values: "token", "shared_key", "sas_token", "ak_sk".`,
											Computed:    true,
											Validators: []validator.String{
												stringvalidator.OneOfCaseInsensitive(
													"token",
													"shared_key",
													"sas_token",
													"ak_sk",
												),
											},
										},
									},
								},
								"content_type": schema.StringAttribute{
									Description: `Available values: "json", "text".`,
									Computed:    true,
									Validators: []validator.String{
										stringvalidator.OneOfCaseInsensitive("json", "text"),
									},
								},
								"retry": schema.SingleNestedAttribute{
									Computed:   true,
									CustomType: customfield.NewNestedObjectType[CDNLogsUploaderTargetsConfigRetryDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"headers": schema.MapAttribute{
											Computed:    true,
											CustomType:  customfield.NewMapType[types.String](ctx),
											ElementType: types.StringType,
										},
										"method": schema.StringAttribute{
											Description: `Available values: "POST", "PUT".`,
											Computed:    true,
											Validators: []validator.String{
												stringvalidator.OneOfCaseInsensitive("POST", "PUT"),
											},
										},
										"response_actions": schema.ListNestedAttribute{
											Computed:   true,
											CustomType: customfield.NewNestedObjectListType[CDNLogsUploaderTargetsConfigRetryResponseActionsDataSourceModel](ctx),
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"action": schema.StringAttribute{
														Description: `Available values: "drop", "retry", "append".`,
														Computed:    true,
														Validators: []validator.String{
															stringvalidator.OneOfCaseInsensitive(
																"drop",
																"retry",
																"append",
															),
														},
													},
													"description": schema.StringAttribute{
														Computed: true,
													},
													"match_payload": schema.StringAttribute{
														Computed: true,
													},
													"match_status_code": schema.Int64Attribute{
														Computed: true,
													},
												},
											},
										},
										"timeout_seconds": schema.Int64Attribute{
											Computed: true,
										},
										"url": schema.StringAttribute{
											Computed: true,
										},
										"use_compression": schema.BoolAttribute{
											Computed: true,
										},
									},
								},
								"upload": schema.SingleNestedAttribute{
									Computed:   true,
									CustomType: customfield.NewNestedObjectType[CDNLogsUploaderTargetsConfigUploadDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"headers": schema.MapAttribute{
											Computed:    true,
											CustomType:  customfield.NewMapType[types.String](ctx),
											ElementType: types.StringType,
										},
										"method": schema.StringAttribute{
											Description: `Available values: "POST", "PUT".`,
											Computed:    true,
											Validators: []validator.String{
												stringvalidator.OneOfCaseInsensitive("POST", "PUT"),
											},
										},
										"response_actions": schema.ListNestedAttribute{
											Computed:   true,
											CustomType: customfield.NewNestedObjectListType[CDNLogsUploaderTargetsConfigUploadResponseActionsDataSourceModel](ctx),
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"action": schema.StringAttribute{
														Description: `Available values: "drop", "retry", "append".`,
														Computed:    true,
														Validators: []validator.String{
															stringvalidator.OneOfCaseInsensitive(
																"drop",
																"retry",
																"append",
															),
														},
													},
													"description": schema.StringAttribute{
														Computed: true,
													},
													"match_payload": schema.StringAttribute{
														Computed: true,
													},
													"match_status_code": schema.Int64Attribute{
														Computed: true,
													},
												},
											},
										},
										"timeout_seconds": schema.Int64Attribute{
											Computed: true,
										},
										"url": schema.StringAttribute{
											Computed: true,
										},
										"use_compression": schema.BoolAttribute{
											Computed: true,
										},
									},
								},
								"account_name": schema.StringAttribute{
									Description: "Azure Blob Storage account name.",
									Computed:    true,
								},
								"container_name": schema.StringAttribute{
									Description: "Azure Blob Storage container name.",
									Computed:    true,
								},
								"log_store": schema.StringAttribute{
									Description: "SLS logstore name. 3-36 characters; lowercase letters, digits, hyphens, and underscores.",
									Computed:    true,
								},
								"project": schema.StringAttribute{
									Description: "SLS project name. 3-63 characters; lowercase letters, digits, and hyphens.",
									Computed:    true,
								},
								"topic": schema.StringAttribute{
									Description: "Optional SLS topic (0-128 characters).",
									Computed:    true,
								},
							},
						},
						"created": schema.StringAttribute{
							Description: "Time when logs uploader target was created.",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
						"description": schema.StringAttribute{
							Description: "Description of the target.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the target.",
							Computed:    true,
						},
						"related_uploader_configs": schema.ListAttribute{
							Description: "List of logs uploader configs that use this target.",
							Computed:    true,
							CustomType:  customfield.NewListType[types.Int64](ctx),
							ElementType: types.Int64Type,
						},
						"status": schema.StringAttribute{
							Description: "Validation status of the logs uploader target. Informs if the specified target is reachable.",
							Computed:    true,
							CustomType:  jsontypes.NormalizedType{},
						},
						"storage_type": schema.StringAttribute{
							Description: "Type of storage for logs.\nAvailable values: \"s3_gcore\", \"s3_amazon\", \"s3_oss\", \"s3_other\", \"s3_v1\", \"ftp\", \"sftp\", \"http\", \"azure_blob\", \"sls\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"s3_gcore",
									"s3_amazon",
									"s3_oss",
									"s3_other",
									"s3_v1",
									"ftp",
									"sftp",
									"http",
									"azure_blob",
									"sls",
								),
							},
						},
						"updated": schema.StringAttribute{
							Description: "Time when logs uploader target was updated.",
							Computed:    true,
							CustomType:  timetypes.RFC3339Type{},
						},
					},
				},
			},
		},
	}
}

func (d *CDNLogsUploaderTargetsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CDNLogsUploaderTargetsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
