variable "project_name" {
  description = "Name of the project to test against"
  type        = string
  default     = "default"
}

variable "region_id" {
  description = "Numeric region identifier"
  type        = number
  default     = 76
}

variable "lb_name" {
  description = "Friendly name for the load balancer"
  type        = string
  default     = "test-lb-advanced-tls"
}

variable "lb_flavor" {
  description = "Flavor to provision"
  type        = string
  default     = "lb1-2-4"
}

variable "vip_network_cidr" {
  description = "CIDR for the VIP network"
  type        = string
  default     = "10.0.50.0/24"
}

variable "backend_subnet_cidr" {
  description = "CIDR for backend instances"
  type        = string
  default     = "10.0.51.0/24"
}

variable "backend_dns" {
  description = "DNS servers for the backend subnet"
  type        = list(string)
  default     = ["1.1.1.1", "8.8.8.8"]
}

variable "allowed_cidrs" {
  description = "CIDRs allowed to reach the HTTPS listener"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "connection_limit" {
  description = "Listener connection limit"
  type        = number
  default     = 50000
}

variable "https_secret_id" {
  description = "Secret ID containing the PKCS12 bundle for TLS termination"
  type        = string
}

variable "https_sni_secret_ids" {
  description = "Optional list of extra TLS certificates"
  type        = list(string)
  default     = []
}

variable "prometheus_secret_id" {
  description = "Secret ID containing PKCS12 bundle for the Prometheus listener"
  type        = string
}

variable "prometheus_port" {
  description = "Frontend port for the Prometheus listener"
  type        = number
  default     = 9090
}

variable "prometheus_backend_port" {
  description = "Backend port where Prometheus scrapers listen"
  type        = number
  default     = 9091
}

variable "prometheus_users" {
  description = "List of allowed Prometheus basic-auth users"
  type = list(object({
    username            = string
    encrypted_password  = string
  }))
  default = [
    {
      username           = "prom-admin"
      encrypted_password = "$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5"
    }
  ]
}

variable "backend_ca_secret_id" {
  description = "Optional CA bundle to validate member certificates"
  type        = string
  default     = null
}

variable "backend_client_secret_id" {
  description = "Optional secret with client cert to reach backends"
  type        = string
  default     = null
}

variable "healthcheck_path" {
  description = "HTTPS health-check path"
  type        = string
  default     = "/health"
}

variable "logging_topic_name" {
  description = "Name of the log topic"
  type        = string
  default     = "lb-advanced-logs"
}

variable "logging_retention_days" {
  description = "Retention period for LB logs"
  type        = number
  default     = 30
}

variable "lb_tags" {
  description = "Tags applied to the load balancer"
  type        = map(string)
  default     = {
    "environment" = "qa"
    "suite"       = "advanced-tls"
  }
}
