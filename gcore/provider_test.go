package gcore

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	dnssdk "github.com/G-Core/gcore-dns-sdk-go"
	storageSDK "github.com/G-Core/gcore-storage-sdk-go"
	gcdn "github.com/G-Core/gcorelabscdn-go"
	gcdnProvider "github.com/G-Core/gcorelabscdn-go/gcore/provider"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	gc "github.com/G-Core/gcorelabscloud-go/gcore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type VarName string

const (
	GCORE_USERNAME_VAR             VarName = "GCORE_USERNAME"
	GCORE_PASSWORD_VAR             VarName = "GCORE_PASSWORD"
	GCORE_PERMANENT_TOKEN_VAR      VarName = "GCORE_PERMANENT_TOKEN"
	GCORE_CDN_URL_VAR              VarName = "GCORE_CDN_URL"
	GCORE_STORAGE_URL_VAR          VarName = "GCORE_STORAGE_API"
	GCORE_DNS_URL_VAR              VarName = "GCORE_DNS_API"
	GCORE_IMAGE_VAR                VarName = "GCORE_IMAGE"
	GCORE_SECGROUP_VAR             VarName = "GCORE_SECGROUP"
	GCORE_EXT_NET_VAR              VarName = "GCORE_EXT_NET"
	GCORE_PRIV_NET_VAR             VarName = "GCORE_PRIV_NET"
	GCORE_PRIV_SUBNET_VAR          VarName = "GCORE_PRIV_SUBNET"
	GCORE_LB_ID_VAR                VarName = "GCORE_LB_ID"
	GCORE_LBLISTENER_ID_VAR        VarName = "GCORE_LBLISTENER_ID"
	GCORE_LBPOOL_ID_VAR            VarName = "GCORE_LBPOOL_ID"
	GCORE_VOLUME_ID_VAR            VarName = "GCORE_VOLUME_ID"
	GCORE_CDN_ORIGINGROUP_ID_VAR   VarName = "GCORE_CDN_ORIGINGROUP_ID"
	GCORE_CDN_RESOURCE_ID_VAR      VarName = "GCORE_CDN_RESOURCE_ID"
	GCORE_CDN_PRESET_ID_VAR        VarName = "GCORE_CDN_PRESET_ID"
	GCORE_CDN_PRESET_OBJECT_ID_VAR VarName = "GCORE_CDN_PRESET_OBJECT_ID"
	GCORE_NETWORK_ID_VAR           VarName = "GCORE_NETWORK_ID"
	GCORE_SUBNET_ID_VAR            VarName = "GCORE_SUBNET_ID"
	GCORE_CLUSTER_ID_VAR           VarName = "GCORE_CLUSTER_ID"
	GCORE_CLUSTER_POOL_ID_VAR      VarName = "GCORE_CLUSTER_POOL_ID"
)

func getEnv(name VarName) string {
	return os.Getenv(string(name))
}

var (
	GCORE_USERNAME             = getEnv(GCORE_USERNAME_VAR)
	GCORE_PASSWORD             = getEnv(GCORE_PASSWORD_VAR)
	GCORE_PERMANENT_TOKEN      = getEnv(GCORE_PERMANENT_TOKEN_VAR)
	GCORE_CDN_URL              = getEnv(GCORE_CDN_URL_VAR)
	GCORE_IMAGE                = getEnv(GCORE_IMAGE_VAR)
	GCORE_SECGROUP             = getEnv(GCORE_SECGROUP_VAR)
	GCORE_EXT_NET              = getEnv(GCORE_EXT_NET_VAR)
	GCORE_PRIV_NET             = getEnv(GCORE_PRIV_NET_VAR)
	GCORE_PRIV_SUBNET          = getEnv(GCORE_PRIV_SUBNET_VAR)
	GCORE_LB_ID                = getEnv(GCORE_LB_ID_VAR)
	GCORE_LBLISTENER_ID        = getEnv(GCORE_LBLISTENER_ID_VAR)
	GCORE_LBPOOL_ID            = getEnv(GCORE_LBPOOL_ID_VAR)
	GCORE_VOLUME_ID            = getEnv(GCORE_VOLUME_ID_VAR)
	GCORE_CDN_ORIGINGROUP_ID   = getEnv(GCORE_CDN_ORIGINGROUP_ID_VAR)
	GCORE_CDN_RESOURCE_ID      = getEnv(GCORE_CDN_RESOURCE_ID_VAR)
	GCORE_CDN_PRESET_ID        = getEnv(GCORE_CDN_PRESET_ID_VAR)
	GCORE_CDN_PRESET_OBJECT_ID = getEnv(GCORE_CDN_PRESET_OBJECT_ID_VAR)
	GCORE_STORAGE_API          = getEnv(GCORE_STORAGE_URL_VAR)
	GCORE_DNS_API              = getEnv(GCORE_DNS_URL_VAR)
	GCORE_NETWORK_ID           = getEnv(GCORE_NETWORK_ID_VAR)
	GCORE_SUBNET_ID            = getEnv(GCORE_SUBNET_ID_VAR)
	GCORE_CLUSTER_ID           = getEnv(GCORE_CLUSTER_ID_VAR)
	GCORE_CLUSTER_POOL_ID      = getEnv(GCORE_CLUSTER_POOL_ID_VAR)
)

var varsMap = map[VarName]string{
	GCORE_USERNAME_VAR:             GCORE_USERNAME,
	GCORE_PASSWORD_VAR:             GCORE_PASSWORD,
	GCORE_PERMANENT_TOKEN_VAR:      GCORE_PERMANENT_TOKEN,
	GCORE_CDN_URL_VAR:              GCORE_CDN_URL,
	GCORE_IMAGE_VAR:                GCORE_IMAGE,
	GCORE_SECGROUP_VAR:             GCORE_SECGROUP,
	GCORE_EXT_NET_VAR:              GCORE_EXT_NET,
	GCORE_PRIV_NET_VAR:             GCORE_PRIV_NET,
	GCORE_PRIV_SUBNET_VAR:          GCORE_PRIV_SUBNET,
	GCORE_LB_ID_VAR:                GCORE_LB_ID,
	GCORE_LBLISTENER_ID_VAR:        GCORE_LBLISTENER_ID,
	GCORE_LBPOOL_ID_VAR:            GCORE_LBPOOL_ID,
	GCORE_VOLUME_ID_VAR:            GCORE_VOLUME_ID,
	GCORE_CDN_ORIGINGROUP_ID_VAR:   GCORE_CDN_ORIGINGROUP_ID,
	GCORE_CDN_RESOURCE_ID_VAR:      GCORE_CDN_RESOURCE_ID,
	GCORE_CDN_PRESET_ID_VAR:        GCORE_CDN_PRESET_ID,
	GCORE_CDN_PRESET_OBJECT_ID_VAR: GCORE_CDN_PRESET_OBJECT_ID,
	GCORE_STORAGE_URL_VAR:          GCORE_STORAGE_API,
	GCORE_DNS_URL_VAR:              GCORE_DNS_API,
	GCORE_NETWORK_ID_VAR:           GCORE_NETWORK_ID,
	GCORE_SUBNET_ID_VAR:            GCORE_SUBNET_ID,
	GCORE_CLUSTER_ID_VAR:           GCORE_CLUSTER_ID,
	GCORE_CLUSTER_POOL_ID_VAR:      GCORE_CLUSTER_POOL_ID,
}

func testAccPreCheckVars(t *testing.T, vars ...VarName) {
	for _, name := range vars {
		if val := varsMap[name]; val == "" {
			t.Fatalf("'%s' must be set for acceptance test", name)
		}
	}
}

var testAccProvider *schema.Provider
var testAccProviders map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"gcore": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// nolint: deadcode,unused
func testAccPreCheck(t *testing.T) {
	vars := map[string]interface{}{
		"GCORE_USERNAME": GCORE_USERNAME,
		"GCORE_PASSWORD": GCORE_PASSWORD,
	}
	for k, v := range vars {
		if v == "" {
			t.Fatalf("'%s' must be set for acceptance test", k)
		}
	}
	checkNameAndID("PROJECT", t)
	checkNameAndID("REGION", t)
}

// nolint: unused
func checkNameAndID(resourceType string, t *testing.T) {
	// resourceType is a word in capital letters
	keyID := fmt.Sprintf("TEST_%s_ID", resourceType)
	keyName := fmt.Sprintf("TEST_%s_NAME", resourceType)
	_, haveID := os.LookupEnv(keyID)
	_, haveName := os.LookupEnv(keyName)
	if !haveID && !haveName {
		t.Fatalf("%s or %s must be set for acceptance tests", keyID, keyName)
	}
	if haveID && haveName {
		t.Fatalf("Use only one from environment variables: %s or %s", keyID, keyName)
	}
}

// nolint: deadcode,unused
func regionInfo() string {
	return objectInfo("REGION")
}

// nolint: deadcode,unused
func projectInfo() string {
	return objectInfo("PROJECT")
}

// nolint: unused
func objectInfo(resourceType string) string {
	// resourceType is a word in capital letters
	keyID := fmt.Sprintf("TEST_%s_ID", resourceType)
	keyName := fmt.Sprintf("TEST_%s_NAME", resourceType)
	if objectID, exists := os.LookupEnv(keyID); exists {
		return fmt.Sprintf(`%s_id = %s`, strings.ToLower(resourceType), objectID)
	}
	return fmt.Sprintf(`%s_name = "%s"`, strings.ToLower(resourceType), os.Getenv(keyName))
}

func CreateTestClient(provider *gcorecloud.ProviderClient, endpoint, version string) (*gcorecloud.ServiceClient, error) {
	projectID := 0
	var err error
	if strProjectID, exists := os.LookupEnv("TEST_PROJECT_ID"); exists {
		projectID, err = strconv.Atoi(strProjectID)
		if err != nil {
			return nil, err
		}
	} else {
		projectID, err = GetProject(provider, 0, os.Getenv("TEST_PROJECT_NAME"))
		if err != nil {
			return nil, err
		}
	}
	regionID := 0
	if strRegionID, exists := os.LookupEnv("TEST_REGION_ID"); exists {
		regionID, err = strconv.Atoi(strRegionID)
		if err != nil {
			return nil, err
		}
	} else {
		regionID, err = GetProject(provider, 0, os.Getenv("TEST_REGION_NAME"))
		if err != nil {
			return nil, err
		}
	}

	client, err := gc.ClientServiceFromProvider(provider, gcorecloud.EndpointOpts{
		Name:    endpoint,
		Region:  regionID,
		Project: projectID,
		Version: version,
	})

	if err != nil {
		return nil, err
	}
	return client, nil
}

func createTestConfig() (*Config, error) {
	provider, err := gc.AuthenticatedClient(gcorecloud.AuthOptions{
		APIURL:      os.Getenv("GCORE_API"),
		AuthURL:     os.Getenv("GCORE_PLATFORM"),
		Username:    os.Getenv("GCORE_USERNAME"),
		Password:    os.Getenv("GCORE_PASSWORD"),
		AllowReauth: true,
	})
	if err != nil {
		return nil, err
	}

	cdnProvider := gcdnProvider.NewClient(GCORE_CDN_URL, gcdnProvider.WithSignerFunc(func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+provider.AccessToken())
		return nil
	}))
	cdnService := gcdn.NewService(cdnProvider)

	storageAPI := GCORE_STORAGE_API
	stHost, stPath, err := ExtractHostAndPath(storageAPI)
	var storageClient *storageSDK.SDK
	if err == nil {
		storageClient = storageSDK.NewSDK(stHost, stPath, storageSDK.WithBearerAuth(provider.AccessToken))
	}

	var dnsClient *dnssdk.Client
	if GCORE_DNS_API != "" {
		baseUrl, err := url.Parse(GCORE_DNS_API)
		if err == nil {
			authorizer := dnssdk.BearerAuth(provider.AccessToken())
			dnsClient = dnssdk.NewClient(authorizer, func(client *dnssdk.Client) {
				client.BaseURL = baseUrl
			})
		}

	}

	config := Config{
		Provider:      provider,
		CDNClient:     cdnService,
		StorageClient: storageClient,
		DNSClient:     dnsClient,
	}

	return &config, nil
}

func testAccCheckResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Widget ID is not set")
		}
		return nil
	}
}
