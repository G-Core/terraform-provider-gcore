package gcore

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
				ValidateFunc: func(v any, k string) ([]string, []error) {
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
				ForceNew:    true, // new resource on file content change
			},
		},
		CreateContext: resourceFastEdgeBinaryUpload,
		ReadContext:   resourceFastEdgeBinaryRead,
		DeleteContext: resourceFastEdgeBinaryDelete,
		Description:   "WebAssembly binary to use in FastEdge applications.",
		// calculate file checksum to detect file content change
		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, meta any) error {
			checksum, err := fileChecksum(diff.Get("filename").(string))
			if err != nil {
				return nil
			}
			return diff.SetNew("checksum", checksum)
		},
	}
}

func resourceFastEdgeBinaryUpload(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge binary upload")
	config := m.(*Config)
	client := config.FastEdgeClient

	filename := d.Get("filename").(string)
	wasmFile, err := os.Open(filename)
	if err != nil {
		return diag.Errorf("opening file %s: %v", filename, err)
	}
	defer wasmFile.Close()

	rsp, err := client.StoreBinaryWithBodyWithResponse(ctx, "application/octet-stream", wasmFile)
	if err != nil {
		return diag.Errorf("calling StoreBinary API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		return diag.Errorf("calling StoreBinary API: %s", extractErrorMessage(rsp.Body))
	}

	// make sure binary was not damaged in transit
	expectedChecksum := d.Get("checksum").(string)
	if *rsp.JSON200.Checksum != expectedChecksum {
		// binary damaged in transit
		return diag.Errorf("uploaded binary checksum (%s) does not match expected (%s), please retry", *rsp.JSON200.Checksum, expectedChecksum)
	}

	d.SetId(strconv.FormatInt(rsp.JSON200.Id, 10))

	log.Printf("[DEBUG] Finish FastEdge binary upload (id=%d)\n", rsp.JSON200.Id)
	return nil
}

func resourceFastEdgeBinaryDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Println("[DEBUG] Start FastEdge binary deletion")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("converting id to number: %v", err)
	}

	rsp, err := client.DelBinaryWithResponse(ctx, id)
	if err != nil {
		return diag.Errorf("calling DelBinary API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() == http.StatusConflict {
			diags = diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Wasm binary (%d) is referenced so cannot be deleted", id),
				},
			}
		} else {
			return diag.Errorf("calling DelBinary API: %s", extractErrorMessage(rsp.Body))
		}
	} else {
		d.SetId("")
	}

	log.Println("[DEBUG] Finish FastEdge binary deletion")
	return diags
}

func resourceFastEdgeBinaryRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	log.Println("[DEBUG] Start FastEdge binary read")
	config := m.(*Config)
	client := config.FastEdgeClient

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("converting id to number: %v", err)
	}

	rsp, err := client.GetBinaryWithResponse(ctx, id)
	if err != nil {
		return diag.Errorf("calling GetBinary API: %v", err)
	}
	if !statusOK(rsp.StatusCode()) {
		if rsp.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("[FastEdge binary (%d) was not found, removed from TF state", id),
				},
			}
		}
		return diag.Errorf("calling GetBinary API: %s", extractErrorMessage(rsp.Body))
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
