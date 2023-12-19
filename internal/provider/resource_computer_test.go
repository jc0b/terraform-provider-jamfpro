package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputerResource(t *testing.T) {
	Name := acctest.RandString(12)
	newName := acctest.RandString(12)
	serialNumber := randomSerialNumber()
	newSerialNumber := randomSerialNumber()
	resourceName := "jamfpro_computer.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccComputerResourceConfig(Name, serialNumber),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", Name),
					resource.TestCheckResourceAttr(
						resourceName, "serial_number", serialNumber),
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
				Config: testAccComputerResourceConfig(newName, newSerialNumber),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", newName),
					resource.TestCheckResourceAttr(
						resourceName, "serial_number", newSerialNumber),
				),
			},
		},
	})
}

func testAccComputerResourceConfig(name string, serial_number string) string {
	return fmt.Sprintf(`
resource "jamfpro_computer" "test" {
  name     		= %q
  serial_number = %q
}
`, name, serial_number)
}
