---
page_title: "gcore_gpu_virtual_image Resource - terraform-provider-gcore"
subcategory: ""
description: |-
  Manages GPU-accelerated virtual machine images for cloud instances
---

# gcore_gpu_virtual_image (Resource)

Manages virtual machine images with GPU acceleration support, optimized for cloud-based machine learning, rendering workloads, and other GPU-intensive applications in virtualized environments.

## Example Usage

```terraform
provider "gcore" {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

data "gcore_region" "region" {
  name = "Luxembourg-2"
}

resource "gcore_gpu_virtual_image" "example" {
  project_id   = data.gcore_project.project.id
  region_id    = data.gcore_region.region.id
  name         = "my-cirros-image"
  url          = "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img"
  architecture = "x86_64"
  os_type      = "linux"
  ssh_key      = "allow"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Optional) The ID of the project. Conflicts with `project_name`.
* `project_name` - (Optional) The name of the project. Conflicts with `project_id`.
* `region_id` - (Optional) The ID of the region. Conflicts with `region_name`.
* `region_name` - (Optional) The name of the region. Conflicts with `region_id`.
* `name` - (Required) The name of the image. Must be unique within the project.
* `url` - (Required) The URL from which to download the image.
* `ssh_key` - (Optional) SSH key permission setting. Valid values are:
  * `allow` - (Default) Allow SSH key usage
  * `deny` - Deny SSH key usage
  * `required` - Require SSH key
* `cow_format` - (Optional) When set to `true`, the image cannot be deleted until all volumes created from it are deleted.
* `architecture` - (Optional) CPU architecture type. Valid values are:
  * `x86_64` - (Default) x86 64-bit architecture
  * `aarch64` - ARM 64-bit architecture
* `os_type` - (Optional) The type of operating system.
* `os_distro` - (Optional) The distribution of the operating system (e.g., "ubuntu", "centos").
* `os_version` - (Optional) The version of the operating system.
* `hw_firmware_type` - (Optional) The type of firmware used for booting.
* `metadata` - (Optional) A map of metadata key-value pairs to associate with the image.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for the image.
* `status` - The current status of the image.
* `metadata` - A map of metadata associated with the image.

## Import

GPU virtual images can be imported using the `project_id` and image name, separated by a slash, e.g.,

```shell
terraform import gcore_gpu_virtual_image.example project_id/image_name
```
