//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/keypair/v2/keypairs"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/availablenetworks"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/routers"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	testK8sClusterName               = "test-cluster"
	testK8sClusterVersion            = "v1.28.1"
	testK8sClusterPoolName           = "test-pool"
	testK8sClusterPoolFlavor         = "g1-standard-1-2"
	testK8sClusterPoolMinNodeCount   = 1
	testK8sClusterPoolMaxNodeCount   = 1
	testK8sClusterPoolBootVolumeSize = 10
	testK8sClusterPoolBootVolumeType = volumes.Standard
	testK8sKeypairName               = "testkp"
)

func TestAccK8sV2DataSource(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	k8sClient, err := CreateTestClient(cfg.Provider, K8sPoint, versionPointV2)
	if err != nil {
		t.Fatal(err)
	}

	netClient, err := CreateTestClient(cfg.Provider, networksPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	subnetClient, err := CreateTestClient(cfg.Provider, subnetPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	kpClient, err := CreateTestClient(cfg.Provider, keypairsPoint, versionPointV2)
	if err != nil {
		t.Fatal(err)
	}

	netOpts := networks.CreateOpts{
		Name:         networkTestName,
		CreateRouter: true,
	}
	networkID, err := createTestNetwork(netClient, netOpts)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteTestNetwork(netClient, networkID)

	gw := net.ParseIP("")
	subnetOpts := subnets.CreateOpts{
		Name:                   subnetTestName,
		NetworkID:              networkID,
		ConnectToNetworkRouter: true,
		EnableDHCP:             true,
		GatewayIP:              &gw,
	}

	subnetID, err := CreateTestSubnet(subnetClient, subnetOpts)
	if err != nil {
		t.Fatal(err)
	}

	// update our new network router so that the k8s nodes will have access to the Nexus
	// registry to download images
	if err := patchRouterForK8S(cfg.Provider, networkID); err != nil {
		t.Fatal(err)
	}

	pid, err := strconv.Atoi(os.Getenv("TEST_PROJECT_ID"))
	if err != nil {
		t.Fatal(err)
	}

	kpOpts := keypairs.CreateOpts{
		Name:      testK8sKeypairName,
		PublicKey: pkTest,
		ProjectID: pid,
	}
	keyPair, err := keypairs.Create(kpClient, kpOpts).Extract()
	if err != nil {
		t.Fatal(err)
	}
	defer keypairs.Delete(kpClient, keyPair.ID)

	k8sOpts := clusters.CreateOpts{
		Name:         testK8sClusterName,
		FixedNetwork: networkID,
		FixedSubnet:  subnetID,
		KeyPair:      keyPair.ID,
		Version:      testK8sClusterVersion,
		Pools: []pools.CreateOpts{
			{
				Name:           testK8sClusterPoolName,
				FlavorID:       testK8sClusterPoolFlavor,
				BootVolumeSize: testK8sClusterPoolBootVolumeSize,
				BootVolumeType: testK8sClusterPoolBootVolumeType,
				MinNodeCount:   testK8sClusterPoolMinNodeCount,
				MaxNodeCount:   testK8sClusterPoolMaxNodeCount,
			},
		},
	}
	clusterName, err := createTestClusterV2(k8sClient, k8sOpts)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteTestClusterV2(k8sClient, clusterName)

	fullName := "data.gcore_k8sv2.acctest"
	ipTemplate := fmt.Sprintf(`
			data "gcore_k8sv2" "acctest" {
			  %s
              %s
              name = "%s"
			}
		`, projectInfo(), regionInfo(), clusterName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ipTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", clusterName),
				),
			},
		},
	})
}

func createTestClusterV2(client *gcorecloud.ServiceClient, opts clusters.CreateOpts) (string, error) {
	res, err := clusters.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	tasksClient, err := CreateTestClient(client.ProviderClient, tasksPoint, versionPointV1)
	if err != nil {
		return "", err
	}
	taskID := res.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		return nil, nil
	})
	if err != nil {
		return "", err
	}

	return opts.Name, nil
}

func deleteTestClusterV2(client *gcorecloud.ServiceClient, clusterID string) error {
	results, err := clusters.Delete(client, clusterID).Extract()
	if err != nil {
		return err
	}

	tasksClient, err := CreateTestClient(client.ProviderClient, tasksPoint, versionPointV1)
	if err != nil {
		return err
	}
	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(tasksClient, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := clusters.Get(client, clusterID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete k8s cluster with ID: %s", clusterID)
		}
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
	return err
}

func patchRouterForK8S(provider *gcorecloud.ProviderClient, networkID string) error {
	routersClient, err := CreateTestClient(provider, RouterPoint, versionPointV1)
	if err != nil {
		return err
	}

	aNetClient, err := CreateTestClient(provider, sharedNetworksPoint, versionPointV1)
	if err != nil {
		return err
	}

	availableNetworks, err := availablenetworks.ListAll(aNetClient, nil)
	if err != nil {
		return err
	}
	var extNet availablenetworks.Network
	for _, an := range availableNetworks {
		if an.External {
			extNet = an
			break
		}
	}

	rs, err := routers.ListAll(routersClient, nil)
	if err != nil {
		return err
	}

	var router routers.Router
	for _, r := range rs {
		if strings.Contains(r.Name, networkID) {
			router = r
			break
		}
	}

	extSubnet := extNet.Subnets[0]
	routerOpts := routers.UpdateOpts{Routes: extSubnet.HostRoutes}
	_, err = routers.Update(routersClient, router.ID, routerOpts).Extract()
	if err != nil {
		return err
	}
	return nil
}
