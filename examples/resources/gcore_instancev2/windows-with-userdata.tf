variable "second_user_userdata" {
 description = "This is a variable of type string"
 type        = string
 default     = <<EOF
<powershell>
# Be sure to set the username and password on these two lines. Of course this is not a good
# security practice to include a password at command line.
$User = "SecondUser"
$Password = ConvertTo-SecureString "s3cR3tP@ssw0rd" -AsPlainText -Force
New-LocalUser $User -Password $Password
Add-LocalGroupMember -Group "Remote Desktop Users" -Member $User
Add-LocalGroupMember -Group "Administrators" -Member $User
</powershell>
EOF
}

data "gcore_image" "windows" {
  name       = "windows-server-2022"
  region_id  = data.gcore_region.region.id
  project_id = data.gcore_project.project.id
}

resource "gcore_volume" "boot_volume_windows" {
  name       = "my-windows-boot-volume"
  type_name  = "ssd_hiiops"
  size       = 50
  image_id   = data.gcore_image.windows.id
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_instancev2" "instance" {
  flavor_id     = "g1w-standard-4-8"
  name          = "my-windows-instance"
  password      = "my-s3cR3tP@ssw0rd"
  user_data     = base64encode(var.second_user_userdata)

  volume {
    volume_id  = gcore_volume.boot_volume_windows.id
    boot_index = 0
  }

  interface {
    type = "external"
    name = "my-external-interface"
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}