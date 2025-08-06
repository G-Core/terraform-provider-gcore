// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package internal

import (
	"context"
	"os"
	"strconv"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_audit_log"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_baremetal_server"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_billing_reservation"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_file_share"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_file_share_access_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_floating_ip"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_gpu_baremetal_cluster"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_gpu_baremetal_cluster_image"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_api_key"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_deployment"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_deployment_log"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_flavor"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_model"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_registry_credential"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_secret"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_instance"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_instance_image"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_l7_policy"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_l7_policy_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_listener"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_pool"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_pool_health_monitor"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_status"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_network"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_network_router"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_network_subnet"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_placement_group"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_project"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_quota_request"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_region"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_registry"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_registry_user"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_reserved_fixed_ip"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_secret"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_security_group"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_security_group_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_ssh_key"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_task"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_usage_report"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_user_role_assignment"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_volume"
	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_app"
	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_app_log"
	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_binary"
	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_kv_store"
	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_secret"
	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_template"
	"github.com/stainless-sdks/gcore-terraform/internal/services/iam_api_token"
	"github.com/stainless-sdks/gcore-terraform/internal/services/iam_user"
	"github.com/stainless-sdks/gcore-terraform/internal/services/security_event"
	"github.com/stainless-sdks/gcore-terraform/internal/services/security_profile"
	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_ai_task"
	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_broadcast"
	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_directory"
	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_player"
	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_playlist"
	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_restream"
	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_stream"
	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_stream_overlay"
	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_video"
	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_video_subtitle"
	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_custom_page_set"
	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain"
	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_advanced_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_api_path"
	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_custom_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_firewall_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_insight"
	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_insight_silence"
	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_setting"
	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_organization"
	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_tag"
)

var _ provider.ProviderWithConfigValidators = (*GcoreProvider)(nil)

// GcoreProvider defines the provider implementation.
type GcoreProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// GcoreProviderModel describes the provider data model.
type GcoreProviderModel struct {
	BaseURL                     types.String `tfsdk:"base_url" json:"base_url,optional"`
	APIKey                      types.String `tfsdk:"api_key" json:"api_key,optional"`
	CloudProjectID              types.Int64  `tfsdk:"cloud_project_id" json:"cloud_project_id,optional"`
	CloudRegionID               types.Int64  `tfsdk:"cloud_region_id" json:"cloud_region_id,optional"`
	CloudPollingIntervalSeconds types.Int64  `tfsdk:"cloud_polling_interval_seconds" json:"cloud_polling_interval_seconds,optional"`
}

func (p *GcoreProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gcore"
	resp.Version = p.version
}

func ProviderSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Description: "Set the base url that the provider connects to.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Optional: true,
			},
			"cloud_project_id": schema.Int64Attribute{
				Optional: true,
			},
			"cloud_region_id": schema.Int64Attribute{
				Optional: true,
			},
			"cloud_polling_interval_seconds": schema.Int64Attribute{
				Optional: true,
			},
		},
	}
}

func (p *GcoreProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = ProviderSchema(ctx)
}

func (p *GcoreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var data GcoreProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	opts := []option.RequestOption{}

	if !data.BaseURL.IsNull() && !data.BaseURL.IsUnknown() {
		opts = append(opts, option.WithBaseURL(data.BaseURL.ValueString()))
	} else if o, ok := os.LookupEnv("GCORE_BASE_URL"); ok {
		opts = append(opts, option.WithBaseURL(o))
	}

	if !data.APIKey.IsNull() && !data.APIKey.IsUnknown() {
		opts = append(opts, option.WithAPIKey(data.APIKey.ValueString()))
	} else if o, ok := os.LookupEnv("GCORE_API_KEY"); ok {
		opts = append(opts, option.WithAPIKey(o))
	} else {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing api_key value",
			"The api_key field is required. Set it in provider configuration or via the \"GCORE_API_KEY\" environment variable.",
		)
		return
	}

	if !data.CloudProjectID.IsNull() && !data.CloudProjectID.IsUnknown() {
		opts = append(opts, option.WithCloudProjectID(data.CloudProjectID.ValueInt64()))
	} else if o, ok := os.LookupEnv("GCORE_CLOUD_PROJECT_ID"); ok {
		parsed, err := strconv.ParseInt(o, 10, 64)
		if err != nil {
			resp.Diagnostics.Append(diag.NewErrorDiagnostic("failed to parse environment variable: GCORE_CLOUD_PROJECT_ID", err.Error()))
			return
		}
		opts = append(opts, option.WithCloudProjectID(parsed))
	}

	if !data.CloudRegionID.IsNull() && !data.CloudRegionID.IsUnknown() {
		opts = append(opts, option.WithCloudRegionID(data.CloudRegionID.ValueInt64()))
	} else if o, ok := os.LookupEnv("GCORE_CLOUD_REGION_ID"); ok {
		parsed, err := strconv.ParseInt(o, 10, 64)
		if err != nil {
			resp.Diagnostics.Append(diag.NewErrorDiagnostic("failed to parse environment variable: GCORE_CLOUD_REGION_ID", err.Error()))
			return
		}
		opts = append(opts, option.WithCloudRegionID(parsed))
	}

	if !data.CloudPollingIntervalSeconds.IsNull() && !data.CloudPollingIntervalSeconds.IsUnknown() {
		opts = append(opts, option.WithCloudPollingIntervalSeconds(data.CloudPollingIntervalSeconds.ValueInt64()))
	} else {
		opts = append(opts, option.WithCloudPollingIntervalSeconds(3))
	}

	client := gcore.NewClient(
		opts...,
	)

	resp.DataSourceData = &client
	resp.ResourceData = &client
}

func (p *GcoreProvider) ConfigValidators(_ context.Context) []provider.ConfigValidator {
	return []provider.ConfigValidator{}
}

func (p *GcoreProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		cloud_project.NewResource,
		cloud_quota_request.NewResource,
		cloud_ssh_key.NewResource,
		cloud_load_balancer.NewResource,
		cloud_load_balancer_l7_policy.NewResource,
		cloud_load_balancer_l7_policy_rule.NewResource,
		cloud_load_balancer_listener.NewResource,
		cloud_load_balancer_pool.NewResource,
		cloud_load_balancer_pool_health_monitor.NewResource,
		cloud_reserved_fixed_ip.NewResource,
		cloud_network.NewResource,
		cloud_network_subnet.NewResource,
		cloud_network_router.NewResource,
		cloud_volume.NewResource,
		cloud_floating_ip.NewResource,
		cloud_security_group.NewResource,
		cloud_security_group_rule.NewResource,
		cloud_user_role_assignment.NewResource,
		cloud_inference_deployment.NewResource,
		cloud_inference_registry_credential.NewResource,
		cloud_inference_secret.NewResource,
		cloud_inference_api_key.NewResource,
		cloud_placement_group.NewResource,
		cloud_baremetal_server.NewResource,
		cloud_registry.NewResource,
		cloud_registry_user.NewResource,
		cloud_file_share.NewResource,
		cloud_file_share_access_rule.NewResource,
		cloud_gpu_baremetal_cluster.NewResource,
		cloud_instance.NewResource,
		cloud_instance_image.NewResource,
		waap_domain.NewResource,
		waap_domain_setting.NewResource,
		waap_domain_api_path.NewResource,
		waap_domain_insight_silence.NewResource,
		waap_domain_custom_rule.NewResource,
		waap_domain_firewall_rule.NewResource,
		waap_domain_advanced_rule.NewResource,
		waap_custom_page_set.NewResource,
		iam_api_token.NewResource,
		iam_user.NewResource,
		fastedge_template.NewResource,
		fastedge_secret.NewResource,
		fastedge_binary.NewResource,
		fastedge_app.NewResource,
		fastedge_kv_store.NewResource,
		streaming_ai_task.NewResource,
		streaming_broadcast.NewResource,
		streaming_directory.NewResource,
		streaming_player.NewResource,
		streaming_playlist.NewResource,
		streaming_video.NewResource,
		streaming_video_subtitle.NewResource,
		streaming_stream.NewResource,
		streaming_stream_overlay.NewResource,
		streaming_restream.NewResource,
		security_profile.NewResource,
	}
}

func (p *GcoreProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		cloud_project.NewCloudProjectDataSource,
		cloud_project.NewCloudProjectsDataSource,
		cloud_task.NewCloudTaskDataSource,
		cloud_task.NewCloudTasksDataSource,
		cloud_region.NewCloudRegionDataSource,
		cloud_region.NewCloudRegionsDataSource,
		cloud_quota_request.NewCloudQuotaRequestDataSource,
		cloud_quota_request.NewCloudQuotaRequestsDataSource,
		cloud_secret.NewCloudSecretDataSource,
		cloud_secret.NewCloudSecretsDataSource,
		cloud_ssh_key.NewCloudSSHKeyDataSource,
		cloud_ssh_key.NewCloudSSHKeysDataSource,
		cloud_load_balancer.NewCloudLoadBalancerDataSource,
		cloud_load_balancer.NewCloudLoadBalancersDataSource,
		cloud_load_balancer_l7_policy.NewCloudLoadBalancerL7PolicyDataSource,
		cloud_load_balancer_l7_policy_rule.NewCloudLoadBalancerL7PolicyRuleDataSource,
		cloud_load_balancer_listener.NewCloudLoadBalancerListenerDataSource,
		cloud_load_balancer_pool.NewCloudLoadBalancerPoolDataSource,
		cloud_load_balancer_status.NewCloudLoadBalancerStatusDataSource,
		cloud_reserved_fixed_ip.NewCloudReservedFixedIPDataSource,
		cloud_reserved_fixed_ip.NewCloudReservedFixedIPsDataSource,
		cloud_network.NewCloudNetworkDataSource,
		cloud_network.NewCloudNetworksDataSource,
		cloud_network_subnet.NewCloudNetworkSubnetDataSource,
		cloud_network_subnet.NewCloudNetworkSubnetsDataSource,
		cloud_network_router.NewCloudNetworkRouterDataSource,
		cloud_network_router.NewCloudNetworkRoutersDataSource,
		cloud_volume.NewCloudVolumeDataSource,
		cloud_volume.NewCloudVolumesDataSource,
		cloud_floating_ip.NewCloudFloatingIPDataSource,
		cloud_floating_ip.NewCloudFloatingIPsDataSource,
		cloud_security_group.NewCloudSecurityGroupDataSource,
		cloud_security_group.NewCloudSecurityGroupsDataSource,
		cloud_user_role_assignment.NewCloudUserRoleAssignmentsDataSource,
		cloud_inference_flavor.NewCloudInferenceFlavorDataSource,
		cloud_inference_flavor.NewCloudInferenceFlavorsDataSource,
		cloud_inference_model.NewCloudInferenceModelDataSource,
		cloud_inference_model.NewCloudInferenceModelsDataSource,
		cloud_inference_deployment.NewCloudInferenceDeploymentDataSource,
		cloud_inference_deployment.NewCloudInferenceDeploymentsDataSource,
		cloud_inference_deployment_log.NewCloudInferenceDeploymentLogsDataSource,
		cloud_inference_registry_credential.NewCloudInferenceRegistryCredentialDataSource,
		cloud_inference_registry_credential.NewCloudInferenceRegistryCredentialsDataSource,
		cloud_inference_secret.NewCloudInferenceSecretDataSource,
		cloud_inference_secret.NewCloudInferenceSecretsDataSource,
		cloud_inference_api_key.NewCloudInferenceAPIKeyDataSource,
		cloud_inference_api_key.NewCloudInferenceAPIKeysDataSource,
		cloud_placement_group.NewCloudPlacementGroupDataSource,
		cloud_baremetal_server.NewCloudBaremetalServersDataSource,
		cloud_registry.NewCloudRegistryDataSource,
		cloud_file_share.NewCloudFileShareDataSource,
		cloud_file_share.NewCloudFileSharesDataSource,
		cloud_billing_reservation.NewCloudBillingReservationDataSource,
		cloud_billing_reservation.NewCloudBillingReservationsDataSource,
		cloud_gpu_baremetal_cluster.NewCloudGPUBaremetalClusterDataSource,
		cloud_gpu_baremetal_cluster.NewCloudGPUBaremetalClustersDataSource,
		cloud_gpu_baremetal_cluster_image.NewCloudGPUBaremetalClusterImageDataSource,
		cloud_instance.NewCloudInstanceDataSource,
		cloud_instance.NewCloudInstancesDataSource,
		cloud_instance_image.NewCloudInstanceImageDataSource,
		cloud_audit_log.NewCloudAuditLogsDataSource,
		cloud_usage_report.NewCloudUsageReportDataSource,
		waap_domain.NewWaapDomainDataSource,
		waap_domain.NewWaapDomainsDataSource,
		waap_domain_setting.NewWaapDomainSettingDataSource,
		waap_domain_api_path.NewWaapDomainAPIPathDataSource,
		waap_domain_api_path.NewWaapDomainAPIPathsDataSource,
		waap_domain_insight.NewWaapDomainInsightDataSource,
		waap_domain_insight.NewWaapDomainInsightsDataSource,
		waap_domain_insight_silence.NewWaapDomainInsightSilenceDataSource,
		waap_domain_insight_silence.NewWaapDomainInsightSilencesDataSource,
		waap_domain_custom_rule.NewWaapDomainCustomRuleDataSource,
		waap_domain_custom_rule.NewWaapDomainCustomRulesDataSource,
		waap_domain_firewall_rule.NewWaapDomainFirewallRuleDataSource,
		waap_domain_firewall_rule.NewWaapDomainFirewallRulesDataSource,
		waap_domain_advanced_rule.NewWaapDomainAdvancedRuleDataSource,
		waap_domain_advanced_rule.NewWaapDomainAdvancedRulesDataSource,
		waap_custom_page_set.NewWaapCustomPageSetDataSource,
		waap_custom_page_set.NewWaapCustomPageSetsDataSource,
		waap_tag.NewWaapTagsDataSource,
		waap_organization.NewWaapOrganizationsDataSource,
		iam_api_token.NewIamAPITokenDataSource,
		iam_user.NewIamUserDataSource,
		iam_user.NewIamUsersDataSource,
		fastedge_template.NewFastedgeTemplateDataSource,
		fastedge_template.NewFastedgeTemplatesDataSource,
		fastedge_secret.NewFastedgeSecretDataSource,
		fastedge_binary.NewFastedgeBinaryDataSource,
		fastedge_app.NewFastedgeAppDataSource,
		fastedge_app.NewFastedgeAppsDataSource,
		fastedge_app_log.NewFastedgeAppLogsDataSource,
		fastedge_kv_store.NewFastedgeKvStoreDataSource,
		streaming_ai_task.NewStreamingAITaskDataSource,
		streaming_ai_task.NewStreamingAITasksDataSource,
		streaming_broadcast.NewStreamingBroadcastDataSource,
		streaming_broadcast.NewStreamingBroadcastsDataSource,
		streaming_directory.NewStreamingDirectoryDataSource,
		streaming_player.NewStreamingPlayerDataSource,
		streaming_player.NewStreamingPlayersDataSource,
		streaming_playlist.NewStreamingPlaylistDataSource,
		streaming_playlist.NewStreamingPlaylistsDataSource,
		streaming_video.NewStreamingVideoDataSource,
		streaming_video.NewStreamingVideosDataSource,
		streaming_video_subtitle.NewStreamingVideoSubtitleDataSource,
		streaming_stream.NewStreamingStreamDataSource,
		streaming_stream.NewStreamingStreamsDataSource,
		streaming_stream_overlay.NewStreamingStreamOverlayDataSource,
		streaming_restream.NewStreamingRestreamDataSource,
		streaming_restream.NewStreamingRestreamsDataSource,
		security_event.NewSecurityEventsDataSource,
		security_profile.NewSecurityProfileDataSource,
	}
}

func NewProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GcoreProvider{
			version: version,
		}
	}
}
