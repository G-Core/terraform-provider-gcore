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

const GPUImageCreatingTimeout = 1200

type GPUImageType string

const (
	GPUImageTypeBaremetal GPUImageType = "baremetal"
	GPUImageTypeVirtual   GPUImageType = "virtual"
)

func resourceGPUImageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Description: "Project ID",
		},
		"region_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Description: "Region ID",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Image name",
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z 0-9._\-]{1,61}[a-zA-Z0-9]$`), "Invalid image name format"),
		},
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Image URL",
		},
		"ssh_key": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "allow",
			ForceNew:     true,
			Description:  "Permission to use a ssh key in instances",
			ValidateFunc: validation.StringInSlice([]string{"allow", "deny", "required"}, false),
		},
		"cow_format": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: "When True, image cannot be deleted unless all volumes, created from it, are deleted",
		},
		"architecture": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "x86_64",
			ForceNew:     true,
			Description:  "Image architecture type: aarch64, x86_64",
			ValidateFunc: validation.StringInSlice([]string{"aarch64", "x86_64"}, false),
		},
		"os_type": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The operating system installed on the image",
		},
		"os_distro": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "OS Distribution, i.e. Debian, CentOS, Ubuntu, CoreOS etc",
		},
		"os_version": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "OS version, i.e. 19.04 (for Ubuntu) or 9.4 for Debian",
		},
		"hw_firmware_type": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Specifies the type of firmware with which to boot the guest",
		},
		"metadata": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Create one or more metadata items for a cluster",
		},
		"project_name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Project name",
		},
		"region_name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Region name",
		},
	}
}

func resourceGPUImage(imageType GPUImageType) *schema.Resource {
	return &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			return resourceGPUImageCreate(ctx, d, m, imageType)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			return resourceGPUImageRead(ctx, d, m, imageType)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			return resourceGPUImageDelete(ctx, d, m, imageType)
		},
		Description: fmt.Sprintf("Manages a %s GPU image", imageType),
		Schema:      resourceGPUImageSchema(),
	}
}

func getGPUServicePath(imageType GPUImageType) string {
	return fmt.Sprintf("gpu/%s", imageType)
}

func resourceGPUImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}, imageType GPUImageType) diag.Diagnostics {
	log.Printf("[DEBUG] Start %s GPU image creation", imageType)
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, getGPUServicePath(imageType), "v3")
	if err != nil {
		return diag.FromErr(err)
	}

	metadata := make(map[string]interface{})
	if v, ok := d.GetOk("metadata"); ok {
		metadata = v.(map[string]interface{})
	}

	// Add additional metadata from specific fields
	for _, field := range []string{"os_type", "os_distro", "os_version", "hw_firmware_type", "ssh_key"} {
		if v, ok := d.GetOk(field); ok {
			metadata[field] = v.(string)
		}
	}

	opts := images.ImageOpts{
		Name:     d.Get("name").(string),
		URL:      d.Get("url").(string),
		Metadata: metadata,
	}

	if v, ok := d.GetOk("cow_format"); ok {
		cowFormat := v.(bool)
		opts.CowFormat = &cowFormat
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

	result := images.UploadImage(client, opts)
	if result.Err != nil {
		return diag.FromErr(result.Err)
	}

	taskID, err := result.Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	firstTaskID := taskID.Tasks[0]

	taskClient, err := CreateClient(provider, d, "tasks", "v1")
	if err != nil {
		return diag.FromErr(err)
	}

	waitResult, err := tasks.WaitTaskAndReturnResult(taskClient, firstTaskID, true, GPUImageCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(taskClient, string(task)).Extract()
		if err != nil {
			return nil, err
		}
		return taskInfo, nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	resultTask := waitResult.(*tasks.Task)
	imageID, err := images.ExtractImageIDFromTask(resultTask)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(imageID)
	return resourceGPUImageRead(ctx, d, m, imageType)
}

func resourceGPUImageRead(ctx context.Context, d *schema.ResourceData, m interface{}, imageType GPUImageType) diag.Diagnostics {
	log.Printf("[DEBUG] Start %s GPU image reading", imageType)
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, getGPUServicePath(imageType), "v3")
	if err != nil {
		return diag.FromErr(err)
	}

	var image *images.Image
	err = resource.RetryContext(ctx, 20*time.Second, func() *resource.RetryError {
		getResult := images.Get(client, d.Id())
		if getResult.Err != nil {
			if _, ok := getResult.Err.(gcorecloud.ErrDefault404); ok {
				return resource.RetryableError(getResult.Err)
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
		return diag.FromErr(err)
	}

	d.Set("name", image.Name)
	if image.URL != "" {
		d.Set("url", image.URL)
	}

	if image.Metadata != nil {
		d.Set("metadata", image.Metadata)

		for _, field := range []string{"hw_architecture", "os_type", "os_distro", "os_version", "hw_firmware_type", "ssh_key"} {
			if v, ok := image.Metadata[field]; ok {
				d.Set(field, v)
			}
		}
	}

	return nil
}

func resourceGPUImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}, imageType GPUImageType) diag.Diagnostics {
	log.Printf("[DEBUG] Start %s GPU image deletion", imageType)
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, getGPUServicePath(imageType), "v3")
	if err != nil {
		return diag.FromErr(err)
	}

	result := images.Delete(client, d.Id())
	if result.Err != nil {
		return diag.FromErr(result.Err)
	}

	d.SetId("")
	return nil
}
