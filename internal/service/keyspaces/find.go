package keyspaces

import (
	"context"

	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/keyspaces"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/PixarV/terraform-provider-ritt/internal/tfresource"
)

func FindKeyspaceByName(ctx context.Context, conn *keyspaces.Keyspaces, name string) (*keyspaces.GetKeyspaceOutput, error) {
	input := keyspaces.GetKeyspaceInput{
		KeyspaceName: aws.String(name),
	}

	output, err := conn.GetKeyspaceWithContext(ctx, &input)

	if tfawserr.ErrCodeEquals(err, keyspaces.ErrCodeResourceNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}
