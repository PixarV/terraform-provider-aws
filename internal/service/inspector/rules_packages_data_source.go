package inspector

import (
	"fmt"
	"sort"

	"github.com/PixarV/aws-sdk-go/aws"
	"github.com/PixarV/aws-sdk-go/service/inspector"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/PixarV/terraform-provider-ritt/internal/conns"
)

func DataSourceRulesPackages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRulesPackagesRead,

		Schema: map[string]*schema.Schema{
			"arns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceRulesPackagesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).InspectorConn

	output, err := findRulesPackageArns(conn)

	if err != nil {
		return fmt.Errorf("error reading Inspector Rules Packages: %w", err)
	}

	arns := aws.StringValueSlice(output)
	sort.Strings(arns)

	d.SetId(meta.(*conns.AWSClient).Region)
	d.Set("arns", arns)

	return nil
}

func findRulesPackageArns(conn *inspector.Inspector) ([]*string, error) {
	input := &inspector.ListRulesPackagesInput{}
	var output []*string

	err := conn.ListRulesPackagesPages(input, func(page *inspector.ListRulesPackagesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.RulesPackageArns {
			if v != nil {
				output = append(output, v)
			}
		}

		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
