// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_deployment

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudInferenceDeploymentDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Inference deployments run containerized ML models with configurable scaling, health probes, and GPU flavors.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Inference instance name.",
				Computed:    true,
			},
			"deployment_name": schema.StringAttribute{
				Description: "Inference instance name.",
				Required:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID",
				Optional:    true,
			},
			"address": schema.StringAttribute{
				Description: "Address of the inference instance",
				Computed:    true,
			},
			"auth_enabled": schema.BoolAttribute{
				Description:        "`true` if instance uses API key authentication. `\"Authorization\": \"Bearer *****\"` or `\"X-Api-Key\": \"*****\"` header is required for the requests to the instance if enabled.",
				Computed:           true,
				DeprecationMessage: "This attribute is deprecated.",
			},
			"command": schema.StringAttribute{
				Description: "Command to be executed when running a container from an image.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Inference instance creation date in ISO 8601 format.",
				Computed:    true,
			},
			"credentials_name": schema.StringAttribute{
				Description: "Registry credentials name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Inference instance description.",
				Computed:    true,
			},
			"flavor_name": schema.StringAttribute{
				Description: "Flavor name for the inference instance",
				Computed:    true,
			},
			"image": schema.StringAttribute{
				Description: "Docker image for the inference instance. This field should contain the image name and tag in the format 'name:tag', e.g., 'nginx:latest'. It defaults to Docker Hub as the image registry, but any accessible Docker image URL can be specified.",
				Computed:    true,
			},
			"listening_port": schema.Int64Attribute{
				Description: "Listening port for the inference instance.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Inference instance name.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Inference instance status.\n\nValue can be one of the following:\n- `DEPLOYING` - The instance is being deployed. Containers are not yet created.\n- `PARTIALLYDEPLOYED` - All containers have been created, but some may not be ready yet. Instances stuck in this state typically indicate either image being pulled, or a failure of some kind. In the latter case, the `error_message` field of the respective container object in the `containers` collection explains the failure reason.\n- `ACTIVE` - The instance is running and ready to accept requests.\n- `DISABLED` - The instance is disabled and not accepting any requests.\n- `PENDING` - The instance is running but scaled to zero. It will be automatically scaled up when a request is made.\n- `DELETING` - The instance is being deleted.\nAvailable values: \"ACTIVE\", \"DELETING\", \"DEPLOYING\", \"DISABLED\", \"PARTIALLYDEPLOYED\", \"PENDING\".",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"ACTIVE",
						"DELETING",
						"DEPLOYING",
						"DISABLED",
						"PARTIALLYDEPLOYED",
						"PENDING",
					),
				},
			},
			"timeout": schema.Int64Attribute{
				Description: "Specifies the duration in seconds without any requests after which the containers will be downscaled to their minimum scale value as defined by `scale.min`. If set, this helps in optimizing resource usage by reducing the number of container instances during periods of inactivity.",
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"api_keys": schema.ListAttribute{
				Description: "List of API keys for the inference instance",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"envs": schema.MapAttribute{
				Description: "Environment variables for the inference instance",
				Computed:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"containers": schema.ListNestedAttribute{
				Description: "List of containers for the inference instance",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudInferenceDeploymentContainersDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Description: "Address of the inference instance",
							Computed:    true,
						},
						"deploy_status": schema.SingleNestedAttribute{
							Description: "Status of the containers deployment",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentContainersDeployStatusDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"ready": schema.Int64Attribute{
									Description: "Number of ready instances",
									Computed:    true,
								},
								"total": schema.Int64Attribute{
									Description: "Total number of instances",
									Computed:    true,
								},
							},
						},
						"error_message": schema.StringAttribute{
							Description: "Error message if the container deployment failed",
							Computed:    true,
						},
						"region_id": schema.Int64Attribute{
							Description: "Region name for the container",
							Computed:    true,
						},
						"scale": schema.SingleNestedAttribute{
							Description: "Scale for the container",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentContainersScaleDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"cooldown_period": schema.Int64Attribute{
									Description: "Cooldown period between scaling actions in seconds",
									Computed:    true,
								},
								"max": schema.Int64Attribute{
									Description: "Maximum scale for the container",
									Computed:    true,
								},
								"min": schema.Int64Attribute{
									Description: "Minimum scale for the container",
									Computed:    true,
								},
								"polling_interval": schema.Int64Attribute{
									Description: "Polling interval for scaling triggers in seconds",
									Computed:    true,
								},
								"triggers": schema.SingleNestedAttribute{
									Description: "Triggers for scaling actions",
									Computed:    true,
									CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentContainersScaleTriggersDataSourceModel](ctx),
									Attributes: map[string]schema.Attribute{
										"cpu": schema.SingleNestedAttribute{
											Description: "CPU trigger configuration",
											Computed:    true,
											CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentContainersScaleTriggersCPUDataSourceModel](ctx),
											Attributes: map[string]schema.Attribute{
												"threshold": schema.Int64Attribute{
													Description: "Threshold value for the trigger in percentage",
													Computed:    true,
												},
											},
										},
										"gpu_memory": schema.SingleNestedAttribute{
											Description: "GPU memory trigger configuration. Calculated by `DCGM_FI_DEV_MEM_COPY_UTIL` metric",
											Computed:    true,
											CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentContainersScaleTriggersGPUMemoryDataSourceModel](ctx),
											Attributes: map[string]schema.Attribute{
												"threshold": schema.Int64Attribute{
													Description: "Threshold value for the trigger in percentage",
													Computed:    true,
												},
											},
										},
										"gpu_utilization": schema.SingleNestedAttribute{
											Description: "GPU utilization trigger configuration. Calculated by `DCGM_FI_DEV_GPU_UTIL` metric",
											Computed:    true,
											CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentContainersScaleTriggersGPUUtilizationDataSourceModel](ctx),
											Attributes: map[string]schema.Attribute{
												"threshold": schema.Int64Attribute{
													Description: "Threshold value for the trigger in percentage",
													Computed:    true,
												},
											},
										},
										"http": schema.SingleNestedAttribute{
											Description: "HTTP trigger configuration",
											Computed:    true,
											CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentContainersScaleTriggersHTTPDataSourceModel](ctx),
											Attributes: map[string]schema.Attribute{
												"rate": schema.Int64Attribute{
													Description: "Request count per 'window' seconds for the http trigger",
													Computed:    true,
												},
												"window": schema.Int64Attribute{
													Description: "Time window for rate calculation in seconds",
													Computed:    true,
												},
											},
										},
										"memory": schema.SingleNestedAttribute{
											Description: "Memory trigger configuration",
											Computed:    true,
											CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentContainersScaleTriggersMemoryDataSourceModel](ctx),
											Attributes: map[string]schema.Attribute{
												"threshold": schema.Int64Attribute{
													Description: "Threshold value for the trigger in percentage",
													Computed:    true,
												},
											},
										},
										"sqs": schema.SingleNestedAttribute{
											Description: "SQS trigger configuration",
											Computed:    true,
											CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentContainersScaleTriggersSqsDataSourceModel](ctx),
											Attributes: map[string]schema.Attribute{
												"activation_queue_length": schema.Int64Attribute{
													Description: "Number of messages for activation",
													Computed:    true,
												},
												"aws_endpoint": schema.StringAttribute{
													Description: "Custom AWS endpoint",
													Computed:    true,
												},
												"aws_region": schema.StringAttribute{
													Description: "AWS region",
													Computed:    true,
												},
												"queue_length": schema.Int64Attribute{
													Description: "Number of messages for one replica",
													Computed:    true,
												},
												"queue_url": schema.StringAttribute{
													Description: "SQS queue URL",
													Computed:    true,
												},
												"scale_on_delayed": schema.BoolAttribute{
													Description: "Scale on delayed messages",
													Computed:    true,
												},
												"scale_on_flight": schema.BoolAttribute{
													Description: "Scale on in-flight messages",
													Computed:    true,
												},
												"secret_name": schema.StringAttribute{
													Description: "Auth secret name",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"ingress_opts": schema.SingleNestedAttribute{
				Description: "Ingress options for the inference instance",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentIngressOptsDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"disable_response_buffering": schema.BoolAttribute{
						Description: "Disable response buffering if true. A client usually has a much slower connection and can not consume the response data as fast as it is produced by an upstream application. Ingress tries to buffer the whole response in order to release the upstream application as soon as possible.By default, the response buffering is enabled.",
						Computed:    true,
					},
				},
			},
			"logging": schema.SingleNestedAttribute{
				Description: "Logging configuration for the inference instance",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentLoggingDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"destination_region_id": schema.Int64Attribute{
						Description: "ID of the region in which the logs will be stored",
						Computed:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "Indicates if log streaming is enabled or disabled",
						Computed:    true,
					},
					"topic_name": schema.StringAttribute{
						Description: "The topic name to stream logs to",
						Computed:    true,
					},
					"retention_policy": schema.SingleNestedAttribute{
						Description: "Logs retention policy",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentLoggingRetentionPolicyDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"period": schema.Int64Attribute{
								Description: "Duration of days for which logs must be kept.",
								Computed:    true,
								Validators: []validator.Int64{
									int64validator.AtMost(1825),
								},
							},
						},
					},
				},
			},
			"object_references": schema.ListNestedAttribute{
				Description: "Indicates to which parent object this inference belongs to.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudInferenceDeploymentObjectReferencesDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"kind": schema.StringAttribute{
							Description: "Kind of the inference object to be referenced\nAvailable values: \"AppDeployment\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive("AppDeployment"),
							},
						},
						"name": schema.StringAttribute{
							Description: "Name of the inference object to be referenced",
							Computed:    true,
						},
					},
				},
			},
			"probes": schema.SingleNestedAttribute{
				Description: "Probes configured for all containers of the inference instance.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesDataSourceModel](ctx),
				Attributes: map[string]schema.Attribute{
					"liveness_probe": schema.SingleNestedAttribute{
						Description: "Liveness probe configuration",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesLivenessProbeDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Whether the probe is enabled or not.",
								Computed:    true,
							},
							"probe": schema.SingleNestedAttribute{
								Description: "Probe configuration (exec, `http_get` or `tcp_socket`)",
								Computed:    true,
								CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesLivenessProbeProbeDataSourceModel](ctx),
								Attributes: map[string]schema.Attribute{
									"exec": schema.SingleNestedAttribute{
										Description: "Exec probe configuration",
										Computed:    true,
										CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesLivenessProbeProbeExecDataSourceModel](ctx),
										Attributes: map[string]schema.Attribute{
											"command": schema.ListAttribute{
												Description: "Command to be executed inside the running container.",
												Computed:    true,
												CustomType:  customfield.NewListType[types.String](ctx),
												ElementType: types.StringType,
											},
										},
									},
									"failure_threshold": schema.Int64Attribute{
										Description: "The number of consecutive probe failures that mark the container as unhealthy.",
										Computed:    true,
									},
									"http_get": schema.SingleNestedAttribute{
										Description: "HTTP GET probe configuration",
										Computed:    true,
										CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesLivenessProbeProbeHTTPGetDataSourceModel](ctx),
										Attributes: map[string]schema.Attribute{
											"headers": schema.MapAttribute{
												Description: "HTTP headers to be sent with the request.",
												Computed:    true,
												CustomType:  customfield.NewMapType[types.String](ctx),
												ElementType: types.StringType,
											},
											"host": schema.StringAttribute{
												Description: "Host name to send HTTP request to.",
												Computed:    true,
											},
											"path": schema.StringAttribute{
												Description: "The endpoint to send the HTTP request to.",
												Computed:    true,
											},
											"port": schema.Int64Attribute{
												Description: "Port number the probe should connect to.",
												Computed:    true,
											},
											"schema": schema.StringAttribute{
												Description: "Schema to use for the HTTP request.",
												Computed:    true,
											},
										},
									},
									"initial_delay_seconds": schema.Int64Attribute{
										Description: "The initial delay before starting the first probe.",
										Computed:    true,
									},
									"period_seconds": schema.Int64Attribute{
										Description: "How often (in seconds) to perform the probe.",
										Computed:    true,
									},
									"success_threshold": schema.Int64Attribute{
										Description: "The number of consecutive successful probes that mark the container as healthy.",
										Computed:    true,
									},
									"tcp_socket": schema.SingleNestedAttribute{
										Description: "TCP socket probe configuration",
										Computed:    true,
										CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesLivenessProbeProbeTcpSocketDataSourceModel](ctx),
										Attributes: map[string]schema.Attribute{
											"port": schema.Int64Attribute{
												Description: "Port number to check if it's open.",
												Computed:    true,
											},
										},
									},
									"timeout_seconds": schema.Int64Attribute{
										Description: "The timeout for each probe.",
										Computed:    true,
									},
								},
							},
						},
					},
					"readiness_probe": schema.SingleNestedAttribute{
						Description: "Readiness probe configuration",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesReadinessProbeDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Whether the probe is enabled or not.",
								Computed:    true,
							},
							"probe": schema.SingleNestedAttribute{
								Description: "Probe configuration (exec, `http_get` or `tcp_socket`)",
								Computed:    true,
								CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesReadinessProbeProbeDataSourceModel](ctx),
								Attributes: map[string]schema.Attribute{
									"exec": schema.SingleNestedAttribute{
										Description: "Exec probe configuration",
										Computed:    true,
										CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesReadinessProbeProbeExecDataSourceModel](ctx),
										Attributes: map[string]schema.Attribute{
											"command": schema.ListAttribute{
												Description: "Command to be executed inside the running container.",
												Computed:    true,
												CustomType:  customfield.NewListType[types.String](ctx),
												ElementType: types.StringType,
											},
										},
									},
									"failure_threshold": schema.Int64Attribute{
										Description: "The number of consecutive probe failures that mark the container as unhealthy.",
										Computed:    true,
									},
									"http_get": schema.SingleNestedAttribute{
										Description: "HTTP GET probe configuration",
										Computed:    true,
										CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesReadinessProbeProbeHTTPGetDataSourceModel](ctx),
										Attributes: map[string]schema.Attribute{
											"headers": schema.MapAttribute{
												Description: "HTTP headers to be sent with the request.",
												Computed:    true,
												CustomType:  customfield.NewMapType[types.String](ctx),
												ElementType: types.StringType,
											},
											"host": schema.StringAttribute{
												Description: "Host name to send HTTP request to.",
												Computed:    true,
											},
											"path": schema.StringAttribute{
												Description: "The endpoint to send the HTTP request to.",
												Computed:    true,
											},
											"port": schema.Int64Attribute{
												Description: "Port number the probe should connect to.",
												Computed:    true,
											},
											"schema": schema.StringAttribute{
												Description: "Schema to use for the HTTP request.",
												Computed:    true,
											},
										},
									},
									"initial_delay_seconds": schema.Int64Attribute{
										Description: "The initial delay before starting the first probe.",
										Computed:    true,
									},
									"period_seconds": schema.Int64Attribute{
										Description: "How often (in seconds) to perform the probe.",
										Computed:    true,
									},
									"success_threshold": schema.Int64Attribute{
										Description: "The number of consecutive successful probes that mark the container as healthy.",
										Computed:    true,
									},
									"tcp_socket": schema.SingleNestedAttribute{
										Description: "TCP socket probe configuration",
										Computed:    true,
										CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesReadinessProbeProbeTcpSocketDataSourceModel](ctx),
										Attributes: map[string]schema.Attribute{
											"port": schema.Int64Attribute{
												Description: "Port number to check if it's open.",
												Computed:    true,
											},
										},
									},
									"timeout_seconds": schema.Int64Attribute{
										Description: "The timeout for each probe.",
										Computed:    true,
									},
								},
							},
						},
					},
					"startup_probe": schema.SingleNestedAttribute{
						Description: "Startup probe configuration",
						Computed:    true,
						CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesStartupProbeDataSourceModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Whether the probe is enabled or not.",
								Computed:    true,
							},
							"probe": schema.SingleNestedAttribute{
								Description: "Probe configuration (exec, `http_get` or `tcp_socket`)",
								Computed:    true,
								CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesStartupProbeProbeDataSourceModel](ctx),
								Attributes: map[string]schema.Attribute{
									"exec": schema.SingleNestedAttribute{
										Description: "Exec probe configuration",
										Computed:    true,
										CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesStartupProbeProbeExecDataSourceModel](ctx),
										Attributes: map[string]schema.Attribute{
											"command": schema.ListAttribute{
												Description: "Command to be executed inside the running container.",
												Computed:    true,
												CustomType:  customfield.NewListType[types.String](ctx),
												ElementType: types.StringType,
											},
										},
									},
									"failure_threshold": schema.Int64Attribute{
										Description: "The number of consecutive probe failures that mark the container as unhealthy.",
										Computed:    true,
									},
									"http_get": schema.SingleNestedAttribute{
										Description: "HTTP GET probe configuration",
										Computed:    true,
										CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesStartupProbeProbeHTTPGetDataSourceModel](ctx),
										Attributes: map[string]schema.Attribute{
											"headers": schema.MapAttribute{
												Description: "HTTP headers to be sent with the request.",
												Computed:    true,
												CustomType:  customfield.NewMapType[types.String](ctx),
												ElementType: types.StringType,
											},
											"host": schema.StringAttribute{
												Description: "Host name to send HTTP request to.",
												Computed:    true,
											},
											"path": schema.StringAttribute{
												Description: "The endpoint to send the HTTP request to.",
												Computed:    true,
											},
											"port": schema.Int64Attribute{
												Description: "Port number the probe should connect to.",
												Computed:    true,
											},
											"schema": schema.StringAttribute{
												Description: "Schema to use for the HTTP request.",
												Computed:    true,
											},
										},
									},
									"initial_delay_seconds": schema.Int64Attribute{
										Description: "The initial delay before starting the first probe.",
										Computed:    true,
									},
									"period_seconds": schema.Int64Attribute{
										Description: "How often (in seconds) to perform the probe.",
										Computed:    true,
									},
									"success_threshold": schema.Int64Attribute{
										Description: "The number of consecutive successful probes that mark the container as healthy.",
										Computed:    true,
									},
									"tcp_socket": schema.SingleNestedAttribute{
										Description: "TCP socket probe configuration",
										Computed:    true,
										CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesStartupProbeProbeTcpSocketDataSourceModel](ctx),
										Attributes: map[string]schema.Attribute{
											"port": schema.Int64Attribute{
												Description: "Port number to check if it's open.",
												Computed:    true,
											},
										},
									},
									"timeout_seconds": schema.Int64Attribute{
										Description: "The timeout for each probe.",
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

func (d *CloudInferenceDeploymentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudInferenceDeploymentDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
