package provider

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApiRoleResource(t *testing.T) {
	Name := acctest.RandString(12)
	NewName := acctest.RandString(12)
	Privileges := []string{"Create Packages", "Read Static Mobile Device Groups", "Read eBooks"}
	NewPrivileges := []string{"Read Teacher App Settings", "Read SMTP Server", "Read PKI", "Read iBeacon"}
	resourceName := "jamfpro_api_role.test_role"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccApiRoleResourceConfig(Name, Privileges),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", Name),
					resource.TestCheckResourceAttr(
						resourceName, "privileges.#", strconv.Itoa(len(Privileges))),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "privileges.*", Privileges[0]),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "privileges.*", Privileges[1]),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "privileges.*", Privileges[2]),
				),
			},
			// ImportState
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccApiRoleResourceConfig(NewName, NewPrivileges),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", NewName),
					resource.TestCheckResourceAttr(
						resourceName, "privileges.#", strconv.Itoa(len(NewPrivileges))),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "privileges.*", NewPrivileges[0]),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "privileges.*", NewPrivileges[1]),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "privileges.*", NewPrivileges[2]),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "privileges.*", NewPrivileges[3]),
				),
			},
		},
	})
}

func testAccApiRoleResourceConfig(name string, privileges []string) string {
	b, _ := json.Marshal(privileges)
	return fmt.Sprintf(`
resource "jamfpro_api_role" "test_role" {
  name       = %q
  privileges = %+v
}
`, name, string(b))
}
