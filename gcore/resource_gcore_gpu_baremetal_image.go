package gcore

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/gpu/v3/images"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBaremetalImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBaremetalImageCreate,
		ReadContext:   resourceBaremetalImageRead,
		DeleteContext: resourceBaremetalImageDelete,
		Description:   "Manages a baremetal custom image",

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Project ID",
			},
			"region_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Region ID",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Image name",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z 0-9._\-]{1,61}[a-zA-Z0-9]$`), "Invalid image name format"),
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Image URL",
			},
			"ssh_key": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "allow",
				ForceNew:     true,
				Description:  "Permission to use a ssh key in instances",
				ValidateFunc: validation.StringInSlice([]string{"allow", "deny", "required"}, false),
			},
			"cow_format": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "When True, image cannot be deleted unless all volumes, created from it, are deleted",
			},
			"architecture": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "x86_64",
				ForceNew:     true,
				Description:  "Image architecture type: aarch64, x86_64",
				ValidateFunc: validation.StringInSlice([]string{"aarch64", "x86_64"}, false),
			},
			"os_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The operating system installed on the image",
			},
			"os_distro": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "OS Distribution, i.e. Debian, CentOS, Ubuntu, CoreOS etc",
			},
			"os_version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "OS version, i.e. 19.04 (for Ubuntu) or 9.4 for Debian",
			},
			"hw_firmware_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the type of firmware with which to boot the guest",
			},
			"metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Create one or more metadata items for a cluster",
			},
			"project_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Project name",
			},
			"region_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Region name",
			},
		},
	}
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

	task, err := images.UploadImage(client, opts)
	if err != nil {
		log.Printf("[ERROR] Error uploading baremetal image: %s", err)
		return diag.FromErr(fmt.Errorf("error uploading baremetal image: %w", err))
	}
	log.Printf("[DEBUG] Upload task started with ID: %s", task.Tasks[0])

	taskID := tasks.TaskID(task.Tasks[0])

	taskClient, err := CreateClient(config.Provider, d, "tasks", "v1")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating task client: %w", err))
	}

	result, err := tasks.WaitTaskAndReturnResult(taskClient, taskID, true, BaremetalImageCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
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

	resultTask := result.(*tasks.Task)
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
		image, err = images.Get(client, d.Id())
		if err != nil {
			log.Printf("[TRACE] HTTP Response Code: %d", err)
			log.Printf("[TRACE] Response Body: %s", err)
			if _, ok := err.(gcorecloud.ErrDefault404); ok {
				return resource.RetryableError(fmt.Errorf("baremetal image not found, retrying: %w", err))
			}
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

	_, err = images.Delete(client, d.Id())
	if err != nil {
		log.Printf("[ERROR] Error deleting baremetal image %s: %s", d.Id(), err)
		return diag.FromErr(fmt.Errorf("error deleting baremetal image: %w", err))
	}
	log.Printf("[DEBUG] Successfully deleted baremetal image: %s", d.Id())

	d.SetId("")
	return nil
}

const BaremetalImageCreatingTimeout = 1200
