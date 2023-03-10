//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	flavorID         = "g1-standard-1-2"
	instanceTestName = "test-vm"
	testOsDistro     = "ubuntu"
)

func TestAccInstanceDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	clientVolume, err := CreateTestClient(cfg.Provider, volumesPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	clientImage, err := CreateTestClient(cfg.Provider, imagesPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	imgs, err := images.ListAll(clientImage, nil)
	if err != nil {
		t.Fatal(err)
	}

	var img images.Image
	for _, i := range imgs {
		if i.OsDistro == testOsDistro {
			img = i
			break
		}
	}
	if img.ID == "" {
		t.Fatalf("images with os_distro='%s' does not exist", testOsDistro)
	}

	optsV := volumes.CreateOpts{
		Name:     volumeTestName,
		Size:     volumeTestSize * 5,
		Source:   volumes.Image,
		TypeName: volumes.Standard,
		ImageID:  img.ID,
	}
	volumeID, err := createTestVolume(clientVolume, optsV)
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, InstancePoint, versionPointV2)
	if err != nil {
		t.Fatal(err)
	}

	clientV1, err := CreateTestClient(cfg.Provider, InstancePoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	opts := instances.CreateOpts{
		Names:  []string{instanceTestName},
		Flavor: flavorID,
		Volumes: []instances.CreateVolumeOpts{{
			Source:    types.ExistingVolume,
			BootIndex: 0,
			VolumeID:  volumeID,
		}},
		Interfaces:     []instances.InterfaceInstanceCreateOpts{{InterfaceOpts: instances.InterfaceOpts{Type: types.ExternalInterfaceType}}},
		SecurityGroups: []gcorecloud.ItemID{},
	}

	res, err := instances.Create(client, opts).Extract()
	if err != nil {
		t.Fatal(err)
	}

	taskID := res.Tasks[0]
	instanceID, err := tasks.WaitTaskAndReturnResult(clientVolume, taskID, true, InstanceCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(clientVolume, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		id, err := instances.ExtractInstanceIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve instance ID from task info: %w", err)
		}
		return id, nil
	},
	)
	if err != nil {
		t.Fatal(err)
	}

	defer instances.Delete(clientV1, instanceID.(string), instances.DeleteOpts{Volumes: []string{volumeID}})

	fullName := "data.gcore_instance.acctest"
	tpl := func(name string) string {
		return fmt.Sprintf(`
			data "gcore_instance" "acctest" {
			  %s
              %s
              name = "%s"
			}
		`, projectInfo(), regionInfo(), name)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(instanceTestName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", instanceTestName),
					resource.TestCheckResourceAttr(fullName, "id", instanceID.(string)),
				),
			},
		},
	})
}
