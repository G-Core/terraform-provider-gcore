variable "second_user_userdata" {
  description = "PowerShell script to create a second Windows user"
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

resource "gcore_cloud_volume" "boot_volume_windows_userdata" {
  project_id = 1
  region_id  = 1
  name       = "my-windows-boot-volume"
  source     = "image"
  image_id   = "a2c1681c-94e0-4aab-8fa3-09a8e662d4c0"
  size       = 50
  type_name  = "ssd_hiiops"
}

resource "gcore_cloud_instance" "windows_with_userdata" {
  project_id = 1
  region_id  = 1
  flavor     = "g1w-standard-4-8"
  name       = "my-windows-instance"
  password   = "my-s3cR3tP@ssw0rd"
  user_data  = base64encode(var.second_user_userdata)

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume_windows_userdata.id }]

  interfaces = [{
    type      = "external"
    ip_family = "ipv4"
  }]
}
