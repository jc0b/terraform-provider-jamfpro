package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
)

type computergroup struct {
	Id        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Computers types.Set    `tfsdk:"computers"`
}

var computerAttrTypes = map[string]attr.Type{
	"id":            types.Int64Type,
	"name":          types.StringType,
	"serial_number": types.StringType,
}

func computerGroupForState(c *jamfpro.ComputerGroup) computergroup {
	computers := make([]attr.Value, 0)
	for _, machine := range c.Computers {
		computers = append(
			computers,
			types.ObjectValueMust(
				computerAttrTypes,
				map[string]attr.Value{
					"id":            types.Int64Value(int64(machine.Id)),
					"name":          types.StringValue(machine.Name),
					"serial_number": types.StringValue(machine.SerialNumber),
				},
			),
		)
	}
	return computergroup{
		Id:        types.Int64Value(int64(c.Id)),
		Name:      types.StringValue(c.Name),
		Computers: types.SetValueMust(types.ObjectType{AttrTypes: computerAttrTypes}, computers),
	}
}

func computerGroupRequestWithState(data computergroup) *jamfpro.ComputerGroupRequest {
	computers := make([]jamfpro.Computer, 0)
	for _, machine := range data.Computers.Elements() {
		machineMap := machine.(types.Object).Attributes()
		if machineMap != nil {
			computers = append(
				computers,
				jamfpro.Computer{
					Id:           int(machineMap["id"].(types.Int64).ValueInt64()),
					Name:         machineMap["name"].(types.String).ValueString(),
					SerialNumber: machineMap["serial_number"].(types.String).ValueString(),
				})
		}
	}
	return &jamfpro.ComputerGroupRequest{
		Name:      data.Name.ValueString(),
		Computers: computers,
	}
}
