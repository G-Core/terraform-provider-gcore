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

# Create a Windows baremetal server with userdata to add a second user
resource "gcore_cloud_baremetal_server" "windows_with_userdata" {
  project_id          = 1
  region_id           = 1
  flavor              = "bm1-infrastructure-small"
  name                = "my-windows-bare-metal"
  image_id            = "408a0e4d-6a28-4bae-93fa-f738d964f555"
  password_wo         = "my-s3cR3tP@ssw0rd"
  password_wo_version = 1
  user_data           = base64encode(var.second_user_userdata)

  interfaces = [{
    type      = "external"
    ip_family = "ipv4"
  }]
}
