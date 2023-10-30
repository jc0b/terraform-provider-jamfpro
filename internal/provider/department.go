package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
	"strconv"
)

type department struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func departmentForState(b *jamfpro.Department) department {
	parsedIntId, _ := strconv.ParseInt(b.Id, 10, 64)

	return department{
		Id:   types.Int64Value(parsedIntId),
		Name: types.StringValue(b.Name),
	}
}
