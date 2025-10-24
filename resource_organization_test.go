package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOrganizationResourceConfig("Test Org", "test-org"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clerk_organization.test", "name", "Test Org"),
					resource.TestCheckResourceAttr("clerk_organization.test", "slug", "test-org"),
					resource.TestCheckResourceAttrSet("clerk_organization.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "clerk_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccOrganizationResourceConfig("Test Org Updated", "test-org-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clerk_organization.test", "name", "Test Org Updated"),
					resource.TestCheckResourceAttr("clerk_organization.test", "slug", "test-org-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccOrganizationResource_withMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with metadata
			{
				Config: testAccOrganizationResourceConfigWithMetadata("Metadata Org", "metadata-org"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clerk_organization.test", "name", "Metadata Org"),
					resource.TestCheckResourceAttrSet("clerk_organization.test", "public_metadata"),
					resource.TestCheckResourceAttrSet("clerk_organization.test", "private_metadata"),
				),
			},
		},
	})
}

func TestAccOrganizationResource_minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with only required fields
			{
				Config: testAccOrganizationResourceConfigMinimal("Minimal Org"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clerk_organization.test", "name", "Minimal Org"),
					resource.TestCheckResourceAttrSet("clerk_organization.test", "id"),
					resource.TestCheckResourceAttrSet("clerk_organization.test", "slug"), // Should be auto-generated
				),
			},
		},
	})
}

func TestAccOrganizationResource_withMaxMemberships(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with max_allowed_memberships
			{
				Config: testAccOrganizationResourceConfigWithMax("Max Org", "max-org", 50),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clerk_organization.test", "name", "Max Org"),
					resource.TestCheckResourceAttr("clerk_organization.test", "max_allowed_memberships", "50"),
				),
			},
			// Update max_allowed_memberships
			{
				Config: testAccOrganizationResourceConfigWithMax("Max Org", "max-org", 100),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clerk_organization.test", "max_allowed_memberships", "100"),
				),
			},
		},
	})
}

// Test configuration functions

func testAccOrganizationResourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "clerk_organization" "test" {
  name = %[1]q
  slug = %[2]q
}
`, name, slug)
}

func testAccOrganizationResourceConfigMinimal(name string) string {
	return fmt.Sprintf(`
resource "clerk_organization" "test" {
  name = %[1]q
}
`, name)
}

func testAccOrganizationResourceConfigWithMetadata(name, slug string) string {
	return fmt.Sprintf(`
resource "clerk_organization" "test" {
  name = %[1]q
  slug = %[2]q

  public_metadata = jsonencode({
    environment = "test"
    region      = "us-west-1"
  })

  private_metadata = jsonencode({
    test_id = "acc-test-123"
  })
}
`, name, slug)
}

func testAccOrganizationResourceConfigWithMax(name, slug string, max int) string {
	return fmt.Sprintf(`
resource "clerk_organization" "test" {
  name                    = %[1]q
  slug                    = %[2]q
  max_allowed_memberships = %[3]d
}
`, name, slug, max)
}
