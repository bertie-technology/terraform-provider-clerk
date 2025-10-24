terraform {
  required_providers {
    clerk = {
      source = "registry.terraform.io/bertie-technology/clerk"
    }
  }
}

provider "clerk" {
  # Configuration via CLERK_API_KEY environment variable is recommended
  # api_key = "sk_test_your_api_key_here"
}
