# Testing the Clerk Terraform Provider

This document describes how to test the Terraform provider for Clerk.

## Types of Tests

### 1. Unit Tests

Unit tests verify individual functions and components without making real API calls.

```bash
# Run unit tests
make test

# Or directly with go
go test -v ./...
```

### 2. Acceptance Tests

Acceptance tests run actual Terraform operations against the Clerk API to verify end-to-end functionality.

**⚠️ Warning:** Acceptance tests create real resources in your Clerk account. Make sure you're using a test account or are prepared for these resources to be created and destroyed.

## Running Acceptance Tests

### Prerequisites

1. **Clerk API Key**: You need a valid Clerk API key for a test account
2. **Clean Test Environment**: Use a dedicated Clerk account for testing
3. **Network Access**: Tests need to reach the Clerk API

### Setup

Set your Clerk API key as an environment variable:

```bash
export CLERK_API_KEY="sk_test_your_test_api_key_here"
```

**Important:** Use a test/development API key, not your production key!

### Running Tests

```bash
# Run all acceptance tests
make testacc

# Or run with go directly
TF_ACC=1 go test -v -count=1 -timeout 30m ./...

# Run a specific test
TF_ACC=1 go test -v -run TestAccOrganizationResource_minimal ./...
```

### Test Options

```bash
# Run tests in parallel (faster)
TF_ACC=1 go test -v -parallel=4 ./...

# Run with more verbose output
TF_ACC=1 go test -v -count=1 ./... 2>&1 | tee test.log

# Run a specific test file
TF_ACC=1 go test -v ./resource_organization_test.go ./provider_test.go
```

## Available Acceptance Tests

### TestAccOrganizationResource

Basic CRUD (Create, Read, Update, Delete) operations with name and slug.

### TestAccOrganizationResource_withMetadata

Tests creating an organization with public and private metadata.

### TestAccOrganizationResource_minimal

Tests creating an organization with only the required name field.

### TestAccOrganizationResource_withMaxMemberships

Tests setting and updating max_allowed_memberships.

## Writing New Tests

When adding new features, add corresponding tests:

1. **Create a test function** in `resource_organization_test.go`
2. **Use descriptive names** like `TestAccOrganizationResource_yourFeature`
3. **Test multiple scenarios**: create, update, import
4. **Clean up properly**: The test framework handles cleanup automatically

Example test structure:

```go
func TestAccOrganizationResource_newFeature(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccConfigForNewFeature(),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("clerk_organization.test", "field", "value"),
                ),
            },
        },
    })
}
```

## Debugging Failed Tests

### Enable Detailed Logging

```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=./test-debug.log
TF_ACC=1 go test -v -run TestAccOrganizationResource ./...
```

### Check API Responses

Failed tests will show the full error from the Clerk API. Common issues:

- **401 Unauthorized**: Check your API key
- **409 Conflict**: Organization slug already exists (tests should use unique names)
- **Rate Limiting**: Add delays between tests if needed

### Manual Cleanup

If tests fail and leave resources behind:

```bash
# List organizations
curl -H "Authorization: Bearer $CLERK_API_KEY" \
  https://api.clerk.com/v1/organizations

# Delete test organizations manually
curl -X DELETE \
  -H "Authorization: Bearer $CLERK_API_KEY" \
  https://api.clerk.com/v1/organizations/org_xxxxx
```

## Continuous Integration

For CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run acceptance tests
  env:
    CLERK_API_KEY: ${{ secrets.CLERK_TEST_API_KEY }}
  run: make testacc
```

**Best Practice:** Use separate Clerk accounts for:

- Development
- Testing/CI
- Production

Never use production API keys in tests!

## Test Coverage

Check test coverage:

```bash
# Generate coverage report
go test -v -cover -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out
```

## Common Issues

### Test Timeout

If tests timeout, increase the timeout:

```bash
TF_ACC=1 go test -v -timeout 60m ./...
```

### Parallel Test Conflicts

If tests interfere with each other (e.g., slug conflicts), run sequentially:

```bash
TF_ACC=1 go test -v -parallel=1 ./...
```

### API Rate Limits

Add delays between test steps if you hit rate limits:

```go
time.Sleep(1 * time.Second)
```

## Resources

- [Terraform Plugin Testing Guide](https://developer.hashicorp.com/terraform/plugin/testing)
- [Clerk API Documentation](https://clerk.com/docs/reference/backend-api)
- [Go Testing Package](https://pkg.go.dev/testing)
