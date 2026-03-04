// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_deployment

import (
	"context"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*CloudInferenceDeploymentResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Inference deployments run containerized ML models with configurable scaling, health probes, and GPU flavors.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "Inference instance name.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:   "Inference instance name.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.Int64Attribute{
				Description:   "Project ID",
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"flavor_name": schema.StringAttribute{
				Description: "Flavor name for the inference instance.",
				Required:    true,
			},
			"image": schema.StringAttribute{
				Description: "Docker image for the inference instance. This field should contain the image name and tag in the format 'name:tag', e.g., 'nginx:latest'. It defaults to Docker Hub as the image registry, but any accessible Docker image URL can be specified.",
				Required:    true,
			},
			"listening_port": schema.Int64Attribute{
				Description: "Listening port for the inference instance.",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"containers": schema.ListNestedAttribute{
				Description: "List of containers for the inference instance.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region_id": schema.Int64Attribute{
							Description: "Region id for the container",
							Required:    true,
						},
						"scale": schema.SingleNestedAttribute{
							Description: "Scale for the container",
							Required:    true,
							Attributes: map[string]schema.Attribute{
								"max": schema.Int64Attribute{
									Description: "Maximum scale for the container",
									Required:    true,
									Validators: []validator.Int64{
										int64validator.AtMost(300),
									},
								},
								"min": schema.Int64Attribute{
									Description: "Minimum scale for the container",
									Required:    true,
									Validators: []validator.Int64{
										int64validator.AtLeast(0),
									},
								},
								"cooldown_period": schema.Int64Attribute{
									Description: "Cooldown period between scaling actions in seconds",
									Computed:    true,
									Optional:    true,
									Validators: []validator.Int64{
										int64validator.Between(1, 3600),
									},
									Default: int64default.StaticInt64(60),
								},
								"polling_interval": schema.Int64Attribute{
									Description: "Polling interval for scaling triggers in seconds",
									Computed:    true,
									Optional:    true,
									Validators: []validator.Int64{
										int64validator.Between(5, 3600),
									},
									Default: int64default.StaticInt64(30),
								},
								"triggers": schema.SingleNestedAttribute{
									Description: "Triggers for scaling actions",
									Computed:    true,
									Optional:    true,
									CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentContainersScaleTriggersModel](ctx),
									Attributes: map[string]schema.Attribute{
										"cpu": schema.SingleNestedAttribute{
											Description: "CPU trigger configuration",
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												"threshold": schema.Int64Attribute{
													Description: "Threshold value for the trigger in percentage",
													Required:    true,
													Validators: []validator.Int64{
														int64validator.Between(1, 100),
													},
												},
											},
										},
										"gpu_memory": schema.SingleNestedAttribute{
											Description: "GPU memory trigger configuration. Calculated by `DCGM_FI_DEV_MEM_COPY_UTIL` metric",
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												"threshold": schema.Int64Attribute{
													Description: "Threshold value for the trigger in percentage",
													Required:    true,
													Validators: []validator.Int64{
														int64validator.Between(1, 100),
													},
												},
											},
										},
										"gpu_utilization": schema.SingleNestedAttribute{
											Description: "GPU utilization trigger configuration. Calculated by `DCGM_FI_DEV_GPU_UTIL` metric",
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												"threshold": schema.Int64Attribute{
													Description: "Threshold value for the trigger in percentage",
													Required:    true,
													Validators: []validator.Int64{
														int64validator.Between(1, 100),
													},
												},
											},
										},
										"http": schema.SingleNestedAttribute{
											Description: "HTTP trigger configuration",
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												"rate": schema.Int64Attribute{
													Description: "Request count per 'window' seconds for the http trigger",
													Required:    true,
													Validators: []validator.Int64{
														int64validator.Between(1, 1000),
													},
												},
												"window": schema.Int64Attribute{
													Description: "Time window for rate calculation in seconds",
													Required:    true,
													Validators: []validator.Int64{
														int64validator.Between(1, 3600),
													},
												},
											},
										},
										"memory": schema.SingleNestedAttribute{
											Description: "Memory trigger configuration",
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												"threshold": schema.Int64Attribute{
													Description: "Threshold value for the trigger in percentage",
													Required:    true,
													Validators: []validator.Int64{
														int64validator.Between(1, 100),
													},
												},
											},
										},
										"sqs": schema.SingleNestedAttribute{
											Description: "SQS trigger configuration",
											Optional:    true,
											Attributes: map[string]schema.Attribute{
												"activation_queue_length": schema.Int64Attribute{
													Description: "Number of messages for activation",
													Required:    true,
													Validators: []validator.Int64{
														int64validator.AtLeast(1),
													},
												},
												"aws_region": schema.StringAttribute{
													Description: "AWS region",
													Required:    true,
												},
												"queue_length": schema.Int64Attribute{
													Description: "Number of messages for one replica",
													Required:    true,
													Validators: []validator.Int64{
														int64validator.AtLeast(1),
													},
												},
												"queue_url": schema.StringAttribute{
													Description: "SQS queue URL",
													Required:    true,
												},
												"secret_name": schema.StringAttribute{
													Description: "Auth secret name",
													Required:    true,
												},
												"aws_endpoint": schema.StringAttribute{
													Description: "Custom AWS endpoint",
													Optional:    true,
												},
												"scale_on_delayed": schema.BoolAttribute{
													Description: "Scale on delayed messages",
													Computed:    true,
													Optional:    true,
													Default:     booldefault.StaticBool(false),
												},
												"scale_on_flight": schema.BoolAttribute{
													Description: "Scale on in-flight messages",
													Computed:    true,
													Optional:    true,
													Default:     booldefault.StaticBool(false),
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
			"api_keys": schema.ListAttribute{
				Description: "List of API keys for the inference instance. Multiple keys can be attached to one deployment.If `auth_enabled` and `api_keys` are both specified, a ValidationError will be raised.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"command": schema.ListAttribute{
				Description: "Command to be executed when running a container from an image.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"auth_enabled": schema.BoolAttribute{
				Description:        "Set to `true` to enable API key authentication for the inference instance. `\"Authorization\": \"Bearer *****\"` or `\"X-Api-Key\": \"*****\"` header is required for the requests to the instance if enabled. This field is deprecated and will be removed in the future. Use `api_keys` field instead.If `auth_enabled` and `api_keys` are both specified, a ValidationError will be raised.",
				Computed:           true,
				Optional:           true,
				DeprecationMessage: "This attribute is deprecated.",
				Default:            booldefault.StaticBool(false),
			},
			"credentials_name": schema.StringAttribute{
				Description: "Registry credentials name",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: "Inference instance description.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
			},
			"timeout": schema.Int64Attribute{
				Description: "Specifies the duration in seconds without any requests after which the containers will be downscaled to their minimum scale value as defined by `scale.min`. If set, this helps in optimizing resource usage by reducing the number of container instances during periods of inactivity. The default value when the parameter is not set is 120.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
				Default: int64default.StaticInt64(120),
			},
			"envs": schema.MapAttribute{
				Description: "Environment variables for the inference instance.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewMapType[types.String](ctx),
				ElementType: types.StringType,
			},
			"ingress_opts": schema.SingleNestedAttribute{
				Description: "Ingress options for the inference instance",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentIngressOptsModel](ctx),
				Attributes: map[string]schema.Attribute{
					"disable_response_buffering": schema.BoolAttribute{
						Description: "Disable response buffering if true. A client usually has a much slower connection and can not consume the response data as fast as it is produced by an upstream application. Ingress tries to buffer the whole response in order to release the upstream application as soon as possible.By default, the response buffering is enabled.",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
			},
			"logging": schema.SingleNestedAttribute{
				Description: "Logging configuration for the inference instance",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentLoggingModel](ctx),
				Attributes: map[string]schema.Attribute{
					"destination_region_id": schema.Int64Attribute{
						Description: "ID of the region in which the logs will be stored",
						Optional:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "Enable or disable log streaming",
						Computed:    true,
						Optional:    true,
						Default:     booldefault.StaticBool(false),
					},
					"retention_policy": schema.SingleNestedAttribute{
						Description: "Logs retention policy",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"period": schema.Int64Attribute{
								Description: "Duration of days for which logs must be kept.",
								Required:    true,
								Validators: []validator.Int64{
									int64validator.AtMost(1825),
								},
							},
						},
					},
					"topic_name": schema.StringAttribute{
						Description: "The topic name to stream logs to",
						Optional:    true,
					},
				},
			},
			"probes": schema.SingleNestedAttribute{
				Description: "Probes configured for all containers of the inference instance. If probes are not provided, and the `image_name` is from a the Model Catalog registry, the default probes will be used.",
				Computed:    true,
				Optional:    true,
				CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesModel](ctx),
				Attributes: map[string]schema.Attribute{
					"liveness_probe": schema.SingleNestedAttribute{
						Description: "Liveness probe configuration",
						Computed:    true,
						Optional:    true,
						CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesLivenessProbeModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Whether the probe is enabled or not.",
								Required:    true,
							},
							"probe": schema.SingleNestedAttribute{
								Description: "Probe configuration (exec, `http_get` or `tcp_socket`)",
								Computed:    true,
								Optional:    true,
								CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesLivenessProbeProbeModel](ctx),
								Attributes: map[string]schema.Attribute{
									"exec": schema.SingleNestedAttribute{
										Description: "Exec probe configuration",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"command": schema.ListAttribute{
												Description: "Command to be executed inside the running container.",
												Required:    true,
												ElementType: types.StringType,
											},
										},
									},
									"failure_threshold": schema.Int64Attribute{
										Description: "The number of consecutive probe failures that mark the container as unhealthy.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
										Default: int64default.StaticInt64(3),
									},
									"http_get": schema.SingleNestedAttribute{
										Description: "HTTP GET probe configuration",
										Computed:    true,
										Optional:    true,
										CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesLivenessProbeProbeHTTPGetModel](ctx),
										Attributes: map[string]schema.Attribute{
											"port": schema.Int64Attribute{
												Description: "Port number the probe should connect to.",
												Required:    true,
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
											"headers": schema.MapAttribute{
												Description: "HTTP headers to be sent with the request.",
												Optional:    true,
												ElementType: types.StringType,
											},
											"host": schema.StringAttribute{
												Description: "Host name to send HTTP request to.",
												Optional:    true,
											},
											"path": schema.StringAttribute{
												Description: "The endpoint to send the HTTP request to.",
												Computed:    true,
												Optional:    true,
												Default:     stringdefault.StaticString("/"),
											},
											"schema": schema.StringAttribute{
												Description: "Schema to use for the HTTP request.",
												Computed:    true,
												Optional:    true,
												Default:     stringdefault.StaticString("HTTP"),
											},
										},
									},
									"initial_delay_seconds": schema.Int64Attribute{
										Description: "The initial delay before starting the first probe.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(0),
										},
										Default: int64default.StaticInt64(0),
									},
									"period_seconds": schema.Int64Attribute{
										Description: "How often (in seconds) to perform the probe.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
										Default: int64default.StaticInt64(10),
									},
									"success_threshold": schema.Int64Attribute{
										Description: "The number of consecutive successful probes that mark the container as healthy.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
										Default: int64default.StaticInt64(1),
									},
									"tcp_socket": schema.SingleNestedAttribute{
										Description: "TCP socket probe configuration",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"port": schema.Int64Attribute{
												Description: "Port number to check if it's open.",
												Required:    true,
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
										},
									},
									"timeout_seconds": schema.Int64Attribute{
										Description: "The timeout for each probe.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
										Default: int64default.StaticInt64(3),
									},
								},
							},
						},
					},
					"readiness_probe": schema.SingleNestedAttribute{
						Description: "Readiness probe configuration",
						Computed:    true,
						Optional:    true,
						CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesReadinessProbeModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Whether the probe is enabled or not.",
								Required:    true,
							},
							"probe": schema.SingleNestedAttribute{
								Description: "Probe configuration (exec, `http_get` or `tcp_socket`)",
								Computed:    true,
								Optional:    true,
								CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesReadinessProbeProbeModel](ctx),
								Attributes: map[string]schema.Attribute{
									"exec": schema.SingleNestedAttribute{
										Description: "Exec probe configuration",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"command": schema.ListAttribute{
												Description: "Command to be executed inside the running container.",
												Required:    true,
												ElementType: types.StringType,
											},
										},
									},
									"failure_threshold": schema.Int64Attribute{
										Description: "The number of consecutive probe failures that mark the container as unhealthy.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
										Default: int64default.StaticInt64(3),
									},
									"http_get": schema.SingleNestedAttribute{
										Description: "HTTP GET probe configuration",
										Computed:    true,
										Optional:    true,
										CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesReadinessProbeProbeHTTPGetModel](ctx),
										Attributes: map[string]schema.Attribute{
											"port": schema.Int64Attribute{
												Description: "Port number the probe should connect to.",
												Required:    true,
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
											"headers": schema.MapAttribute{
												Description: "HTTP headers to be sent with the request.",
												Optional:    true,
												ElementType: types.StringType,
											},
											"host": schema.StringAttribute{
												Description: "Host name to send HTTP request to.",
												Optional:    true,
											},
											"path": schema.StringAttribute{
												Description: "The endpoint to send the HTTP request to.",
												Computed:    true,
												Optional:    true,
												Default:     stringdefault.StaticString("/"),
											},
											"schema": schema.StringAttribute{
												Description: "Schema to use for the HTTP request.",
												Computed:    true,
												Optional:    true,
												Default:     stringdefault.StaticString("HTTP"),
											},
										},
									},
									"initial_delay_seconds": schema.Int64Attribute{
										Description: "The initial delay before starting the first probe.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(0),
										},
										Default: int64default.StaticInt64(0),
									},
									"period_seconds": schema.Int64Attribute{
										Description: "How often (in seconds) to perform the probe.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
										Default: int64default.StaticInt64(10),
									},
									"success_threshold": schema.Int64Attribute{
										Description: "The number of consecutive successful probes that mark the container as healthy.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
										Default: int64default.StaticInt64(1),
									},
									"tcp_socket": schema.SingleNestedAttribute{
										Description: "TCP socket probe configuration",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"port": schema.Int64Attribute{
												Description: "Port number to check if it's open.",
												Required:    true,
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
										},
									},
									"timeout_seconds": schema.Int64Attribute{
										Description: "The timeout for each probe.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
										Default: int64default.StaticInt64(3),
									},
								},
							},
						},
					},
					"startup_probe": schema.SingleNestedAttribute{
						Description: "Startup probe configuration",
						Computed:    true,
						Optional:    true,
						CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesStartupProbeModel](ctx),
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Whether the probe is enabled or not.",
								Required:    true,
							},
							"probe": schema.SingleNestedAttribute{
								Description: "Probe configuration (exec, `http_get` or `tcp_socket`)",
								Computed:    true,
								Optional:    true,
								CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesStartupProbeProbeModel](ctx),
								Attributes: map[string]schema.Attribute{
									"exec": schema.SingleNestedAttribute{
										Description: "Exec probe configuration",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"command": schema.ListAttribute{
												Description: "Command to be executed inside the running container.",
												Required:    true,
												ElementType: types.StringType,
											},
										},
									},
									"failure_threshold": schema.Int64Attribute{
										Description: "The number of consecutive probe failures that mark the container as unhealthy.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
										Default: int64default.StaticInt64(3),
									},
									"http_get": schema.SingleNestedAttribute{
										Description: "HTTP GET probe configuration",
										Computed:    true,
										Optional:    true,
										CustomType:  customfield.NewNestedObjectType[CloudInferenceDeploymentProbesStartupProbeProbeHTTPGetModel](ctx),
										Attributes: map[string]schema.Attribute{
											"port": schema.Int64Attribute{
												Description: "Port number the probe should connect to.",
												Required:    true,
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
											"headers": schema.MapAttribute{
												Description: "HTTP headers to be sent with the request.",
												Optional:    true,
												ElementType: types.StringType,
											},
											"host": schema.StringAttribute{
												Description: "Host name to send HTTP request to.",
												Optional:    true,
											},
											"path": schema.StringAttribute{
												Description: "The endpoint to send the HTTP request to.",
												Computed:    true,
												Optional:    true,
												Default:     stringdefault.StaticString("/"),
											},
											"schema": schema.StringAttribute{
												Description: "Schema to use for the HTTP request.",
												Computed:    true,
												Optional:    true,
												Default:     stringdefault.StaticString("HTTP"),
											},
										},
									},
									"initial_delay_seconds": schema.Int64Attribute{
										Description: "The initial delay before starting the first probe.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(0),
										},
										Default: int64default.StaticInt64(0),
									},
									"period_seconds": schema.Int64Attribute{
										Description: "How often (in seconds) to perform the probe.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
										Default: int64default.StaticInt64(10),
									},
									"success_threshold": schema.Int64Attribute{
										Description: "The number of consecutive successful probes that mark the container as healthy.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
										Default: int64default.StaticInt64(1),
									},
									"tcp_socket": schema.SingleNestedAttribute{
										Description: "TCP socket probe configuration",
										Optional:    true,
										Attributes: map[string]schema.Attribute{
											"port": schema.Int64Attribute{
												Description: "Port number to check if it's open.",
												Required:    true,
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
										},
									},
									"timeout_seconds": schema.Int64Attribute{
										Description: "The timeout for each probe.",
										Computed:    true,
										Optional:    true,
										Validators: []validator.Int64{
											int64validator.AtLeast(1),
										},
										Default: int64default.StaticInt64(3),
									},
								},
							},
						},
					},
				},
			},
			"address": schema.StringAttribute{
				Description: "Address of the inference instance",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Inference instance creation date in ISO 8601 format.",
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
			"tasks": schema.ListAttribute{
				Description: "List of task IDs representing asynchronous operations. Use these IDs to monitor operation progress:\n- `GET /v1/tasks/{task_id}` - Check individual task status and details\nPoll task status until completion (`FINISHED`/`ERROR`) before proceeding with dependent operations.",
				Computed:    true,
				CustomType:  customfield.NewListType[types.String](ctx),
				ElementType: types.StringType,
			},
			"object_references": schema.ListNestedAttribute{
				Description: "Indicates to which parent object this inference belongs to.",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectListType[CloudInferenceDeploymentObjectReferencesModel](ctx),
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
		},
	}
}

func (r *CloudInferenceDeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *CloudInferenceDeploymentResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
