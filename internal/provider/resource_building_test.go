package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBuildingResource(t *testing.T) {
	Name := acctest.RandString(12)
	StreetAddress1 := acctest.RandString(12)
	StreetAddress2 := acctest.RandString(12)
	City := acctest.RandString(12)
	StateProvince := acctest.RandString(12)
	ZipPostalCode := acctest.RandString(12)
	Country := acctest.RandString(12)
	newCity := acctest.RandString(12)
	newCountry := acctest.RandString(12)
	resourceName := "jamfpro_building.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccBuildingResourceConfig(Name, StreetAddress1, StreetAddress2, City, StateProvince, ZipPostalCode, Country),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "name", Name),
					resource.TestCheckResourceAttr(
						resourceName, "street_address1", StreetAddress1),
					resource.TestCheckResourceAttr(
						resourceName, "street_address2", StreetAddress2),
					resource.TestCheckResourceAttr(
						resourceName, "city", City),
					resource.TestCheckResourceAttr(
						resourceName, "state_province", StateProvince),
					resource.TestCheckResourceAttr(
						resourceName, "zip_postal_code", ZipPostalCode),
					resource.TestCheckResourceAttr(
						resourceName, "country", Country),
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
				Config: testAccBuildingResourceConfig(Name, StreetAddress1, StreetAddress2, newCity, StateProvince, ZipPostalCode, newCountry),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "city", newCity),
					resource.TestCheckResourceAttr(
						resourceName, "country", newCountry),
				),
			},
		},
	})
}

func testAccBuildingResourceConfig(name string, streetAddress1 string, streetAddress2 string, city string, stateProvince string, zipPostalCode string, country string) string {
	return fmt.Sprintf(`
resource "jamfpro_building" "test" {
  name            = %[1]q
  street_address1 = %[2]q
  street_address2 = %[3]q
  city            = %[4]q
  state_province  = %[5]q
  zip_postal_code = %[6]q
  country         = %[7]q
}
`, name, streetAddress1, streetAddress2, city, stateProvince, zipPostalCode, country)
}
