package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jc0b/go-jamfpro-api/jamfpro"
	"strconv"
)

func resourceImportStatePassthroughJamfProID(ctx context.Context, name string, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	jamfProId, err := strconv.ParseInt(request.ID, 10, 64)
	if err != nil {
		response.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Jamf Pro %s ID must be an integer", name),
		)
	} else {
		response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(jamfProId))...)
	}
}

func AreGroupsEquivalent(planned, actual *jamfpro.ComputerGroup) bool {
	if actual == nil {
		return false
	}

	if planned.Name != actual.Name {
		return false
	}
	if planned.Id != actual.Id {
		return false
	}
	for i, v := range planned.Computers {
		if v != actual.Computers[i] {
			return false
		}
	}
	for i, v := range planned.Criteria {
		if v != actual.Criteria[i] {
			return false
		}
	}

	return true
}
