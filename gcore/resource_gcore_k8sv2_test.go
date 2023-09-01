//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/keypair/v2/keypairs"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccK8sV2(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cfg, err := createTestConfig()
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
		Name:      kpName,
		PublicKey: pkTest,
		ProjectID: pid,
	}
	keyPair, err := keypairs.Create(kpClient, kpOpts).Extract()
	if err != nil {
		t.Fatal(err)
	}
	defer keypairs.Delete(kpClient, keyPair.ID)

	fullName := "gcore_k8sv2.acctest"

	ipTemplate := fmt.Sprintf(`
			resource "gcore_k8sv2" "acctest" {
			  %s
              %s
              name = "tf-k8s"
			  fixed_network = "%s"
			  fixed_subnet = "%s"
              keypair = "%s"
			  version = "%s"
			  pool {
				name = "tf-pool1"
				flavor_id = "g1-standard-1-2"
				min_node_count = 1
				max_node_count = 1
				boot_volume_size = 10
				boot_volume_type = "standard"
			  }
			}
		`, projectInfo(), regionInfo(), networkID, subnetID, keyPair.ID, testK8sClusterVersion)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccK8sDestroy,
		Steps: []resource.TestStep{
			{
				Config: ipTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "tf-k8s"),
				),
			},
		},
	})
}

func testAccK8sV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, K8sPoint, versionPointV2)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_k8s" {
			continue
		}

		_, err := clusters.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("k8s cluster still exists")
		}
	}

	return nil
}

func TestDiffK8sV2ClusterPoolChange(t *testing.T) {
	tests := []struct {
		name                      string
		old, new                  interface{}
		wantAdd, wantUpd, wantDel []map[string]interface{}
	}{
		{
			name: "no change",
			old: []map[string]interface{}{
				{"name": "pool-1", "max_node_count": 1},
			},
			new: []map[string]interface{}{
				{"name": "pool-1", "max_node_count": 1},
			},
			wantUpd: []map[string]interface{}{
				{"name": "pool-1", "max_node_count": 1},
			},
		},
		{
			name: "remove pool",
			old: []map[string]interface{}{
				{"name": "pool-1", "max_node_count": 1},
				{"name": "pool-2", "max_node_count": 2},
			},
			new: []map[string]interface{}{
				{"name": "pool-2", "max_node_count": 2},
			},
			wantUpd: []map[string]interface{}{
				{"name": "pool-2", "max_node_count": 2},
			},
			wantDel: []map[string]interface{}{
				{"name": "pool-1", "max_node_count": 1},
			},
		},
		{
			name: "add pool",
			old: []map[string]interface{}{
				{"name": "pool-1", "max_node_count": 1},
			},
			new: []map[string]interface{}{
				{"name": "pool-1", "max_node_count": 1},
				{"name": "pool-2", "max_node_count": 2},
			},
			wantAdd: []map[string]interface{}{
				{"name": "pool-2", "max_node_count": 2},
			},
			wantUpd: []map[string]interface{}{
				{"name": "pool-1", "max_node_count": 1},
			},
		},
		{
			name: "add remove pool",
			old: []map[string]interface{}{
				{"name": "pool-1", "max_node_count": 1},
				{"name": "pool-2", "max_node_count": 2},
			},
			new: []map[string]interface{}{
				{"name": "pool-2", "max_node_count": 2},
				{"name": "pool-3", "max_node_count": 3},
			},
			wantAdd: []map[string]interface{}{
				{"name": "pool-3", "max_node_count": 3},
			},
			wantUpd: []map[string]interface{}{
				{"name": "pool-2", "max_node_count": 2},
			},
			wantDel: []map[string]interface{}{
				{"name": "pool-1", "max_node_count": 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			add, upd, del := diffK8sV2ClusterPoolChange(tt.old, tt.new)
			if !reflect.DeepEqual(add, tt.wantAdd) {
				t.Errorf("diffClusterPoolChange() add got: %v, want: %v", add, tt.wantAdd)
			}
			if !reflect.DeepEqual(upd, tt.wantUpd) {
				t.Errorf("diffClusterPoolChange() upd got: %v, want: %v", upd, tt.wantUpd)
			}
			if !reflect.DeepEqual(del, tt.wantDel) {
				t.Errorf("diffClusterPoolChange() del got: %v, want: %v", del, tt.wantDel)
			}
		})
	}
}
