package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
	"strconv"
)

type apirole struct {
	Id         types.Int64  `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Privileges types.Set    `tfsdk:"privileges"`
}

func apiRoleForState(a *jamfpro.ApiRole) apirole {
	parsedInt, _ := strconv.ParseInt(*a.Id, 10, 64)

	privileges := make([]attr.Value, 0)
	for _, pv := range *a.Privileges {
		privileges = append(privileges, types.StringValue(pv))
	}

	return apirole{
		Id:         types.Int64Value(parsedInt),
		Name:       types.StringValue(*a.DisplayName),
		Privileges: types.SetValueMust(types.StringType, privileges),
	}
}
