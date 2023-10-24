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

var _ resource.Resource = &BuildingResource{}
var _ resource.ResourceWithImportState = &BuildingResource{}

func NewBuildingResource() resource.Resource {
	return &BuildingResource{}
}

type BuildingResource struct {
	client *jamfpro.Client
}

func (b BuildingResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (b BuildingResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_building"
}

func (b BuildingResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Represents a building resource in Jamf Pro",
		MarkdownDescription: "This resource (`jamfpro_building" + "`) manages buildings in Jamf Pro",

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

func (b BuildingResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
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
			fmt.Sprintf("Unable to create building, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a tag")

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, buildingForState(building))...)
}

func (b BuildingResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (b BuildingResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (b BuildingResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}
