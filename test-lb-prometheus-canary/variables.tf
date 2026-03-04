variable "project_name" {
  description = "Project name"
  type        = string
  default     = "default"
}

variable "region_id" {
  description = "Region ID"
  type        = number
  default     = 76
}

variable "lb_name" {
  description = "Load balancer name"
  type        = string
  default     = "test-lb-prometheus-canary"
}

variable "lb_flavor" {
  description = "Flavor"
  type        = string
  default     = "lb1-2-4"
}

variable "vip_network_cidr" {
  description = "CIDR for VIP network"
  type        = string
  default     = "10.0.60.0/24"
}

variable "backend_subnet_cidr" {
  description = "CIDR for backend network"
  type        = string
  default     = "10.0.61.0/24"
}

variable "allowed_cidrs" {
  description = "CIDRs permitted on prod listener"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "canary_port" {
  description = "Frontend HTTP port for the canary listener"
  type        = number
  default     = 8080
}

variable "canary_backend_port" {
  description = "Backend port for canary traffic"
  type        = number
  default     = 8081
}

variable "canary_health_path" {
  description = "Health endpoint for the canary pool"
  type        = string
  default     = "/status"
}

variable "prometheus_port" {
  description = "Frontend Prometheus listener port"
  type        = number
  default     = 9090
}

variable "prometheus_backend_port" {
  description = "Backend metrics port"
  type        = number
  default     = 9100
}

variable "prometheus_secret_id" {
  description = "Secret ID with PKCS12 bundle for the Prometheus listener"
  type        = string
}

variable "prometheus_users" {
  description = "Basic-auth credentials for Prometheus listener"
  type = list(object({
    username           = string
    encrypted_password = string
  }))
  default = [
    {
      username           = "metrics"
      encrypted_password = "$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5"
    }
  ]
}

variable "lb_tags" {
  description = "Tags applied to the LB"
  type        = map(string)
  default     = {
    "environment" = "qa"
    "suite"       = "prom-canary"
  }
}
