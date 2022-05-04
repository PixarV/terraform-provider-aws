package ses_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/PixarV/aws-sdk-go/service/ses"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/PixarV/terraform-provider-ritt/internal/acctest"
)

func TestAccSESEmailIdentityDataSource_basic(t *testing.T) {
	email := acctest.DefaultEmailAddress

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ses.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckEmailIdentityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailIdentityDataSourceConfig(email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityExists("aws_ses_email_identity.test"),
					acctest.MatchResourceAttrRegionalARN("data.aws_ses_email_identity.test", "arn", "ses", regexp.MustCompile(fmt.Sprintf("identity/%s$", regexp.QuoteMeta(email)))),
				),
			},
		},
	})
}

func TestAccSESEmailIdentityDataSource_trailingPeriod(t *testing.T) {
	email := fmt.Sprintf("%s.", acctest.DefaultEmailAddress)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, ses.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckEmailIdentityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailIdentityDataSourceConfig(email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailIdentityExists("aws_ses_email_identity.test"),
					acctest.MatchResourceAttrRegionalARN("data.aws_ses_email_identity.test", "arn", "ses", regexp.MustCompile(fmt.Sprintf("identity/%s$", regexp.QuoteMeta(strings.TrimSuffix(email, "."))))),
				),
			},
		},
	})
}

func testAccEmailIdentityDataSourceConfig(email string) string {
	return fmt.Sprintf(`
resource "aws_ses_email_identity" "test" {
  email = %q
}
data "aws_ses_email_identity" "test" {
  depends_on = [aws_ses_email_identity.test]
  email      = %q
}
`, email, email)
}
