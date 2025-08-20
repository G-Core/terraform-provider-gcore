// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_task

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudTaskDataSourceModel struct {
	TaskID            types.String                                                       `tfsdk:"task_id" path:"task_id,required"`
	AcknowledgedAt    types.String                                                       `tfsdk:"acknowledged_at" json:"acknowledged_at,computed"`
	AcknowledgedBy    types.Int64                                                        `tfsdk:"acknowledged_by" json:"acknowledged_by,computed"`
	ClientID          types.Int64                                                        `tfsdk:"client_id" json:"client_id,computed"`
	CreatedOn         types.String                                                       `tfsdk:"created_on" json:"created_on,computed"`
	DetailedState     types.String                                                       `tfsdk:"detailed_state" json:"detailed_state,computed"`
	Error             types.String                                                       `tfsdk:"error" json:"error,computed"`
	FinishedOn        types.String                                                       `tfsdk:"finished_on" json:"finished_on,computed"`
	ID                types.String                                                       `tfsdk:"id" json:"id,computed"`
	JobID             types.String                                                       `tfsdk:"job_id" json:"job_id,computed"`
	LifecyclePolicyID types.Int64                                                        `tfsdk:"lifecycle_policy_id" json:"lifecycle_policy_id,computed"`
	ProjectID         types.Int64                                                        `tfsdk:"project_id" json:"project_id,computed"`
	RegionID          types.Int64                                                        `tfsdk:"region_id" json:"region_id,computed"`
	RequestID         types.String                                                       `tfsdk:"request_id" json:"request_id,computed"`
	ScheduleID        types.String                                                       `tfsdk:"schedule_id" json:"schedule_id,computed"`
	State             types.String                                                       `tfsdk:"state" json:"state,computed"`
	TaskType          types.String                                                       `tfsdk:"task_type" json:"task_type,computed"`
	UpdatedOn         types.String                                                       `tfsdk:"updated_on" json:"updated_on,computed"`
	UserClientID      types.Int64                                                        `tfsdk:"user_client_id" json:"user_client_id,computed"`
	UserID            types.Int64                                                        `tfsdk:"user_id" json:"user_id,computed"`
	CreatedResources  customfield.NestedObject[CloudTaskCreatedResourcesDataSourceModel] `tfsdk:"created_resources" json:"created_resources,computed"`
	Data              jsontypes.Normalized                                               `tfsdk:"data" json:"data,computed"`
}

type CloudTaskCreatedResourcesDataSourceModel struct {
	AIClusters         customfield.List[types.String] `tfsdk:"ai_clusters" json:"ai_clusters,computed"`
	APIKeys            customfield.List[types.String] `tfsdk:"api_keys" json:"api_keys,computed"`
	CaasContainers     customfield.List[types.String] `tfsdk:"caas_containers" json:"caas_containers,computed"`
	DDOSProfiles       customfield.List[types.Int64]  `tfsdk:"ddos_profiles" json:"ddos_profiles,computed"`
	FaasFunctions      customfield.List[types.String] `tfsdk:"faas_functions" json:"faas_functions,computed"`
	FaasNamespaces     customfield.List[types.String] `tfsdk:"faas_namespaces" json:"faas_namespaces,computed"`
	FileShares         customfield.List[types.String] `tfsdk:"file_shares" json:"file_shares,computed"`
	Floatingips        customfield.List[types.String] `tfsdk:"floatingips" json:"floatingips,computed"`
	Healthmonitors     customfield.List[types.String] `tfsdk:"healthmonitors" json:"healthmonitors,computed"`
	Images             customfield.List[types.String] `tfsdk:"images" json:"images,computed"`
	InferenceInstances customfield.List[types.String] `tfsdk:"inference_instances" json:"inference_instances,computed"`
	Instances          customfield.List[types.String] `tfsdk:"instances" json:"instances,computed"`
	K8sClusters        customfield.List[types.String] `tfsdk:"k8s_clusters" json:"k8s_clusters,computed"`
	K8sPools           customfield.List[types.String] `tfsdk:"k8s_pools" json:"k8s_pools,computed"`
	L7polices          customfield.List[types.String] `tfsdk:"l7polices" json:"l7polices,computed"`
	L7rules            customfield.List[types.String] `tfsdk:"l7rules" json:"l7rules,computed"`
	LaasTopic          customfield.List[types.String] `tfsdk:"laas_topic" json:"laas_topic,computed"`
	Listeners          customfield.List[types.String] `tfsdk:"listeners" json:"listeners,computed"`
	Loadbalancers      customfield.List[types.String] `tfsdk:"loadbalancers" json:"loadbalancers,computed"`
	Members            customfield.List[types.String] `tfsdk:"members" json:"members,computed"`
	Networks           customfield.List[types.String] `tfsdk:"networks" json:"networks,computed"`
	Pools              customfield.List[types.String] `tfsdk:"pools" json:"pools,computed"`
	Ports              customfield.List[types.String] `tfsdk:"ports" json:"ports,computed"`
	PostgreSQLClusters customfield.List[types.String] `tfsdk:"postgresql_clusters" json:"postgresql_clusters,computed"`
	Projects           customfield.List[types.Int64]  `tfsdk:"projects" json:"projects,computed"`
	RegistryRegistries customfield.List[types.String] `tfsdk:"registry_registries" json:"registry_registries,computed"`
	RegistryUsers      customfield.List[types.String] `tfsdk:"registry_users" json:"registry_users,computed"`
	Routers            customfield.List[types.String] `tfsdk:"routers" json:"routers,computed"`
	Secrets            customfield.List[types.String] `tfsdk:"secrets" json:"secrets,computed"`
	Servergroups       customfield.List[types.String] `tfsdk:"servergroups" json:"servergroups,computed"`
	Snapshots          customfield.List[types.String] `tfsdk:"snapshots" json:"snapshots,computed"`
	Subnets            customfield.List[types.String] `tfsdk:"subnets" json:"subnets,computed"`
	Volumes            customfield.List[types.String] `tfsdk:"volumes" json:"volumes,computed"`
}
