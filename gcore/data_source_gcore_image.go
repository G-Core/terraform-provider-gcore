package gcore

import (
	"context"
	"fmt"
	"log"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	imagesPoint        = "images"
	bmImagesPoint      = "bmimages"
	downloadImagePoint = "downloadimage"
	ImageUploadTimeout = 1200
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageRead,
		Description: "Represent image data",
		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
				DiffSuppressFunc: suppressDiffProjectID,
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
				DiffSuppressFunc: suppressDiffRegionID,
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "use 'os-version', for example 'ubuntu-20.04'",
				Optional:     true,
				ExactlyOneOf: []string{"name", "image_id"},
			},
			"image_id": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "use 'image_id' if you know it, for example 'f4b1b1b1-1b1b-1b1b-1b1b-1b1b1b1b1b1b'",
				Optional:     true,
				ExactlyOneOf: []string{"name", "image_id"},
			},
			"is_baremetal": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "set to true if you need to get baremetal image",
				Optional:    true,
			},
			"min_disk": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"min_ram": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"os_distro": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata_k": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata_kv": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"metadata_read_only": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Image reading")
	name := d.Get("name").(string)
	imageID := d.Get("image_id").(string)

	config := m.(*Config)
	provider := config.Provider

	point := imagesPoint
	if isBm, _ := d.Get("is_baremetal").(bool); isBm {
		point = bmImagesPoint
	}
	client, err := CreateClient(provider, d, point, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	var image *images.Image
	if imageID != "" {
		image, err = images.Get(client, imageID).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		image, err = findImageByNameAndMetadata(client, d, name)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(image.ID)
	d.Set("project_id", d.Get("project_id").(int))
	d.Set("region_id", d.Get("region_id").(int))
	d.Set("min_disk", image.MinDisk)
	d.Set("min_ram", image.MinRAM)
	d.Set("os_distro", image.OsDistro)
	d.Set("os_version", image.OsVersion)
	d.Set("description", image.Description)

	metadataReadOnly := make([]map[string]interface{}, 0, len(image.Metadata))
	if len(image.Metadata) > 0 {
		for _, metadataItem := range image.Metadata {
			metadataReadOnly = append(metadataReadOnly, map[string]interface{}{
				"key":       metadataItem.Key,
				"value":     metadataItem.Value,
				"read_only": metadataItem.ReadOnly,
			})
		}
	}

	if err := d.Set("metadata_read_only", metadataReadOnly); err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] Finish Image reading")
	return nil
}

func imagesNames(images []images.Image) []string {
	result := make([]string, 0, len(images))
	for _, img := range images {
		result = append(result, fmt.Sprintf("%s (%s)", img.Name, img.ID))
	}
	return result

}

func findImageByNameAndMetadata(client *gcorecloud.ServiceClient, d *schema.ResourceData, name string) (*images.Image, error) {
	listOpts := &images.ListOpts{}
	if metadataK, ok := d.GetOk("metadata_k"); ok {
		listOpts.MetadataK = metadataK.(string)
	}

	if metadataRaw, ok := d.GetOk("metadata_kv"); ok {
		typedMetadataKV := make(map[string]string, len(metadataRaw.(map[string]interface{})))
		for k, v := range metadataRaw.(map[string]interface{}) {
			typedMetadataKV[k] = v.(string)
		}
		listOpts.MetadataKV = typedMetadataKV
	}

	allImages, err := images.ListAll(client, *listOpts)
	if err != nil {
		return nil, err
	}

	collectedImages := make([]images.Image, 0)
	for _, img := range allImages {
		if strings.HasPrefix(strings.ToLower(img.Name), strings.ToLower(name)) {
			collectedImages = append(collectedImages, img)
		}
	}

	if len(collectedImages) == 0 {
		return nil, fmt.Errorf("image with name %s not found", name)
	}

	if len(collectedImages) > 1 {
		return nil, fmt.Errorf(
			"found more than one image with name %s - %s, pls choose image by iamge_id",
			name, strings.Join(imagesNames(collectedImages), ", "),
		)
	}

	return &collectedImages[0], nil
}
