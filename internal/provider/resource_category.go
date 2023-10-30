package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
)

var _ resource.Resource = &CategoryResource{}
var _ resource.ResourceWithImportState = &CategoryResource{}

func NewCategoryResource() resource.Resource {
	return &CategoryResource{}
}

type CategoryResource struct {
	client *jamfpro.Client
}

func (c *CategoryResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*jamfpro.Client)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *jamfpro.Client, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	c.client = client
}

func (c *CategoryResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_category"
}

func (c *CategoryResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Represents a category resource in Jamf Pro",
		MarkdownDescription: "This resource (`jamfpro_category`) manages Categories in Jamf Pro",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the Category",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the Category",
			},
			"priority": schema.Int64Attribute{
				Optional:    true,
				Description: "The Category priority",
			},
		},
	}
}

func (c *CategoryResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data category

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	CategoryCreateRequest := &jamfpro.CategoryCreateRequest{
		Name:     data.Name.ValueString(),
		Priority: int(data.Priority.ValueInt64()),
	}
	category, _, err := c.client.Categories.Create(ctx, CategoryCreateRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create category, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a tag")

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, categoryForState(category))...)
}

func (c *CategoryResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data category

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	category, _, err := c.client.Categories.GetByID(ctx, int(data.Id.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read category with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Category")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, categoryForState(category))...)
}

func (c *CategoryResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data category

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	CategoryUpdateRequest := &jamfpro.CategoryUpdateRequest{
		Name:     data.Name.ValueString(),
		Priority: int(data.Priority.ValueInt64()),
	}
	category, _, err := c.client.Categories.Update(ctx, int(data.Id.ValueInt64()), CategoryUpdateRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update category with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Category")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, categoryForState(category))...)
}

func (c *CategoryResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data category

	diags := request.State.Get(ctx, &data)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.client.Categories.Delete(ctx, int(data.Id.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete category with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Category")
}

func (c *CategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughJamfProID(ctx, "category", req, resp)
}
