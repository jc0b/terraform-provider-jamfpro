package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputerDataSource(t *testing.T) {
	computerName := "Jacob’s MacBook Air"
	computerSerial := "C02J6KH9Q6LR"
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
data "jamfpro_computer" "test1_by_name" {
  name = %[1]q
}

data "jamfpro_computer" "test2_by_serial" {
  serial_number = %[2]q
}
`, computerName, computerSerial)
}
