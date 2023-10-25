package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
