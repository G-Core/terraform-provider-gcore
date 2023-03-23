package gcore

import (
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceStorageSFTP() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			StorageSchemaId: {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					StorageSchemaId,
					StorageSchemaName,
				},
				Description: "An id of storage resource.",
			},
			StorageSchemaClientId: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "An client id of storage resource.",
			},
			StorageSchemaName: {
				Type:     schema.TypeString,
				Optional: true,
				AtLeastOneOf: []string{
					StorageSchemaId,
					StorageSchemaName,
				},
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					storageName := i.(string)
					if !regexp.MustCompile(`^[a-z0-9\-]+$`).MatchString(storageName) || len(storageName) > 26 {
						return diag.Errorf("sftp storage name can't be empty and can have only lowercase letters, numbers and dashes; it also must be less than 27 characters length")
					}
					return nil
				},
				Description: "A name of storage resource.",
			},
			StorageSFTPSchemaServerAlias: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An alias of storage resource.",
			},
			StorageSFTPSchemaExpires: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A expires date of storage resource.",
			},
			StorageSchemaLocation: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A location of new storage resource. One of (ams, sin, fra, mia)",
			},
			StorageSFTPSchemaKeyId: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed:    true,
				Description: "An ssh keys IDs to link with new sftp storage resource only. https://storage.gcore.com/ssh-key/list",
			},
			StorageSchemaGenerateHTTPEndpoint: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A http sftp entry point for new storage resource.",
			},
			StorageSchemaGenerateSFTPEndpoint: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A ssh sftp entry point for new storage resource.",
			},
		},
		ReadContext: resourceStorageSFTPRead,
		Description: "Represent sftp storage resource. https://storage.gcore.com/storage/list",
	}
}
