package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
	"log"
	"os"
	"time"
)

var providerConfigurationError = "Jamf Pro provider configuration error"
var localSessionStorageFile = "/tmp/jamf-tf-session"

var _ provider.Provider = &JamfProProvider{}

type JamfProProvider struct {
	version string
}

type JamfProProviderModel struct {
	InstanceURL  types.String `tfsdk:"instance_url"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (j JamfProProvider) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "jamfpro"
	response.Version = j.version
}

func (j JamfProProvider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"instance_url": schema.StringAttribute{
				Optional:    true,
				Sensitive:   false,
				Description: "The url of your Jamf Pro instance.",
				MarkdownDescription: "The url of your Jamf Pro instance (e.g. myinstance.jamfcloud.com)." +
					"Can also be set with the `JAMF_INSTANCE_URL` environment variable.",
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "A Jamf Pro API Client ID.",
				MarkdownDescription: "The Client ID of an API Client. Can also be set with the `JAMF_CLIENT_ID`" +
					"environment variable. Must be used in conjunction with a matching Client Secret.",
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "A Jamf Pro API Client Secret.",
				MarkdownDescription: "The Client Secret of an API Client. Can also be set with the `JAMF_CLIENT_SECRET`" +
					"environment variable. Must be used in conjunction with a matching Client ID.",
			},
		},
	}
}

func (j JamfProProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	fmt.Println("Configuring Jamf Pro Provider...")
	var data JamfProProviderModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	// Instance URL
	var InstanceURL string
	if data.InstanceURL.IsUnknown() {
		response.Diagnostics.AddWarning(
			providerConfigurationError,
			"Cannot use unknown value as Instance URL",
		)
		return
	}

	if data.InstanceURL.IsNull() {
		InstanceURL = os.Getenv("JAMF_INSTANCE_URL")
	} else {
		InstanceURL = data.InstanceURL.ValueString()
	}

	if InstanceURL == "" {
		response.Diagnostics.AddError(
			providerConfigurationError,
			"Instance URL cannot be an empty string",
		)
		return
	}

	// Client ID & Secret
	var clientId string
	var clientSecret string
	if data.ClientID.IsUnknown() {
		response.Diagnostics.AddWarning(
			providerConfigurationError,
			"Cannot use unknown value as Client ID",
		)
		return
	}

	if data.ClientSecret.IsUnknown() {
		response.Diagnostics.AddWarning(
			providerConfigurationError,
			"Cannot use unknown value as Client Secret",
		)
		return
	}

	if data.ClientID.IsNull() {
		clientId = os.Getenv("JAMF_CLIENT_ID")
	} else {
		clientId = data.ClientID.ValueString()
	}

	if data.ClientSecret.IsNull() {
		tflog.Info(ctx, "Client Secret is NULL")
		clientSecret = os.Getenv("JAMF_CLIENT_SECRET")
	} else {
		clientSecret = data.ClientSecret.ValueString()
	}

	var apiClient = (clientId != "") == (clientSecret != "")

	if !apiClient {
		response.Diagnostics.AddError(
			providerConfigurationError,
			"You must supply API Client credentials to authenticate.")
	}

	userAgent := fmt.Sprintf("terraform-provider-jamfpro/%s", j.version)

	jamfSession := ""
	token, err := os.ReadFile(localSessionStorageFile)

	// if the file is populated, read it and pass it
	// else, its empty string

	if err == nil {
		info, err := os.Stat(localSessionStorageFile)
		if info.ModTime().Before(time.Now().Add(time.Minute*30)) && err == nil {
			jamfSession = string(token)
		}
	}

	var c *jamfpro.Client

	c, err = jamfpro.NewClient(clientId, clientSecret, InstanceURL, jamfSession)
	if err != nil {
		response.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create OAuth client:\n\n"+err.Error())
	}

	if c.GetSessionToken() != jamfSession {
		f, err := os.Create(localSessionStorageFile)
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.WriteString(c.GetSessionToken() + "\n")
		if err != nil {
			log.Fatal(err)
		}
		err = f.Sync()
		if err != nil {
			return
		}
	}

	c.ExtraHeader["User-Agent"] = userAgent

	response.DataSourceData = c
	response.ResourceData = c
}

func (j JamfProProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCategoryDataSource,
		NewComputerDataSource,
	}
}

func (j JamfProProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApiRoleResource,
		NewBuildingResource,
		NewCategoryResource,
		NewComputerGroupResource,
		NewComputerResource,
		NewDepartmentResource,
		NewSmartComputerGroupResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JamfProProvider{
			version: version,
		}
	}
}
