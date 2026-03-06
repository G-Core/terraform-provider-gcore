// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_deployment

import (
	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudInferenceDeploymentModel struct {
	ID               types.String                                                                `tfsdk:"id" json:"-,computed"`
	Name             types.String                                                                `tfsdk:"name" json:"name,required"`
	ProjectID        types.Int64                                                                 `tfsdk:"project_id" path:"project_id,optional"`
	FlavorName       types.String                                                                `tfsdk:"flavor_name" json:"flavor_name,required"`
	Image            types.String                                                                `tfsdk:"image" json:"image,required"`
	ListeningPort    types.Int64                                                                 `tfsdk:"listening_port" json:"listening_port,required"`
	Containers       *[]*CloudInferenceDeploymentContainersModel                                 `tfsdk:"containers" json:"containers,required"`
	APIKeys          *[]types.String                                                             `tfsdk:"api_keys" json:"api_keys,optional"`
	Command          *[]types.String                                                             `tfsdk:"command" json:"command,optional,no_refresh"`
	AuthEnabled      types.Bool                                                                  `tfsdk:"auth_enabled" json:"auth_enabled,computed_optional"`
	CredentialsName  types.String                                                                `tfsdk:"credentials_name" json:"credentials_name,computed_optional"`
	Description      types.String                                                                `tfsdk:"description" json:"description,computed_optional"`
	Timeout          types.Int64                                                                 `tfsdk:"timeout" json:"timeout,computed_optional"`
	Envs             customfield.Map[types.String]                                               `tfsdk:"envs" json:"envs,computed_optional"`
	IngressOpts      customfield.NestedObject[CloudInferenceDeploymentIngressOptsModel]          `tfsdk:"ingress_opts" json:"ingress_opts,computed_optional"`
	Logging          customfield.NestedObject[CloudInferenceDeploymentLoggingModel]              `tfsdk:"logging" json:"logging,computed_optional"`
	Probes           customfield.NestedObject[CloudInferenceDeploymentProbesModel]               `tfsdk:"probes" json:"probes,computed_optional"`
	Address          types.String                                                                `tfsdk:"address" json:"address,computed"`
	CreatedAt        types.String                                                                `tfsdk:"created_at" json:"created_at,computed"`
	Status           types.String                                                                `tfsdk:"status" json:"status,computed"`
	Tasks            customfield.List[types.String]                                              `tfsdk:"tasks" json:"tasks,computed,no_refresh"`
	ObjectReferences customfield.NestedObjectList[CloudInferenceDeploymentObjectReferencesModel] `tfsdk:"object_references" json:"object_references,computed"`
}

func (m CloudInferenceDeploymentModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m CloudInferenceDeploymentModel) MarshalJSONForUpdate(state CloudInferenceDeploymentModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}

type CloudInferenceDeploymentContainersModel struct {
	RegionID types.Int64                                   `tfsdk:"region_id" json:"region_id,required"`
	Scale    *CloudInferenceDeploymentContainersScaleModel `tfsdk:"scale" json:"scale,required"`
}

type CloudInferenceDeploymentContainersScaleModel struct {
	Max             types.Int64                                                                    `tfsdk:"max" json:"max,required"`
	Min             types.Int64                                                                    `tfsdk:"min" json:"min,required"`
	CooldownPeriod  types.Int64                                                                    `tfsdk:"cooldown_period" json:"cooldown_period,computed_optional"`
	PollingInterval types.Int64                                                                    `tfsdk:"polling_interval" json:"polling_interval,computed_optional"`
	Triggers        customfield.NestedObject[CloudInferenceDeploymentContainersScaleTriggersModel] `tfsdk:"triggers" json:"triggers,computed_optional"`
}

type CloudInferenceDeploymentContainersScaleTriggersModel struct {
	CPU            *CloudInferenceDeploymentContainersScaleTriggersCPUModel            `tfsdk:"cpu" json:"cpu,optional"`
	GPUMemory      *CloudInferenceDeploymentContainersScaleTriggersGPUMemoryModel      `tfsdk:"gpu_memory" json:"gpu_memory,optional"`
	GPUUtilization *CloudInferenceDeploymentContainersScaleTriggersGPUUtilizationModel `tfsdk:"gpu_utilization" json:"gpu_utilization,optional"`
	HTTP           *CloudInferenceDeploymentContainersScaleTriggersHTTPModel           `tfsdk:"http" json:"http,optional"`
	Memory         *CloudInferenceDeploymentContainersScaleTriggersMemoryModel         `tfsdk:"memory" json:"memory,optional"`
	Sqs            *CloudInferenceDeploymentContainersScaleTriggersSqsModel            `tfsdk:"sqs" json:"sqs,optional"`
}

type CloudInferenceDeploymentContainersScaleTriggersCPUModel struct {
	Threshold types.Int64 `tfsdk:"threshold" json:"threshold,required"`
}

type CloudInferenceDeploymentContainersScaleTriggersGPUMemoryModel struct {
	Threshold types.Int64 `tfsdk:"threshold" json:"threshold,required"`
}

type CloudInferenceDeploymentContainersScaleTriggersGPUUtilizationModel struct {
	Threshold types.Int64 `tfsdk:"threshold" json:"threshold,required"`
}

type CloudInferenceDeploymentContainersScaleTriggersHTTPModel struct {
	Rate   types.Int64 `tfsdk:"rate" json:"rate,required"`
	Window types.Int64 `tfsdk:"window" json:"window,required"`
}

type CloudInferenceDeploymentContainersScaleTriggersMemoryModel struct {
	Threshold types.Int64 `tfsdk:"threshold" json:"threshold,required"`
}

type CloudInferenceDeploymentContainersScaleTriggersSqsModel struct {
	ActivationQueueLength types.Int64  `tfsdk:"activation_queue_length" json:"activation_queue_length,required"`
	AwsRegion             types.String `tfsdk:"aws_region" json:"aws_region,required"`
	QueueLength           types.Int64  `tfsdk:"queue_length" json:"queue_length,required"`
	QueueURL              types.String `tfsdk:"queue_url" json:"queue_url,required"`
	SecretName            types.String `tfsdk:"secret_name" json:"secret_name,required"`
	AwsEndpoint           types.String `tfsdk:"aws_endpoint" json:"aws_endpoint,optional"`
	ScaleOnDelayed        types.Bool   `tfsdk:"scale_on_delayed" json:"scale_on_delayed,computed_optional"`
	ScaleOnFlight         types.Bool   `tfsdk:"scale_on_flight" json:"scale_on_flight,computed_optional"`
}

type CloudInferenceDeploymentIngressOptsModel struct {
	DisableResponseBuffering types.Bool `tfsdk:"disable_response_buffering" json:"disable_response_buffering,computed_optional"`
}

type CloudInferenceDeploymentLoggingModel struct {
	DestinationRegionID types.Int64                                          `tfsdk:"destination_region_id" json:"destination_region_id,optional"`
	Enabled             types.Bool                                           `tfsdk:"enabled" json:"enabled,computed_optional"`
	RetentionPolicy     *CloudInferenceDeploymentLoggingRetentionPolicyModel `tfsdk:"retention_policy" json:"retention_policy,optional"`
	TopicName           types.String                                         `tfsdk:"topic_name" json:"topic_name,optional"`
}

type CloudInferenceDeploymentLoggingRetentionPolicyModel struct {
	Period types.Int64 `tfsdk:"period" json:"period,required"`
}

type CloudInferenceDeploymentProbesModel struct {
	LivenessProbe  customfield.NestedObject[CloudInferenceDeploymentProbesLivenessProbeModel]  `tfsdk:"liveness_probe" json:"liveness_probe,computed_optional"`
	ReadinessProbe customfield.NestedObject[CloudInferenceDeploymentProbesReadinessProbeModel] `tfsdk:"readiness_probe" json:"readiness_probe,computed_optional"`
	StartupProbe   customfield.NestedObject[CloudInferenceDeploymentProbesStartupProbeModel]   `tfsdk:"startup_probe" json:"startup_probe,computed_optional"`
}

type CloudInferenceDeploymentProbesLivenessProbeModel struct {
	Enabled types.Bool                                                                      `tfsdk:"enabled" json:"enabled,required"`
	Probe   customfield.NestedObject[CloudInferenceDeploymentProbesLivenessProbeProbeModel] `tfsdk:"probe" json:"probe,computed_optional"`
}

type CloudInferenceDeploymentProbesLivenessProbeProbeModel struct {
	Exec                *CloudInferenceDeploymentProbesLivenessProbeProbeExecModel                             `tfsdk:"exec" json:"exec,optional"`
	FailureThreshold    types.Int64                                                                            `tfsdk:"failure_threshold" json:"failure_threshold,computed_optional"`
	HTTPGet             customfield.NestedObject[CloudInferenceDeploymentProbesLivenessProbeProbeHTTPGetModel] `tfsdk:"http_get" json:"http_get,computed_optional"`
	InitialDelaySeconds types.Int64                                                                            `tfsdk:"initial_delay_seconds" json:"initial_delay_seconds,computed_optional"`
	PeriodSeconds       types.Int64                                                                            `tfsdk:"period_seconds" json:"period_seconds,computed_optional"`
	SuccessThreshold    types.Int64                                                                            `tfsdk:"success_threshold" json:"success_threshold,computed_optional"`
	TcpSocket           *CloudInferenceDeploymentProbesLivenessProbeProbeTcpSocketModel                        `tfsdk:"tcp_socket" json:"tcp_socket,optional"`
	TimeoutSeconds      types.Int64                                                                            `tfsdk:"timeout_seconds" json:"timeout_seconds,computed_optional"`
}

type CloudInferenceDeploymentProbesLivenessProbeProbeExecModel struct {
	Command *[]types.String `tfsdk:"command" json:"command,required"`
}

type CloudInferenceDeploymentProbesLivenessProbeProbeHTTPGetModel struct {
	Port    types.Int64              `tfsdk:"port" json:"port,required"`
	Headers *map[string]types.String `tfsdk:"headers" json:"headers,optional"`
	Host    types.String             `tfsdk:"host" json:"host,optional"`
	Path    types.String             `tfsdk:"path" json:"path,computed_optional"`
	Schema  types.String             `tfsdk:"schema" json:"schema,computed_optional"`
}

type CloudInferenceDeploymentProbesLivenessProbeProbeTcpSocketModel struct {
	Port types.Int64 `tfsdk:"port" json:"port,required"`
}

type CloudInferenceDeploymentProbesReadinessProbeModel struct {
	Enabled types.Bool                                                                       `tfsdk:"enabled" json:"enabled,required"`
	Probe   customfield.NestedObject[CloudInferenceDeploymentProbesReadinessProbeProbeModel] `tfsdk:"probe" json:"probe,computed_optional"`
}

type CloudInferenceDeploymentProbesReadinessProbeProbeModel struct {
	Exec                *CloudInferenceDeploymentProbesReadinessProbeProbeExecModel                             `tfsdk:"exec" json:"exec,optional"`
	FailureThreshold    types.Int64                                                                             `tfsdk:"failure_threshold" json:"failure_threshold,computed_optional"`
	HTTPGet             customfield.NestedObject[CloudInferenceDeploymentProbesReadinessProbeProbeHTTPGetModel] `tfsdk:"http_get" json:"http_get,computed_optional"`
	InitialDelaySeconds types.Int64                                                                             `tfsdk:"initial_delay_seconds" json:"initial_delay_seconds,computed_optional"`
	PeriodSeconds       types.Int64                                                                             `tfsdk:"period_seconds" json:"period_seconds,computed_optional"`
	SuccessThreshold    types.Int64                                                                             `tfsdk:"success_threshold" json:"success_threshold,computed_optional"`
	TcpSocket           *CloudInferenceDeploymentProbesReadinessProbeProbeTcpSocketModel                        `tfsdk:"tcp_socket" json:"tcp_socket,optional"`
	TimeoutSeconds      types.Int64                                                                             `tfsdk:"timeout_seconds" json:"timeout_seconds,computed_optional"`
}

type CloudInferenceDeploymentProbesReadinessProbeProbeExecModel struct {
	Command *[]types.String `tfsdk:"command" json:"command,required"`
}

type CloudInferenceDeploymentProbesReadinessProbeProbeHTTPGetModel struct {
	Port    types.Int64              `tfsdk:"port" json:"port,required"`
	Headers *map[string]types.String `tfsdk:"headers" json:"headers,optional"`
	Host    types.String             `tfsdk:"host" json:"host,optional"`
	Path    types.String             `tfsdk:"path" json:"path,computed_optional"`
	Schema  types.String             `tfsdk:"schema" json:"schema,computed_optional"`
}

type CloudInferenceDeploymentProbesReadinessProbeProbeTcpSocketModel struct {
	Port types.Int64 `tfsdk:"port" json:"port,required"`
}

type CloudInferenceDeploymentProbesStartupProbeModel struct {
	Enabled types.Bool                                                                     `tfsdk:"enabled" json:"enabled,required"`
	Probe   customfield.NestedObject[CloudInferenceDeploymentProbesStartupProbeProbeModel] `tfsdk:"probe" json:"probe,computed_optional"`
}

type CloudInferenceDeploymentProbesStartupProbeProbeModel struct {
	Exec                *CloudInferenceDeploymentProbesStartupProbeProbeExecModel                             `tfsdk:"exec" json:"exec,optional"`
	FailureThreshold    types.Int64                                                                           `tfsdk:"failure_threshold" json:"failure_threshold,computed_optional"`
	HTTPGet             customfield.NestedObject[CloudInferenceDeploymentProbesStartupProbeProbeHTTPGetModel] `tfsdk:"http_get" json:"http_get,computed_optional"`
	InitialDelaySeconds types.Int64                                                                           `tfsdk:"initial_delay_seconds" json:"initial_delay_seconds,computed_optional"`
	PeriodSeconds       types.Int64                                                                           `tfsdk:"period_seconds" json:"period_seconds,computed_optional"`
	SuccessThreshold    types.Int64                                                                           `tfsdk:"success_threshold" json:"success_threshold,computed_optional"`
	TcpSocket           *CloudInferenceDeploymentProbesStartupProbeProbeTcpSocketModel                        `tfsdk:"tcp_socket" json:"tcp_socket,optional"`
	TimeoutSeconds      types.Int64                                                                           `tfsdk:"timeout_seconds" json:"timeout_seconds,computed_optional"`
}

type CloudInferenceDeploymentProbesStartupProbeProbeExecModel struct {
	Command *[]types.String `tfsdk:"command" json:"command,required"`
}

type CloudInferenceDeploymentProbesStartupProbeProbeHTTPGetModel struct {
	Port    types.Int64              `tfsdk:"port" json:"port,required"`
	Headers *map[string]types.String `tfsdk:"headers" json:"headers,optional"`
	Host    types.String             `tfsdk:"host" json:"host,optional"`
	Path    types.String             `tfsdk:"path" json:"path,computed_optional"`
	Schema  types.String             `tfsdk:"schema" json:"schema,computed_optional"`
}

type CloudInferenceDeploymentProbesStartupProbeProbeTcpSocketModel struct {
	Port types.Int64 `tfsdk:"port" json:"port,required"`
}

type CloudInferenceDeploymentObjectReferencesModel struct {
	Kind types.String `tfsdk:"kind" json:"kind,computed"`
	Name types.String `tfsdk:"name" json:"name,computed"`
}
