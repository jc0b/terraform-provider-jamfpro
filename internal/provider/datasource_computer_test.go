package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputerDataSource(t *testing.T) {
	computerName := "Datasource Test Computer"
	computerSerial := randomSerialNumber()
	c1DataSourceName := "data.jamfpro_computer.test1_by_name"
	c2DataSourceName := "data.jamfpro_computer.test2_by_serial"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComputerDataSourceConfig(computerName, computerSerial),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name, no taxonomy, default color
					resource.TestCheckResourceAttr(
						c1DataSourceName, "name", computerName),
					resource.TestCheckResourceAttr(
						c1DataSourceName, "serial_number", computerSerial),
					resource.TestCheckResourceAttr(
						c2DataSourceName, "name", computerName),
					resource.TestCheckResourceAttr(
						c2DataSourceName, "serial_number", computerSerial),
				),
			},
		},
	})
}

func testAccComputerDataSourceConfig(computerName string, computerSerial string) string {
	return fmt.Sprintf(`
resource "jamfpro_computer" "test1" {
  name     		= %[1]q
  serial_number = %[2]q
}

resource "jamfpro_computer" "test2" {
  name     		= %[1]q
  serial_number = %[2]q
}

data "jamfpro_computer" "test1_by_name" {
  name = jamfpro_computer.test1.name
}

data "jamfpro_computer" "test2_by_serial" {
  serial_number = jamfpro_computer.test2.serial_number
}
`, computerName, computerSerial)
}
