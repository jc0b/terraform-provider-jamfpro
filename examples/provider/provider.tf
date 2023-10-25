terraform {
  required_providers {
    jamfpro = {
      source = "jc0b/jamfpro"
    }
  }
}

// configure the provider
provider "jamfpro" {
  // Base URL of your Jamf Pro instance.
  // The `JAMF_INSTANCE_URL` environment variable can be used instead.
  instance_url = "https://jc0b.jamfcloud.com"

  // Jamf Pro API Client ID.
  // This is a secret, it must be managed using a variable.
  // The `JAMF_CLIENT_ID` environment variable can be used instead.
  client_id = var.client_id

  // Jamf Pro API Client Secret.
  // This is a secret, it must be managed using a variable.
  // The `JAMF_CLIENT_SECRET` environment variable can be used instead.
  client_secret = var.client_secret
}