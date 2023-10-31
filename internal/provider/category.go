package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
	"strconv"
)

type category struct {
	Id       types.Int64  `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Priority types.Int64  `tfsdk:"priority"`
}

func categoryForState(c *jamfpro.Category) category {
	parsedIntId, _ := strconv.ParseInt(c.Id, 10, 64)

	return category{
		Id:       types.Int64Value(parsedIntId),
		Name:     types.StringValue(c.Name),
		Priority: types.Int64Value(int64(c.Priority)),
	}
}
