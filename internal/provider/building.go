package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
	"strconv"
)

type building struct {
	Id             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	StreetAddress1 types.String `tfsdk:"streetAddress1"`
	StreetAddress2 types.String `tfsdk:"streetAddress1"`
	City           types.String `tfsdk:"city"`
	StateProvince  types.String `tfsdk:"stateProvince"`
	ZipPostalCode  types.String `tfsdk:"zipPostalCode"`
	Country        types.String `tfsdk:"country"`
}

func buildingForState(b *jamfpro.Building) building {
	parsedInt, _ := strconv.ParseInt(*b.Id, 10, 64)

	return building{
		Id:             types.Int64Value(parsedInt),
		Name:           types.StringValue(*b.Name),
		StreetAddress1: types.StringValue(*b.StreetAddress1),
		StreetAddress2: types.StringValue(*b.StreetAddress2),
		City:           types.StringValue(*b.City),
		StateProvince:  types.StringValue(*b.StateProvince),
		ZipPostalCode:  types.StringValue(*b.ZipPostalCode),
		Country:        types.StringValue(*b.Country),
	}
}
