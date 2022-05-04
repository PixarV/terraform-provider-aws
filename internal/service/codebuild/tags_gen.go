// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package codebuild

import (
	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/codebuild"
	tftags "github.com/PixarV/terraform-provider-ritt/internal/tags"
)

// []*SERVICE.Tag handling

// Tags returns codebuild service tags.
func Tags(tags tftags.KeyValueTags) []*codebuild.Tag {
	result := make([]*codebuild.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &codebuild.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from codebuild service tags.
func KeyValueTags(tags []*codebuild.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}
