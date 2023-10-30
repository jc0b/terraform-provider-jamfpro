package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCategoryResource(t *testing.T) {
	Name := acctest.RandString(12)
	newName := acctest.RandString(12)
	Priority := acctest.RandIntRange(1, 20)
	newPriority := acctest.RandIntRange(1, 20)
	resourceName := "jamfpro_category.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccCategoryResourceConfig(Name, Priority),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", Name),
					resource.TestCheckResourceAttr(
						resourceName, "priority", strconv.FormatInt(int64(Priority), 10)),
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
				Config: testAccCategoryResourceConfig(newName, newPriority),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", newName),
					resource.TestCheckResourceAttr(
						resourceName, "priority", strconv.FormatInt(int64(newPriority), 10)),
				),
			},
		},
	})
}

func testAccCategoryResourceConfig(name string, priority int) string {
	return fmt.Sprintf(`
resource "jamfpro_category" "test" {
  name     = %q
  priority = %d
}
`, name, priority)
}
