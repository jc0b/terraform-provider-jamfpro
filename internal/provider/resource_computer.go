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

var _ resource.Resource = &ComputerResource{}
var _ resource.ResourceWithImportState = &ComputerResource{}

func NewComputerResource() resource.Resource {
	return &ComputerResource{}
}

type ComputerResource struct {
	client *jamfpro.Client
}

func (c ComputerResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_computer"

}

func (c ComputerResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Represents a Computer resource in Jamf Pro. Used primarily as a vehicle for testing the Computer datasource",
		MarkdownDescription: "This resource (`jamfpro_computer`) manages Computer records in Jamf Pro",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the Computer",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the Computer",
			},
			"serial_number": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Serial Number of the Computer",
			},
			"udid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Hardware UDID of the Computer",
			},
		},
	}
}

func (c *ComputerResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*jamfpro.Client)

	if !ok {
		response.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *jamfpro.Client, got: #{request.ProviderData}. Please report this issue to the provider developers."),
		)

		return
	}

	c.client = client
}

func (c *ComputerResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data computer

	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	computerCreateRequest := &jamfpro.ComputerCreateRequest{
		General: jamfpro.ComputerCreateGeneral{
			Name:         data.Name.ValueString(),
			SerialNumber: data.SerialNumber.ValueString(),
			Udid:         data.Udid.ValueString(),
		},
	}
	computer, _, err := c.client.Computers.Create(ctx, computerCreateRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create computer, got error: %q)", err.Error()))
		return
	}

	tflog.Trace(ctx, "created a computer")

	response.Diagnostics.Append(response.State.Set(ctx, computerForState(computer))...)

}

func (c *ComputerResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data computer
	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	computer, _, err := c.client.Computers.GetByID(ctx, int(data.Id.ValueInt64()))

	if err != nil {
		response.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read computer with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a computer")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, computerForState(computer))...)
}

func (c *ComputerResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data computer

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	computerUpdateRequest := &jamfpro.ComputerUpdateRequest{
		General: jamfpro.ComputerCreateGeneral{
			Name:         data.Name.ValueString(),
			SerialNumber: data.SerialNumber.ValueString(),
			Udid:         data.Udid.ValueString(),
		},
	}

	computer, _, err := c.client.Computers.Update(ctx, int(data.Id.ValueInt64()), computerUpdateRequest)

	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update computer with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a computer")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, computerForState(computer))...)
}

func (c *ComputerResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data computer

	diags := request.State.Get(ctx, &data)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.client.Computers.Delete(ctx, int(data.Id.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete computer with ID %d, got error: %s", data.Id.ValueInt64(), err.Error()),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Computer")
}

func (c *ComputerResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resourceImportStatePassthroughJamfProID(ctx, "computer", request, response)
}

func (c *ComputerResource) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {

}
