// Code generated by "internal/generate/listpages/main.go -ListOps=ListIPSets,ListRegexPatternSets,ListRuleGroups,ListWebACLs -Paginator=NextMarker"; DO NOT EDIT.

package wafv2

import (
	"context"

	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/wafv2"
)

func listIPSetsPages(conn *wafv2.WAFV2, input *wafv2.ListIPSetsInput, fn func(*wafv2.ListIPSetsOutput, bool) bool) error {
	return listIPSetsPagesWithContext(context.Background(), conn, input, fn)
}

func listIPSetsPagesWithContext(ctx context.Context, conn *wafv2.WAFV2, input *wafv2.ListIPSetsInput, fn func(*wafv2.ListIPSetsOutput, bool) bool) error {
	for {
		output, err := conn.ListIPSetsWithContext(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.StringValue(output.NextMarker) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextMarker = output.NextMarker
	}
	return nil
}

func listRegexPatternSetsPages(conn *wafv2.WAFV2, input *wafv2.ListRegexPatternSetsInput, fn func(*wafv2.ListRegexPatternSetsOutput, bool) bool) error {
	return listRegexPatternSetsPagesWithContext(context.Background(), conn, input, fn)
}

func listRegexPatternSetsPagesWithContext(ctx context.Context, conn *wafv2.WAFV2, input *wafv2.ListRegexPatternSetsInput, fn func(*wafv2.ListRegexPatternSetsOutput, bool) bool) error {
	for {
		output, err := conn.ListRegexPatternSetsWithContext(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.StringValue(output.NextMarker) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextMarker = output.NextMarker
	}
	return nil
}

func listRuleGroupsPages(conn *wafv2.WAFV2, input *wafv2.ListRuleGroupsInput, fn func(*wafv2.ListRuleGroupsOutput, bool) bool) error {
	return listRuleGroupsPagesWithContext(context.Background(), conn, input, fn)
}

func listRuleGroupsPagesWithContext(ctx context.Context, conn *wafv2.WAFV2, input *wafv2.ListRuleGroupsInput, fn func(*wafv2.ListRuleGroupsOutput, bool) bool) error {
	for {
		output, err := conn.ListRuleGroupsWithContext(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.StringValue(output.NextMarker) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextMarker = output.NextMarker
	}
	return nil
}

func listWebACLsPages(conn *wafv2.WAFV2, input *wafv2.ListWebACLsInput, fn func(*wafv2.ListWebACLsOutput, bool) bool) error {
	return listWebACLsPagesWithContext(context.Background(), conn, input, fn)
}

func listWebACLsPagesWithContext(ctx context.Context, conn *wafv2.WAFV2, input *wafv2.ListWebACLsInput, fn func(*wafv2.ListWebACLsOutput, bool) bool) error {
	for {
		output, err := conn.ListWebACLsWithContext(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.StringValue(output.NextMarker) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextMarker = output.NextMarker
	}
	return nil
}
