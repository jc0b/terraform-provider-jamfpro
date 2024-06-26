---
page_title: "jamfpro_building Resource - terraform-provider-jamfpro"
description: |-
  This resource (`jamfpro_building`) manages buildings in Jamf Pro
---

# jamfpro_building (Resource)
This resource (`jamfpro_building`) manages buildings in Jamf Pro

## Example Usage
```terraform
resource "jamfpro_building" "30_rock" {
    city            = "New York"
    country         = "United States of America"
    state_province  = "New York"
    street_address1 = "30 Rockefeller Plaza"
    zip_postal_code = "NY 10112"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the building

### Optional

- `city` (String) City of the building
- `country` (String) Country of the building
- `state_province` (String) State/province of the building
- `street_address1` (String) A street address for the building
- `street_address2` (String) A second street address for the building
- `zip_postal_code` (String) ZIP/Postal code of the building

### Read-Only

- `id` (Number) ID of the building