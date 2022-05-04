// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package opensearch

import (
	"fmt"

	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/opensearchservice"
	tftags "github.com/PixarV/terraform-provider-ritt/internal/tags"
)

// ListTags lists opensearch service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(conn *opensearchservice.OpenSearchService, identifier string) (tftags.KeyValueTags, error) {
	input := &opensearchservice.ListTagsInput{
		ARN: aws.String(identifier),
	}

	output, err := conn.ListTags(input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.TagList), nil
}

// []*SERVICE.Tag handling

// Tags returns opensearch service tags.
func Tags(tags tftags.KeyValueTags) []*opensearchservice.Tag {
	result := make([]*opensearchservice.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &opensearchservice.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from opensearchservice service tags.
func KeyValueTags(tags []*opensearchservice.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}

// UpdateTags updates opensearch service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn *opensearchservice.OpenSearchService, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &opensearchservice.RemoveTagsInput{
			ARN:     aws.String(identifier),
			TagKeys: aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.RemoveTags(input)

		if err != nil {
			return fmt.Errorf("error untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &opensearchservice.AddTagsInput{
			ARN:     aws.String(identifier),
			TagList: Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.AddTags(input)

		if err != nil {
			return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
