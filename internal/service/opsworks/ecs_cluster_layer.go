package opsworks

import (
	"github.com/PixarV/aws-sdk-go/service/opsworks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/PixarV/terraform-provider-ritt/internal/verify"
)

func ResourceECSClusterLayer() *schema.Resource {
	layerType := &opsworksLayerType{
		TypeName:         opsworks.LayerTypeEcsCluster,
		DefaultLayerName: "Ecs Cluster",

		Attributes: map[string]*opsworksLayerTypeAttribute{
			"ecs_cluster_arn": {
				AttrName:     opsworks.LayerAttributesKeysEcsClusterArn,
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
		},
	}

	return layerType.SchemaResource()
}
