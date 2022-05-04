package synthetics

import (
	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/synthetics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/PixarV/terraform-provider-ritt/internal/tfresource"
)

func statusCanaryState(conn *synthetics.Synthetics, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindCanaryByName(conn, name)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status.State), nil
	}
}
