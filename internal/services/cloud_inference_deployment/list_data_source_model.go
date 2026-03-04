// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_deployment

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudInferenceDeploymentsResultsListDataSourceEnvelope struct {
	Results customfield.NestedObjectList[CloudInferenceDeploymentsItemsDataSourceModel] `json:"results,computed"`
}

type CloudInferenceDeploymentsDataSourceModel struct {
	ProjectID types.Int64                                                                 `tfsdk:"project_id" path:"project_id,optional"`
	MaxItems  types.Int64                                                                 `tfsdk:"max_items"`
	Items     customfield.NestedObjectList[CloudInferenceDeploymentsItemsDataSourceModel] `tfsdk:"items"`
}

func (m *CloudInferenceDeploymentsDataSourceModel) toListParams(_ context.Context) (params cloud.InferenceDeploymentListParams, diags diag.Diagnostics) {
	params = cloud.InferenceDeploymentListParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

type CloudInferenceDeploymentsItemsDataSourceModel struct {
	ID               types.String                                                                           `tfsdk:"id" json:"name,computed"`
	Address          types.String                                                                           `tfsdk:"address" json:"address,computed"`
	AuthEnabled      types.Bool                                                                             `tfsdk:"auth_enabled" json:"auth_enabled,computed"`
	Command          types.String                                                                           `tfsdk:"command" json:"command,computed"`
	Containers       customfield.NestedObjectList[CloudInferenceDeploymentsContainersDataSourceModel]       `tfsdk:"containers" json:"containers,computed"`
	CreatedAt        types.String                                                                           `tfsdk:"created_at" json:"created_at,computed"`
	CredentialsName  types.String                                                                           `tfsdk:"credentials_name" json:"credentials_name,computed"`
	Description      types.String                                                                           `tfsdk:"description" json:"description,computed"`
	Envs             customfield.Map[types.String]                                                          `tfsdk:"envs" json:"envs,computed"`
	FlavorName       types.String                                                                           `tfsdk:"flavor_name" json:"flavor_name,computed"`
	Image            types.String                                                                           `tfsdk:"image" json:"image,computed"`
	IngressOpts      customfield.NestedObject[CloudInferenceDeploymentsIngressOptsDataSourceModel]          `tfsdk:"ingress_opts" json:"ingress_opts,computed"`
	ListeningPort    types.Int64                                                                            `tfsdk:"listening_port" json:"listening_port,computed"`
	Logging          customfield.NestedObject[CloudInferenceDeploymentsLoggingDataSourceModel]              `tfsdk:"logging" json:"logging,computed"`
	Name             types.String                                                                           `tfsdk:"name" json:"name,computed"`
	ObjectReferences customfield.NestedObjectList[CloudInferenceDeploymentsObjectReferencesDataSourceModel] `tfsdk:"object_references" json:"object_references,computed"`
	Probes           customfield.NestedObject[CloudInferenceDeploymentsProbesDataSourceModel]               `tfsdk:"probes" json:"probes,computed"`
	ProjectID        types.Int64                                                                            `tfsdk:"project_id" json:"project_id,computed"`
	Status           types.String                                                                           `tfsdk:"status" json:"status,computed"`
	Timeout          types.Int64                                                                            `tfsdk:"timeout" json:"timeout,computed"`
	APIKeys          customfield.List[types.String]                                                         `tfsdk:"api_keys" json:"api_keys,computed"`
}

type CloudInferenceDeploymentsContainersDataSourceModel struct {
	Address      types.String                                                                             `tfsdk:"address" json:"address,computed"`
	DeployStatus customfield.NestedObject[CloudInferenceDeploymentsContainersDeployStatusDataSourceModel] `tfsdk:"deploy_status" json:"deploy_status,computed"`
	ErrorMessage types.String                                                                             `tfsdk:"error_message" json:"error_message,computed"`
	RegionID     types.Int64                                                                              `tfsdk:"region_id" json:"region_id,computed"`
	Scale        customfield.NestedObject[CloudInferenceDeploymentsContainersScaleDataSourceModel]        `tfsdk:"scale" json:"scale,computed"`
}

type CloudInferenceDeploymentsContainersDeployStatusDataSourceModel struct {
	Ready types.Int64 `tfsdk:"ready" json:"ready,computed"`
	Total types.Int64 `tfsdk:"total" json:"total,computed"`
}

type CloudInferenceDeploymentsContainersScaleDataSourceModel struct {
	CooldownPeriod  types.Int64                                                                               `tfsdk:"cooldown_period" json:"cooldown_period,computed"`
	Max             types.Int64                                                                               `tfsdk:"max" json:"max,computed"`
	Min             types.Int64                                                                               `tfsdk:"min" json:"min,computed"`
	PollingInterval types.Int64                                                                               `tfsdk:"polling_interval" json:"polling_interval,computed"`
	Triggers        customfield.NestedObject[CloudInferenceDeploymentsContainersScaleTriggersDataSourceModel] `tfsdk:"triggers" json:"triggers,computed"`
}

type CloudInferenceDeploymentsContainersScaleTriggersDataSourceModel struct {
	CPU            customfield.NestedObject[CloudInferenceDeploymentsContainersScaleTriggersCPUDataSourceModel]            `tfsdk:"cpu" json:"cpu,computed"`
	GPUMemory      customfield.NestedObject[CloudInferenceDeploymentsContainersScaleTriggersGPUMemoryDataSourceModel]      `tfsdk:"gpu_memory" json:"gpu_memory,computed"`
	GPUUtilization customfield.NestedObject[CloudInferenceDeploymentsContainersScaleTriggersGPUUtilizationDataSourceModel] `tfsdk:"gpu_utilization" json:"gpu_utilization,computed"`
	HTTP           customfield.NestedObject[CloudInferenceDeploymentsContainersScaleTriggersHTTPDataSourceModel]           `tfsdk:"http" json:"http,computed"`
	Memory         customfield.NestedObject[CloudInferenceDeploymentsContainersScaleTriggersMemoryDataSourceModel]         `tfsdk:"memory" json:"memory,computed"`
	Sqs            customfield.NestedObject[CloudInferenceDeploymentsContainersScaleTriggersSqsDataSourceModel]            `tfsdk:"sqs" json:"sqs,computed"`
}

type CloudInferenceDeploymentsContainersScaleTriggersCPUDataSourceModel struct {
	Threshold types.Int64 `tfsdk:"threshold" json:"threshold,computed"`
}

type CloudInferenceDeploymentsContainersScaleTriggersGPUMemoryDataSourceModel struct {
	Threshold types.Int64 `tfsdk:"threshold" json:"threshold,computed"`
}

type CloudInferenceDeploymentsContainersScaleTriggersGPUUtilizationDataSourceModel struct {
	Threshold types.Int64 `tfsdk:"threshold" json:"threshold,computed"`
}

type CloudInferenceDeploymentsContainersScaleTriggersHTTPDataSourceModel struct {
	Rate   types.Int64 `tfsdk:"rate" json:"rate,computed"`
	Window types.Int64 `tfsdk:"window" json:"window,computed"`
}

type CloudInferenceDeploymentsContainersScaleTriggersMemoryDataSourceModel struct {
	Threshold types.Int64 `tfsdk:"threshold" json:"threshold,computed"`
}

type CloudInferenceDeploymentsContainersScaleTriggersSqsDataSourceModel struct {
	ActivationQueueLength types.Int64  `tfsdk:"activation_queue_length" json:"activation_queue_length,computed"`
	AwsEndpoint           types.String `tfsdk:"aws_endpoint" json:"aws_endpoint,computed"`
	AwsRegion             types.String `tfsdk:"aws_region" json:"aws_region,computed"`
	QueueLength           types.Int64  `tfsdk:"queue_length" json:"queue_length,computed"`
	QueueURL              types.String `tfsdk:"queue_url" json:"queue_url,computed"`
	ScaleOnDelayed        types.Bool   `tfsdk:"scale_on_delayed" json:"scale_on_delayed,computed"`
	ScaleOnFlight         types.Bool   `tfsdk:"scale_on_flight" json:"scale_on_flight,computed"`
	SecretName            types.String `tfsdk:"secret_name" json:"secret_name,computed"`
}

type CloudInferenceDeploymentsIngressOptsDataSourceModel struct {
	DisableResponseBuffering types.Bool `tfsdk:"disable_response_buffering" json:"disable_response_buffering,computed"`
}

type CloudInferenceDeploymentsLoggingDataSourceModel struct {
	DestinationRegionID types.Int64                                                                              `tfsdk:"destination_region_id" json:"destination_region_id,computed"`
	Enabled             types.Bool                                                                               `tfsdk:"enabled" json:"enabled,computed"`
	TopicName           types.String                                                                             `tfsdk:"topic_name" json:"topic_name,computed"`
	RetentionPolicy     customfield.NestedObject[CloudInferenceDeploymentsLoggingRetentionPolicyDataSourceModel] `tfsdk:"retention_policy" json:"retention_policy,computed"`
}

type CloudInferenceDeploymentsLoggingRetentionPolicyDataSourceModel struct {
	Period types.Int64 `tfsdk:"period" json:"period,computed"`
}

type CloudInferenceDeploymentsObjectReferencesDataSourceModel struct {
	Kind types.String `tfsdk:"kind" json:"kind,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}

type CloudInferenceDeploymentsProbesDataSourceModel struct {
	LivenessProbe  customfield.NestedObject[CloudInferenceDeploymentsProbesLivenessProbeDataSourceModel]  `tfsdk:"liveness_probe" json:"liveness_probe,computed"`
	ReadinessProbe customfield.NestedObject[CloudInferenceDeploymentsProbesReadinessProbeDataSourceModel] `tfsdk:"readiness_probe" json:"readiness_probe,computed"`
	StartupProbe   customfield.NestedObject[CloudInferenceDeploymentsProbesStartupProbeDataSourceModel]   `tfsdk:"startup_probe" json:"startup_probe,computed"`
}

type CloudInferenceDeploymentsProbesLivenessProbeDataSourceModel struct {
	Enabled types.Bool                                                                                 `tfsdk:"enabled" json:"enabled,computed"`
	Probe   customfield.NestedObject[CloudInferenceDeploymentsProbesLivenessProbeProbeDataSourceModel] `tfsdk:"probe" json:"probe,computed"`
}

type CloudInferenceDeploymentsProbesLivenessProbeProbeDataSourceModel struct {
	Exec                customfield.NestedObject[CloudInferenceDeploymentsProbesLivenessProbeProbeExecDataSourceModel]      `tfsdk:"exec" json:"exec,computed"`
	FailureThreshold    types.Int64                                                                                         `tfsdk:"failure_threshold" json:"failure_threshold,computed"`
	HTTPGet             customfield.NestedObject[CloudInferenceDeploymentsProbesLivenessProbeProbeHTTPGetDataSourceModel]   `tfsdk:"http_get" json:"http_get,computed"`
	InitialDelaySeconds types.Int64                                                                                         `tfsdk:"initial_delay_seconds" json:"initial_delay_seconds,computed"`
	PeriodSeconds       types.Int64                                                                                         `tfsdk:"period_seconds" json:"period_seconds,computed"`
	SuccessThreshold    types.Int64                                                                                         `tfsdk:"success_threshold" json:"success_threshold,computed"`
	TcpSocket           customfield.NestedObject[CloudInferenceDeploymentsProbesLivenessProbeProbeTcpSocketDataSourceModel] `tfsdk:"tcp_socket" json:"tcp_socket,computed"`
	TimeoutSeconds      types.Int64                                                                                         `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
}

type CloudInferenceDeploymentsProbesLivenessProbeProbeExecDataSourceModel struct {
	Command customfield.List[types.String] `tfsdk:"command" json:"command,computed"`
}

type CloudInferenceDeploymentsProbesLivenessProbeProbeHTTPGetDataSourceModel struct {
	Headers customfield.Map[types.String] `tfsdk:"headers" json:"headers,computed"`
	Host    types.String                  `tfsdk:"host" json:"host,computed"`
	Path    types.String                  `tfsdk:"path" json:"path,computed"`
	Port    types.Int64                   `tfsdk:"port" json:"port,computed"`
	Schema  types.String                  `tfsdk:"schema" json:"schema,computed"`
}

type CloudInferenceDeploymentsProbesLivenessProbeProbeTcpSocketDataSourceModel struct {
	Port types.Int64 `tfsdk:"port" json:"port,computed"`
}

type CloudInferenceDeploymentsProbesReadinessProbeDataSourceModel struct {
	Enabled types.Bool                                                                                  `tfsdk:"enabled" json:"enabled,computed"`
	Probe   customfield.NestedObject[CloudInferenceDeploymentsProbesReadinessProbeProbeDataSourceModel] `tfsdk:"probe" json:"probe,computed"`
}

type CloudInferenceDeploymentsProbesReadinessProbeProbeDataSourceModel struct {
	Exec                customfield.NestedObject[CloudInferenceDeploymentsProbesReadinessProbeProbeExecDataSourceModel]      `tfsdk:"exec" json:"exec,computed"`
	FailureThreshold    types.Int64                                                                                          `tfsdk:"failure_threshold" json:"failure_threshold,computed"`
	HTTPGet             customfield.NestedObject[CloudInferenceDeploymentsProbesReadinessProbeProbeHTTPGetDataSourceModel]   `tfsdk:"http_get" json:"http_get,computed"`
	InitialDelaySeconds types.Int64                                                                                          `tfsdk:"initial_delay_seconds" json:"initial_delay_seconds,computed"`
	PeriodSeconds       types.Int64                                                                                          `tfsdk:"period_seconds" json:"period_seconds,computed"`
	SuccessThreshold    types.Int64                                                                                          `tfsdk:"success_threshold" json:"success_threshold,computed"`
	TcpSocket           customfield.NestedObject[CloudInferenceDeploymentsProbesReadinessProbeProbeTcpSocketDataSourceModel] `tfsdk:"tcp_socket" json:"tcp_socket,computed"`
	TimeoutSeconds      types.Int64                                                                                          `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
}

type CloudInferenceDeploymentsProbesReadinessProbeProbeExecDataSourceModel struct {
	Command customfield.List[types.String] `tfsdk:"command" json:"command,computed"`
}

type CloudInferenceDeploymentsProbesReadinessProbeProbeHTTPGetDataSourceModel struct {
	Headers customfield.Map[types.String] `tfsdk:"headers" json:"headers,computed"`
	Host    types.String                  `tfsdk:"host" json:"host,computed"`
	Path    types.String                  `tfsdk:"path" json:"path,computed"`
	Port    types.Int64                   `tfsdk:"port" json:"port,computed"`
	Schema  types.String                  `tfsdk:"schema" json:"schema,computed"`
}

type CloudInferenceDeploymentsProbesReadinessProbeProbeTcpSocketDataSourceModel struct {
	Port types.Int64 `tfsdk:"port" json:"port,computed"`
}

type CloudInferenceDeploymentsProbesStartupProbeDataSourceModel struct {
	Enabled types.Bool                                                                                `tfsdk:"enabled" json:"enabled,computed"`
	Probe   customfield.NestedObject[CloudInferenceDeploymentsProbesStartupProbeProbeDataSourceModel] `tfsdk:"probe" json:"probe,computed"`
}

type CloudInferenceDeploymentsProbesStartupProbeProbeDataSourceModel struct {
	Exec                customfield.NestedObject[CloudInferenceDeploymentsProbesStartupProbeProbeExecDataSourceModel]      `tfsdk:"exec" json:"exec,computed"`
	FailureThreshold    types.Int64                                                                                        `tfsdk:"failure_threshold" json:"failure_threshold,computed"`
	HTTPGet             customfield.NestedObject[CloudInferenceDeploymentsProbesStartupProbeProbeHTTPGetDataSourceModel]   `tfsdk:"http_get" json:"http_get,computed"`
	InitialDelaySeconds types.Int64                                                                                        `tfsdk:"initial_delay_seconds" json:"initial_delay_seconds,computed"`
	PeriodSeconds       types.Int64                                                                                        `tfsdk:"period_seconds" json:"period_seconds,computed"`
	SuccessThreshold    types.Int64                                                                                        `tfsdk:"success_threshold" json:"success_threshold,computed"`
	TcpSocket           customfield.NestedObject[CloudInferenceDeploymentsProbesStartupProbeProbeTcpSocketDataSourceModel] `tfsdk:"tcp_socket" json:"tcp_socket,computed"`
	TimeoutSeconds      types.Int64                                                                                        `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
}

type CloudInferenceDeploymentsProbesStartupProbeProbeExecDataSourceModel struct {
	Command customfield.List[types.String] `tfsdk:"command" json:"command,computed"`
}

type CloudInferenceDeploymentsProbesStartupProbeProbeHTTPGetDataSourceModel struct {
	Headers customfield.Map[types.String] `tfsdk:"headers" json:"headers,computed"`
	Host    types.String                  `tfsdk:"host" json:"host,computed"`
	Path    types.String                  `tfsdk:"path" json:"path,computed"`
	Port    types.Int64                   `tfsdk:"port" json:"port,computed"`
	Schema  types.String                  `tfsdk:"schema" json:"schema,computed"`
}

type CloudInferenceDeploymentsProbesStartupProbeProbeTcpSocketDataSourceModel struct {
	Port types.Int64 `tfsdk:"port" json:"port,computed"`
}
