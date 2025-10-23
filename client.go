package main

import (
	"context"
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/organization"
)

// ClerkClient wraps the Clerk SDK client configuration
type ClerkClient struct {
	APIKey string
}

// CreateOrganization creates a new organization using the Clerk SDK
func (c *ClerkClient) CreateOrganization(ctx context.Context, params *organization.CreateParams) (*clerk.Organization, error) {
	org, err := organization.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}
	return org, nil
}

// GetOrganization retrieves an organization by ID using the Clerk SDK
func (c *ClerkClient) GetOrganization(ctx context.Context, id string) (*clerk.Organization, error) {
	org, err := organization.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return org, nil
}

// UpdateOrganization updates an existing organization using the Clerk SDK
func (c *ClerkClient) UpdateOrganization(ctx context.Context, id string, params *organization.UpdateParams) (*clerk.Organization, error) {
	org, err := organization.Update(ctx, id, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}
	return org, nil
}

// DeleteOrganization deletes an organization using the Clerk SDK
func (c *ClerkClient) DeleteOrganization(ctx context.Context, id string) error {
	_, err := organization.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	return nil
}
