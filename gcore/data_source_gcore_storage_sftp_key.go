package gcore

import (
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceStorageSFTPKey() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			StorageKeySchemaName: {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					name := i.(string)
					if !regexp.MustCompile(`^[\w\-]+$`).MatchString(name) || len(name) > 127 {
						return diag.Errorf("key name can't be empty and can have only letters, numbers, dashes and underscores, it also should be less than 128 symbols")
					}
					return nil
				},
				Description: "A name of storage key resource.",
			},
			StorageKeySchemaId: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "An id of of storage key resource.",
			},
		},
		ReadContext: resourceStorageSFTPKeyRead,
		Description: "Represent storage key resource. https://storage.gcore.com/ssh-key/list",
	}
}
