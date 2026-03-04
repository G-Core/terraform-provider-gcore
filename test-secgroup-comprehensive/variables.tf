variable "sg_name" {
  default = "tf-test-sg-extended"
}

variable "sg_description" {
  default = "Extended security group test"
}

variable "create_rules" {
  default = false
}

variable "create_import_sg" {
  default = false
}
