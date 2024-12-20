---
subcategory: "EC2 Image Builder"
layout: "aws"
page_title: "AWS: aws_imagebuilder_container_recipe"
description: |-
    Provides details about an Image Builder Container Recipe
---

# Data Source: aws_imagebuilder_container_recipe

Provides details about an Image builder Container Recipe.

## Example Usage

```terraform
data "aws_imagebuilder_container_recipe" "example" {
  arn = "arn:aws:imagebuilder:us-east-1:aws:container-recipe/example/1.0.0"
}
```

## Argument Reference

The following arguments are required:

* `arn` - (Required) Amazon Resource Name (ARN) of the container recipe.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `component` - List of objects with components for the container recipe.
    * `component_arn` - Amazon Resource Name (ARN) of the Image Builder Component.
    * `parameter` - Set of parameters that are used to configure the component.
        * `name` - Name of the component parameter.
        * `value` - Value of the component parameter.
* `container_type` - Type of the container.
* `date_created` - Date the container recipe was created.
* `description` - Description of the container recipe.
* `dockerfile_template_data` - Dockerfile template used to build the image.
* `encrypted` - Flag that indicates if the target container is encrypted.
* `instance_configuration` - List of objects with instance configurations for building and testing container images.
    * `block_device_mapping` - Set of objects with block device mappings for the instance configuration.
        * `device_name` - Name of the device. For example, `/dev/sda` or `/dev/xvdb`.
        * `ebs` - Single list of object with Elastic Block Storage (EBS) block device mapping settings.
            * `delete_on_termination` - Whether to delete the volume on termination. Defaults to unset, which is the value inherited from the parent image.
            * `encrypted` - Whether to encrypt the volume. Defaults to unset, which is the value inherited from the parent image.
            * `iops` - Number of Input/Output (I/O) operations per second to provision for an `io1` or `io2` volume.
            * `kms_key_id` - Amazon Resource Name (ARN) of the Key Management Service (KMS) Key for encryption.
            * `snapshot_id` - Identifier of the EC2 Volume Snapshot.
            * `volume_size` - Size of the volume, in GiB.
            * `volume_type` - Type of the volume. For example, `gp2` or `io2`.
        * `no_device` - Whether to remove a mapping from the parent image.
        * `virtual_name` - Virtual device name. For example, `ephemeral0`. Instance store volumes are numbered starting from 0.
    * `image` - AMI ID of the base image for container build and test instance.
* `kms_key_id` - KMS key used to encrypt the container image.
* `name` - Name of the container recipe.
* `owner` - Owner of the container recipe.
* `parent_image` - Base image for the container recipe.
* `platform` - Platform of the container recipe.
* `tags` - Key-value map of resource tags for the container recipe.
* `target_repository` - Destination repository for the container image.
    * `repository_name` - Name of the container repository where the output container image is stored. The name is prefixed by the repository location.
    * `service` - Service in which this image is registered.
* `version` - Version of the container recipe.
* `working_directory` - The working directory used during build and test workflows.