package main

import (
	"context"
	"encoding/json"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/organization"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &organizationResource{}
	_ resource.ResourceWithConfigure   = &organizationResource{}
	_ resource.ResourceWithImportState = &organizationResource{}
)

// NewOrganizationResource is a helper function to simplify the provider implementation
func NewOrganizationResource() resource.Resource {
	return &organizationResource{}
}

// organizationResource is the resource implementation
type organizationResource struct {
	client *ClerkClient
}

// organizationResourceModel describes the resource data model
type organizationResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Slug                  types.String `tfsdk:"slug"`
	MaxAllowedMemberships types.Int64  `tfsdk:"max_allowed_memberships"`
	PublicMetadata        types.String `tfsdk:"public_metadata"`
	PrivateMetadata       types.String `tfsdk:"private_metadata"`
	CreatedBy             types.String `tfsdk:"created_by"`
}

// Metadata returns the resource type name
func (r *organizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Schema defines the schema for the resource
func (r *organizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Clerk organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the organization.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the organization.",
				Required:    true,
			},
			"slug": schema.StringAttribute{
				Description: "The slug of the organization. If not provided, one will be generated from the name.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_allowed_memberships": schema.Int64Attribute{
				Description: "The maximum number of memberships allowed for the organization.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"public_metadata": schema.StringAttribute{
				Description: "Public metadata for the organization (JSON string).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_metadata": schema.StringAttribute{
				Description: "Private metadata for the organization (JSON string).",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Description: "The user ID who created the organization.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *organizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ClerkClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *ClerkClient, got something else. Please report this issue to the provider developers.",
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state
func (r *organizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organizationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the organization parameters
	params := &organization.CreateParams{
		Name: clerk.String(plan.Name.ValueString()),
	}

	if !plan.Slug.IsNull() && !plan.Slug.IsUnknown() {
		params.Slug = clerk.String(plan.Slug.ValueString())
	}

	if !plan.MaxAllowedMemberships.IsNull() && !plan.MaxAllowedMemberships.IsUnknown() {
		params.MaxAllowedMemberships = clerk.Int64(plan.MaxAllowedMemberships.ValueInt64())
	}

	if !plan.CreatedBy.IsNull() && !plan.CreatedBy.IsUnknown() {
		params.CreatedBy = clerk.String(plan.CreatedBy.ValueString())
	}

	// Parse public metadata
	if !plan.PublicMetadata.IsNull() && !plan.PublicMetadata.IsUnknown() {
		var metadata json.RawMessage
		if err := json.Unmarshal([]byte(plan.PublicMetadata.ValueString()), &metadata); err != nil {
			resp.Diagnostics.AddError(
				"Error parsing public_metadata",
				"Could not parse public_metadata as JSON: "+err.Error(),
			)
			return
		}
		params.PublicMetadata = &metadata
	}

	// Parse private metadata
	if !plan.PrivateMetadata.IsNull() && !plan.PrivateMetadata.IsUnknown() {
		var metadata json.RawMessage
		if err := json.Unmarshal([]byte(plan.PrivateMetadata.ValueString()), &metadata); err != nil {
			resp.Diagnostics.AddError(
				"Error parsing private_metadata",
				"Could not parse private_metadata as JSON: "+err.Error(),
			)
			return
		}
		params.PrivateMetadata = &metadata
	}

	// Create the organization
	org, err := r.client.CreateOrganization(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating organization",
			"Could not create organization: "+err.Error(),
		)
		return
	}

	// Set the ID so we can fetch the full resource
	plan.ID = types.StringValue(org.ID)

	// Fetch the organization again to get the complete state from the API
	// This ensures we capture any computed fields or defaults set by the API
	org, err = r.client.GetOrganization(ctx, org.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading organization after create",
			"Could not read organization ID "+org.ID+" after creation: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.Name = types.StringValue(org.Name)
	plan.Slug = types.StringValue(org.Slug)

	// Always set a known value for max_allowed_memberships
	if org.MaxAllowedMemberships > 0 {
		plan.MaxAllowedMemberships = types.Int64Value(org.MaxAllowedMemberships)
	} else {
		plan.MaxAllowedMemberships = types.Int64Null()
	}

	// Always set a known value for public_metadata
	if org.PublicMetadata != nil {
		metadata, err := json.Marshal(org.PublicMetadata)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error serializing public_metadata",
				"Could not serialize public_metadata: "+err.Error(),
			)
			return
		}
		plan.PublicMetadata = types.StringValue(string(metadata))
	} else {
		plan.PublicMetadata = types.StringNull()
	}

	// Always set a known value for private_metadata
	if org.PrivateMetadata != nil {
		metadata, err := json.Marshal(org.PrivateMetadata)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error serializing private_metadata",
				"Could not serialize private_metadata: "+err.Error(),
			)
			return
		}
		plan.PrivateMetadata = types.StringValue(string(metadata))
	} else {
		plan.PrivateMetadata = types.StringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data
func (r *organizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organizationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the organization from Clerk
	org, err := r.client.GetOrganization(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading organization",
			"Could not read organization ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update state with refreshed values
	state.Name = types.StringValue(org.Name)
	state.Slug = types.StringValue(org.Slug)

	// Always set a known value for max_allowed_memberships
	if org.MaxAllowedMemberships > 0 {
		state.MaxAllowedMemberships = types.Int64Value(org.MaxAllowedMemberships)
	} else {
		state.MaxAllowedMemberships = types.Int64Null()
	}

	// Always set a known value for public_metadata
	if org.PublicMetadata != nil {
		metadata, err := json.Marshal(org.PublicMetadata)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error serializing public_metadata",
				"Could not serialize public_metadata: "+err.Error(),
			)
			return
		}
		state.PublicMetadata = types.StringValue(string(metadata))
	} else {
		state.PublicMetadata = types.StringNull()
	}

	// Always set a known value for private_metadata
	if org.PrivateMetadata != nil {
		metadata, err := json.Marshal(org.PrivateMetadata)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error serializing private_metadata",
				"Could not serialize private_metadata: "+err.Error(),
			)
			return
		}
		state.PrivateMetadata = types.StringValue(string(metadata))
	} else {
		state.PrivateMetadata = types.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *organizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan organizationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the organization update parameters
	params := &organization.UpdateParams{
		Name: clerk.String(plan.Name.ValueString()),
	}

	if !plan.Slug.IsNull() && !plan.Slug.IsUnknown() {
		params.Slug = clerk.String(plan.Slug.ValueString())
	}

	if !plan.MaxAllowedMemberships.IsNull() && !plan.MaxAllowedMemberships.IsUnknown() {
		params.MaxAllowedMemberships = clerk.Int64(plan.MaxAllowedMemberships.ValueInt64())
	}

	// Parse public metadata
	if !plan.PublicMetadata.IsNull() && !plan.PublicMetadata.IsUnknown() {
		var metadata json.RawMessage
		if err := json.Unmarshal([]byte(plan.PublicMetadata.ValueString()), &metadata); err != nil {
			resp.Diagnostics.AddError(
				"Error parsing public_metadata",
				"Could not parse public_metadata as JSON: "+err.Error(),
			)
			return
		}
		params.PublicMetadata = &metadata
	}

	// Parse private metadata
	if !plan.PrivateMetadata.IsNull() && !plan.PrivateMetadata.IsUnknown() {
		var metadata json.RawMessage
		if err := json.Unmarshal([]byte(plan.PrivateMetadata.ValueString()), &metadata); err != nil {
			resp.Diagnostics.AddError(
				"Error parsing private_metadata",
				"Could not parse private_metadata as JSON: "+err.Error(),
			)
			return
		}
		params.PrivateMetadata = &metadata
	}

	// Update the organization
	_, err := r.client.UpdateOrganization(ctx, plan.ID.ValueString(), params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating organization",
			"Could not update organization ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Fetch the organization again to get the latest state from the API
	// This ensures we capture any values set by the API (like computed fields)
	org, err := r.client.GetOrganization(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading organization after update",
			"Could not read organization ID "+plan.ID.ValueString()+" after update: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.Name = types.StringValue(org.Name)
	plan.Slug = types.StringValue(org.Slug)

	// Always set a known value for max_allowed_memberships
	if org.MaxAllowedMemberships > 0 {
		plan.MaxAllowedMemberships = types.Int64Value(org.MaxAllowedMemberships)
	} else {
		plan.MaxAllowedMemberships = types.Int64Null()
	}

	// Always set a known value for public_metadata
	if org.PublicMetadata != nil {
		metadata, err := json.Marshal(org.PublicMetadata)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error serializing public_metadata",
				"Could not serialize public_metadata: "+err.Error(),
			)
			return
		}
		plan.PublicMetadata = types.StringValue(string(metadata))
	} else {
		plan.PublicMetadata = types.StringNull()
	}

	// Always set a known value for private_metadata
	if org.PrivateMetadata != nil {
		metadata, err := json.Marshal(org.PrivateMetadata)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error serializing private_metadata",
				"Could not serialize private_metadata: "+err.Error(),
			)
			return
		}
		plan.PrivateMetadata = types.StringValue(string(metadata))
	} else {
		plan.PrivateMetadata = types.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *organizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state organizationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the organization
	err := r.client.DeleteOrganization(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting organization",
			"Could not delete organization ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

// ImportState imports an existing resource into Terraform state
func (r *organizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the ID field for import
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
