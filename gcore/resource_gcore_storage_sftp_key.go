package gcore

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"

	"github.com/G-Core/gcore-storage-sdk-go/swagger/client/key"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	StorageKeySchemaKey  = "key"
	StorageKeySchemaName = "name"
	StorageKeySchemaId   = "key_id"
)

func resourceStorageSFTPKey() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			StorageKeySchemaName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					name := i.(string)
					if !regexp.MustCompile(`^[\w\-]+$`).MatchString(name) || len(name) > 127 {
						return diag.Errorf("key name can't be empty and can have only letters, numbers, dashes and underscores, it also should be less than 128 symbols")
					}
					return nil
				},
				Description: "A name of new storage key resource.",
			},
			StorageKeySchemaKey: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A body of of new storage key resource.",
			},
			StorageKeySchemaId: {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "An id of of new storage key resource.",
			},
		},
		CreateContext: resourceStorageSFTPKeyCreate,
		ReadContext:   resourceStorageSFTPKeyRead,
		DeleteContext: resourceStorageSFTPKeyDelete,
		Description:   "Represent storage key resource. https://storage.gcore.com/ssh-key/list",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceStorageSFTPKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (dErr diag.Diagnostics) {
	id := new(int)
	log.Println("[DEBUG] Start Storage Key Resource creating")
	defer log.Printf("[DEBUG] Finish Storage Key Resource creating (id=%d)\n", *id)
	config := m.(*Config)
	client := config.StorageClient

	opts := make([]func(opt *key.KeyCreateHTTPParams), 0)
	opts = append(opts, func(opt *key.KeyCreateHTTPParams) { opt.Context = ctx })
	name := strings.TrimSpace(d.Get(StorageKeySchemaName).(string))
	if name != "" {
		opts = append(opts, func(opt *key.KeyCreateHTTPParams) { opt.Body.Name = name })
	}
	keyText := strings.TrimSpace(d.Get(StorageKeySchemaKey).(string))
	if keyText != "" {
		opts = append(opts, func(opt *key.KeyCreateHTTPParams) { opt.Body.Key = keyText })
	}

	result, err := client.CreateKey(opts...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("create storage key: %v", err))
	}
	d.SetId(fmt.Sprintf("%d", result.ID))
	return resourceStorageSFTPKeyRead(ctx, d, m)

}

func resourceStorageSFTPKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceId := storageKeyResourceID(d)
	log.Printf("[DEBUG] Start Storage Key Resource reading (id=%s)\n", resourceId)
	defer log.Println("[DEBUG] Finish Storage Key Resource reading")

	config := m.(*Config)
	client := config.StorageClient

	opts := []func(opt *key.KeyListHTTPV2Params){
		func(opt *key.KeyListHTTPV2Params) { opt.Context = ctx },
	}
	if resourceId != "" {
		opts = append(opts, func(opt *key.KeyListHTTPV2Params) { opt.ID = &resourceId })
	}
	name := strings.TrimSpace(d.Get(StorageKeySchemaName).(string))
	if name != "" {
		opts = append(opts, func(opt *key.KeyListHTTPV2Params) { opt.Name = &name })
	}
	if resourceId == "" && name == "" {
		return diag.Errorf("get storage: empty storage id/name")
	}
	result, err := client.KeysList(opts...)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(result) != 1 {
		return diag.Errorf("get storage key: wrong length of search result (%d), want 1", len(result))
	}
	st := result[0]

	d.SetId(fmt.Sprint(st.ID))
	_ = d.Set(StorageKeySchemaId, st.ID)
	_ = d.Set(StorageKeySchemaName, st.Name)

	return nil
}

func resourceStorageSFTPKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceId := storageKeyResourceID(d)
	log.Printf("[DEBUG] Start Storage Key Resource deleting (id=%s)\n", resourceId)
	defer log.Println("[DEBUG] Finish Storage Key Resource deleting")
	if resourceId == "" {
		return diag.Errorf("empty storage id")
	}

	config := m.(*Config)
	client := config.StorageClient

	id, err := strconv.ParseInt(resourceId, 10, 64)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get resource id: %w", err))
	}

	opts := []func(opt *key.KeyDeleteHTTPParams){
		func(opt *key.KeyDeleteHTTPParams) { opt.Context = ctx },
		func(opt *key.KeyDeleteHTTPParams) { opt.ID = id },
	}
	err = client.DeleteKey(opts...)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func storageKeyResourceID(d *schema.ResourceData) string {
	resourceID := d.Id()
	if resourceID == "" {
		id := d.Get(StorageKeySchemaId).(int)
		if id > 0 {
			resourceID = fmt.Sprint(id)
		}
	}
	return resourceID
}
