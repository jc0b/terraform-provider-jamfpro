package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
)

var _ datasource.DataSource = &ComputerDataSource{}

func NewComputerDataSource() datasource.DataSource {
	return &ComputerDataSource{}
}

type ComputerDataSource struct {
	client *jamfpro.Client
}

func (c *ComputerDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_computer"
}

func (c *ComputerDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Allows details of a computer to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `jamfpro_computer` allows details of a computer to be retrieved by its `ID`, name, or serial number.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the computer.",
				MarkdownDescription: "`ID` of the computer.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the computer.",
				MarkdownDescription: "`name` of the computer.",
				Optional:            true,
				Computed:            true,
			},
			"serial_number": schema.StringAttribute{
				Description:         "Serial number of the computer.",
				MarkdownDescription: "`serial_number` of the computer.",
				Optional:            true,
				Computed:            true,
			},
			"udid": schema.StringAttribute{
				Description:         "Hardware UDID of the computer.",
				MarkdownDescription: "`udid` of the computer.",
				Computed:            true,
			},
		},
	}
}

func (c *ComputerDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data computer

	// Read Terraform configuration data into the model
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	var jamfComputer *jamfpro.Computer
	var err error
	if !data.Id.IsNull() && data.Id.ValueInt64() != 0 {
		jamfComputer, _, err = c.client.Computers.GetByID(ctx, int(data.Id.ValueInt64()))
		if err != nil {
			response.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get computer with ID '%d', got error: %s", data.Id.ValueInt64(), err),
			)
		}
	} else if data.SerialNumber.ValueString() != "" {
		jamfComputer, _, err = c.client.Computers.GetBySerialNumber(ctx, data.SerialNumber.ValueString())
		if err != nil {
			response.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get computer with serial number '%s', got error: %s", data.SerialNumber.ValueString(), err),
			)
		}
	} else {
		jamfComputer, _, err = c.client.Computers.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			response.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get computer '%s', got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if jamfComputer != nil {
		response.Diagnostics.Append(response.State.Set(ctx, computerForState(jamfComputer))...)
	}

	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ComputerDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*jamfpro.Client)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jamfpro.Client, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	c.client = client
}

func (c *ComputerDataSource) ValidateConfig(ctx context.Context, request datasource.ValidateConfigRequest, response *datasource.ValidateConfigResponse) {
	var data computer
	diags := request.Config.Get(ctx, &data)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if data.Id.IsNull() && data.Name.IsNull() && data.SerialNumber.IsNull() {
		response.Diagnostics.AddError("Invalid `jamfpro_computer` data source", "`id`, `name`, or `serial_number` missing. At least one is required in order to create the data source.")
	}
}
