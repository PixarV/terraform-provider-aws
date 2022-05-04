package directconnect

import (
	"fmt"

	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/PixarV/terraform-provider-ritt/internal/conns"
)

func DataSourceLocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationsRead,

		Schema: map[string]*schema.Schema{
			"location_codes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceLocationsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DirectConnectConn

	locations, err := FindLocations(conn, &directconnect.DescribeLocationsInput{})

	if err != nil {
		return fmt.Errorf("error reading Direct Connect locations: %w", err)
	}

	var locationCodes []*string

	for _, location := range locations {
		locationCodes = append(locationCodes, location.LocationCode)
	}

	d.SetId(meta.(*conns.AWSClient).Region)
	d.Set("location_codes", aws.StringValueSlice(locationCodes))

	return nil
}
