package gcore

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/G-Core/gcorelabscloud-go/gcore/inference/v3/inferences"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const inferenceDeploymentPoint = "inferences"

func resourceInferenceDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInferenceDeploymentCreate,
		ReadContext:   resourceInferenceDeploymentRead,
		UpdateContext: resourceInferenceDeploymentUpdate,
		DeleteContext: resourceInferenceDeploymentDelete,
		Description:   "Represent inference deployment",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, inferenceName, err := ImportStringParserWithNoRegion(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("name", inferenceName)
				d.SetId(inferenceName)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"image": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"listening_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"containers": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "Region id for the container",
							Required:    true,
						},
						"ready_containers": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "Status of the containers deployment. Number of ready instances",
							Computed:    true,
						},
						"total_containers": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "Status of the containers deployment. Total number of instances",
							Computed:    true,
						},
						"cooldown_period": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "Cooldown period between scaling actions in seconds",
							Required:    true,
						},
						"scale_max": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "Maximum scale for the container",
							Required:    true,
						},
						"scale_min": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "Minimum scale for the container",
							Required:    true,
						},
						"polling_interval": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Polling interval for scaling triggers in seconds",
						},
						"triggers_cpu_threshold": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "CPU trigger threshold configuration",
							Optional:    true,
						},
						"triggers_memory_threshold": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "Memory trigger threshold configuration",
							Optional:    true,
						},
						"triggers_gpu_memory_threshold": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "GPU memory trigger threshold configuration. Calculated by DCGM_FI_DEV_MEM_COPY_UTIL metric",
							Optional:    true,
						},
						"triggers_gpu_utilization_threshold": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "GPU utilization trigger threshold configuration. Calculated by DCGM_FI_DEV_GPU_UTIL metric",
							Optional:    true,
						},
						"triggers_http_rate": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "Request count per 'window' seconds for the http trigger. Required if you use http trigger",
							Optional:    true,
						},
						"triggers_http_window": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "Time window for rate calculation in seconds. Required if you use http trigger",
							Optional:    true,
						},
						"triggers_sqs_secret_name": &schema.Schema{
							Type:        schema.TypeString,
							Description: "Name of the secret with AWS credentials. Required if you use SQS trigger",
							Optional:    true,
						},
						"triggers_sqs_aws_region": &schema.Schema{
							Type:        schema.TypeString,
							Description: "AWS region. Required if you use SQS trigger",
							Optional:    true,
						},
						"triggers_sqs_aws_endpoint": &schema.Schema{
							Type:        schema.TypeString,
							Description: "Custom AWS endpoint, left empty to use default aws endpoint",
							Optional:    true,
						},
						"triggers_sqs_queue_url": &schema.Schema{
							Type:        schema.TypeString,
							Description: "URL of the SQS queue. Required if you use SQS trigger",
							Optional:    true,
						},
						"triggers_sqs_queue_length": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "Number of messages for one replica",
							Optional:    true,
						},
						"triggers_sqs_activation_queue_length": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "Number of messages for activation",
							Optional:    true,
						},
						"triggers_sqs_scale_on_flight": &schema.Schema{
							Type:        schema.TypeBool,
							Description: "Scale on in-flight messages",
							Optional:    true,
						},
						"triggers_sqs_scale_on_delayed": &schema.Schema{
							Type:        schema.TypeBool,
							Description: "Scale on delayed messages",
							Optional:    true,
						},
					},
				},
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"flavor_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"envs": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"command": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"credentials_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"logging": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"destination_region_id": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"topic_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"retention_policy_period": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"liveness_probe": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     probeSchema(),
			},
			"readiness_probe": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     probeSchema(),
			},
			"startup_probe": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     probeSchema(),
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Datetime when the inference deployment was created. The format is 2025-12-28T19:14:44.180394",
				Computed:    true,
			},
		},
	}
}

func resourceInferenceDeploymentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Inference deployment creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := inferences.CreateInferenceDeploymentOpts{
		Name:          d.Get("name").(string),
		Image:         d.Get("image").(string),
		ListeningPort: d.Get("listening_port").(int),
		Description:   d.Get("description").(string),
		AuthEnabled:   d.Get("auth_enabled").(bool),
		FlavorName:    d.Get("flavor_name").(string),
		Containers:    nil,
		Probes:        nil,
	}

	if v, ok := d.GetOk("timeout"); ok {
		timeout := v.(int)
		opts.Timeout = &timeout
	}

	if v, ok := d.GetOk("envs"); ok {
		envs := v.(map[string]string)
		opts.Envs = envs
	}

	if v, ok := d.GetOk("command"); ok {
		command := v.(string)
		opts.Command = strings.Fields(command)
	}

	if v, ok := d.GetOk("credentials_name"); ok {
		credentialsName := v.(string)
		opts.CredentialsName = &credentialsName
	}

	if v, ok := d.GetOk("logging"); ok {
		logs := v.([]interface{})
		if len(logs) > 0 {
			logOpts := logs[0].(map[string]interface{})
			loggingOpts := inferences.CreateLoggingOpts{
				Enabled:             logOpts["enabled"].(bool),
				DestinationRegionID: logOpts["destination_region_id"].(int),
				TopicName:           logOpts["topic_name"].(string),
			}

			retention := logOpts["retention_policy_period"].(int)
			if retention != 0 {
				loggingOpts.RetentionPolicy = inferences.LoggingRetentionPolicy{
					Period: &retention,
				}
			}
			opts.Logging = &loggingOpts
		}
	}

	opts.Probes = &inferences.Probes{}
	if probe, ok := d.GetOk("liveness_probe"); ok && len(probe.([]interface{})) > 0 {
		probeOpts := probe.([]interface{})[0].(map[string]interface{})
		opts.Probes.LivenessProbe = toProbeConfig(probeOpts)
	} else {
		opts.Probes.LivenessProbe = &inferences.ProbeConfiguration{Enabled: false}
	}

	if probe, ok := d.GetOk("readiness_probe"); ok && len(probe.([]interface{})) > 0 {
		probeOpts := probe.([]interface{})[0].(map[string]interface{})
		opts.Probes.ReadinessProbe = toProbeConfig(probeOpts)
	} else {
		opts.Probes.ReadinessProbe = &inferences.ProbeConfiguration{Enabled: false}
	}

	if probe, ok := d.GetOk("startup_probe"); ok && len(probe.([]interface{})) > 0 {
		probeOpts := probe.([]interface{})[0].(map[string]interface{})
		opts.Probes.StartupProbe = toProbeConfig(probeOpts)
	} else {
		opts.Probes.StartupProbe = &inferences.ProbeConfiguration{Enabled: false}
	}

	if v, ok := d.GetOk("containers"); ok {
		containers := v.([]interface{})
		for _, container := range containers {
			opts.Containers = append(opts.Containers, toContainerOpts(container))
		}
	}

	results, err := inferences.CreateInferenceDeployment(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)

	inferenceName := d.Get("name").(string)
	clientV1, err := CreateClient(provider, d, tasksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = tasks.WaitTaskAndReturnResult(clientV1, taskID, true, SecretCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		inference, err := inferences.GetInferenceDeployment(client, inferenceName).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get inference deployment with name: %s. Error: %w", inferenceName, err)
		}
		return inference, nil
	},
	)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Inference deployment name (%s)", inferenceName)
	d.SetId(inferenceName)

	resourceInferenceDeploymentRead(ctx, d, m)

	log.Printf("[DEBUG] Finish inference deployment creating (%s)", inferenceName)
	return diags
}

func resourceInferenceDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start inference deployment reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	inferenceName := d.Id()
	log.Printf("[DEBUG] inference name = %s", inferenceName)

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	inference, err := inferences.GetInferenceDeployment(client, inferenceName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", inference.Name)
	d.Set("description", inference.Description)
	d.Set("image", inference.Image)
	d.Set("listening_port", inference.ListeningPort)
	d.Set("status", inference.Status)
	d.Set("auth_enabled", inference.AuthEnabled)
	d.Set("address", inference.Address)
	d.Set("timeout", inference.Timeout)
	d.Set("flavor_name", inference.FlavorName)
	d.Set("credentials_name", inference.CredentialsName)
	d.Set("logging", []interface{}{})

	if inference.Command != nil {
		if err := d.Set("command", inference.Command); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("envs", inference.Envs); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("containers", containersToSchema(inference.Containers)); err != nil {
		return diag.FromErr(err)
	}

	if inference.Logging != nil {
		if err := d.Set("logging", []interface{}{map[string]interface{}{
			"enabled":                 inference.Logging.Enabled,
			"destination_region_id":   inference.Logging.DestinationRegionID,
			"topic_name":              inference.Logging.TopicName,
			"retention_policy_period": inference.Logging.RetentionPolicy,
		}}); err != nil {
			return diag.FromErr(err)
		}
	}

	if inference.Probes != nil {
		if inference.Probes.LivenessProbe != nil {
			if err := d.Set("liveness_probe", probeToSchema(inference.Probes.LivenessProbe)); err != nil {
				return diag.FromErr(err)
			}
		}

		if inference.Probes.ReadinessProbe != nil {
			if err := d.Set("readiness_probe", probeToSchema(inference.Probes.ReadinessProbe)); err != nil {
				return diag.FromErr(err)
			}
		}

		if inference.Probes.StartupProbe != nil {
			if err := d.Set("startup_probe", probeToSchema(inference.Probes.StartupProbe)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	log.Println("[DEBUG] Finish inference deployment reading")
	return diags
}

func probeToSchema(probe *inferences.ProbeConfiguration) interface{} {
	result := make(map[string]interface{})
	result["enabled"] = probe.Enabled
	if probe.Probe != nil {
		result["failure_threshold"] = probe.Probe.FailureThreshold
		result["initial_delay_seconds"] = probe.Probe.InitialDelaySeconds
		result["period_seconds"] = probe.Probe.PeriodSeconds
		result["timeout_seconds"] = probe.Probe.TimeoutSeconds
		result["success_threshold"] = probe.Probe.SuccessThreshold
		if probe.Probe.Exec != nil {
			result["exec_command"] = strings.Join(probe.Probe.Exec.Command, " ")
		}

		if probe.Probe.TcpSocket != nil {
			result["tcp_socket_port"] = probe.Probe.TcpSocket.Port
		}

		if probe.Probe.HttpGet != nil {
			result["http_get_headers"] = probe.Probe.HttpGet.Headers
			result["http_get_host"] = *probe.Probe.HttpGet.Host
			result["http_get_path"] = probe.Probe.HttpGet.Path
			result["http_get_port"] = probe.Probe.HttpGet.Port
			result["http_get_schema"] = probe.Probe.HttpGet.Schema
		}
	}

	return []interface{}{result}
}

func containersToSchema(containers []inferences.Container) interface{} {
	result := make([]interface{}, 0, len(containers))
	for _, c := range containers {
		container := make(map[string]interface{})
		container["region_id"] = c.RegionID
		container["ready_containers"] = c.DeployStatus.Ready
		container["total_containers"] = c.DeployStatus.Total

		if c.Scale.Triggers.Cpu != nil {
			container["triggers_cpu_threshold"] = c.Scale.Triggers.Cpu.Threshold
		}
		if c.Scale.Triggers.Memory != nil {
			container["triggers_memory_threshold"] = c.Scale.Triggers.Memory.Threshold
		}
		if c.Scale.Triggers.GpuMemory != nil {
			container["triggers_gpu_memory_threshold"] = c.Scale.Triggers.GpuMemory.Threshold
		}
		if c.Scale.Triggers.GpuUtilization != nil {
			container["triggers_gpu_utilization_threshold"] = c.Scale.Triggers.GpuUtilization.Threshold
			container["triggers_http_rate"] = c.Scale.Triggers.Http.Rate
			container["triggers_http_window"] = c.Scale.Triggers.Http.Window
		}
		if c.Scale.Triggers.Http != nil {
			if c.Scale.Triggers.Http.Rate != nil {
				container["triggers_http_rate"] = c.Scale.Triggers.Http.Rate
			}

			if c.Scale.Triggers.Http.Window != nil {
				container["triggers_http_window"] = c.Scale.Triggers.Http.Window
			}
		}

		if c.Scale.Triggers.Sqs != nil {
			container["triggers_sqs_secret_name"] = c.Scale.Triggers.Sqs.SecretName
			container["triggers_sqs_aws_region"] = c.Scale.Triggers.Sqs.AwsRegion
			container["triggers_sqs_aws_endpoint"] = c.Scale.Triggers.Sqs.AwsEndpoint
			container["triggers_sqs_queue_url"] = c.Scale.Triggers.Sqs.QueueURL
			container["triggers_sqs_queue_length"] = c.Scale.Triggers.Sqs.QueueLength
			container["triggers_sqs_activation_queue_length"] = c.Scale.Triggers.Sqs.ActivationQueueLength
			container["triggers_sqs_scale_on_flight"] = c.Scale.Triggers.Sqs.ScaleOnFlight
			container["triggers_sqs_scale_on_delayed"] = c.Scale.Triggers.Sqs.ScaleOnDelayed
		}

		container["cooldown_period"] = 0
		if c.Scale.CooldownPeriod != nil {
			container["cooldown_period"] = c.Scale.CooldownPeriod
		}

		container["polling_interval"] = 0
		if c.Scale.PollingInterval != nil {
			container["polling_interval"] = c.Scale.PollingInterval
		}

		container["scale_max"] = c.Scale.Max
		container["scale_min"] = c.Scale.Min

		result = append(result, container)
	}

	return result
}

func resourceInferenceDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start inference deployment updating")
	config := m.(*Config)
	provider := config.Provider
	inferenceName := d.Id()
	log.Printf("[DEBUG] inference deployment = %s", inferenceName)

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := inferences.UpdateInferenceDeploymentOpts{}
	if d.HasChange("image") {
		image := d.Get("image").(string)
		opts.Image = &image
	}

	if d.HasChange("flavor_name") {
		flavorName := d.Get("flavor_name").(string)
		opts.FlavorName = &flavorName
	}

	if d.HasChange("listening_port") {
		listeningPort := d.Get("listening_port").(int)
		opts.ListeningPort = &listeningPort
	}

	if d.HasChange("auth_enabled") {
		authEnabled := d.Get("auth_enabled").(bool)
		opts.AuthEnabled = &authEnabled
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
	}

	if d.HasChange("timeout") {
		timeout := d.Get("timeout").(int)
		opts.Timeout = &timeout
	}

	if d.HasChange("envs") {
		envs := d.Get("envs").(map[string]string)
		opts.Envs = envs
	}

	if d.HasChange("command") {
		command := d.Get("command").(string)
		opts.Command = strings.Fields(command)
	}

	if d.HasChange("credentials_name") {
		credentialsName := d.Get("credentials_name").(string)
		opts.CredentialsName = &credentialsName
	}

	if d.HasChange("logging") {
		logs := d.Get("logging").([]interface{})
		if len(logs) > 0 {
			logOpts := logs[0].(map[string]interface{})
			loggingOpts := inferences.CreateLoggingOpts{
				Enabled:             logOpts["enabled"].(bool),
				DestinationRegionID: logOpts["destination_region_id"].(int),
				TopicName:           logOpts["topic_name"].(string),
			}

			retention := logOpts["retention_policy_period"].(int)
			if retention != 0 {
				loggingOpts.RetentionPolicy = inferences.LoggingRetentionPolicy{Period: &retention}
			}
			opts.Logging = &loggingOpts
		}
	}

	if d.HasChange("liveness_probe") || d.HasChange("readiness_probe") || d.HasChange("startup_probe") {
		opts.Probes = &inferences.Probes{}
	}

	if d.HasChange("liveness_probe") {
		if probe, ok := d.Get("liveness_probe").([]interface{}); ok && len(probe) > 0 {
			probeOpts := probe[0].(map[string]interface{})
			opts.Probes.LivenessProbe = toProbeConfig(probeOpts)
		} else {
			opts.Probes.LivenessProbe = &inferences.ProbeConfiguration{Enabled: false}
		}
	}

	if d.HasChange("readiness_probe") {
		if probe, ok := d.Get("readiness_probe").([]interface{}); ok && len(probe) > 0 {
			probeOpts := probe[0].(map[string]interface{})
			opts.Probes.ReadinessProbe = toProbeConfig(probeOpts)
		} else {
			opts.Probes.ReadinessProbe = &inferences.ProbeConfiguration{Enabled: false}
		}
	}

	if d.HasChange("startup_probe") {
		if probe, ok := d.Get("startup_probe").([]interface{}); ok && len(probe) > 0 {
			probeOpts := probe[0].(map[string]interface{})
			opts.Probes.StartupProbe = toProbeConfig(probeOpts)
		} else {
			opts.Probes.StartupProbe = &inferences.ProbeConfiguration{Enabled: false}
		}
	}

	if d.HasChange("containers") {
		containers := d.Get("containers").([]interface{})
		for _, container := range containers {
			opts.Containers = append(opts.Containers, toContainerOpts(container))
		}
	}

	results, err := inferences.UpdateInferenceDeployment(client, inferenceName, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)

	clientV1, err := CreateClient(provider, d, tasksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = tasks.WaitTaskAndReturnResult(clientV1, taskID, true, SecretCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		return nil, nil
	},
	)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Finish of inference deployment updating")
	return resourceInferenceDeploymentRead(ctx, d, m)
}

func toContainerOpts(container interface{}) inferences.CreateContainerOpts {
	containerOpts := container.(map[string]interface{})
	containerScale := inferences.ContainerScale{
		Max: containerOpts["scale_max"].(int),
		Min: containerOpts["scale_min"].(int),
	}

	if v, ok := containerOpts["triggers_cpu_threshold"].(int); ok && v != 0 {
		containerScale.Triggers.Cpu = &inferences.ScaleTriggerThreshold{Threshold: v}
	}
	if v, ok := containerOpts["triggers_memory_threshold"].(int); ok && v != 0 {
		containerScale.Triggers.Memory = &inferences.ScaleTriggerThreshold{Threshold: v}
	}
	if v, ok := containerOpts["triggers_gpu_memory_threshold"].(int); ok && v != 0 {
		containerScale.Triggers.GpuMemory = &inferences.ScaleTriggerThreshold{Threshold: v}
	}
	if v, ok := containerOpts["triggers_gpu_utilization_threshold"].(int); ok && v != 0 {
		containerScale.Triggers.GpuUtilization = &inferences.ScaleTriggerThreshold{Threshold: v}
	}

	if cooldownPeriod, ok := containerOpts["cooldown_period"].(int); cooldownPeriod != 0 && ok {
		containerScale.CooldownPeriod = &cooldownPeriod
	}
	if pollingInterval, ok := containerOpts["polling_interval"].(int); pollingInterval != 0 && ok {
		containerScale.PollingInterval = &pollingInterval
	}

	if httpRate, ok := containerOpts["triggers_http_rate"].(int); httpRate != 0 && ok {
		containerScale.Triggers.Http = &inferences.ScaleTriggerHttp{Rate: &httpRate}

		if httpWindow, ok := containerOpts["triggers_http_window"].(int); httpWindow != 0 && ok {
			containerScale.Triggers.Http.Window = &httpWindow
		}
	}

	if secretName, ok := containerOpts["triggers_sqs_secret_name"].(string); ok && len(secretName) > 0 {
		containerScale.Triggers.Sqs = &inferences.ScaleTriggerSqs{
			SecretName:            secretName,
			AwsRegion:             containerOpts["triggers_sqs_aws_region"].(string),
			QueueURL:              containerOpts["triggers_sqs_queue_url"].(string),
			QueueLength:           containerOpts["triggers_sqs_queue_length"].(int),
			ActivationQueueLength: containerOpts["triggers_sqs_activation_queue_length"].(int),
			ScaleOnFlight:         containerOpts["triggers_sqs_scale_on_flight"].(bool),
			ScaleOnDelayed:        containerOpts["triggers_sqs_scale_on_delayed"].(bool),
		}

		if endpoint, ok := containerOpts["triggers_sqs_aws_endpoint"].(string); ok && len(endpoint) > 0 {
			containerScale.Triggers.Sqs.AwsEndpoint = &endpoint
		}
	}

	return inferences.CreateContainerOpts{
		RegionID: containerOpts["region_id"].(int),
		Scale:    containerScale,
	}
}

func resourceInferenceDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start inference deployment deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	inferenceName := d.Id()
	log.Printf("[DEBUG] inference deployment = %s", inferenceName)

	client, err := CreateClient(provider, d, inferenceDeploymentPoint, versionPointV3)
	if err != nil {
		return diag.FromErr(err)
	}

	results, err := inferences.DeleteInferenceDeployment(client, inferenceName).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)

	clientV1, err := CreateClient(provider, d, tasksPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = tasks.WaitTaskAndReturnResult(clientV1, taskID, true, SecretCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := inferences.GetInferenceDeployment(client, inferenceName).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete inference deployment with name: %s. Error: %w", inferenceName, err)
		}
		return nil, nil
	},
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Finish of inference deployment deleting")
	return diags
}

func probeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"failure_threshold": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"initial_delay_seconds": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"period_seconds": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"timeout_seconds": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"success_threshold": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"exec_command": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tcp_socket_port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"http_get_headers": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"http_get_host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_get_path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_get_port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"http_get_schema": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func toProbeConfig(probe map[string]interface{}) *inferences.ProbeConfiguration {
	probeConfig := inferences.ProbeConfiguration{
		Enabled: probe["enabled"].(bool),
	}

	if !probeConfig.Enabled {
		return &probeConfig
	}

	probeConfig.Probe = &inferences.Probe{}

	probeConfig.Probe.FailureThreshold = probe["failure_threshold"].(int)
	probeConfig.Probe.InitialDelaySeconds = probe["initial_delay_seconds"].(int)
	probeConfig.Probe.PeriodSeconds = probe["period_seconds"].(int)
	probeConfig.Probe.TimeoutSeconds = probe["timeout_seconds"].(int)
	probeConfig.Probe.SuccessThreshold = probe["success_threshold"].(int)

	if exec, ok := probe["exec_command"].(string); ok && len(exec) > 0 {
		probeConfig.Probe.Exec = &inferences.ExecProbe{
			Command: strings.Fields(exec),
		}
	}
	if tcpSocket, ok := probe["tcp_socket_port"].(int); ok && tcpSocket != 0 {
		probeConfig.Probe.TcpSocket = &inferences.TcpSocketProbe{
			Port: tcpSocket,
		}
	}

	if host, ok := probe["http_get_host"].(string); ok && len(host) > 0 {
		if probeConfig.Probe.HttpGet == nil {
			probeConfig.Probe.HttpGet = &inferences.HttpGetProbe{}
		}
		probeConfig.Probe.HttpGet.Host = &host
	}

	if path, ok := probe["http_get_path"].(string); ok && len(path) > 0 {
		if probeConfig.Probe.HttpGet == nil {
			probeConfig.Probe.HttpGet = &inferences.HttpGetProbe{}
		}
		probeConfig.Probe.HttpGet.Path = path
	}

	if port, ok := probe["http_get_port"].(int); ok && port > 0 {
		if probeConfig.Probe.HttpGet == nil {
			probeConfig.Probe.HttpGet = &inferences.HttpGetProbe{}
		}
		probeConfig.Probe.HttpGet.Port = port
	}

	if httpSchema, ok := probe["http_get_schema"].(string); ok && len(httpSchema) > 0 {
		if probeConfig.Probe.HttpGet == nil {
			probeConfig.Probe.HttpGet = &inferences.HttpGetProbe{}
		}
		probeConfig.Probe.HttpGet.Schema = httpSchema
	}

	if headersRaw, ok := probe["http_get_headers"].(map[string]interface{}); ok && len(headersRaw) > 0 {
		if probeConfig.Probe.HttpGet == nil {
			probeConfig.Probe.HttpGet = &inferences.HttpGetProbe{}
		}

		headers := make(map[string]string)
		for k, v := range headersRaw {
			headers[k] = v.(string)
		}
		probeConfig.Probe.HttpGet.Headers = headers
	}

	if probeConfig.Probe.HttpGet != nil {
		if probeConfig.Probe.HttpGet.Schema == "" {
			probeConfig.Probe.HttpGet.Schema = "HTTP"
		}
		if probeConfig.Probe.HttpGet.Headers == nil {
			probeConfig.Probe.HttpGet.Headers = map[string]string{}
		}
	}

	return &probeConfig
}
