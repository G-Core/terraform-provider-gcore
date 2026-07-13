// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_target

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CDNLogsUploaderTargetResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	emptyResponseActions := types.ListValueMust(
		types.ObjectType{AttrTypes: map[string]attr.Type{
			"action":            types.StringType,
			"description":       types.StringType,
			"match_payload":     types.StringType,
			"match_status_code": types.Int64Type,
		}},
		[]attr.Value{},
	)

	return schema.Schema{
		MarkdownDescription: "Logs uploader targets define destinations for CDN log delivery, such as S3 buckets or SFTP servers, with associated authentication and configuration settings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseNonNullStateForUnknown()},
			},
			"storage_type": schema.StringAttribute{
				Description: "Type of storage for logs.\nAvailable values: \"s3_gcore\", \"s3_amazon\", \"s3_oss\", \"s3_other\", \"s3_v1\", \"ftp\", \"sftp\", \"http\", \"azure_blob\", \"sls\".",
				Required:    true,
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
			"config": schema.SingleNestedAttribute{
				Description: "Config for specific storage type.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"access_key_id": schema.StringAttribute{
						Optional: true,
					},
					"bucket_name": schema.StringAttribute{
						Optional: true,
					},
					"directory": schema.StringAttribute{
						Optional: true,
					},
					"endpoint": schema.StringAttribute{
						Optional: true,
					},
					"region": schema.StringAttribute{
						Optional: true,
					},
					"secret_access_key": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"use_path_style": schema.BoolAttribute{
						Computed:      true,
						Optional:      true,
						PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
					},
					"hostname": schema.StringAttribute{
						Optional: true,
					},
					"password": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"timeout_seconds": schema.Int64Attribute{
						Optional: true,
						Validators: []validator.Int64{
							int64validator.Between(0, 300),
						},
					},
					"user": schema.StringAttribute{
						Optional: true,
					},
					"key_passphrase": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"private_key": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"upload": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"url": schema.StringAttribute{
								Required: true,
							},
							"headers": schema.MapAttribute{
								Optional:    true,
								ElementType: types.StringType,
							},
							"method": schema.StringAttribute{
								Description: `Available values: "POST", "PUT".`,
								Computed:    true,
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("POST", "PUT"),
								},
								Default: stringdefault.StaticString("POST"),
							},
							"response_actions": schema.ListNestedAttribute{
								Computed: true,
								Optional: true,
								Default:  listdefault.StaticValue(emptyResponseActions),
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"action": schema.StringAttribute{
											Description: `Available values: "drop", "retry", "append".`,
											Required:    true,
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
											Optional: true,
											Default:  stringdefault.StaticString(""),
										},
										"match_payload": schema.StringAttribute{
											Computed: true,
											Optional: true,
											Default:  stringdefault.StaticString(""),
										},
										"match_status_code": schema.Int64Attribute{
											Computed: true,
											Optional: true,
											Default:  int64default.StaticInt64(0),
										},
									},
								},
							},
							"timeout_seconds": schema.Int64Attribute{
								Computed: true,
								Optional: true,
								Default:  int64default.StaticInt64(30),
							},
							"use_compression": schema.BoolAttribute{
								Computed: true,
								Optional: true,
								Default:  booldefault.StaticBool(false),
							},
						},
					},
					"append": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"url": schema.StringAttribute{
								Required: true,
							},
							"headers": schema.MapAttribute{
								Optional:    true,
								ElementType: types.StringType,
							},
							"method": schema.StringAttribute{
								Description: `Available values: "POST", "PUT".`,
								Computed:    true,
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("POST", "PUT"),
								},
								Default: stringdefault.StaticString("POST"),
							},
							"response_actions": schema.ListNestedAttribute{
								Computed: true,
								Optional: true,
								Default:  listdefault.StaticValue(emptyResponseActions),
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"action": schema.StringAttribute{
											Description: `Available values: "drop", "retry", "append".`,
											Required:    true,
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
											Optional: true,
											Default:  stringdefault.StaticString(""),
										},
										"match_payload": schema.StringAttribute{
											Computed: true,
											Optional: true,
											Default:  stringdefault.StaticString(""),
										},
										"match_status_code": schema.Int64Attribute{
											Computed: true,
											Optional: true,
											Default:  int64default.StaticInt64(0),
										},
									},
								},
							},
							"timeout_seconds": schema.Int64Attribute{
								Computed: true,
								Optional: true,
								Default:  int64default.StaticInt64(30),
							},
							"use_compression": schema.BoolAttribute{
								Computed: true,
								Optional: true,
								Default:  booldefault.StaticBool(false),
							},
						},
					},
					"auth": schema.SingleNestedAttribute{
						Optional:   true,
						Validators: []validator.Object{authConfigValidator{}},
						Attributes: map[string]schema.Attribute{
							"config": schema.SingleNestedAttribute{
								Required: true,
								Attributes: map[string]schema.Attribute{
									"token": schema.StringAttribute{
										Optional:  true,
										Sensitive: true,
									},
									"header_name": schema.StringAttribute{
										Optional: true,
									},
									"account_key": schema.StringAttribute{
										Description: "Azure Blob Storage account key.",
										Optional:    true,
										Sensitive:   true,
									},
									"access_key_id": schema.StringAttribute{
										Description: "Alibaba access key ID.",
										Optional:    true,
									},
									"secret_access_key": schema.StringAttribute{
										Description: "Alibaba secret access key.",
										Optional:    true,
										Sensitive:   true,
									},
								},
							},
							"type": schema.StringAttribute{
								Description: `Available values: "token", "shared_key", "sas_token", "ak_sk".`,
								Required:    true,
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
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive("json", "text"),
						},
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
					"retry": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"url": schema.StringAttribute{
								Required: true,
							},
							"headers": schema.MapAttribute{
								Optional:    true,
								ElementType: types.StringType,
							},
							"method": schema.StringAttribute{
								Description: `Available values: "POST", "PUT".`,
								Computed:    true,
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive("POST", "PUT"),
								},
								Default: stringdefault.StaticString("POST"),
							},
							"response_actions": schema.ListNestedAttribute{
								Computed: true,
								Optional: true,
								Default:  listdefault.StaticValue(emptyResponseActions),
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"action": schema.StringAttribute{
											Description: `Available values: "drop", "retry", "append".`,
											Required:    true,
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
											Optional: true,
											Default:  stringdefault.StaticString(""),
										},
										"match_payload": schema.StringAttribute{
											Computed: true,
											Optional: true,
											Default:  stringdefault.StaticString(""),
										},
										"match_status_code": schema.Int64Attribute{
											Computed: true,
											Optional: true,
											Default:  int64default.StaticInt64(0),
										},
									},
								},
							},
							"timeout_seconds": schema.Int64Attribute{
								Computed: true,
								Optional: true,
								Default:  int64default.StaticInt64(30),
							},
							"use_compression": schema.BoolAttribute{
								Computed: true,
								Optional: true,
								Default:  booldefault.StaticBool(false),
							},
						},
					},
					"account_name": schema.StringAttribute{
						Description: "Azure Blob Storage account name.",
						Optional:    true,
					},
					"container_name": schema.StringAttribute{
						Description: "Azure Blob Storage container name.",
						Optional:    true,
					},
					"log_store": schema.StringAttribute{
						Description: "SLS logstore name. 3-36 characters; lowercase letters, digits, hyphens, and underscores.",
						Optional:    true,
					},
					"project": schema.StringAttribute{
						Description: "SLS project name. 3-63 characters; lowercase letters, digits, and hyphens.",
						Optional:    true,
					},
					"topic": schema.StringAttribute{
						Description: "Optional SLS topic (0-128 characters).",
						Optional:    true,
					},
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the target.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
			},
			"name": schema.StringAttribute{
				Description: "Name of the target.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("Target"),
			},
			"client_id": schema.Int64Attribute{
				Description:   "Client that owns the target.",
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"created": schema.StringAttribute{
				Description:   "Time when logs uploader target was created.",
				Computed:      true,
				CustomType:    timetypes.RFC3339Type{},
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *CDNLogsUploaderTargetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CDNLogsUploaderTargetResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
