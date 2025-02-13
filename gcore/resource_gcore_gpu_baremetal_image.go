package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/gpu/v3/images"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBaremetalImage() *schema.Resource {
	return resourceGPUImage(GPUImageTypeBaremetal)
}

func resourceBaremetalImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Start baremetal image creation")
	config := m.(*Config)
	log.Printf("[DEBUG] Config: %+v", config)
	client, err := CreateClient(config.Provider, d, "gpu/baremetal", "v3")
	if err != nil {
		log.Printf("[ERROR] Error creating client: %s", err)
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Client created: %+v", client)

	metadata := make(map[string]interface{})
	if v, ok := d.GetOk("metadata"); ok {
		metadata = v.(map[string]interface{})
		log.Printf("[DEBUG] Using metadata: %v", metadata)
	}

	// Add additional metadata from specific fields
	if v, ok := d.GetOk("os_type"); ok {
		metadata["os_type"] = v.(string)
		log.Printf("[DEBUG] Setting os_type: %s", v.(string))
	}
	if v, ok := d.GetOk("os_distro"); ok {
		metadata["os_distro"] = v.(string)
		log.Printf("[DEBUG] Setting os_distro: %s", v.(string))
	}
	if v, ok := d.GetOk("os_version"); ok {
		metadata["os_version"] = v.(string)
		log.Printf("[DEBUG] Setting os_version: %s", v.(string))
	}
	if v, ok := d.GetOk("hw_firmware_type"); ok {
		metadata["hw_firmware_type"] = v.(string)
		log.Printf("[DEBUG] Setting hw_firmware_type: %s", v.(string))
	}
	if v, ok := d.GetOk("ssh_key"); ok {
		metadata["ssh_key"] = v.(string)
		log.Printf("[DEBUG] Setting ssh_key: %s", v.(string))
	}

	opts := images.ImageOpts{
		Name:     d.Get("name").(string),
		URL:      d.Get("url").(string),
		Metadata: metadata,
	}

	if v, ok := d.GetOk("cow_format"); ok {
		cowFormat := v.(bool)
		opts.CowFormat = &cowFormat
		log.Printf("[DEBUG] Setting cow_format: %t", cowFormat)
	}

	if v, ok := d.GetOk("architecture"); ok {
		arch := v.(string)
		opts.Architecture = &arch
	}

	if v, ok := d.GetOk("os_type"); ok {
		osType := images.ImageOsType(v.(string))
		opts.OsType = &osType
	}

	if v, ok := d.GetOk("os_distro"); ok {
		osDistro := v.(string)
		opts.OsDistro = &osDistro
	}

	if v, ok := d.GetOk("os_version"); ok {
		osVersion := v.(string)
		opts.OsVersion = &osVersion
	}

	if v, ok := d.GetOk("hw_firmware_type"); ok {
		hwType := images.ImageHwFirmwareType(v.(string))
		opts.HwFirmwareType = &hwType
	}

	if v, ok := d.GetOk("ssh_key"); ok {
		sshKey := images.SshKeyType(v.(string))
		opts.SshKey = &sshKey
	}

	log.Printf("[DEBUG] Baremetal image create options: %+v", opts)

	result := images.UploadImage(client, opts)
	if result.Err != nil {
		log.Printf("[ERROR] Error uploading baremetal image: %s", result.Err)
		return diag.FromErr(fmt.Errorf("error uploading baremetal image: %w", result.Err))
	}

	taskResponse := result.Body.(map[string]interface{})
	taskID := tasks.TaskID(taskResponse["tasks"].([]interface{})[0].(string))
	log.Printf("[DEBUG] Upload task started with ID: %s", taskID)

	taskClient, err := CreateClient(config.Provider, d, "tasks", "v1")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating task client: %w", err))
	}

	waitResult, err := tasks.WaitTaskAndReturnResult(taskClient, taskID, true, BaremetalImageCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(taskClient, string(task)).Extract()
		if err != nil {
			log.Printf("[ERROR] Error getting task %s info: %s", task, err)
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		log.Printf("[DEBUG] Task %s state: %s", task, taskInfo.State)
		return taskInfo, nil
	})

	if err != nil {
		log.Printf("[ERROR] Error waiting for baremetal image task %s: %s", taskID, err)
		return diag.FromErr(fmt.Errorf("error waiting for baremetal image %s to become ready: %w", taskID, err))
	}

	resultTask := waitResult.(*tasks.Task)
	imageID, err := images.ExtractImageIDFromTask(resultTask)
	if err != nil {
		log.Printf("[ERROR] Error extracting image ID from task: %s", err)
		return diag.FromErr(fmt.Errorf("error extracting image ID from task: %w", err))
	}
	log.Printf("[DEBUG] Successfully created baremetal image with ID: %s", imageID)

	d.SetId(imageID)
	return resourceBaremetalImageRead(ctx, d, m)
}

func resourceBaremetalImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Start baremetal image reading")
	config := m.(*Config)
	client, err := CreateClient(config.Provider, d, "gpu/baremetal", "v3")
	if err != nil {
		log.Printf("[ERROR] Error creating client: %s", err)
		return diag.FromErr(err)
	}

	// Add retry mechanism for image retrieval
	var image *images.Image
	err = resource.RetryContext(ctx, 20*time.Second, func() *resource.RetryError {
		getResult := images.Get(client, d.Id())
		if getResult.Err != nil {
			log.Printf("[TRACE] HTTP Response Code: %d", getResult.Err)
			log.Printf("[TRACE] Response Body: %s", getResult.Err)
			if _, ok := getResult.Err.(gcorecloud.ErrDefault404); ok {
				return resource.RetryableError(fmt.Errorf("baremetal image not found, retrying: %w", getResult.Err))
			}
			return resource.NonRetryableError(getResult.Err)
		}
		var err error
		image, err = getResult.Extract()
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		log.Printf("[ERROR] Error getting baremetal image: %s", err)
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Retrieved baremetal image: %+v", image)

	d.Set("name", image.Name)
	if image.URL != "" {
		if err := d.Set("url", image.URL); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set url: %w", err))
		}
	} else {
		// Handle cases where API might return empty URL
		storedURL := d.Get("url").(string)
		if storedURL == "" {
			return diag.Errorf("empty URL in both API response and state")
		}
		log.Printf("[WARN] Using stored URL as API returned empty value")
	}

	if image.Metadata != nil {
		d.Set("metadata", image.Metadata)
		log.Printf("[DEBUG] Setting metadata: %v", image.Metadata)

		if v, ok := image.Metadata["hw_architecture"]; ok {
			d.Set("architecture", v)
		}
		if v, ok := image.Metadata["os_type"]; ok {
			d.Set("os_type", v)
		}
		if v, ok := image.Metadata["os_distro"]; ok {
			d.Set("os_distro", v)
		}
		if v, ok := image.Metadata["os_version"]; ok {
			d.Set("os_version", v)
		}
		if v, ok := image.Metadata["hw_firmware_type"]; ok {
			d.Set("hw_firmware_type", v)
		}
		if v, ok := image.Metadata["ssh_key"]; ok {
			d.Set("ssh_key", v)
		}
	}

	return nil
}

func resourceBaremetalImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Start baremetal image deletion")
	config := m.(*Config)
	client, err := CreateClient(config.Provider, d, "gpu/baremetal", "v3")
	if err != nil {
		log.Printf("[ERROR] Error creating client: %s", err)
		return diag.FromErr(err)
	}

	deleteResult := images.Delete(client, d.Id())
	if deleteResult.Err != nil {
		log.Printf("[ERROR] Error deleting baremetal image %s: %s", d.Id(), deleteResult.Err)
		return diag.FromErr(fmt.Errorf("error deleting baremetal image: %w", deleteResult.Err))
	}
	log.Printf("[DEBUG] Successfully deleted baremetal image: %s", d.Id())

	d.SetId("")
	return nil
}

const BaremetalImageCreatingTimeout = 1200
