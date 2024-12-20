---
subcategory: "CloudTrail"
layout: "aws"
page_title: "AWS: aws_cloudtrail_event_data_store"
description: |-
  Provides a CloudTrail Event Data Store resource.
---

# Resource: aws_cloudtrail_event_data_store

Provides a CloudTrail Event Data Store.

More information about event data stores can be found in the [Event Data Store User Guide](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/query-event-data-store.html).

-> **Tip:** For an organization event data store you must create this resource in the management account.

## Example Usage

### Basic

The most simple event data store configuration requires us to only set the `name` and `retention_period` attributes. The event data store will automatically capture all management events. To capture management events from all the regions, `multi_region_enabled` must be `true`.

```terraform
resource "aws_cloudtrail_event_data_store" "example" {
  name             = "example-event-data-store"
  retention_period = 7
}
```

### Data Event Logging

CloudTrail can log [Data Events](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/logging-data-events-with-cloudtrail.html) for certain services such as S3 bucket objects and Lambda function invocations. Additional information about data event configuration can be found in the following links:

- [CloudTrail API AdvancedFieldSelector documentation](https://docs.aws.amazon.com/awscloudtrail/latest/APIReference/API_AdvancedFieldSelector.html)

#### Log all DynamoDB PutEvent actions for a specific DynamoDB table

```terraform
data "aws_dynamodb_table" "table" {
  name = "not-important-dynamodb-table"
}

resource "aws_cloudtrail_event_data_store" "example" {
  # ... other configuration ...

  advanced_event_selector {
    name = "Log all DynamoDB PutEvent actions for a specific DynamoDB table"

    field_selector {
      field  = "eventCategory"
      equals = ["Data"]
    }

    field_selector {
      field = "resources.type"

      equals = [
        "AWS::DynamoDB::Table"
      ]
    }

    field_selector {
      field  = "eventName"
      equals = ["PutItem"]
    }

    field_selector {
      field = "resources.ARN"

      equals = [
        data.aws_dynamodb_table.table.arn
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the event data store.
- `advanced_event_selector` - (Required) The advanced event selectors to use to select the events for the data store. For more information about how to use advanced event selectors, see [Log events by using advanced event selectors](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/logging-data-events-with-cloudtrail.html#creating-data-event-selectors-advanced) in the CloudTrail User Guide.
- `multi_region_enabled` - (Optional) Specifies whether the event data store includes events from all regions, or only from the region in which the event data store is created. Default: `true`.
- `organization_enabled` - (Optional) Specifies whether an event data store collects events logged for an organization in AWS Organizations. Default: `false`.
- `retention_period` - (Optional) The retention period of the event data store, in days. You can set a retention period of up to 2555 days, the equivalent of seven years. Default: `2555`.
- `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block](/docs/providers/aws/index.html#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.
- `termination_protection_enabled` - (Optional) Specifies whether termination protection is enabled for the event data store. If termination protection is enabled, you cannot delete the event data store until termination protection is disabled. Default: `true`.

### Advanced Event Selector Arguments

For **advanced_event_selector** the following attributes are supported.

- `name` (Optional) - Specifies the name of the advanced event selector.
- `field_selector` (Required) - Specifies the selector statements in an advanced event selector. Fields documented below.

#### Field Selector Arguments

For **field_selector** the following attributes are supported.

- `field` (Required) - Specifies a field in an event record on which to filter events to be logged. You can specify only the following values: `readOnly`, `eventSource`, `eventName`, `eventCategory`, `resources.type`, `resources.ARN`.
- `equals` (Optional) - A list of values that includes events that match the exact value of the event record field specified as the value of `field`. This is the only valid operator that you can use with the `readOnly`, `eventCategory`, and `resources.type` fields.
- `not_equals` (Optional) - A list of values that excludes events that match the exact value of the event record field specified as the value of `field`.
- `starts_with` (Optional) - A list of values that includes events that match the first few characters of the event record field specified as the value of `field`.
- `not_starts_with` (Optional) - A list of values that excludes events that match the first few characters of the event record field specified as the value of `field`.
- `ends_with` (Optional) - A list of values that includes events that match the last few characters of the event record field specified as the value of `field`.
- `not_ends_with` (Optional) - A list of values that excludes events that match the last few characters of the event record field specified as the value of `field`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `arn` - ARN of the event data store.
- `id` - Name of the event data store.
- `tags_all` - Map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](/docs/providers/aws/index.html#default_tags-configuration-block).

## Import

Event data stores can be imported using their `arn`, e.g.,

```
$ terraform import aws_cloudtrail_event_data_store.example arn:aws:cloudtrail:us-east-1:123456789123:eventdatastore/22333815-4414-412c-b155-dd254033gfhf
```