// Code generated by "internal/generate/listpages/main.go -ListOps=ListApplications"; DO NOT EDIT.

package kinesisanalyticsv2

import (
	"context"

	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/kinesisanalyticsv2"
)

func listApplicationsPages(conn *kinesisanalyticsv2.KinesisAnalyticsV2, input *kinesisanalyticsv2.ListApplicationsInput, fn func(*kinesisanalyticsv2.ListApplicationsOutput, bool) bool) error {
	return listApplicationsPagesWithContext(context.Background(), conn, input, fn)
}

func listApplicationsPagesWithContext(ctx context.Context, conn *kinesisanalyticsv2.KinesisAnalyticsV2, input *kinesisanalyticsv2.ListApplicationsInput, fn func(*kinesisanalyticsv2.ListApplicationsOutput, bool) bool) error {
	for {
		output, err := conn.ListApplicationsWithContext(ctx, input)
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
