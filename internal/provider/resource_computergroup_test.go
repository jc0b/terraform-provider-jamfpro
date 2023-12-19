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
				Config: testAccComputerGroupResourceConfig(Name, serialNumber),
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
				Config: testAccComputerGroupResourceConfig(newName, serialNumber),
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
  name 			= "Test Mac"
  serial_number = %q
}

resource "jamfpro_computergroup" "test" {
  name     = %q
  computers = [resource.jamfpro_computer.test_computer]
}`, serial_number, name)
}

//
//{
//id = data.jamfpro_computer.test_computer.id
//name = data.jamfpro_computer.test_computer.name
//serial_number = data.jamfpro_computer.test_computer.serial_number
//},
