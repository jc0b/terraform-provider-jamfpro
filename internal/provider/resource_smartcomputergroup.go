package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
	"net/http"
	"time"
)

var _ resource.Resource = &SmartComputerGroupResource{}
var _ resource.ResourceWithImportState = &SmartComputerGroupResource{}

func NewSmartComputerGroupResource() resource.Resource {
	return &SmartComputerGroupResource{}
}

type SmartComputerGroupResource struct {
	client *jamfpro.Client
}

func (c SmartComputerGroupResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_smartcomputergroup"

}

func (c SmartComputerGroupResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Represents a Smart Computer Group resource in Jamf Pro",
		MarkdownDescription: "This resource (`jamfpro_smartcomputergroup`) manages Smart Computer Groups in Jamf Pro",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the Smart Computer Group",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the Smart Computer Group",
			},
			"criteria": schema.SetNestedAttribute{
				Required:    true,
				Computed:    false,
				Description: "Represents criteria by which members of a smart group are defined.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional:    true,
							Description: "Represents the name of a criteria to check against",
						},
						"priority": schema.Int64Attribute{
							Optional:    true,
							Description: "Represents this elements position in the order of criteria. Counting starts at 1.",
							Validators: []validator.Int64{
								int64validator.AtLeast(0),
							},
						},
						"and_or": schema.StringAttribute{
							Optional: true,
							Description: "Whether this criteria will be AND or ORed with the previous criteria. " +
								"Possible values are `and` and `or`.",
							Validators: []validator.String{
								stringvalidator.OneOf("and", "or"),
							},
						},
						"search_type": schema.StringAttribute{
							Optional: true,
							Description: "Represents the operator used to assess the relationship between the criteria " +
								"and the value fields.",
							MarkdownDescription: "Represents the operator used to assess the relationship between the " +
								"`name` and the `value` fields. Possible values are: `is`, `is not`, `has`, and `does " +
								"not have`.",
							Validators: []validator.String{
								stringvalidator.OneOf("is", "is not", "has", "does not have"),
							},
						},
						"value": schema.StringAttribute{
							Optional:            true,
							Description:         "Represents the value that the name criteria is checked against.",
							MarkdownDescription: "Represents the value that the `name` criteria is checked against.",
						},
						"opening_paren": schema.BoolAttribute{
							Optional:    true,
							Description: "Represents whether this criteria contains an opening parenthesis.",
						},
						"closing_paren": schema.BoolAttribute{
							Optional:    true,
							Description: "Represents whether this criteria contains a closing parenthesis.",
						},
					},
				},
			},
		},
	}
}

func (c *SmartComputerGroupResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (c *SmartComputerGroupResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data smartcomputergroup

	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	computergroup, _, err := c.client.ComputerGroups.Create(ctx, smartComputerGroupRequestWithState(data))
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create computergroup, got error: %q)", err.Error()))
		return
	}

	tflog.Trace(ctx, "Waiting for smartcomputergroup to propagate in Jamf")
	createdSmartComputerGroup, resp, err := c.client.ComputerGroups.GetByID(ctx, computergroup.Id)
	interval := 1
	for resp.StatusCode != http.StatusOK && !AreGroupsEquivalent(computergroup, createdSmartComputerGroup) {
		time.Sleep(time.Duration(interval) * time.Second)
		createdSmartComputerGroup, resp, err = c.client.ComputerGroups.GetByID(ctx, computergroup.Id)
		interval = interval * 2
	}

	tflog.Trace(ctx, "created a smartcomputergroup")

	response.Diagnostics.Append(response.State.Set(ctx, smartComputerGroupForState(computergroup))...)

}

func (c *SmartComputerGroupResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data smartcomputergroup
	retryCount := 3

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	smartComputerGroup, resp, err := c.client.ComputerGroups.GetByID(ctx, int(data.Id.ValueInt64()))
	if resp.StatusCode == 404 {
		for resp.StatusCode == 404 && retryCount > 0 {
			time.Sleep(time.Duration(2) * time.Second)
			smartComputerGroup, resp, err = c.client.ComputerGroups.GetByID(ctx, int(data.Id.ValueInt64()))
			retryCount = retryCount - 1
		}
	}

	if err != nil {
		response.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read smartcomputergroup with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a smartcomputergroup")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, smartComputerGroupForState(smartComputerGroup))...)
}

func (c *SmartComputerGroupResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data smartcomputergroup
	retryCount := 5
	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	smartComputerGroupUpdateRequest := smartComputerGroupRequestWithState(data)

	smartComputerGroup, resp, err := c.client.ComputerGroups.Update(ctx, int(data.Id.ValueInt64()), smartComputerGroupUpdateRequest)
	if resp.StatusCode == 404 {
		for resp.StatusCode == 404 && retryCount > 0 {
			time.Sleep(time.Duration(2) * time.Second)
			_, resp, err = c.client.ComputerGroups.Update(ctx, int(data.Id.ValueInt64()), smartComputerGroupUpdateRequest)
			retryCount = retryCount - 1
		}
	}

	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update smartcomputergroup with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "Waiting for computergroup to propagate in Jamf")
	updatedSmartComputerGroup, _, err := c.client.ComputerGroups.GetByID(ctx, smartComputerGroup.Id)
	interval := 1
	for !AreGroupsEquivalent(smartComputerGroup, updatedSmartComputerGroup) {
		time.Sleep(time.Duration(interval) * time.Second)
		updatedSmartComputerGroup, _, err = c.client.ComputerGroups.GetByID(ctx, smartComputerGroup.Id)
		interval = interval * 2
	}
	time.Sleep(time.Duration(interval) * time.Second)

	tflog.Trace(ctx, "updated a smartcomputergroup")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, computerGroupForState(smartComputerGroup))...)
}

func (c *SmartComputerGroupResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data smartcomputergroup

	diags := request.State.Get(ctx, &data)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.client.ComputerGroups.Delete(ctx, int(data.Id.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete smartcomputergroup with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Smart Computer Group")
}

func (c *SmartComputerGroupResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resourceImportStatePassthroughJamfProID(ctx, "smartcomputergroup", request, response)
}
