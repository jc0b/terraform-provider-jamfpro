package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputerGroupResource(t *testing.T) {
	Name := acctest.RandString(12)
	newName := acctest.RandString(12)
	serialNumber := randomSerialNumber()

	resourceName := "jamfpro_computergroup.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccComputerGroupResourceConfig(serialNumber, Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", Name),
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
				Config: testAccComputerGroupResourceConfig(serialNumber, newName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", newName),
				),
			},
		},
	})
}

func testAccComputerGroupResourceConfig(serial_number string, name string) string {
	return fmt.Sprintf(`
resource "jamfpro_computer" "test_computer" {
  name 			= "Static Computergroup Test Mac"
  serial_number = %[1]q
}

data "jamfpro_computer" "test_computer" {
  serial_number = %[1]q
  depends_on = [jamfpro_computer.test_computer]
}

resource "jamfpro_computergroup" "test" {
  name     = %[2]q
  computers = [jamfpro_computer.test_computer]
}`, serial_number, name)
}
