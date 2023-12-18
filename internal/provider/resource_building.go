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

var resourceName = "building"

var _ resource.Resource = &BuildingResource{}
var _ resource.ResourceWithImportState = &BuildingResource{}

func NewBuildingResource() resource.Resource {
	return &BuildingResource{}
}

type BuildingResource struct {
	client *jamfpro.Client
}

func (b *BuildingResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

	b.client = client
}

func (b *BuildingResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_" + resourceName
}

func (b *BuildingResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Represents a " + resourceName + " resource in Jamf Pro",
		MarkdownDescription: "This resource (`jamfpro_" + resourceName + "`) manages buildings in Jamf Pro",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the building",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the building",
			},
			"street_address1": schema.StringAttribute{
				Optional:    true,
				Description: "A street address for the building",
			},
			"street_address2": schema.StringAttribute{
				Optional:    true,
				Description: "A second street address for the building",
			},
			"city": schema.StringAttribute{
				Optional:    true,
				Description: "City of the building",
			},
			"state_province": schema.StringAttribute{
				Optional:    true,
				Description: "State/province of the building",
			},
			"zip_postal_code": schema.StringAttribute{
				Optional:    true,
				Description: "ZIP/Postal code of the building",
			},
			"country": schema.StringAttribute{
				Optional:    true,
				Description: "Country of the building",
			},
		},
	}
}

func (b *BuildingResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data building

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	buildingCreateRequest := &jamfpro.BuildingCreateRequest{
		Name:           data.Name.ValueString(),
		StreetAddress1: data.StreetAddress1.ValueString(),
		StreetAddress2: data.StreetAddress2.ValueString(),
		City:           data.City.ValueString(),
		StateProvince:  data.StateProvince.ValueString(),
		ZipPostalCode:  data.ZipPostalCode.ValueString(),
		Country:        data.Country.ValueString(),
	}
	building, _, err := b.client.Buildings.Create(ctx, buildingCreateRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create %s, got error: %s", resourceName, err),
		)
		return
	}

	tflog.Trace(ctx, "created a building")

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, buildingForState(building))...)
}

func (b *BuildingResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data building

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	building, _, err := b.client.Buildings.GetByID(ctx, int(data.Id.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s with ID %d, got error: %s", resourceName, data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a building")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, buildingForState(building))...)
}

func (b *BuildingResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data building

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	buildingUpdateRequest := &jamfpro.BuildingUpdateRequest{
		Name:           data.Name.ValueString(),
		StreetAddress1: data.StreetAddress1.ValueString(),
		StreetAddress2: data.StreetAddress2.ValueString(),
		City:           data.City.ValueString(),
		StateProvince:  data.StateProvince.ValueString(),
		ZipPostalCode:  data.ZipPostalCode.ValueString(),
		Country:        data.Country.ValueString(),
	}
	building, _, err := b.client.Buildings.Update(ctx, int(data.Id.ValueInt64()), buildingUpdateRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update %s with ID %d, got error: %s", resourceName, data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a building")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, buildingForState(building))...)
}

func (b *BuildingResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data building

	diags := request.State.Get(ctx, &data)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	_, err := b.client.Buildings.Delete(ctx, int(data.Id.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete %s with ID %d, got error: %s", resourceName, data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a building")
}

func (b *BuildingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughJamfProID(ctx, resourceName, req, resp)
}
