package gcore

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFastEdgeBinary() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filename": {
				Description: "WebAssembly binary file to upload.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // new resource on filename change
				// make sure file exists and is readable
				ValidateFunc: func(v interface{}, k string) ([]string, []error) {
					f, err := os.Open(v.(string))
					if err != nil {
						return nil, []error{err}
					}
					f.Close()
					return nil, nil
				},
			},
			"checksum": {
				Description: "Binary checksum.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true, // new resource on file content change
			},
		},
		CreateContext: resourceFastEdgeBinaryUpload,
		ReadContext:   resourceFastEdgeBinaryRead,
		DeleteContext: resourceFastEdgeBinaryDelete,
		Description:   "WebAssembly binary to use in FastEdge applications.",
		// calculate file checksum to detect file content change
		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			checksum, err := fileChecksum(diff.Get("filename").(string))
			if err != nil {
				return nil
			}
			return diff.SetNew("checksum", checksum)
		},
	}
}

func resourceFastEdgeBinaryUpload(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge binary upload")
	config := m.(*Config)
	client := config.FastEdgeClient

	wasmFile, err := os.Open(d.Get("filename").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	defer wasmFile.Close()

	rsp, err := client.StoreBinaryWithBodyWithResponse(ctx, "application/octet-stream", wasmFile)
	if err != nil {
		return diag.FromErr(err)
	}
	if !statusOK(rsp.StatusCode()) {
		return diag.FromErr(errors.New(extractErrorMessage(rsp.Body)))
	}
	d.SetId(strconv.FormatInt(rsp.JSON200.Id, 10))

	log.Printf("[DEBUG] Finish FastEdge binary upload (id=%d)\n", rsp.JSON200.Id)
	return nil
}

func resourceFastEdgeBinaryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge binary deletion")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	rsp, err := client.DelBinaryWithResponse(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() != http.StatusConflict {
			return diag.FromErr(errors.New(extractErrorMessage(rsp.Body)))
		}
		log.Printf("[WARN] FastEdge binary (%d) is referenced, cannot delete but removing from TF state", id)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish FastEdge binary deletion")
	return nil
}

func resourceFastEdgeBinaryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge binary read")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	rsp, err := client.GetBinaryWithResponse(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() == http.StatusNotFound {
			log.Printf("[WARN] FastEdge binary (%d) was not found, removing from TF state", id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(errors.New(extractErrorMessage(rsp.Body)))
	}
	d.Set("checksum", rsp.JSON200.Checksum)

	log.Println("[DEBUG] Finish FastEdge binary read")
	return nil
}

func extractErrorMessage(rspBuf []byte) string {
	var rsp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rspBuf, &rsp); err == nil {
		return rsp.Error
	}
	return string(rspBuf) // "error" is missing - output response as is
}

func statusOK(status int) bool {
	return status == http.StatusOK || status == http.StatusNoContent
}

func fileChecksum(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
