terraform {
  required_providers {
    clerk = {
      source = "registry.terraform.io/bertie-technology/clerk"
    }
  }
}

provider "clerk" {
  api_key = var.clerk_api_key
}

# Example: Create an organization
resource "clerk_organization" "example" {
  name                    = var.organization_name
  slug                    = var.organization_slug
  max_allowed_memberships = var.max_allowed_memberships

  public_metadata  = var.public_metadata != null ? jsonencode(var.public_metadata) : null
  private_metadata = var.private_metadata != null ? jsonencode(var.private_metadata) : null
}

# Outputs - show all organization values
output "organization_id" {
  value       = clerk_organization.example.id
  description = "The ID of the created organization"
}

output "organization_name" {
  value       = clerk_organization.example.name
  description = "The name of the organization"
}

output "organization_slug" {
  value       = clerk_organization.example.slug
  description = "The slug of the organization"
}

output "max_allowed_memberships" {
  value       = clerk_organization.example.max_allowed_memberships
  description = "Maximum number of memberships allowed"
}

output "public_metadata" {
  value       = clerk_organization.example.public_metadata
  description = "Public metadata (JSON string)"
}

output "private_metadata" {
  value       = clerk_organization.example.private_metadata
  description = "Private metadata (JSON string)"
  sensitive   = true
}

output "created_by" {
  value       = clerk_organization.example.created_by
  description = "User ID who created the organization"
}
