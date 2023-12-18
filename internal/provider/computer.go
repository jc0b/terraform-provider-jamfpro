package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
)

type computer struct {
	Id           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	SerialNumber types.String `tfsdk:"serial_number"`
}

func computerForState(c *jamfpro.Computer) computer {
	return computer{
		Id:           types.Int64Value(int64(c.Id)),
		Name:         types.StringValue(c.Name),
		SerialNumber: types.StringValue(c.SerialNumber),
	}
}
