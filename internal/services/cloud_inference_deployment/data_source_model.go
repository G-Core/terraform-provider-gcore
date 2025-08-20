// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_deployment

import (
	"context"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/customfield"
)

type CloudInferenceDeploymentDataSourceModel struct {
	DeploymentName   types.String                                                                          `tfsdk:"deployment_name" path:"deployment_name,required"`
	ProjectID        types.Int64                                                                           `tfsdk:"project_id" path:"project_id,required"`
	Address          types.String                                                                          `tfsdk:"address" json:"address,computed"`
	AuthEnabled      types.Bool                                                                            `tfsdk:"auth_enabled" json:"auth_enabled,computed"`
	Command          types.String                                                                          `tfsdk:"command" json:"command,computed"`
	CreatedAt        types.String                                                                          `tfsdk:"created_at" json:"created_at,computed"`
	CredentialsName  types.String                                                                          `tfsdk:"credentials_name" json:"credentials_name,computed"`
	Description      types.String                                                                          `tfsdk:"description" json:"description,computed"`
	FlavorName       types.String                                                                          `tfsdk:"flavor_name" json:"flavor_name,computed"`
	Image            types.String                                                                          `tfsdk:"image" json:"image,computed"`
	ListeningPort    types.Int64                                                                           `tfsdk:"listening_port" json:"listening_port,computed"`
	Name             types.String                                                                          `tfsdk:"name" json:"name,computed"`
	Status           types.String                                                                          `tfsdk:"status" json:"status,computed"`
	Timeout          types.Int64                                                                           `tfsdk:"timeout" json:"timeout,computed"`
	APIKeys          customfield.List[types.String]                                                        `tfsdk:"api_keys" json:"api_keys,computed"`
	Envs             customfield.Map[types.String]                                                         `tfsdk:"envs" json:"envs,computed"`
	Containers       customfield.NestedObjectList[CloudInferenceDeploymentContainersDataSourceModel]       `tfsdk:"containers" json:"containers,computed"`
	IngressOpts      customfield.NestedObject[CloudInferenceDeploymentIngressOptsDataSourceModel]          `tfsdk:"ingress_opts" json:"ingress_opts,computed"`
	Logging          customfield.NestedObject[CloudInferenceDeploymentLoggingDataSourceModel]              `tfsdk:"logging" json:"logging,computed"`
	ObjectReferences customfield.NestedObjectList[CloudInferenceDeploymentObjectReferencesDataSourceModel] `tfsdk:"object_references" json:"object_references,computed"`
	Probes           customfield.NestedObject[CloudInferenceDeploymentProbesDataSourceModel]               `tfsdk:"probes" json:"probes,computed"`
}

func (m *CloudInferenceDeploymentDataSourceModel) toReadParams(_ context.Context) (params cloud.InferenceDeploymentGetParams, diags diag.Diagnostics) {
	params = cloud.InferenceDeploymentGetParams{}

	if !m.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(m.ProjectID.ValueInt64())
	}

	return
}

type CloudInferenceDeploymentContainersDataSourceModel struct {
	Address      types.String                                                                            `tfsdk:"address" json:"address,computed"`
	DeployStatus customfield.NestedObject[CloudInferenceDeploymentContainersDeployStatusDataSourceModel] `tfsdk:"deploy_status" json:"deploy_status,computed"`
	ErrorMessage types.String                                                                            `tfsdk:"error_message" json:"error_message,computed"`
	RegionID     types.Int64                                                                             `tfsdk:"region_id" json:"region_id,computed"`
	Scale        customfield.NestedObject[CloudInferenceDeploymentContainersScaleDataSourceModel]        `tfsdk:"scale" json:"scale,computed"`
}

type CloudInferenceDeploymentContainersDeployStatusDataSourceModel struct {
	Ready types.Int64 `tfsdk:"ready" json:"ready,computed"`
	Total types.Int64 `tfsdk:"total" json:"total,computed"`
}

type CloudInferenceDeploymentContainersScaleDataSourceModel struct {
	CooldownPeriod  types.Int64                                                                              `tfsdk:"cooldown_period" json:"cooldown_period,computed"`
	Max             types.Int64                                                                              `tfsdk:"max" json:"max,computed"`
	Min             types.Int64                                                                              `tfsdk:"min" json:"min,computed"`
	PollingInterval types.Int64                                                                              `tfsdk:"polling_interval" json:"polling_interval,computed"`
	Triggers        customfield.NestedObject[CloudInferenceDeploymentContainersScaleTriggersDataSourceModel] `tfsdk:"triggers" json:"triggers,computed"`
}

type CloudInferenceDeploymentContainersScaleTriggersDataSourceModel struct {
	CPU            customfield.NestedObject[CloudInferenceDeploymentContainersScaleTriggersCPUDataSourceModel]            `tfsdk:"cpu" json:"cpu,computed"`
	GPUMemory      customfield.NestedObject[CloudInferenceDeploymentContainersScaleTriggersGPUMemoryDataSourceModel]      `tfsdk:"gpu_memory" json:"gpu_memory,computed"`
	GPUUtilization customfield.NestedObject[CloudInferenceDeploymentContainersScaleTriggersGPUUtilizationDataSourceModel] `tfsdk:"gpu_utilization" json:"gpu_utilization,computed"`
	HTTP           customfield.NestedObject[CloudInferenceDeploymentContainersScaleTriggersHTTPDataSourceModel]           `tfsdk:"http" json:"http,computed"`
	Memory         customfield.NestedObject[CloudInferenceDeploymentContainersScaleTriggersMemoryDataSourceModel]         `tfsdk:"memory" json:"memory,computed"`
	Sqs            customfield.NestedObject[CloudInferenceDeploymentContainersScaleTriggersSqsDataSourceModel]            `tfsdk:"sqs" json:"sqs,computed"`
}

type CloudInferenceDeploymentContainersScaleTriggersCPUDataSourceModel struct {
	Threshold types.Int64 `tfsdk:"threshold" json:"threshold,computed"`
}

type CloudInferenceDeploymentContainersScaleTriggersGPUMemoryDataSourceModel struct {
	Threshold types.Int64 `tfsdk:"threshold" json:"threshold,computed"`
}

type CloudInferenceDeploymentContainersScaleTriggersGPUUtilizationDataSourceModel struct {
	Threshold types.Int64 `tfsdk:"threshold" json:"threshold,computed"`
}

type CloudInferenceDeploymentContainersScaleTriggersHTTPDataSourceModel struct {
	Rate   types.Int64 `tfsdk:"rate" json:"rate,computed"`
	Window types.Int64 `tfsdk:"window" json:"window,computed"`
}

type CloudInferenceDeploymentContainersScaleTriggersMemoryDataSourceModel struct {
	Threshold types.Int64 `tfsdk:"threshold" json:"threshold,computed"`
}

type CloudInferenceDeploymentContainersScaleTriggersSqsDataSourceModel struct {
	ActivationQueueLength types.Int64  `tfsdk:"activation_queue_length" json:"activation_queue_length,computed"`
	AwsEndpoint           types.String `tfsdk:"aws_endpoint" json:"aws_endpoint,computed"`
	AwsRegion             types.String `tfsdk:"aws_region" json:"aws_region,computed"`
	QueueLength           types.Int64  `tfsdk:"queue_length" json:"queue_length,computed"`
	QueueURL              types.String `tfsdk:"queue_url" json:"queue_url,computed"`
	ScaleOnDelayed        types.Bool   `tfsdk:"scale_on_delayed" json:"scale_on_delayed,computed"`
	ScaleOnFlight         types.Bool   `tfsdk:"scale_on_flight" json:"scale_on_flight,computed"`
	SecretName            types.String `tfsdk:"secret_name" json:"secret_name,computed"`
}

type CloudInferenceDeploymentIngressOptsDataSourceModel struct {
	DisableResponseBuffering types.Bool `tfsdk:"disable_response_buffering" json:"disable_response_buffering,computed"`
}

type CloudInferenceDeploymentLoggingDataSourceModel struct {
	DestinationRegionID types.Int64                                                                             `tfsdk:"destination_region_id" json:"destination_region_id,computed"`
	Enabled             types.Bool                                                                              `tfsdk:"enabled" json:"enabled,computed"`
	TopicName           types.String                                                                            `tfsdk:"topic_name" json:"topic_name,computed"`
	RetentionPolicy     customfield.NestedObject[CloudInferenceDeploymentLoggingRetentionPolicyDataSourceModel] `tfsdk:"retention_policy" json:"retention_policy,computed"`
}

type CloudInferenceDeploymentLoggingRetentionPolicyDataSourceModel struct {
	Period types.Int64 `tfsdk:"period" json:"period,computed"`
}

type CloudInferenceDeploymentObjectReferencesDataSourceModel struct {
	Kind types.String `tfsdk:"kind" json:"kind,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}

type CloudInferenceDeploymentProbesDataSourceModel struct {
	LivenessProbe  customfield.NestedObject[CloudInferenceDeploymentProbesLivenessProbeDataSourceModel]  `tfsdk:"liveness_probe" json:"liveness_probe,computed"`
	ReadinessProbe customfield.NestedObject[CloudInferenceDeploymentProbesReadinessProbeDataSourceModel] `tfsdk:"readiness_probe" json:"readiness_probe,computed"`
	StartupProbe   customfield.NestedObject[CloudInferenceDeploymentProbesStartupProbeDataSourceModel]   `tfsdk:"startup_probe" json:"startup_probe,computed"`
}

type CloudInferenceDeploymentProbesLivenessProbeDataSourceModel struct {
	Enabled types.Bool                                                                                `tfsdk:"enabled" json:"enabled,computed"`
	Probe   customfield.NestedObject[CloudInferenceDeploymentProbesLivenessProbeProbeDataSourceModel] `tfsdk:"probe" json:"probe,computed"`
}

type CloudInferenceDeploymentProbesLivenessProbeProbeDataSourceModel struct {
	Exec                customfield.NestedObject[CloudInferenceDeploymentProbesLivenessProbeProbeExecDataSourceModel]      `tfsdk:"exec" json:"exec,computed"`
	FailureThreshold    types.Int64                                                                                        `tfsdk:"failure_threshold" json:"failure_threshold,computed"`
	HTTPGet             customfield.NestedObject[CloudInferenceDeploymentProbesLivenessProbeProbeHTTPGetDataSourceModel]   `tfsdk:"http_get" json:"http_get,computed"`
	InitialDelaySeconds types.Int64                                                                                        `tfsdk:"initial_delay_seconds" json:"initial_delay_seconds,computed"`
	PeriodSeconds       types.Int64                                                                                        `tfsdk:"period_seconds" json:"period_seconds,computed"`
	SuccessThreshold    types.Int64                                                                                        `tfsdk:"success_threshold" json:"success_threshold,computed"`
	TcpSocket           customfield.NestedObject[CloudInferenceDeploymentProbesLivenessProbeProbeTcpSocketDataSourceModel] `tfsdk:"tcp_socket" json:"tcp_socket,computed"`
	TimeoutSeconds      types.Int64                                                                                        `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
}

type CloudInferenceDeploymentProbesLivenessProbeProbeExecDataSourceModel struct {
	Command customfield.List[types.String] `tfsdk:"command" json:"command,computed"`
}

type CloudInferenceDeploymentProbesLivenessProbeProbeHTTPGetDataSourceModel struct {
	Headers customfield.Map[types.String] `tfsdk:"headers" json:"headers,computed"`
	Host    types.String                  `tfsdk:"host" json:"host,computed"`
	Path    types.String                  `tfsdk:"path" json:"path,computed"`
	Port    types.Int64                   `tfsdk:"port" json:"port,computed"`
	Schema  types.String                  `tfsdk:"schema" json:"schema,computed"`
}

type CloudInferenceDeploymentProbesLivenessProbeProbeTcpSocketDataSourceModel struct {
	Port types.Int64 `tfsdk:"port" json:"port,computed"`
}

type CloudInferenceDeploymentProbesReadinessProbeDataSourceModel struct {
	Enabled types.Bool                                                                                 `tfsdk:"enabled" json:"enabled,computed"`
	Probe   customfield.NestedObject[CloudInferenceDeploymentProbesReadinessProbeProbeDataSourceModel] `tfsdk:"probe" json:"probe,computed"`
}

type CloudInferenceDeploymentProbesReadinessProbeProbeDataSourceModel struct {
	Exec                customfield.NestedObject[CloudInferenceDeploymentProbesReadinessProbeProbeExecDataSourceModel]      `tfsdk:"exec" json:"exec,computed"`
	FailureThreshold    types.Int64                                                                                         `tfsdk:"failure_threshold" json:"failure_threshold,computed"`
	HTTPGet             customfield.NestedObject[CloudInferenceDeploymentProbesReadinessProbeProbeHTTPGetDataSourceModel]   `tfsdk:"http_get" json:"http_get,computed"`
	InitialDelaySeconds types.Int64                                                                                         `tfsdk:"initial_delay_seconds" json:"initial_delay_seconds,computed"`
	PeriodSeconds       types.Int64                                                                                         `tfsdk:"period_seconds" json:"period_seconds,computed"`
	SuccessThreshold    types.Int64                                                                                         `tfsdk:"success_threshold" json:"success_threshold,computed"`
	TcpSocket           customfield.NestedObject[CloudInferenceDeploymentProbesReadinessProbeProbeTcpSocketDataSourceModel] `tfsdk:"tcp_socket" json:"tcp_socket,computed"`
	TimeoutSeconds      types.Int64                                                                                         `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
}

type CloudInferenceDeploymentProbesReadinessProbeProbeExecDataSourceModel struct {
	Command customfield.List[types.String] `tfsdk:"command" json:"command,computed"`
}

type CloudInferenceDeploymentProbesReadinessProbeProbeHTTPGetDataSourceModel struct {
	Headers customfield.Map[types.String] `tfsdk:"headers" json:"headers,computed"`
	Host    types.String                  `tfsdk:"host" json:"host,computed"`
	Path    types.String                  `tfsdk:"path" json:"path,computed"`
	Port    types.Int64                   `tfsdk:"port" json:"port,computed"`
	Schema  types.String                  `tfsdk:"schema" json:"schema,computed"`
}

type CloudInferenceDeploymentProbesReadinessProbeProbeTcpSocketDataSourceModel struct {
	Port types.Int64 `tfsdk:"port" json:"port,computed"`
}

type CloudInferenceDeploymentProbesStartupProbeDataSourceModel struct {
	Enabled types.Bool                                                                               `tfsdk:"enabled" json:"enabled,computed"`
	Probe   customfield.NestedObject[CloudInferenceDeploymentProbesStartupProbeProbeDataSourceModel] `tfsdk:"probe" json:"probe,computed"`
}

type CloudInferenceDeploymentProbesStartupProbeProbeDataSourceModel struct {
	Exec                customfield.NestedObject[CloudInferenceDeploymentProbesStartupProbeProbeExecDataSourceModel]      `tfsdk:"exec" json:"exec,computed"`
	FailureThreshold    types.Int64                                                                                       `tfsdk:"failure_threshold" json:"failure_threshold,computed"`
	HTTPGet             customfield.NestedObject[CloudInferenceDeploymentProbesStartupProbeProbeHTTPGetDataSourceModel]   `tfsdk:"http_get" json:"http_get,computed"`
	InitialDelaySeconds types.Int64                                                                                       `tfsdk:"initial_delay_seconds" json:"initial_delay_seconds,computed"`
	PeriodSeconds       types.Int64                                                                                       `tfsdk:"period_seconds" json:"period_seconds,computed"`
	SuccessThreshold    types.Int64                                                                                       `tfsdk:"success_threshold" json:"success_threshold,computed"`
	TcpSocket           customfield.NestedObject[CloudInferenceDeploymentProbesStartupProbeProbeTcpSocketDataSourceModel] `tfsdk:"tcp_socket" json:"tcp_socket,computed"`
	TimeoutSeconds      types.Int64                                                                                       `tfsdk:"timeout_seconds" json:"timeout_seconds,computed"`
}

type CloudInferenceDeploymentProbesStartupProbeProbeExecDataSourceModel struct {
	Command customfield.List[types.String] `tfsdk:"command" json:"command,computed"`
}

type CloudInferenceDeploymentProbesStartupProbeProbeHTTPGetDataSourceModel struct {
	Headers customfield.Map[types.String] `tfsdk:"headers" json:"headers,computed"`
	Host    types.String                  `tfsdk:"host" json:"host,computed"`
	Path    types.String                  `tfsdk:"path" json:"path,computed"`
	Port    types.Int64                   `tfsdk:"port" json:"port,computed"`
	Schema  types.String                  `tfsdk:"schema" json:"schema,computed"`
}

type CloudInferenceDeploymentProbesStartupProbeProbeTcpSocketDataSourceModel struct {
	Port types.Int64 `tfsdk:"port" json:"port,computed"`
}
