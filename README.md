# Terraform Provider for Clerk

A Terraform provider for managing Clerk resources. This provider uses:

- The official [Clerk Go SDK](https://github.com/clerk/clerk-sdk-go) for interacting with the Clerk API
- HashiCorp's [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) for a modern, type-safe provider implementation

## Current Scope

This provider currently supports:

- **Organizations** - Create, read, update, and delete Clerk organizations

Additional resources may be added in future versions.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25 (for development)
- A Clerk account with an API key

## Building the Provider

Clone the repository and build the provider:

```bash
make build
# Or manually:
mkdir -p bin
go build -o bin/terraform-provider-clerk
```

## Installing the Provider

### Local Development

For local development, you can use the following `~/.terraformrc` configuration to tell Terraform where to find your local provider:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/bertie-technology/clerk" = "/Users/bertie-technology/Documents/work/terraform-clerk-organization/bin"
  }

  direct {}
}
```

Replace the path with the absolute path to your `bin/` directory containing the built provider binary.

## Using the Provider

### Provider Configuration

```hcl
terraform {
  required_providers {
    clerk = {
      source = "bertie-technology/clerk"
      version = "~> 0.1"
    }
  }
}

provider "clerk" {
  api_key = "your-clerk-api-key"  # Or set CLERK_API_KEY environment variable
}
```

### Resources

The following resources are currently available:

#### `clerk_organization`

Manages a Clerk organization. This is currently the only resource provided by this provider.

**Example Usage:**

```hcl
resource "clerk_organization" "example" {
  name                    = "My Organization"
  slug                    = "my-org"
  max_allowed_memberships = 100

  public_metadata = jsonencode({
    environment = "production"
    region      = "us-west-2"
  })

  private_metadata = jsonencode({
    billing_id = "cus_123456"
  })
}
```

**Argument Reference:**

- `name` - (Required) The name of the organization.
- `slug` - (Optional) The slug of the organization. If not provided, one will be generated from the name.
- `max_allowed_memberships` - (Optional) The maximum number of memberships allowed for the organization.
- `public_metadata` - (Optional) Public metadata for the organization as a JSON string. Defaults to `{}`.
- `private_metadata` - (Optional, Sensitive) Private metadata for the organization as a JSON string. Defaults to `{}`.
- `created_by` - (Optional) The user ID who created the organization.

**Attribute Reference:**

In addition to all arguments above, the following attributes are exported:

- `id` - The unique identifier for the organization.

All arguments are also available as attributes and can be referenced in outputs or other resources.

## Environment Variables

- `CLERK_API_KEY` - Your Clerk API key (can be used instead of `api_key` in provider configuration)

## Development

### Building

```bash
make build
# Or manually:
go build -o bin/terraform-provider-clerk
```

### Testing

#### Unit Tests

```bash
make test
```

#### Acceptance Tests

Acceptance tests create real resources in Clerk and require a valid API key.

**⚠️ Warning:** Use a test/development Clerk account, not production!

```bash
# Set your test API key
export CLERK_API_KEY="sk_test_your_test_api_key"

# Run acceptance tests
make testacc
```

See [TESTING.md](TESTING.md) for detailed testing documentation.

### Documentation

Generate provider documentation:

```bash
make docs
```

This uses [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs) to automatically generate documentation from:

- Provider schema
- Resource schemas
- Example files in `examples/`
- Templates in `templates/`

Generated docs are placed in the `docs/` directory.

### Running with Debug Mode

```bash
go run main.go -debug
```

## Documentation

Full provider documentation is available in the [docs/](docs/) directory:

- [Provider Configuration](docs/index.md)
- [clerk_organization Resource](docs/resources/organization.md)

## License

MIT
