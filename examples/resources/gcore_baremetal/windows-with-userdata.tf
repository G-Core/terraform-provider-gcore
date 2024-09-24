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

resource "gcore_baremetal" "baremetal_windows_with_userdata" {
  flavor_id     = "bm1-infrastructure-small"
  name          = "my-windows-baremetal"
  password      = "my-s3cR3tP@ssw0rd"
  user_data     = base64encode(var.second_user_userdata)
  image_id      = data.gcore_image.windows.id

  interface {
    type = "external"
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}