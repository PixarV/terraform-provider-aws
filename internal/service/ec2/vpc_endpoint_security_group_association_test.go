package ec2_test

import (
	"fmt"
	"testing"

	"github.com/PixarV/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/PixarV/terraform-provider-ritt/internal/acctest"
	"github.com/PixarV/terraform-provider-ritt/internal/conns"
	tfec2 "github.com/PixarV/terraform-provider-ritt/internal/service/ec2"
	"github.com/PixarV/terraform-provider-ritt/internal/tfresource"
)

func TestAccEC2VPCEndpointSecurityGroupAssociation_basic(t *testing.T) {
	var v ec2.VpcEndpoint
	resourceName := "aws_vpc_endpoint_security_group_association.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckVpcEndpointSecurityGroupAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndpointSecurityGroupAssociationConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcEndpointSecurityGroupAssociationExists(resourceName, &v),
					testAccCheckVpcEndpointSecurityGroupAssociationNumAssociations(&v, 2),
				),
			},
		},
	})
}

func TestAccEC2VPCEndpointSecurityGroupAssociation_disappears(t *testing.T) {
	var v ec2.VpcEndpoint
	resourceName := "aws_vpc_endpoint_security_group_association.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckVpcEndpointSecurityGroupAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndpointSecurityGroupAssociationConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcEndpointSecurityGroupAssociationExists(resourceName, &v),
					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceVPCEndpointSecurityGroupAssociation(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccEC2VPCEndpointSecurityGroupAssociation_multiple(t *testing.T) {
	var v ec2.VpcEndpoint
	resourceName0 := "aws_vpc_endpoint_security_group_association.test.0"
	resourceName1 := "aws_vpc_endpoint_security_group_association.test.1"
	resourceName2 := "aws_vpc_endpoint_security_group_association.test.2"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckVpcEndpointSecurityGroupAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndpointSecurityGroupAssociationConfigMultiple(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcEndpointSecurityGroupAssociationExists(resourceName0, &v),
					testAccCheckVpcEndpointSecurityGroupAssociationExists(resourceName1, &v),
					testAccCheckVpcEndpointSecurityGroupAssociationExists(resourceName2, &v),
					testAccCheckVpcEndpointSecurityGroupAssociationNumAssociations(&v, 4),
				),
			},
		},
	})
}

func TestAccEC2VPCEndpointSecurityGroupAssociation_replaceDefaultAssociation(t *testing.T) {
	var v ec2.VpcEndpoint
	resourceName := "aws_vpc_endpoint_security_group_association.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckVpcEndpointSecurityGroupAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndpointSecurityGroupAssociationConfigReplaceDefaultAssociation(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcEndpointSecurityGroupAssociationExists(resourceName, &v),
					testAccCheckVpcEndpointSecurityGroupAssociationNumAssociations(&v, 1),
				),
			},
		},
	})
}

func testAccCheckVpcEndpointSecurityGroupAssociationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_vpc_endpoint_security_group_association" {
			continue
		}

		err := tfec2.FindVPCEndpointSecurityGroupAssociationExists(conn, rs.Primary.Attributes["vpc_endpoint_id"], rs.Primary.Attributes["security_group_id"])

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("VPC Endpoint Security Group Association %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckVpcEndpointSecurityGroupAssociationExists(n string, v *ec2.VpcEndpoint) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC Endpoint Security Group Association ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		output, err := tfec2.FindVPCEndpointByID(conn, rs.Primary.Attributes["vpc_endpoint_id"])

		if err != nil {
			return err
		}

		err = tfec2.FindVPCEndpointSecurityGroupAssociationExists(conn, rs.Primary.Attributes["vpc_endpoint_id"], rs.Primary.Attributes["security_group_id"])

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckVpcEndpointSecurityGroupAssociationNumAssociations(v *ec2.VpcEndpoint, n int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len := len(v.Groups); len != n {
			return fmt.Errorf("got %d associations; wanted %d", len, n)
		}

		return nil
	}
}

func testAccVpcEndpointSecurityGroupAssociationConfigBase(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

data "aws_region" "current" {}

resource "aws_security_group" "test" {
  count = 3

  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc_endpoint" "test" {
  vpc_id            = aws_vpc.test.id
  service_name      = "com.amazonaws.${data.aws_region.current.name}.ec2"
  vpc_endpoint_type = "Interface"

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccVpcEndpointSecurityGroupAssociationConfigBasic(rName string) string {
	return acctest.ConfigCompose(
		testAccVpcEndpointSecurityGroupAssociationConfigBase(rName),
		`
resource "aws_vpc_endpoint_security_group_association" "test" {
  vpc_endpoint_id   = aws_vpc_endpoint.test.id
  security_group_id = aws_security_group.test[0].id
}
`)
}

func testAccVpcEndpointSecurityGroupAssociationConfigMultiple(rName string) string {
	return acctest.ConfigCompose(
		testAccVpcEndpointSecurityGroupAssociationConfigBase(rName),
		`
resource "aws_vpc_endpoint_security_group_association" "test" {
  count = length(aws_security_group.test)

  vpc_endpoint_id   = aws_vpc_endpoint.test.id
  security_group_id = aws_security_group.test[count.index].id
}
`)
}

func testAccVpcEndpointSecurityGroupAssociationConfigReplaceDefaultAssociation(rName string) string {
	return acctest.ConfigCompose(
		testAccVpcEndpointSecurityGroupAssociationConfigBase(rName),
		`
resource "aws_vpc_endpoint_security_group_association" "test" {
  vpc_endpoint_id   = aws_vpc_endpoint.test.id
  security_group_id = aws_security_group.test[0].id

  replace_default_association = true
}
`)
}
