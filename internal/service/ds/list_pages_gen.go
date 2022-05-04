// Code generated by "internal/generate/listpages/main.go -ListOps=DescribeDirectories"; DO NOT EDIT.

package ds

import (
	"context"

	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/directoryservice"
)

func describeDirectoriesPages(conn *directoryservice.DirectoryService, input *directoryservice.DescribeDirectoriesInput, fn func(*directoryservice.DescribeDirectoriesOutput, bool) bool) error {
	return describeDirectoriesPagesWithContext(context.Background(), conn, input, fn)
}

func describeDirectoriesPagesWithContext(ctx context.Context, conn *directoryservice.DirectoryService, input *directoryservice.DescribeDirectoriesInput, fn func(*directoryservice.DescribeDirectoriesOutput, bool) bool) error {
	for {
		output, err := conn.DescribeDirectoriesWithContext(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.StringValue(output.NextToken) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextToken = output.NextToken
	}
	return nil
}
