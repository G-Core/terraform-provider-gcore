# ===== LOAD BALANCER VARIABLES =====
variable "lb_name" {
  description = "Load balancer name - test rename drift"
  type        = string
  default     = "qa-lb-comprehensive"
}

variable "lb_flavor" {
  description = "Load balancer flavor"
  type        = string
  default     = "lb1-2-4"
}

variable "lb_tags" {
  description = "LB tags - test tag operations including removal"
  type        = list(string)
  default     = []
}

# ===== LISTENER VARIABLES =====
variable "listener_name" {
  description = "Listener name"
  type        = string
  default     = "qa-listener"
}

variable "listener_protocol" {
  description = "Listener protocol"
  type        = string
  default     = "HTTP"
}

variable "listener_port" {
  description = "Listener port"
  type        = number
  default     = 80
}

# Test computed_optional timeout fields (Kirill's issue)
variable "timeout_client_data" {
  description = "Frontend client inactivity timeout - test drift"
  type        = number
  default     = null
}

variable "timeout_member_connect" {
  description = "Backend member connection timeout - test drift"
  type        = number
  default     = null
}

variable "timeout_member_data" {
  description = "Backend member inactivity timeout - test drift"
  type        = number
  default     = null
}

# Test list fields
variable "sni_secret_ids" {
  description = "SNI secret IDs - test list clearing"
  type        = list(string)
  default     = []
}

variable "user_list" {
  description = "User list for basic auth - test list clearing"
  type = list(object({
    username           = string
    encrypted_password = string
  }))
  default = []
}

# ===== POOL VARIABLES =====
variable "create_pool" {
  description = "Whether to create pool - test timing with listener"
  type        = bool
  default     = true
}

variable "pool_name" {
  description = "Pool name"
  type        = string
  default     = "qa-pool"
}

variable "pool_algorithm" {
  description = "Pool load balancing algorithm"
  type        = string
  default     = "ROUND_ROBIN"
}

variable "pool_protocol" {
  description = "Pool protocol"
  type        = string
  default     = "HTTP"
}

variable "pool_timeout_client_data" {
  description = "Pool frontend client timeout"
  type        = number
  default     = null
}

variable "pool_timeout_member_connect" {
  description = "Pool backend connection timeout"
  type        = number
  default     = null
}

variable "pool_timeout_member_data" {
  description = "Pool backend inactivity timeout"
  type        = number
  default     = null
}

# Health monitor with computed_optional fields
variable "pool_healthmonitor" {
  description = "Pool health monitor configuration"
  type = object({
    type             = string
    delay            = number
    max_retries      = number
    timeout          = number
    http_method      = optional(string)
    max_retries_down = optional(number)
    url_path         = optional(string)
    expected_codes   = optional(string)
  })
  default = null
}

# Session persistence
variable "pool_session_persistence" {
  description = "Pool session persistence configuration"
  type = object({
    type                    = string
    cookie_name             = optional(string)
    persistence_granularity = optional(string)
    persistence_timeout     = optional(number)
  })
  default = null
}

# Pool members inline
variable "pool_members" {
  description = "Pool members defined inline"
  type = list(object({
    address         = string
    protocol_port   = number
    subnet_id       = string
    weight          = optional(number)
    admin_state_up  = optional(bool)
    backup          = optional(bool)
    monitor_address = optional(string)
    monitor_port    = optional(number)
  }))
  default = []
}

# ===== POOL MEMBER (SEPARATE RESOURCE) VARIABLES =====
variable "create_separate_member" {
  description = "Whether to create separate pool member resource"
  type        = bool
  default     = false
}

variable "separate_member_address" {
  description = "Separate member IP address"
  type        = string
  default     = "10.0.1.10"
}

variable "separate_member_port" {
  description = "Separate member port"
  type        = number
  default     = 80
}

variable "separate_member_weight" {
  description = "Separate member weight"
  type        = number
  default     = null
}

variable "separate_member_admin_state" {
  description = "Separate member admin state"
  type        = bool
  default     = null
}

variable "separate_member_backup" {
  description = "Separate member backup flag"
  type        = bool
  default     = null
}

variable "separate_member_monitor_address" {
  description = "Separate member monitor address"
  type        = string
  default     = null
}

variable "separate_member_monitor_port" {
  description = "Separate member monitor port"
  type        = number
  default     = null
}
