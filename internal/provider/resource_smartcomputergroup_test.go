package provider

import (
	"fmt"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSmartComputerGroupResource(t *testing.T) {
	Name := acctest.RandString(12)
	newName := acctest.RandString(12)

	testCriteria := jamfpro.ComputerGroupCriteria{
		Name:         "Application Title",
		Priority:     0,
		AndOr:        "and",
		SearchType:   "is",
		Value:        "Safari.app",
		OpeningParen: false,
		ClosingParen: false,
	}
	resourceName := "jamfpro_smartcomputergroup.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccSmartComputerGroupResourceConfig(Name, testCriteria),
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
				Config: testAccSmartComputerGroupResourceConfig(newName, testCriteria),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", newName),
				),
			},
		},
	})
}

func testAccSmartComputerGroupResourceConfig(name string, criteria jamfpro.ComputerGroupCriteria) string {
	return fmt.Sprintf(`
resource "jamfpro_smartcomputergroup" "test" {
  name     = %q
  criteria = [
	{
		name = %q
		priority = %d
		and_or = %q
		search_type = %q
		value = %q
		opening_paren = %t
		closing_paren = %t
	},
  ]
}`, name, criteria.Name, criteria.Priority, criteria.AndOr, criteria.SearchType, criteria.Value, criteria.OpeningParen, criteria.ClosingParen)
}
