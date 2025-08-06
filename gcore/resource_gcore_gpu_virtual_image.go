package gcore

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceVirtualImage() *schema.Resource {
	return resourceGPUImage(GPUNodeTypeVirtual)
}
