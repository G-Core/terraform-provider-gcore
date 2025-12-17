package gcore

import (
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceStorageS3() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			StorageSchemaId: {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					StorageSchemaId,
					StorageSchemaName,
				},
				Description: "An id of new storage resource.",
			},
			StorageSchemaClientId: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "An client id of new storage resource.",
			},
			StorageSchemaName: {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					storageName := i.(string)
					if !regexp.MustCompile(`^[\w\-]+$`).MatchString(storageName) || len(storageName) > 255 {
						return diag.Errorf("storage name can't be empty and can have only letters, numbers, dashes and underscores, it also should be less than 256 symbols")
					}
					return nil
				},
				AtLeastOneOf: []string{
					StorageSchemaId,
					StorageSchemaName,
				},
				Description: "A name of new storage resource.",
			},
			StorageSchemaLocation: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A location of new storage resource.",
			},
			StorageSchemaGenerateHTTPEndpoint: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A http s3 entry point for new storage resource.",
			},
			StorageSchemaGenerateS3Endpoint: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A s3 endpoint for new storage resource.",
			},
			StorageSchemaGenerateEndpoint: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A s3 entry point for new storage resource.",
			},
		},
		ReadContext: resourceStorageS3Read,
		Description: "Represent s3 storage resource. https://storage.gcore.com/storage/list",
	}
}
