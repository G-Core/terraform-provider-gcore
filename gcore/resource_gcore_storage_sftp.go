package gcore

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	gstorage "github.com/G-Core/gcore-storage-sdk-go"
	"github.com/G-Core/gcore-storage-sdk-go/swagger/client/key"

	"github.com/G-Core/gcore-storage-sdk-go/swagger/client/storage"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	StorageSFTPSchemaGenerateSftpPassword = "generated_password"
	StorageSchemaGenerateSFTPEndpoint     = "generated_sftp_endpoint"
	StorageSFTPSchemaSftpPassword         = "password"
	StorageSFTPSchemaKeyId                = "ssh_key_id"
	StorageSFTPSchemaExpires              = "http_expires_header_value"
	StorageSFTPSchemaServerAlias          = "http_servername_alias"

	StorageSFTPSchemaUpdateAfterCreate = "update_after_create"
)

func resourceStorageSFTP() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			StorageSchemaId: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "An id of new storage resource.",
			},
			StorageSFTPSchemaUpdateAfterCreate: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "A temporary flag. An internal cheat, to skip update ssh keys. Skip it.",
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
					if !regexp.MustCompile(`^[a-z0-9\-]+$`).MatchString(storageName) || len(storageName) > 26 {
						return diag.Errorf("sftp storage name can't be empty and can have only lowercase letters, numbers and dashes; it also must be less than 27 characters length")
					}
					return nil
				},
				Description: "A name of new storage resource.",
			},
			StorageSFTPSchemaServerAlias: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An alias of storage resource.",
			},
			StorageSFTPSchemaExpires: {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					v := i.(string)
					if !regexp.MustCompile(`^(\d+\s(years))?\s?(\d+\s(months))?\s?(\d+\s(weeks))?\s?(\d+\s(days))?\s?(\d+\s(hours))?\s?(\d+\s(minutes))?\s?(\d+\s(seconds))?$`).MatchString(v) || len(v) > 255 {
						return diag.Errorf("storage expires must be matched by " +
							"^([0-9]+\\s(years|months|weeks|days|hours|minutes|seconds)[\\s]?)+$ regexp, " +
							"it also should be less than 256 symbols")
					}
					return nil
				},
				Description: "A expires date of storage resource.",
			},
			StorageSchemaLocation: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A location of new storage resource.",
			},
			StorageSFTPSchemaSftpPassword: {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					v := i.(string)
					if len(v) > 63 || len(v) < 8 {
						return diag.Errorf("storage password should be more than 8 and less than 64 symbols")
					}
					return nil
				},
				Description: "A sftp password for new storage resource.",
			},
			StorageSFTPSchemaGenerateSftpPassword: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "An auto generated sftp password for new storage resource.",
			},
			StorageSFTPSchemaKeyId: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "An ssh keys IDs to link with new sftp storage resource only. https://storage.gcore.com/ssh-key/list",
			},
			StorageSchemaGenerateHTTPEndpoint: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A http sftp entry point for new storage resource.",
			},
			StorageSchemaGenerateSFTPEndpoint: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A ssh sftp entry point for new storage resource.",
			},
		},
		CreateContext: resourceStorageSFTPCreate,
		ReadContext:   resourceStorageSFTPRead,
		UpdateContext: resourceStorageSFTPUpdate,
		DeleteContext: resourceStorageSFTPDelete,
		Description:   "Represent sftp storage resource. https://storage.gcore.com/storage/list",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceStorageValidateKeys(ctx context.Context, sdk *gstorage.SDK, d *schema.ResourceData) error {
	keyIds := d.Get(StorageSFTPSchemaKeyId).([]interface{})
	if len(keyIds) == 0 {
		return nil
	}
	for _, v := range keyIds {
		id, ok := v.(int)
		keyID := fmt.Sprint(id)
		if !ok {
			return fmt.Errorf("key %v is not int", v)
		}
		opts := []func(opt *key.KeyListHTTPV2Params){
			func(opt *key.KeyListHTTPV2Params) { opt.Context = ctx },
			func(opt *key.KeyListHTTPV2Params) { opt.ID = &keyID },
		}
		result, err := sdk.KeysList(opts...)
		if err != nil {
			return fmt.Errorf("problems with key %v: %w", v, err)
		}
		if len(result) == 0 {
			return fmt.Errorf("key %v is not found", v)
		}
	}
	return nil
}

func resourceStorageLinkKeys(ctx context.Context, sdk *gstorage.SDK, d *schema.ResourceData, storageID int64) error {
	keyIds := d.Get(StorageSFTPSchemaKeyId).([]interface{})
	if len(keyIds) == 0 {
		return nil
	}
	for _, v := range keyIds {
		keyId, ok := v.(int)
		if !ok {
			log.Printf("[ERROR] ssh key id should be int: %+v is %T\n", v, v)
			continue
		}
		keyOpts := []func(opt *storage.KeyLinkHTTPParams){
			func(opt *storage.KeyLinkHTTPParams) { opt.Context = ctx },
			func(opt *storage.KeyLinkHTTPParams) { opt.ID = storageID; opt.KeyID = int64(keyId) },
		}
		err := sdk.LinkKeyToStorage(keyOpts...)
		if err != nil {
			return fmt.Errorf("link key #%d to storage: %w", keyId, err)
		}
	}
	return nil
}

func resourceStorageRelinkKeys(ctx context.Context, sdk *gstorage.SDK, d *schema.ResourceData, storageID int64) error {
	if !d.HasChange(StorageSFTPSchemaKeyId) {
		return nil
	}
	oldKeys, newKeys := d.GetChange(StorageSFTPSchemaKeyId)
	oldIDs, newIDs := oldKeys.([]interface{}), newKeys.([]interface{})
	decision := map[interface{}]int{}
	for _, v := range oldIDs {
		decision[v]--
	}
	for _, v := range newIDs {
		decision[v]++
	}
	unlink, link := make([]int, 0), make([]int, 0)
	for id, v := range decision {
		if v > 0 {
			link = append(link, id.(int))
		}
		if v < 0 {
			unlink = append(unlink, id.(int))
		}
	}
	for _, keyId := range link {
		keyOpts := []func(opt *storage.KeyLinkHTTPParams){
			func(opt *storage.KeyLinkHTTPParams) { opt.Context = ctx },
			func(opt *storage.KeyLinkHTTPParams) { opt.ID = storageID; opt.KeyID = int64(keyId) },
		}
		err := sdk.LinkKeyToStorage(keyOpts...)
		if err != nil {
			return fmt.Errorf("link key #%d to storage: %w", keyId, err)
		}
	}
	for _, keyId := range unlink {
		keyOpts := []func(opt *storage.KeyUnlinkHTTPParams){
			func(opt *storage.KeyUnlinkHTTPParams) { opt.Context = ctx },
			func(opt *storage.KeyUnlinkHTTPParams) { opt.ID = storageID; opt.KeyID = int64(keyId) },
		}
		err := sdk.UnlinkKeyFromStorage(keyOpts...)
		if err != nil {
			return fmt.Errorf("unlink key #%d to storage: %w", keyId, err)
		}
	}
	return nil
}

func resourceStorageSFTPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := new(int)
	log.Println("[DEBUG] Start SFTP Storage Resource creating")
	defer log.Printf("[DEBUG] Finish SFTP Storage Resource creating (id=%d)\n", *id)
	config := m.(*Config)
	client := config.StorageClient

	err := resourceStorageValidateKeys(ctx, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := []func(opt *storage.StorageCreateHTTPParams){
		func(opt *storage.StorageCreateHTTPParams) { opt.Context = ctx },
		func(opt *storage.StorageCreateHTTPParams) { opt.Body.Type = "sftp" },
	}
	if d.Get(StorageSFTPSchemaGenerateSftpPassword).(bool) {
		opts = append(opts, func(opt *storage.StorageCreateHTTPParams) { opt.Body.GenerateSftpPassword = true })
	}
	location := strings.TrimSpace(d.Get(StorageSchemaLocation).(string))
	if location != "" {
		opts = append(opts, func(opt *storage.StorageCreateHTTPParams) { opt.Body.Location = location })
	}
	name := strings.TrimSpace(d.Get(StorageSchemaName).(string))
	if name != "" {
		opts = append(opts, func(opt *storage.StorageCreateHTTPParams) { opt.Body.Name = name })
	}
	sftpPassword := strings.TrimSpace(d.Get(StorageSFTPSchemaSftpPassword).(string))
	if sftpPassword != "" {
		opts = append(opts, func(opt *storage.StorageCreateHTTPParams) { opt.Body.SftpPassword = sftpPassword })
	}

	result, err := client.CreateStorage(opts...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("create storage: %v", err))
	}
	d.SetId(fmt.Sprintf("%d", result.ID))
	*id = int(result.ID)
	if result.Credentials.SftpPassword != "" {
		_ = d.Set(StorageSFTPSchemaSftpPassword, result.Credentials.SftpPassword)
	}

	err = resourceStorageLinkKeys(ctx, client, d, result.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	if strings.TrimSpace(d.Get(StorageSFTPSchemaExpires).(string)) != "" ||
		strings.TrimSpace(d.Get(StorageSFTPSchemaServerAlias).(string)) != "" {
		_ = d.Set(StorageSFTPSchemaUpdateAfterCreate, true)
		return resourceStorageSFTPUpdate(ctx, d, m)
	}

	return resourceStorageSFTPRead(ctx, d, m)
}

func resourceStorageSFTPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceId := storageResourceID(d)
	log.Printf("[DEBUG] Start SFTP Storage Resource reading (id=%s)\n", resourceId)
	defer log.Println("[DEBUG] Finish SFTP Storage Resource reading")

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
	if len(st.Credentials.Keys) > 0 {
		keys := make([]int, 0, len(st.Credentials.Keys))
		for _, k := range st.Credentials.Keys {
			if k == nil {
				continue
			}
			keys = append(keys, int(k.ID))
		}
		_ = d.Set(StorageKeySchemaId, keys)
	}
	_ = d.Set(StorageSFTPSchemaServerAlias, st.ServerAlias)
	_ = d.Set(StorageSFTPSchemaExpires, st.Expires)
	_ = d.Set(StorageSchemaId, st.ID)
	_ = d.Set(StorageSchemaLocation, st.Location)
	_ = d.Set(StorageSchemaGenerateHTTPEndpoint,
		fmt.Sprintf("http://%s.%s.origin.gcdn.co", st.Name, st.Location))
	_ = d.Set(StorageSchemaGenerateSFTPEndpoint,
		fmt.Sprintf("ssh://%s@%s.origin.gcdn.co:2200", st.Name, st.Location))

	return nil
}

func resourceStorageSFTPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceId := storageResourceID(d)
	log.Printf("[DEBUG] Start SFTP Storage Resource updating (id=%s)\n", resourceId)
	defer log.Println("[DEBUG] Finish SFTP Storage Resource updating")
	if resourceId == "" {
		return diag.Errorf("empty storage id")
	}

	config := m.(*Config)
	client := config.StorageClient

	id, err := strconv.ParseInt(resourceId, 10, 64)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get resource id: %w", err))
	}

	opts := make([]func(opt *storage.StorageUpdateHTTPParams), 0)
	opts = append(opts, func(opt *storage.StorageUpdateHTTPParams) { opt.Context = ctx })
	opts = append(opts, func(opt *storage.StorageUpdateHTTPParams) { opt.ID = id })
	expires := strings.TrimSpace(d.Get(StorageSFTPSchemaExpires).(string))
	if expires != "" {
		opts = append(opts, func(opt *storage.StorageUpdateHTTPParams) { opt.Body.Expires = expires })
	}
	alias := strings.TrimSpace(d.Get(StorageSFTPSchemaServerAlias).(string))
	if alias != "" {
		opts = append(opts, func(opt *storage.StorageUpdateHTTPParams) { opt.Body.ServerAlias = alias })
	}
	_, err = client.ModifyStorage(opts...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("update storage: %w", err))
	}
	if d.Get(StorageSFTPSchemaUpdateAfterCreate).(bool) {
		_ = d.Set(StorageSFTPSchemaUpdateAfterCreate, false)
		return nil
	}
	if d.HasChange(StorageSFTPSchemaSftpPassword) {
		pass := d.Get(StorageSFTPSchemaSftpPassword).(string)
		deletePass := false
		if pass == "" {
			deletePass = true
		}
		passOpts := []func(*storage.StorageUpdateCredentialsHTTPParams){
			func(params *storage.StorageUpdateCredentialsHTTPParams) {
				params.ID = id
				params.Context = ctx
				params.Body.SftpPassword = pass
				params.Body.DeleteSftpPassword = deletePass
			},
		}
		_, err = client.UpdateStorageCredentials(passOpts...)
		if err != nil {
			return diag.FromErr(fmt.Errorf("update creds: %w", err))
		}
	}
	err = resourceStorageRelinkKeys(ctx, client, d, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("update keys: %w", err))
	}

	return resourceStorageSFTPRead(ctx, d, m)
}

func resourceStorageSFTPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceId := storageResourceID(d)
	log.Printf("[DEBUG] Start SFTP Storage Resource deleting (id=%s)\n", resourceId)
	defer log.Println("[DEBUG] Finish SFTP Storage Resource deleting")
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
