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
	"time"
)

var _ resource.Resource = &ComputerGroupResource{}
var _ resource.ResourceWithImportState = &ComputerGroupResource{}

func NewComputerGroupResource() resource.Resource {
	return &ComputerGroupResource{}
}

type ComputerGroupResource struct {
	client *jamfpro.Client
}

func (c ComputerGroupResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_computergroup"

}

func (c ComputerGroupResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Represents a Computer Group resource in Jamf Pro",
		MarkdownDescription: "This resource (`jamfpro_computergroup`) manages Computer Groups in Jamf Pro",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the Computer Group",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the Computer Group",
			},
			"computers": schema.SetNestedAttribute{
				Required:    true,
				Computed:    false,
				Description: "Represents computers that are members of a static group.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description:         "ID of the computer.",
							MarkdownDescription: "`ID` of the computer.",
							Optional:            true,
						},
						"name": schema.StringAttribute{
							Description:         "Name of the computer.",
							MarkdownDescription: "`name` of the computer.",
							Optional:            true,
						},
						"serial_number": schema.StringAttribute{
							Description:         "Serial number of the computer.",
							MarkdownDescription: "`serial_number` of the computer.",
							Optional:            true,
						},
						"udid": schema.StringAttribute{
							Description:         "Hardware UDID of the computer.",
							MarkdownDescription: "`udid` of the computer.",
							Computed:            true,
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

func (c *ComputerGroupResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (c *ComputerGroupResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data computergroup

	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	computerRequest := computerGroupRequestWithState(data)
	computergroup, _, err := c.client.ComputerGroups.Create(ctx, computerRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create computergroup, got error: %q)", err.Error()))
		return
	}

	tflog.Trace(ctx, "created a computergroup")

	response.Diagnostics.Append(response.State.Set(ctx, computerGroupForState(computergroup))...)

}

func (c *ComputerGroupResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data computergroup
	retryCount := 5
	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	computergroup, resp, err := c.client.ComputerGroups.GetByID(ctx, int(data.Id.ValueInt64()))
	if resp.StatusCode == 404 {
		for resp.StatusCode == 404 && retryCount > 0 {
			time.Sleep(time.Duration(4) * time.Second)
			computergroup, resp, err = c.client.ComputerGroups.GetByID(ctx, int(data.Id.ValueInt64()))
			retryCount = retryCount - 1
		}
	}

	if err != nil {
		response.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read computergroup with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a computergroup")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, computerGroupForState(computergroup))...)
}

func (c *ComputerGroupResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data computergroup
	//retryCount := 5

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	computerGroupUpdateRequest := computerGroupRequestWithState(data)

	computerGroup, _, err := c.client.ComputerGroups.Update(ctx, int(data.Id.ValueInt64()), computerGroupUpdateRequest)
	//if resp.StatusCode == 404 {
	//	for resp.StatusCode == 404 && retryCount > 0 {
	//		time.Sleep(time.Duration(2) * time.Second)
	//		computerGroup, resp, err = c.client.ComputerGroups.Update(ctx, int(data.Id.ValueInt64()), computerGroupUpdateRequest)
	//		retryCount = retryCount - 1
	//	}
	//}

	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update computergroup with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	//tflog.Trace(ctx, "Waiting for computergroup to propagate in Jamf")
	//updatedComputerGroup, _, err := c.client.ComputerGroups.GetByID(ctx, computerGroup.Id)
	//interval := 1
	//for !AreGroupsEquivalent(computerGroup, updatedComputerGroup) {
	//	time.Sleep(time.Duration(interval) * time.Second)
	//	updatedComputerGroup, _, err = c.client.ComputerGroups.GetByID(ctx, computerGroup.Id)
	//	interval = interval * 2
	//}
	//time.Sleep(time.Duration(interval) * time.Second)

	tflog.Trace(ctx, "updated a computergroup")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, computerGroupForState(computerGroup))...)
}

func (c *ComputerGroupResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data computergroup

	diags := request.State.Get(ctx, &data)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.client.ComputerGroups.Delete(ctx, int(data.Id.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete computergroup with ID %d, got error: %s", data.Id.ValueInt64(), err.Error()),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Computer Group")
}

func (c *ComputerGroupResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resourceImportStatePassthroughJamfProID(ctx, "computergroup", request, response)
}

func (c *ComputerGroupResource) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {

}
