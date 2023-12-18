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

var _ resource.Resource = &DepartmentResource{}
var _ resource.ResourceWithImportState = &DepartmentResource{}

func NewDepartmentResource() resource.Resource {
	return &DepartmentResource{}
}

type DepartmentResource struct {
	client *jamfpro.Client
}

func (c *DepartmentResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (c *DepartmentResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_department"
}

func (c *DepartmentResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Represents a department resource in Jamf Pro",
		MarkdownDescription: "This resource (`jamfpro_department`) manages Departments in Jamf Pro",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the Department",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the Department",
			},
		},
	}
}

func (c *DepartmentResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data department

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	DepartmentCreateRequest := &jamfpro.DepartmentCreateRequest{
		Name: data.Name.ValueString(),
	}
	department, _, err := c.client.Departments.Create(ctx, DepartmentCreateRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create department, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a building")

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, departmentForState(department))...)
}

func (c *DepartmentResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data department

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	department, _, err := c.client.Departments.GetByID(ctx, int(data.Id.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read department with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Department")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, departmentForState(department))...)
}

func (c *DepartmentResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data department

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	DepartmentUpdateRequest := &jamfpro.DepartmentUpdateRequest{
		Name: data.Name.ValueString(),
	}
	department, _, err := c.client.Departments.Update(ctx, int(data.Id.ValueInt64()), DepartmentUpdateRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update department with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Department")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, departmentForState(department))...)
}

func (c *DepartmentResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data department

	diags := request.State.Get(ctx, &data)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.client.Departments.Delete(ctx, int(data.Id.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete department with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Department")
}

func (c *DepartmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughJamfProID(ctx, "department", req, resp)
}
