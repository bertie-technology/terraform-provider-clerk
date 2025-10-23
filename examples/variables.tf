variable "clerk_api_key" {
  description = "Clerk API Key for authentication"
  type        = string
  sensitive   = true
}

variable "organization_name" {
  description = "Name of the organization to create"
  type        = string
}

variable "organization_slug" {
  description = "Slug for the organization. If not provided, one will be generated from the name."
  type        = string
  default     = null
}

variable "max_allowed_memberships" {
  description = "Maximum number of memberships allowed"
  type        = number
  default     = null
}

variable "public_metadata" {
  description = "Public metadata for the organization"
  type        = map(string)
  default     = null
}

variable "private_metadata" {
  description = "Private metadata for the organization"
  type        = map(string)
  sensitive   = true
  default     = null
}
