variable "user_list" {
  type = list(object({
    username = string
    encrypted_password = string
   }))
  default = [
    {
      username = "admin"
      encrypted_password = "$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5"
    },
  ]
}