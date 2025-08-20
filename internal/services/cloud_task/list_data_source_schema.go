// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_task

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudTasksDataSource)(nil)

func ListDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"from_timestamp": schema.StringAttribute{
				Description: "ISO formatted datetime string. Filter the tasks by creation date greater than or equal to `from_timestamp`",
				Optional:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"is_acknowledged": schema.BoolAttribute{
				Description: "Filter the tasks by their acknowledgement status",
				Optional:    true,
			},
			"task_type": schema.StringAttribute{
				Description: "Filter the tasks by their type one of ['`activate_ddos_profile`', '`attach_bm_to_reserved_fixed_ip`', '`attach_vm_interface`', '`attach_vm_to_reserved_fixed_ip`', '`attach_volume`', '`create_ai_cluster_gpu`', '`create_bm`', '`create_caas_container`', '`create_dbaas_postgres_cluster`', '`create_ddos_profile`', '`create_faas_function`', '`create_faas_namespace`', '`create_fip`', '`create_gpu_virtual_cluster`', '`create_image`', '`create_inference_application`', '`create_inference_instance`', '`create_k8s_cluster_pool_v2`', '`create_k8s_cluster_v2`', '`create_l7policy`', '`create_l7rule`', '`create_lblistener`', '`create_lbmember`', '`create_lbpool`', '`create_lbpool_health_monitor`', '`create_loadbalancer`', '`create_network`', '`create_reserved_fixed_ip`', '`create_router`', '`create_secret`', '`create_servergroup`', '`create_sfs`', '`create_snapshot`', '`create_subnet`', '`create_vm`', '`create_volume`', '`deactivate_ddos_profile`', '`delete_ai_cluster_gpu`', '`delete_caas_container`', '`delete_dbaas_postgres_cluster`', '`delete_ddos_profile`', '`delete_faas_function`', '`delete_faas_namespace`', '`delete_fip`', '`delete_gpu_virtual_cluster`', '`delete_gpu_virtual_server`', '`delete_image`', '`delete_inference_application`', '`delete_inference_instance`', '`delete_k8s_cluster_pool_v2`', '`delete_k8s_cluster_v2`', '`delete_l7policy`', '`delete_l7rule`', '`delete_lblistener`', '`delete_lbmember`', '`delete_lbmetadata`', '`delete_lbpool`', '`delete_loadbalancer`', '`delete_network`', '`delete_reserved_fixed_ip`', '`delete_router`', '`delete_secret`', '`delete_servergroup`', '`delete_sfs`', '`delete_snapshot`', '`delete_subnet`', '`delete_vm`', '`delete_volume`', '`detach_vm_interface`', '`detach_volume`', '`download_image`', '`downscale_ai_cluster_gpu`', '`downscale_gpu_virtual_cluster`', '`extend_sfs`', '`extend_volume`', '`failover_loadbalancer`', '`hard_reboot_gpu_baremetal_server`', '`hard_reboot_gpu_virtual_cluster`', '`hard_reboot_gpu_virtual_server`', '`hard_reboot_vm`', '`patch_caas_container`', '`patch_dbaas_postgres_cluster`', '`patch_faas_function`', '`patch_faas_namespace`', '`patch_lblistener`', '`patch_lbpool`', '`put_into_server_group`', '`put_l7policy`', '`put_l7rule`', '`rebuild_bm`', '`rebuild_gpu_baremetal_node`', '`remove_from_server_group`', '`replace_lbmetadata`', '`resize_k8s_cluster_v2`', '`resize_loadbalancer`', '`resize_vm`', '`resume_vm`', '`revert_volume`', '`soft_reboot_gpu_baremetal_server`', '`soft_reboot_gpu_virtual_cluster`', '`soft_reboot_gpu_virtual_server`', '`soft_reboot_vm`', '`start_gpu_baremetal_server`', '`start_gpu_virtual_cluster`', '`start_gpu_virtual_server`', '`start_vm`', '`stop_gpu_baremetal_server`', '`stop_gpu_virtual_cluster`', '`stop_gpu_virtual_server`', '`stop_vm`', '`suspend_vm`', '`sync_private_flavors`', '`update_ddos_profile`', '`update_inference_application`', '`update_inference_instance`', '`update_k8s_cluster_v2`', '`update_lbmetadata`', '`update_port_allowed_address_pairs`', '`update_tags_gpu_virtual_cluster`', '`upgrade_k8s_cluster_v2`', '`upscale_ai_cluster_gpu`', '`upscale_gpu_virtual_cluster`']",
				Optional:    true,
			},
			"to_timestamp": schema.StringAttribute{
				Description: "ISO formatted datetime string. Filter the tasks by creation date less than or equal to `to_timestamp`",
				Optional:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"project_id": schema.ListAttribute{
				Description: "The project ID to filter the tasks by project. Supports multiple values of kind key=value1&key=value2",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"region_id": schema.ListAttribute{
				Description: "The region ID to filter the tasks by region. Supports multiple values of kind key=value1&key=value2",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"state": schema.ListAttribute{
				Description: "Filter the tasks by state. Supports multiple values of kind key=value1&key=value2",
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOfCaseInsensitive(
							"ERROR",
							"FINISHED",
							"NEW",
							"RUNNING",
						),
					),
				},
				ElementType: types.StringType,
			},
			"limit": schema.Int64Attribute{
				Description: "Limit the number of returned tasks. Falls back to default of 10 if not specified. Limited by max limit value of 1000",
				Computed:    true,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtMost(1000),
				},
			},
			"order_by": schema.StringAttribute{
				Description: "Sorting by creation date. Oldest first, or most recent first\nAvailable values: \"asc\", \"desc\".",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("asc", "desc"),
				},
			},
			"sorting": schema.StringAttribute{
				Description:        "(DEPRECATED Use '`order_by`' instead) Sorting by creation date. Oldest first, or most recent first\nAvailable values: \"asc\", \"desc\".",
				Computed:           true,
				Optional:           true,
				DeprecationMessage: "This attribute is deprecated.",
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("asc", "desc"),
				},
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
				CustomType:  customfield.NewNestedObjectListType[CloudTasksItemsDataSourceModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The task ID",
							Computed:    true,
						},
						"created_on": schema.StringAttribute{
							Description: "Created timestamp",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "The task state\nAvailable values: \"ERROR\", \"FINISHED\", \"NEW\", \"RUNNING\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"ERROR",
									"FINISHED",
									"NEW",
									"RUNNING",
								),
							},
						},
						"task_type": schema.StringAttribute{
							Description: "The task type",
							Computed:    true,
						},
						"user_id": schema.Int64Attribute{
							Description: "The user ID that initiated the task",
							Computed:    true,
						},
						"acknowledged_at": schema.StringAttribute{
							Description: "If task was acknowledged, this field stores acknowledge timestamp",
							Computed:    true,
						},
						"acknowledged_by": schema.Int64Attribute{
							Description: "If task was acknowledged, this field stores `user_id` of the person",
							Computed:    true,
						},
						"client_id": schema.Int64Attribute{
							Description: "The client ID",
							Computed:    true,
						},
						"created_resources": schema.SingleNestedAttribute{
							Description: "If the task creates resources, this field will contain their IDs",
							Computed:    true,
							CustomType:  customfield.NewNestedObjectType[CloudTasksCreatedResourcesDataSourceModel](ctx),
							Attributes: map[string]schema.Attribute{
								"ai_clusters": schema.ListAttribute{
									Description: "IDs of created AI clusters",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"api_keys": schema.ListAttribute{
									Description: "IDs of created API keys",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"caas_containers": schema.ListAttribute{
									Description: "IDs of created CaaS containers",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"ddos_profiles": schema.ListAttribute{
									Description: "IDs of created ddos protection profiles",
									Computed:    true,
									CustomType:  customfield.NewListType[types.Int64](ctx),
									ElementType: types.Int64Type,
								},
								"faas_functions": schema.ListAttribute{
									Description: "IDs of created FaaS functions",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"faas_namespaces": schema.ListAttribute{
									Description: "IDs of created FaaS namespaces",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"file_shares": schema.ListAttribute{
									Description: "IDs of created file shares",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"floatingips": schema.ListAttribute{
									Description: "IDs of created floating IPs",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"healthmonitors": schema.ListAttribute{
									Description: "IDs of created health monitors",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"images": schema.ListAttribute{
									Description: "IDs of created images",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"inference_instances": schema.ListAttribute{
									Description: "IDs of created inference instances",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"instances": schema.ListAttribute{
									Description: "IDs of created instances",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"k8s_clusters": schema.ListAttribute{
									Description: "IDs of created Kubernetes clusters",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"k8s_pools": schema.ListAttribute{
									Description: "IDs of created Kubernetes pools",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"l7polices": schema.ListAttribute{
									Description: "IDs of created L7 policies",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"l7rules": schema.ListAttribute{
									Description: "IDs of created L7 rules",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"laas_topic": schema.ListAttribute{
									Description: "IDs of created LaaS topics",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"listeners": schema.ListAttribute{
									Description: "IDs of created load balancer listeners",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"loadbalancers": schema.ListAttribute{
									Description: "IDs of created load balancers",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"members": schema.ListAttribute{
									Description: "IDs of created pool members",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"networks": schema.ListAttribute{
									Description: "IDs of created networks",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"pools": schema.ListAttribute{
									Description: "IDs of created load balancer pools",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"ports": schema.ListAttribute{
									Description: "IDs of created ports",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"postgresql_clusters": schema.ListAttribute{
									Description: "IDs of created postgres clusters",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"projects": schema.ListAttribute{
									Description: "IDs of created projects",
									Computed:    true,
									CustomType:  customfield.NewListType[types.Int64](ctx),
									ElementType: types.Int64Type,
								},
								"registry_registries": schema.ListAttribute{
									Description: "IDs of created registry registries",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"registry_users": schema.ListAttribute{
									Description: "IDs of created registry users",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"routers": schema.ListAttribute{
									Description: "IDs of created routers",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"secrets": schema.ListAttribute{
									Description: "IDs of created secrets",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"servergroups": schema.ListAttribute{
									Description: "IDs of created server groups",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"snapshots": schema.ListAttribute{
									Description: "IDs of created volume snapshots",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"subnets": schema.ListAttribute{
									Description: "IDs of created subnets",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
								"volumes": schema.ListAttribute{
									Description: "IDs of created volumes",
									Computed:    true,
									CustomType:  customfield.NewListType[types.String](ctx),
									ElementType: types.StringType,
								},
							},
						},
						"data": schema.StringAttribute{
							Description: "Task parameters",
							Computed:    true,
							CustomType:  jsontypes.NormalizedType{},
						},
						"detailed_state": schema.StringAttribute{
							Description: "Task detailed state that is more specific to task type\nAvailable values: \"CLUSTER_CLEAN_UP\", \"CLUSTER_RESIZE\", \"CLUSTER_RESUME\", \"CLUSTER_SUSPEND\", \"ERROR\", \"FINISHED\", \"IPU_SERVERS\", \"NETWORK\", \"POPLAR_SERVERS\", \"POST_DEPLOY_SETUP\", \"VIPU_CONTROLLER\".",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"CLUSTER_CLEAN_UP",
									"CLUSTER_RESIZE",
									"CLUSTER_RESUME",
									"CLUSTER_SUSPEND",
									"ERROR",
									"FINISHED",
									"IPU_SERVERS",
									"NETWORK",
									"POPLAR_SERVERS",
									"POST_DEPLOY_SETUP",
									"VIPU_CONTROLLER",
								),
							},
						},
						"error": schema.StringAttribute{
							Description: "The error value",
							Computed:    true,
						},
						"finished_on": schema.StringAttribute{
							Description: "Finished timestamp",
							Computed:    true,
						},
						"job_id": schema.StringAttribute{
							Description: "Job ID",
							Computed:    true,
						},
						"lifecycle_policy_id": schema.Int64Attribute{
							Description: "Lifecycle policy ID",
							Computed:    true,
						},
						"project_id": schema.Int64Attribute{
							Description: "The project ID",
							Computed:    true,
						},
						"region_id": schema.Int64Attribute{
							Description: "The region ID",
							Computed:    true,
						},
						"request_id": schema.StringAttribute{
							Description: "The request ID",
							Computed:    true,
						},
						"schedule_id": schema.StringAttribute{
							Description: "Schedule ID",
							Computed:    true,
						},
						"updated_on": schema.StringAttribute{
							Description: "Last updated timestamp",
							Computed:    true,
						},
						"user_client_id": schema.Int64Attribute{
							Description: "Client, specified in user's JWT",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CloudTasksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *CloudTasksDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
