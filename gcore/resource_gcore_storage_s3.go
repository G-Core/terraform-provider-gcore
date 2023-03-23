package gcore

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/G-Core/gcore-storage-sdk-go/swagger/client/storage"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	StorageS3SchemaGenerateAccessKey  = "generated_access_key"
	StorageS3SchemaGenerateSecretKey  = "generated_secret_key"
	StorageSchemaGenerateHTTPEndpoint = "generated_http_endpoint"
	StorageSchemaGenerateS3Endpoint   = "generated_s3_endpoint"
	StorageSchemaGenerateEndpoint     = "generated_endpoint"

	StorageSchemaLocation = "location"
	StorageSchemaName     = "name"
	StorageSchemaId       = "storage_id"
	StorageSchemaClientId = "client_id"
)

func resourceStorageS3() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			StorageSchemaId: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "An id of new storage resource.",
			},
			StorageSchemaClientId: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "An client id of new storage resource.",
			},
			StorageSchemaName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					storageName := i.(string)
					if !regexp.MustCompile(`^[\w\-]+$`).MatchString(storageName) || len(storageName) > 255 {
						return diag.Errorf("storage name can't be empty and can have only letters, numbers, dashes and underscores, it also should be less than 256 symbols")
					}
					return nil
				},
				Description: "A name of new storage resource.",
			},
			StorageSchemaLocation: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					val := v.(string)
					allowed := []string{"s-ed1", "s-darz1", "s-ws1", "s-dt2", "s-drc2"}
					for _, el := range allowed {
						if el == val {
							return nil
						}
					}
					return diag.Errorf(`must be one of %+v`, allowed)
				},
				Description: "A location of new storage resource. One of (s-ed1, s-darz1, s-ws1, s-dt2, s-drc2)",
			},
			StorageS3SchemaGenerateAccessKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A s3 access key for new storage resource.",
			},
			StorageS3SchemaGenerateSecretKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A s3 secret key for new storage resource.",
			},
			StorageSchemaGenerateHTTPEndpoint: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A http s3 entry point for new storage resource.",
			},
			StorageSchemaGenerateS3Endpoint: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A s3 endpoint for new storage resource.",
			},
			StorageSchemaGenerateEndpoint: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A s3 entry point for new storage resource.",
			},
		},
		CreateContext: resourceStorageS3Create,
		ReadContext:   resourceStorageS3Read,
		DeleteContext: resourceStorageS3Delete,
		Description:   "Represent s3 storage resource. https://storage.gcore.com/storage/list",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceStorageS3Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := new(int)
	log.Println("[DEBUG] Start S3 Storage Resource creating")
	defer log.Printf("[DEBUG] Finish S3 Storage Resource creating (id=%d)\n", *id)
	config := m.(*Config)
	client := config.StorageClient

	opts := []func(opt *storage.StorageCreateHTTPParams){
		func(opt *storage.StorageCreateHTTPParams) { opt.Context = ctx },
		func(opt *storage.StorageCreateHTTPParams) { opt.Body.Type = "s3" },
	}
	location := strings.TrimSpace(d.Get(StorageSchemaLocation).(string))
	if location != "" {
		opts = append(opts, func(opt *storage.StorageCreateHTTPParams) { opt.Body.Location = location })
	}
	name := strings.TrimSpace(d.Get(StorageSchemaName).(string))
	if name != "" {
		opts = append(opts, func(opt *storage.StorageCreateHTTPParams) { opt.Body.Name = name })
	}

	result, err := client.CreateStorage(opts...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("create storage: %v", err))
	}
	d.SetId(fmt.Sprintf("%d", result.ID))
	*id = int(result.ID)
	if result.Credentials.S3.AccessKey != "" {
		_ = d.Set(StorageS3SchemaGenerateAccessKey, result.Credentials.S3.AccessKey)
	}
	if result.Credentials.S3.SecretKey != "" {
		_ = d.Set(StorageS3SchemaGenerateSecretKey, result.Credentials.S3.SecretKey)
	}

	return resourceStorageS3Read(ctx, d, m)
}

func resourceStorageS3Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceId := storageResourceID(d)
	log.Printf("[DEBUG] Start S3 Storage Resource reading (id=%s)\n", resourceId)
	defer log.Println("[DEBUG] Finish S3 Storage Resource reading")

	config := m.(*Config)
	client := config.StorageClient

	opts := []func(opt *storage.StorageListHTTPV2Params){
		func(opt *storage.StorageListHTTPV2Params) { opt.Context = ctx },
		func(opt *storage.StorageListHTTPV2Params) { opt.ShowDeleted = new(bool) },
	}
	if resourceId != "" {
		opts = append(opts, func(opt *storage.StorageListHTTPV2Params) { opt.ID = &resourceId })
	}
	name := d.Get(StorageSchemaName).(string)
	if name != "" {
		opts = append(opts, func(opt *storage.StorageListHTTPV2Params) { opt.Name = &name })
	}
	if resourceId == "" && name == "" {
		return diag.Errorf("get storage: empty storage id/name")
	}

	result, err := client.StoragesList(opts...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("storages list: %w", err))
	}

	if (len(result) == 0) || (name == "" && len(result) != 1) {
		return diag.Errorf("get storage: wrong length of search result (%d), want 1", len(result))
	}
	st := result[0]

	d.SetId(fmt.Sprint(st.ID))
	nameParts := strings.Split(st.Name, "-")
	if len(nameParts) > 1 {
		clientID, _ := strconv.ParseInt(nameParts[0], 10, 64)
		_ = d.Set(StorageSchemaClientId, int(clientID))
		_ = d.Set(StorageSchemaName, strings.Join(nameParts[1:], "-"))
	} else {
		_ = d.Set(StorageSchemaName, st.Name)
	}
	_ = d.Set(StorageSchemaId, st.ID)
	_ = d.Set(StorageSchemaLocation, st.Location)
	_ = d.Set(StorageSchemaGenerateEndpoint, fmt.Sprintf("%s.cloud.gcore.lu/%s", st.Location, st.Name))
	_ = d.Set(StorageSchemaGenerateHTTPEndpoint, fmt.Sprintf("https://%s.cloud.gcore.lu/{bucket_name}", st.Location))
	_ = d.Set(StorageSchemaGenerateS3Endpoint, fmt.Sprintf("https://%s.cloud.gcore.lu", st.Location))

	return nil
}

func resourceStorageS3Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceId := storageResourceID(d)
	log.Printf("[DEBUG] Start S3 Storage Resource deleting (id=%s)\n", resourceId)
	defer log.Println("[DEBUG] Finish S3 Storage Resource deleting")
	if resourceId == "" {
		return diag.Errorf("empty storage id")
	}

	config := m.(*Config)
	client := config.StorageClient

	id, err := strconv.ParseInt(resourceId, 10, 64)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get resource id: %w", err))
	}

	opts := []func(opt *storage.StorageDeleteHTTPParams){
		func(opt *storage.StorageDeleteHTTPParams) { opt.Context = ctx },
		func(opt *storage.StorageDeleteHTTPParams) { opt.ID = id },
	}
	err = client.DeleteStorage(opts...)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func storageResourceID(d *schema.ResourceData) string {
	resourceID := d.Id()
	if resourceID == "" {
		id := d.Get(StorageSchemaId).(int)
		if id > 0 {
			resourceID = fmt.Sprint(id)
		}
	}
	return resourceID
}
