package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"jamfpro": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if !isClientIdSet() {
		t.Fatal("JAMF_CLIENT_ID environment variable must be set for acceptance tests")
	}
	if !isClientSecretSet() {
		t.Fatal("JAMF_CLIENT_SECRET environment variable must be set for acceptance tests")
	}
	if !isInstanceURLSet() {
		t.Fatal("JAMF_INSTANCE_URL environment variable must be set for acceptance tests")
	}
}

func isClientIdSet() bool {
	if os.Getenv("JAMF_CLIENT_ID") != "" {
		return true
	}
	return false
}

func isClientSecretSet() bool {
	if os.Getenv("JAMF_CLIENT_SECRET") != "" {
		return true
	}
	return false
}

func isInstanceURLSet() bool {
	if os.Getenv("JAMF_INSTANCE_URL") != "" {
		return true
	}
	return false
}
