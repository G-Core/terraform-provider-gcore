//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/aiimages"
	ai "github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/ais"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	inst_types "github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/keypair/v2/keypairs"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
)

const (
	testAIClusterName       = "acctestds"
	testAIClusterFlavor     = "g2a-ai-fake-v1pod-8"
	testAIClusterVolumeType = volumes.Standard
	testAIClusterKeypair    = "testkp"
)

func TestAccAIClusterDataSource(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	aiClusterClient, err := CreateTestClient(cfg.Provider, AIClusterPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	kpClient, err := CreateTestClient(cfg.Provider, keypairsPoint, versionPointV2)
	if err != nil {
		t.Fatal(err)
	}

	ProjectID, err := strconv.Atoi(os.Getenv("TEST_PROJECT_ID"))
	if err != nil {
		t.Fatal(err)
	}

	kpOpts := keypairs.CreateOpts{
		Name:      testAIClusterKeypair,
		PublicKey: pkTest,
		ProjectID: ProjectID,
	}
	keyPair, err := keypairs.Create(kpClient, kpOpts).Extract()
	if err != nil {
		t.Fatal(err)
	}
	defer keypairs.Delete(kpClient, keyPair.ID)

	sgClient, err := CreateTestClient(cfg.Provider, securityGroupPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	sgList, err := securitygroups.ListAll(sgClient, securitygroups.ListOpts{})
	if err != nil {
		t.Fatal(err)
	}
	var defaultSgID string
	for _, sg := range sgList {
		if sg.Name == "default" {
			defaultSgID = sg.ID
			break
		}
	}

	aiImageClient, err := CreateTestClient(cfg.Provider, AIImagesPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}
	opts := aiimages.AIImageListOpts{}
	aiImages, err := aiimages.ListAll(aiImageClient, opts)
	if err != nil {
		t.Fatal(err)
	}

	testAIImage := aiImages[0].ID
	AIClusterInterfaces := []instances.InterfaceInstanceCreateOpts{
		{
			InterfaceOpts: instances.InterfaceOpts{
				Type: inst_types.ExternalInterfaceType,
			},
		},
	}

	aiClusterOpts := ai.CreateOpts{
		Flavor:     testAIClusterFlavor,
		Name:       testAIClusterName,
		ImageID:    testAIImage,
		Interfaces: AIClusterInterfaces,
		Keypair:    kpOpts.Name,
		Volumes: []instances.CreateVolumeOpts{
			{
				Source:              inst_types.Image,
				ImageID:             testAIImage,
				Size:                30,
				TypeName:            volumes.Standard,
				DeleteOnTermination: true,
			},
		},
		SecurityGroups: []gcorecloud.ItemID{{ID: defaultSgID}},
		Metadata:       map[string]string{"qqqq": "qqtest"},
	}

	clusterID, err := createTestAICluster(aiClusterClient, aiClusterOpts)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteTestAICluster(aiClusterClient, clusterID)

	fullName := "data.gcore_ai_cluster.acctestds"
	aiDataTemplate := fmt.Sprintf(`
			data "gcore_ai_cluster" "acctestds" {
			  %s
              %s
              cluster_id = "%s"
			}
		`, projectInfo(), regionInfo(), clusterID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aiDataTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "cluster_id", clusterID),
				),
			},
		},
	})
}

func createTestAICluster(client *gcorecloud.ServiceClient, opts ai.CreateOpts) (string, error) {
	results, err := ai.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskClient, err := CreateTestClient(client.ProviderClient, tasksPoint, versionPointV1)
	if err != nil {
		return "", err
	}
	taskID := results.Tasks[0]
	clusterID, err := tasks.WaitTaskAndReturnResult(taskClient, taskID, true, AIClusterCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(taskClient, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		clusterID, err := ai.ExtractAIClusterIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve AI cluster ID from task info: %w", err)
		}
		return clusterID, nil
	},
	)
	if err != nil {
		return "", err
	}
	return clusterID.(string), nil
}

func deleteTestAICluster(client *gcorecloud.ServiceClient, clusterID string) error {
	results, err := ai.Delete(client, clusterID, nil).Extract()
	if err != nil {
		return err
	}

	taskClient, err := CreateTestClient(client.ProviderClient, tasksPoint, versionPointV1)
	if err != nil {
		return err
	}
	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(taskClient, taskID, true, AIClusterDeletingTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := ai.Get(client, clusterID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete AI cluster with ID: %s", clusterID)
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
