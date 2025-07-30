// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_task

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

var _ datasource.DataSourceWithConfigValidators = (*CloudTaskDataSource)(nil)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"task_id": schema.StringAttribute{
				Description: "Task ID",
				Required:    true,
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
			"created_on": schema.StringAttribute{
				Description: "Created timestamp",
				Computed:    true,
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
			"id": schema.StringAttribute{
				Description: "The task ID",
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
			"updated_on": schema.StringAttribute{
				Description: "Last updated timestamp",
				Computed:    true,
			},
			"user_client_id": schema.Int64Attribute{
				Description: "Client, specified in user's JWT",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "The user ID that initiated the task",
				Computed:    true,
			},
			"created_resources": schema.SingleNestedAttribute{
				Description: "If the task creates resources, this field will contain their IDs",
				Computed:    true,
				CustomType:  customfield.NewNestedObjectType[CloudTaskCreatedResourcesDataSourceModel](ctx),
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
		},
	}
}

func (d *CloudTaskDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *CloudTaskDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
