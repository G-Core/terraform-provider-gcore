package gcore

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	fastedge "github.com/G-Core/FastEdge-client-sdk-go"
	dnssdk "github.com/G-Core/gcore-dns-sdk-go"
	storageSDK "github.com/G-Core/gcore-storage-sdk-go"
	gcdn "github.com/G-Core/gcorelabscdn-go"
	gcdnProvider "github.com/G-Core/gcorelabscdn-go/gcore/provider"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	gc "github.com/G-Core/gcorelabscloud-go/gcore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform/version"
)

const (
	ProviderOptPermanentToken    = "permanent_api_token"
	ProviderOptSkipCredsAuthErr  = "ignore_creds_auth_error"
	ProviderOptSingleApiEndpoint = "api_endpoint"
	DefaultUserAgent             = "terraform-provider/%s"

	lifecyclePolicyResource = "gcore_lifecyclepolicy"
)

var AppVersion = "dev"

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
				// commented because it's broke all tests
				//AtLeastOneOf: []string{ProviderOptPermanentToken, "user_name"},
				//RequiredWith: []string{"user_name", "password"},
				Deprecated:  fmt.Sprintf("Use `%s` instead", ProviderOptPermanentToken),
				Description: "Gcore account username. Can also be set with the GCORE_USERNAME environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_USERNAME", nil),
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				// commented because it's broke all tests
				//RequiredWith: []string{"user_name", "password"},
				Deprecated:  fmt.Sprintf("Use `%s` instead", ProviderOptPermanentToken),
				Description: "Gcore account password. Can also be set with the GCORE_PASSWORD environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_PASSWORD", nil),
			},
			ProviderOptPermanentToken: {
				Type:     schema.TypeString,
				Optional: true,
				// commented because it's broke all tests
				//AtLeastOneOf: []string{ProviderOptPermanentToken, "user_name"},
				Sensitive:   true,
				Description: "A permanent [API-token](https://gcore.com/docs/account-settings/create-use-or-delete-a-permanent-api-token). Can also be set with the GCORE_PERMANENT_TOKEN environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_PERMANENT_TOKEN", nil),
			},
			ProviderOptSingleApiEndpoint: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A single API endpoint for all products. Will be used when specific product API url is not defined. Can also be set with the GCORE_API_ENDPOINT environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_API_ENDPOINT", "https://api.gcore.com"),
			},
			ProviderOptSkipCredsAuthErr: {
				Type:        schema.TypeBool,
				Optional:    true,
				Deprecated:  "Does not have any effect anymore.",
				Description: "Should be set to true when you are gonna to use storage resource with permanent API-token only.",
			},
			"gcore_platform": {
				Type:          schema.TypeString,
				Optional:      true,
				Deprecated:    "Use `gcore_platform_api` instead.",
				ConflictsWith: []string{"gcore_platform_api"},
				Description:   "Platform URL is used for generate JWT.",
				DefaultFunc:   schema.EnvDefaultFunc("GCORE_PLATFORM", nil),
			},
			"gcore_platform_api": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Platform URL is used for generate JWT (define only if you want to override Platform API endpoint). Can also be set with the GCORE_PLATFORM_API environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_PLATFORM_API", nil),
			},
			"gcore_api": {
				Type:          schema.TypeString,
				Optional:      true,
				Deprecated:    "Use `gcore_cloud_api` instead.",
				ConflictsWith: []string{"gcore_cloud_api"},
				Description:   "Region API.",
				DefaultFunc:   schema.EnvDefaultFunc("GCORE_API", nil),
			},
			"gcore_cloud_api": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region API (define only if you want to override Region API endpoint). Can also be set with the GCORE_CLOUD_API environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_CLOUD_API", nil),
			},
			"gcore_cdn_api": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CDN API (define only if you want to override CDN API endpoint). Can also be set with the GCORE_CDN_API environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_CDN_API", ""),
			},
			"gcore_storage_api": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Storage API (define only if you want to override Storage API endpoint). Can also be set with the GCORE_STORAGE_API environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_STORAGE_API", ""),
			},
			"gcore_dns_api": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS API (define only if you want to override DNS API endpoint). Can also be set with the GCORE_DNS_API environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_DNS_API", ""),
			},
			"gcore_fastedge_api": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "FastEdge API (define only if you want to override FastEdge API endpoint). Can also be set with the GCORE_FASTEDGE_API environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_FASTEDGE_API", ""),
			},
			"gcore_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client ID. Can also be set with the GCORE_CLIENT_ID environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("GCORE_CLIENT_ID", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gcore_ai_cluster":           resourceAICluster(),
			"gcore_volume":               resourceVolume(),
			"gcore_network":              resourceNetwork(),
			"gcore_subnet":               resourceSubnet(),
			"gcore_router":               resourceRouter(),
			"gcore_instance":             resourceInstance(),
			"gcore_instancev2":           resourceInstanceV2(),
			"gcore_keypair":              resourceKeypair(),
			"gcore_reservedfixedip":      resourceReservedFixedIP(),
			"gcore_floatingip":           resourceFloatingIP(),
			"gcore_loadbalancer":         resourceLoadBalancer(),
			"gcore_loadbalancerv2":       resourceLoadBalancerV2(),
			"gcore_lblistener":           resourceLbListener(),
			"gcore_lbpool":               resourceLBPool(),
			"gcore_lbmember":             resourceLBMember(),
			"gcore_securitygroup":        resourceSecurityGroup(),
			"gcore_baremetal":            resourceBmInstance(),
			"gcore_snapshot":             resourceSnapshot(),
			"gcore_servergroup":          resourceServerGroup(),
			"gcore_k8sv2":                resourceK8sV2(),
			"gcore_secret":               resourceSecret(),
			"gcore_laas_topic":           resourceLaaSTopic(),
			"gcore_faas_namespace":       resourceFaaSNamespace(),
			"gcore_faas_function":        resourceFaaSFunction(),
			"gcore_faas_key":             resourceFaaSKey(),
			"gcore_storage_s3":           resourceStorageS3(),
			"gcore_storage_s3_bucket":    resourceStorageS3Bucket(),
			DNSZoneResource:              resourceDNSZone(),
			DNSZoneRecordResource:        resourceDNSZoneRecord(),
			"gcore_storage_sftp":         resourceStorageSFTP(),
			"gcore_storage_sftp_key":     resourceStorageSFTPKey(),
			"gcore_cdn_resource":         resourceCDNResource(),
			"gcore_cdn_origingroup":      resourceCDNOriginGroup(),
			"gcore_cdn_originshielding":  resourceCDNOriginShielding(),
			"gcore_cdn_applied_preset":   resourceCDNAppliedPreset(),
			"gcore_cdn_rule":             resourceCDNRule(),
			"gcore_cdn_sslcert":          resourceCDNCert(),
			"gcore_cdn_rule_template":    resourceRuleTemplate(),
			"gcore_cdn_cacert":           resourceCDNCACert(),
			lifecyclePolicyResource:      resourceLifecyclePolicy(),
			"gcore_ddos_protection":      resourceDDoSProtection(),
			"gcore_inference_deployment": resourceInferenceDeployment(),
			"gcore_inference_secret":     resourceInferenceSecrets(),
			"gcore_registry_credential":  resourceRegistryCredential(),
			"gcore_gpu_baremetal_image":  resourceBaremetalImage(),
			"gcore_gpu_virtual_image":    resourceVirtualImage(),
			"gcore_fastedge_binary":      resourceFastEdgeBinary(),
			"gcore_fastedge_app":         resourceFastEdgeApp(),
			"gcore_fastedge_template":    resourceFastEdgeTemplate(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gcore_ai_cluster":             dataSourceAICluster(),
			"gcore_project":                dataSourceProject(),
			"gcore_region":                 dataSourceRegion(),
			"gcore_securitygroup":          dataSourceSecurityGroup(),
			"gcore_image":                  dataSourceImage(),
			"gcore_volume":                 dataSourceVolume(),
			"gcore_network":                dataSourceNetwork(),
			"gcore_subnet":                 dataSourceSubnet(),
			"gcore_router":                 dataSourceRouter(),
			"gcore_loadbalancer":           dataSourceLoadBalancer(),
			"gcore_loadbalancerv2":         dataSourceLoadBalancerV2(),
			"gcore_lblistener":             dataSourceLBListener(),
			"gcore_lbpool":                 dataSourceLBPool(),
			"gcore_instance":               dataSourceInstance(),
			"gcore_floatingip":             dataSourceFloatingIP(),
			"gcore_storage_s3":             dataSourceStorageS3(),
			"gcore_storage_s3_bucket":      dataSourceStorageS3Bucket(),
			"gcore_storage_sftp":           dataSourceStorageSFTP(),
			"gcore_storage_sftp_key":       dataSourceStorageSFTPKey(),
			"gcore_reservedfixedip":        dataSourceReservedFixedIP(),
			"gcore_servergroup":            dataSourceServerGroup(),
			"gcore_k8sv2":                  dataSourceK8sV2(),
			"gcore_k8sv2_kubeconfig":       dataSourceK8sV2KubeConfig(),
			"gcore_secret":                 dataSourceSecret(),
			"gcore_laas_hosts":             dataSourceLaaSHosts(),
			"gcore_laas_status":            dataSourceLaaSStatus(),
			"gcore_faas_namespace":         dataSourceFaaSNamespace(),
			"gcore_faas_key":               dataSourceFaaSKey(),
			"gcore_faas_function":          dataSourceFaaSFunction(),
			"gcore_ddos_profile_template":  dataSourceDDoSProfileTemplate(),
			"gcore_cdn_shielding_location": dataOriginShieldingLocation(),
			"gcore_cdn_preset":             dataPreset(),
			"gcore_cdn_client":             dataClient(),
			"gcore_inference_flavor":       dataSourceInferenceFlavor(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("user_name").(string)
	password := d.Get("password").(string)
	permanentToken := d.Get(ProviderOptPermanentToken).(string)
	apiEndpoint := d.Get(ProviderOptSingleApiEndpoint).(string)

	var diags diag.Diagnostics
	if permanentToken == "" &&
		(username == "" || password == "") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Field 'permanent_api_token' or 'username' and 'password' are required.",
			Detail:   "To use provider your should fill field 'permanent_api_token' or 'username' and 'password'.",
		})
	}

	cloudApi := d.Get("gcore_cloud_api").(string)
	if cloudApi == "" {
		cloudApi = d.Get("gcore_api").(string)
	}
	if cloudApi == "" {
		cloudApi = apiEndpoint + "/cloud"
	}

	cdnAPI := d.Get("gcore_cdn_api").(string)
	if cdnAPI == "" {
		cdnAPI = apiEndpoint
	}

	storageAPI := d.Get("gcore_storage_api").(string)
	if storageAPI == "" {
		storageAPI = apiEndpoint + "/storage"
	}

	dnsAPI := d.Get("gcore_dns_api").(string)
	if dnsAPI == "" {
		dnsAPI = apiEndpoint + "/dns"
	}

	fastedgeAPI := d.Get("gcore_fastedge_api").(string)
	if fastedgeAPI == "" {
		fastedgeAPI = apiEndpoint + "/fastedge"
	}

	platform := d.Get("gcore_platform_api").(string)
	if platform == "" {
		platform = d.Get("gcore_platform").(string)
	}
	if platform == "" {
		platform = apiEndpoint + "/iam"
	}

	clientID := d.Get("gcore_client_id").(string)

	var err error
	var provider *gcorecloud.ProviderClient
	if permanentToken != "" {
		provider, err = gc.APITokenClient(gcorecloud.APITokenOptions{
			APIURL:   cloudApi,
			APIToken: permanentToken,
		})
	} else {
		provider, err = gc.AuthenticatedClient(gcorecloud.AuthOptions{
			APIURL:      cloudApi,
			AuthURL:     platform,
			Username:    username,
			Password:    password,
			AllowReauth: true,
			ClientID:    clientID,
		})
	}
	if err != nil {
		provider = &gcorecloud.ProviderClient{}
		log.Printf("[ERROR] init auth client: %s\n", err)
	}
	provider.UserAgent.Prepend(fmt.Sprintf(DefaultUserAgent, AppVersion))
	// enable retries on Get requests with 5XX status codes using exponential backoff strategy with a maximum of 3
	// retries and a base interval of 2 seconds
	provider.EnableGetRetriesOn5XX(3, 2)

	cdnProvider := gcdnProvider.NewClient(cdnAPI, gcdnProvider.WithSignerFunc(func(req *http.Request) error {
		for k, v := range provider.AuthenticatedHeaders() {
			req.Header.Set(k, v)
		}

		return nil
	}))
	cdnService := gcdn.NewService(cdnProvider)

	provider.SetDebug(os.Getenv("TF_LOG") == "DEBUG")
	config := Config{
		Provider:  provider,
		CDNClient: cdnService,
		CDNMutex:  &sync.Mutex{},
	}

	userAgent := fmt.Sprintf("terraform/%s", version.Version)
	if storageAPI != "" {
		stHost, stPath, err := ExtractHostAndPath(storageAPI)
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("storage api url: %w", err))
		}
		config.StorageClient = storageSDK.NewSDK(
			stHost,
			stPath,
			storageSDK.WithBearerAuth(provider.AccessToken),
			storageSDK.WithPermanentTokenAuth(func() string { return permanentToken }),
			storageSDK.WithUserAgent(userAgent),
		)
	}
	if dnsAPI != "" {
		baseUrl, err := url.Parse(dnsAPI)
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("dns api url: %w", err))
		}
		authorizer := dnssdk.BearerAuth(provider.AccessToken())
		if permanentToken != "" {
			authorizer = dnssdk.PermanentAPIKeyAuth(permanentToken)
		}
		config.DNSClient = dnssdk.NewClient(
			authorizer,
			func(client *dnssdk.Client) {
				client.BaseURL = baseUrl
				client.Debug = os.Getenv("TF_LOG") == "DEBUG"
			},
			func(client *dnssdk.Client) {
				client.UserAgent = userAgent
			})
	}

	if fastedgeAPI != "" {
		authFunc := func(ctx context.Context, req *http.Request) error {
			if permanentToken != "" {
				req.Header.Set("Authorization", "APIKey "+permanentToken)
			} else {
				req.Header.Set("Authorization", "Bearer "+provider.AccessToken())
			}
			return nil
		}

		config.FastEdgeClient, err = fastedge.NewClientWithVersionCheck(
			fastedgeAPI,
			userAgent,
			"Terraform provider",
			fastedge.WithRequestEditorFn(authFunc),
		)
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("fastedge api init: %w", err))
		}
	}

	return &config, diags
}
