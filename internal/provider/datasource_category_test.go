package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCategoryDataSource(t *testing.T) {
	categoryName := acctest.RandString(12)
	priority1 := acctest.RandIntRange(1, 20)
	category2Name := acctest.RandString(12)
	priority2 := acctest.RandIntRange(1, 20)
	c1ResourceName := "jamfpro_category.test1"
	c2ResourceName := "jamfpro_category.test2"
	c1DataSourceName := "data.jamfpro_category.test1_by_name"
	c2DataSourceName := "data.jamfpro_category.test2_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagDataSourceConfig(categoryName, category2Name, priority1, priority2),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Read by name, no taxonomy, default color
					resource.TestCheckResourceAttrPair(
						c1DataSourceName, "id", c1ResourceName, "id"),
					resource.TestCheckResourceAttr(
						c1DataSourceName, "name", categoryName),
					resource.TestCheckResourceAttr(
						c1DataSourceName, "priority", strconv.FormatInt(int64(priority1), 10)),
					resource.TestCheckResourceAttrPair(
						c2DataSourceName, "id", c2ResourceName, "id"),
					resource.TestCheckResourceAttr(
						c2DataSourceName, "name", category2Name),
					resource.TestCheckResourceAttr(
						c2DataSourceName, "priority", strconv.FormatInt(int64(priority2), 10)),
				),
			},
		},
	})
}

func testAccTagDataSourceConfig(categoryOne string, categoryTwo string, priorityOne int, priorityTwo int) string {
	return fmt.Sprintf(`
resource "jamfpro_category" "test1" {
  name     = %q
  priority = %d
}

resource "jamfpro_category" "test2" {
  name     = %q
  priority = %d
}

data "jamfpro_category" "test1_by_name" {
  name = jamfpro_category.test1.name
}

data "jamfpro_category" "test2_by_id" {
  id = jamfpro_category.test2.id
}
`, categoryOne, priorityOne, categoryTwo, priorityTwo)
}
