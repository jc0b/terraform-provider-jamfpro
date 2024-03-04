package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
)

var _ resource.Resource = &ApiRoleResource{}
var _ resource.ResourceWithImportState = &ApiRoleResource{}

func NewApiRoleResource() resource.Resource {
	return &ApiRoleResource{}
}

type ApiRoleResource struct {
	client *jamfpro.Client
}

func (a *ApiRoleResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*jamfpro.Client)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure type",
			fmt.Sprintf("Expected *jamfpro.Client, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	a.client = client
}

func (a *ApiRoleResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_api_role"
}

func (a *ApiRoleResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Represents an API role resource in Jamf Pro",
		MarkdownDescription: "This resource (`jamfpro_api_role`) manages API roles in Jamf Pro",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the API role",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the API role",
			},
			"privileges": schema.SetAttribute{
				Description: "The privileges granted to the API role",
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

func (a *ApiRoleResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data apirole

	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	privileges := make([]string, 0)
	for _, priv := range data.Privileges.Elements() { // nil if null or unknown → no iterations
		privileges = append(privileges, priv.(types.String).ValueString())
	}

	apiRoleCreateRequest := &jamfpro.ApiRoleCreateRequest{
		DisplayName: data.Name.ValueString(),
		Privileges:  privileges,
	}

	apirole, _, err := a.client.ApiRoles.Create(ctx, apiRoleCreateRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create API role, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created an API role")

	response.Diagnostics.Append(response.State.Set(ctx, apiRoleForState(apirole))...)
}

func (a *ApiRoleResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data apirole

	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	apirole, _, err := a.client.ApiRoles.GetByID(ctx, int(data.Id.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read API role with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "Read an API role")

	response.Diagnostics.Append(response.State.Set(ctx, apiRoleForState(apirole))...)
}

func (a *ApiRoleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data apirole

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	privileges := make([]string, 0)
	for _, priv := range data.Privileges.Elements() { // nil if null or unknown → no iterations
		privileges = append(privileges, priv.(types.String).ValueString())
	}

	apiRoleUpdateRequest := &jamfpro.ApiRoleUpdateRequest{
		DisplayName: data.Name.ValueString(),
		Privileges:  privileges,
	}
	apirole, _, err := a.client.ApiRoles.Update(ctx, int(data.Id.ValueInt64()), apiRoleUpdateRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update API role with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated an API role")

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, apiRoleForState(apirole))...)
}

func (a *ApiRoleResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data apirole

	diags := request.State.Get(ctx, &data)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	_, err := a.client.ApiRoles.Delete(ctx, int(data.Id.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete API role with ID %d, got error: %s", data.Id.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted an API role")
}

func (a *ApiRoleResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resourceImportStatePassthroughJamfProID(ctx, "api_role", request, response)
}
