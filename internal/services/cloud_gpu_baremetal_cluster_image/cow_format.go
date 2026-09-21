package cloud_gpu_baremetal_cluster_image

import "github.com/hashicorp/terraform-plugin-framework/types"

// cowFormatFromDiskFormat recovers the create-only cow_format flag from the
// disk_format the API returns instead of it: an image uploaded with
// cow_format=true is stored as raw, one uploaded with cow_format=false as
// qcow2. Any other disk format leaves the flag null.
//
// Without it an imported image has cow_format null in state, so a config that
// sets cow_format = true forces replacement on the first plan, and a config
// that omits it silently records the default false on a raw image.
func cowFormatFromDiskFormat(diskFormat types.String) types.Bool {
	switch diskFormat.ValueString() {
	case "raw":
		return types.BoolValue(true)
	case "qcow2":
		return types.BoolValue(false)
	}
	return types.BoolNull()
}
