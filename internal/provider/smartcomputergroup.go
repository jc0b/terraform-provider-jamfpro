package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
)

type smartcomputergroup struct {
	Id       types.Int64  `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Criteria types.Set    `tfsdk:"criteria"`
}

var criteriaAttrTypes = map[string]attr.Type{
	"name":          types.StringType,
	"priority":      types.Int64Type,
	"and_or":        types.StringType,
	"search_type":   types.StringType,
	"value":         types.StringType,
	"opening_paren": types.BoolType,
	"closing_paren": types.BoolType,
}

func smartComputerGroupForState(c *jamfpro.ComputerGroup) smartcomputergroup {
	criteria := make([]attr.Value, 0)
	for _, criterion := range c.Criteria {
		criteria = append(
			criteria,
			types.ObjectValueMust(
				criteriaAttrTypes,
				map[string]attr.Value{
					"name":          types.StringValue(criterion.Name),
					"priority":      types.Int64Value(int64(criterion.Priority)),
					"and_or":        types.StringValue(criterion.AndOr),
					"search_type":   types.StringValue(criterion.SearchType),
					"value":         types.StringValue(criterion.Value),
					"opening_paren": types.BoolValue(criterion.OpeningParen),
					"closing_paren": types.BoolValue(criterion.ClosingParen),
				},
			),
		)
	}
	return smartcomputergroup{
		Id:       types.Int64Value(int64(c.Id)),
		Name:     types.StringValue(c.Name),
		Criteria: types.SetValueMust(types.ObjectType{AttrTypes: criteriaAttrTypes}, criteria),
	}
}

func smartComputerGroupRequestWithState(data smartcomputergroup) *jamfpro.ComputerGroupRequest {
	criteria := make([]jamfpro.ComputerGroupCriteria, 0)
	for _, criterion := range data.Criteria.Elements() {
		criterionMap := criterion.(types.Object).Attributes()
		if criterionMap != nil {
			criteria = append(
				criteria,
				jamfpro.ComputerGroupCriteria{
					Name:         criterionMap["name"].(types.String).ValueString(),
					Priority:     int(criterionMap["priority"].(types.Int64).ValueInt64()),
					AndOr:        criterionMap["and_or"].(types.String).ValueString(),
					SearchType:   criterionMap["search_type"].(types.String).ValueString(),
					Value:        criterionMap["value"].(types.String).ValueString(),
					OpeningParen: criterionMap["opening_paren"].(types.Bool).ValueBool(),
					ClosingParen: criterionMap["closing_paren"].(types.Bool).ValueBool(),
				})
		}
	}
	return &jamfpro.ComputerGroupRequest{
		Name:     data.Name.ValueString(),
		Criteria: criteria,
	}
}
