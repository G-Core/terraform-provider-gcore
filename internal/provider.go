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
	"github.com/stainless-sdks/gcore-terraform/internal/services/cdn_cdn_resource_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cdn_certificate"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cdn_origin_group"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cdn_trusted_ca_certificate"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_file_share"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_file_share_access_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_floating_ip"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_gpu_baremetal_cluster_image"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_gpu_virtual_cluster"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_gpu_virtual_cluster_image"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_deployment"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_registry_credential"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_secret"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_instance"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_instance_image"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_k8s_cluster"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_listener"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_pool"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_pool_member"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_network"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_network_router"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_network_subnet"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_placement_group"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_project"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_region"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_reserved_fixed_ip"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_secret"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_security_group"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_security_group_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_ssh_key"
	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_volume"
	"github.com/stainless-sdks/gcore-terraform/internal/services/dns_network_mapping"
	"github.com/stainless-sdks/gcore-terraform/internal/services/dns_zone"
	"github.com/stainless-sdks/gcore-terraform/internal/services/dns_zone_rrset"
	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_app"
	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_binary"
	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_secret"
	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain"
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
	CloudPollingTimeoutSeconds  types.Int64  `tfsdk:"cloud_polling_timeout_seconds" json:"cloud_polling_timeout_seconds,optional"`
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
			"cloud_polling_timeout_seconds": schema.Int64Attribute{
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

	if !data.CloudPollingTimeoutSeconds.IsNull() && !data.CloudPollingTimeoutSeconds.IsUnknown() {
		opts = append(opts, option.WithCloudPollingTimeoutSeconds(data.CloudPollingTimeoutSeconds.ValueInt64()))
	} else {
		opts = append(opts, option.WithCloudPollingTimeoutSeconds(7200))
	}

	// Override Go SDK max retries to 4 from 2 which is the default.
	// The max delay is capped at 8 secs, so the maximum value for max retries is 4.
	opts = append(opts, option.WithMaxRetries(4))

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
		cloud_secret.NewResource,
		cloud_ssh_key.NewResource,
		cloud_load_balancer.NewResource,
		cloud_load_balancer_listener.NewResource,
		cloud_load_balancer_pool.NewResource,
		cloud_load_balancer_pool_member.NewResource,
		cloud_reserved_fixed_ip.NewResource,
		cloud_network.NewResource,
		cloud_network_subnet.NewResource,
		cloud_network_router.NewResource,
		cloud_volume.NewResource,
		cloud_floating_ip.NewResource,
		cloud_security_group.NewResource,
		cloud_security_group_rule.NewResource,
		cloud_inference_deployment.NewResource,
		cloud_inference_registry_credential.NewResource,
		cloud_inference_secret.NewResource,
		cloud_placement_group.NewResource,
		cloud_file_share.NewResource,
		cloud_file_share_access_rule.NewResource,
		cloud_gpu_baremetal_cluster_image.NewResource,
		cloud_gpu_virtual_cluster.NewResource,
		cloud_gpu_virtual_cluster_image.NewResource,
		cloud_instance.NewResource,
		cloud_instance_image.NewResource,
		cloud_k8s_cluster.NewResource,
		waap_domain.NewResource,
		fastedge_secret.NewResource,
		fastedge_binary.NewResource,
		fastedge_app.NewResource,
		dns_zone.NewResource,
		dns_zone_rrset.NewResource,
		dns_network_mapping.NewResource,
		cdn_cdn_resource_rule.NewResource,
		cdn_origin_group.NewResource,
		cdn_certificate.NewResource,
		cdn_trusted_ca_certificate.NewResource,
	}
}

func (p *GcoreProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		cloud_project.NewCloudProjectDataSource,
		cloud_project.NewCloudProjectsDataSource,
		cloud_region.NewCloudRegionDataSource,
		cloud_region.NewCloudRegionsDataSource,
		cloud_secret.NewCloudSecretDataSource,
		cloud_secret.NewCloudSecretsDataSource,
		cloud_ssh_key.NewCloudSSHKeyDataSource,
		cloud_ssh_key.NewCloudSSHKeysDataSource,
		cloud_load_balancer.NewCloudLoadBalancerDataSource,
		cloud_load_balancer.NewCloudLoadBalancersDataSource,
		cloud_load_balancer_listener.NewCloudLoadBalancerListenerDataSource,
		cloud_load_balancer_pool.NewCloudLoadBalancerPoolDataSource,
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
		cloud_inference_deployment.NewCloudInferenceDeploymentDataSource,
		cloud_inference_deployment.NewCloudInferenceDeploymentsDataSource,
		cloud_inference_registry_credential.NewCloudInferenceRegistryCredentialDataSource,
		cloud_inference_registry_credential.NewCloudInferenceRegistryCredentialsDataSource,
		cloud_inference_secret.NewCloudInferenceSecretDataSource,
		cloud_inference_secret.NewCloudInferenceSecretsDataSource,
		cloud_placement_group.NewCloudPlacementGroupDataSource,
		cloud_file_share.NewCloudFileShareDataSource,
		cloud_file_share.NewCloudFileSharesDataSource,
		cloud_gpu_baremetal_cluster_image.NewCloudGPUBaremetalClusterImageDataSource,
		cloud_gpu_virtual_cluster.NewCloudGPUVirtualClusterDataSource,
		cloud_gpu_virtual_cluster.NewCloudGPUVirtualClustersDataSource,
		cloud_gpu_virtual_cluster_image.NewCloudGPUVirtualClusterImageDataSource,
		cloud_instance.NewCloudInstanceDataSource,
		cloud_instance.NewCloudInstancesDataSource,
		cloud_instance_image.NewCloudInstanceImageDataSource,
		cloud_k8s_cluster.NewCloudK8SClusterDataSource,
		waap_domain.NewWaapDomainDataSource,
		waap_domain.NewWaapDomainsDataSource,
		fastedge_secret.NewFastedgeSecretDataSource,
		fastedge_binary.NewFastedgeBinaryDataSource,
		fastedge_app.NewFastedgeAppDataSource,
		fastedge_app.NewFastedgeAppsDataSource,
		cdn_origin_group.NewCDNOriginGroupDataSource,
		dns_zone.NewDNSZoneDataSource,
		dns_zone_rrset.NewDNSZoneRrsetDataSource,
		dns_network_mapping.NewDNSNetworkMappingDataSource,
		cdn_cdn_resource_rule.NewCDNCDNResourceRuleDataSource,
		cdn_certificate.NewCDNCertificateDataSource,
		cdn_trusted_ca_certificate.NewCDNTrustedCaCertificateDataSource,
	}
}

func NewProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GcoreProvider{
			version: version,
		}
	}
}
