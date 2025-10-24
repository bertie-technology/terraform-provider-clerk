# Basic organization
resource "clerk_organization" "example" {
  name = "Example Organization"
  slug = "example-org"
}

# Organization with metadata
resource "clerk_organization" "with_metadata" {
  name                    = "Production Organization"
  slug                    = "prod-org"
  max_allowed_memberships = 100

  public_metadata = jsonencode({
    environment = "production"
    region      = "us-west-2"
    tier        = "premium"
  })

  private_metadata = jsonencode({
    billing_id    = "cus_123456789"
    internal_code = "ORG-PROD-001"
  })
}

# Minimal organization (auto-generated slug)
resource "clerk_organization" "minimal" {
  name = "Simple Organization"
}
