package gcore

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBaremetalImage() *schema.Resource {
	return resourceGPUImage(GPUImageTypeBaremetal)
}
