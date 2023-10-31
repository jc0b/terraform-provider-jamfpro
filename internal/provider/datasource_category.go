package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
)

var _ datasource.DataSource = &CategoryDataSource{}

func NewCategoryDataSource() datasource.DataSource {
	return &CategoryDataSource{}
}

type CategoryDataSource struct {
	client *jamfpro.Client
}

func (c *CategoryDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_category"
}

func (c *CategoryDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Allows details of a category to be retrieved by its ID or name.",
		MarkdownDescription: "The data source `jamfpro_category` allows details of a category to be retrieved by its `ID` or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the category.",
				MarkdownDescription: "`ID` of the category.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the category.",
				MarkdownDescription: "`name` of the category.",
				Optional:            true,
			},
			"priority": schema.Int64Attribute{
				Description:         "Priority of the category.",
				MarkdownDescription: "`priority` of the category.",
				Computed:            true,
			},
		},
	}
}

func (c *CategoryDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data category

	// Read Terraform configuration data into the model
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	var jamfCategory *jamfpro.Category
	var err error
	if data.Id.ValueInt64() > 0 {
		jamfCategory, _, err = c.client.Categories.GetByID(ctx, int(data.Id.ValueInt64()))
		if err != nil {
			response.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get category with ID '%d', got error: %s", data.Id.ValueInt64(), err),
			)
		}
	} else {
		jamfCategory, _, err = c.client.Categories.GetByName(ctx, data.Name.ValueString())
		if err != nil {
			response.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get category '%s', got error: %s", data.Name.ValueString(), err),
			)
		}
	}

	if jamfCategory != nil {
		response.Diagnostics.Append(response.State.Set(ctx, categoryForState(jamfCategory))...)
	}
}

func (c *CategoryDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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
