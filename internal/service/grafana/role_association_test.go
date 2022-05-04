package grafana_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/PixarV/aws-sdk-go/service/managedgrafana"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/PixarV/terraform-provider-ritt/internal/acctest"
	"github.com/PixarV/terraform-provider-ritt/internal/conns"
	tfgrafana "github.com/PixarV/terraform-provider-ritt/internal/service/grafana"
	"github.com/PixarV/terraform-provider-ritt/internal/tfresource"
)

func testAccGrafanaRoleAssociation_usersAdmin(t *testing.T) {
	key := "GRAFANA_SSO_USER_ID"
	userID := os.Getenv(key)
	if userID == "" {
		t.Skipf("Environment variable %s is not set", key)
	}

	role := managedgrafana.RoleAdmin
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_grafana_role_association.test"
	workspaceResourceName := "aws_grafana_workspace.test"

	resource.Test(t,
		resource.TestCase{
			PreCheck: func() {
				acctest.PreCheck(t)
				acctest.PreCheckPartitionHasService(managedgrafana.EndpointsID, t)
				acctest.PreCheckSSOAdminInstances(t)
			},
			ErrorCheck:   acctest.ErrorCheck(t, managedgrafana.EndpointsID),
			CheckDestroy: testAccCheckRoleAssociationDestroy,
			Providers:    acctest.Providers,
			Steps: []resource.TestStep{
				{
					Config: testAccWorkspaceGrafanaRoleAssociationUsers(rName, role, userID),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccCheckRoleAssociationExists(resourceName),
						resource.TestCheckResourceAttr(resourceName, "role", role),
						resource.TestCheckResourceAttr(resourceName, "user_ids.#", "1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "user_ids.*", userID),
						resource.TestCheckResourceAttrPair(resourceName, "workspace_id", workspaceResourceName, "id"),
					),
				},
			},
		})
}

func testAccGrafanaRoleAssociation_usersEditor(t *testing.T) {
	key := "GRAFANA_SSO_USER_ID"
	userID := os.Getenv(key)
	if userID == "" {
		t.Skipf("Environment variable %s is not set", key)
	}

	role := managedgrafana.RoleEditor
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_grafana_role_association.test"
	workspaceResourceName := "aws_grafana_workspace.test"

	resource.Test(t,
		resource.TestCase{
			PreCheck: func() {
				acctest.PreCheck(t)
				acctest.PreCheckPartitionHasService(managedgrafana.EndpointsID, t)
				acctest.PreCheckSSOAdminInstances(t)
			},
			ErrorCheck:   acctest.ErrorCheck(t, managedgrafana.EndpointsID),
			CheckDestroy: testAccCheckRoleAssociationDestroy,
			Providers:    acctest.Providers,
			Steps: []resource.TestStep{
				{
					Config: testAccWorkspaceGrafanaRoleAssociationUsers(rName, role, userID),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccCheckRoleAssociationExists(resourceName),
						resource.TestCheckResourceAttr(resourceName, "role", role),
						resource.TestCheckResourceAttr(resourceName, "user_ids.#", "1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "user_ids.*", userID),
						resource.TestCheckResourceAttrPair(resourceName, "workspace_id", workspaceResourceName, "id"),
					),
				},
			},
		})
}

func testAccGrafanaRoleAssociation_groupsAdmin(t *testing.T) {
	key := "GRAFANA_SSO_GROUP_ID"
	groupID := os.Getenv(key)
	if groupID == "" {
		t.Skipf("Environment variable %s is not set", key)
	}

	role := managedgrafana.RoleAdmin
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_grafana_role_association.test"
	workspaceResourceName := "aws_grafana_workspace.test"

	resource.Test(t,
		resource.TestCase{
			PreCheck: func() {
				acctest.PreCheck(t)
				acctest.PreCheckPartitionHasService(managedgrafana.EndpointsID, t)
				acctest.PreCheckSSOAdminInstances(t)
			},
			ErrorCheck:   acctest.ErrorCheck(t, managedgrafana.EndpointsID),
			CheckDestroy: testAccCheckRoleAssociationDestroy,
			Providers:    acctest.Providers,
			Steps: []resource.TestStep{
				{
					Config: testAccWorkspaceGrafanaRoleAssociationGroups(rName, role, groupID),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccCheckRoleAssociationExists(resourceName),
						resource.TestCheckResourceAttr(resourceName, "group_ids.#", "1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "group_ids.*", groupID),
						resource.TestCheckResourceAttr(resourceName, "role", role),
						resource.TestCheckResourceAttrPair(resourceName, "workspace_id", workspaceResourceName, "id"),
					),
				},
			},
		})
}

func testAccGrafanaRoleAssociation_groupsEditor(t *testing.T) {
	key := "GRAFANA_SSO_GROUP_ID"
	groupID := os.Getenv(key)
	if groupID == "" {
		t.Skipf("Environment variable %s is not set", key)
	}

	role := managedgrafana.RoleEditor
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_grafana_role_association.test"
	workspaceResourceName := "aws_grafana_workspace.test"

	resource.Test(t,
		resource.TestCase{
			PreCheck: func() {
				acctest.PreCheck(t)
				acctest.PreCheckPartitionHasService(managedgrafana.EndpointsID, t)
				acctest.PreCheckSSOAdminInstances(t)
			},
			ErrorCheck:   acctest.ErrorCheck(t, managedgrafana.EndpointsID),
			CheckDestroy: testAccCheckRoleAssociationDestroy,
			Providers:    acctest.Providers,
			Steps: []resource.TestStep{
				{
					Config: testAccWorkspaceGrafanaRoleAssociationGroups(rName, role, groupID),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccCheckRoleAssociationExists(resourceName),
						resource.TestCheckResourceAttr(resourceName, "group_ids.#", "1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "group_ids.*", groupID),
						resource.TestCheckResourceAttr(resourceName, "role", role),
						resource.TestCheckResourceAttrPair(resourceName, "workspace_id", workspaceResourceName, "id"),
					),
				},
			},
		})
}

func testAccGrafanaRoleAssociation_usersAndGroupsAdmin(t *testing.T) {
	key := "GRAFANA_SSO_USER_ID"
	userID := os.Getenv(key)
	if userID == "" {
		t.Skipf("Environment variable %s is not set", key)
	}
	key = "GRAFANA_SSO_GROUP_ID"
	groupID := os.Getenv(key)
	if groupID == "" {
		t.Skipf("Environment variable %s is not set", key)
	}

	role := managedgrafana.RoleAdmin
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_grafana_role_association.test"
	workspaceResourceName := "aws_grafana_workspace.test"

	resource.Test(t,
		resource.TestCase{
			PreCheck: func() {
				acctest.PreCheck(t)
				acctest.PreCheckPartitionHasService(managedgrafana.EndpointsID, t)
				acctest.PreCheckSSOAdminInstances(t)
			},
			ErrorCheck:   acctest.ErrorCheck(t, managedgrafana.EndpointsID),
			CheckDestroy: testAccCheckRoleAssociationDestroy,
			Providers:    acctest.Providers,
			Steps: []resource.TestStep{
				{
					Config: testAccWorkspaceGrafanaRoleAssociationUsersAndGroups(rName, role, userID, groupID),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccCheckRoleAssociationExists(resourceName),
						resource.TestCheckResourceAttr(resourceName, "group_ids.#", "1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "group_ids.*", groupID),
						resource.TestCheckResourceAttr(resourceName, "role", role),
						resource.TestCheckResourceAttr(resourceName, "user_ids.#", "1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "user_ids.*", userID),
						resource.TestCheckResourceAttrPair(resourceName, "workspace_id", workspaceResourceName, "id"),
					),
				},
			},
		})
}

func testAccGrafanaRoleAssociation_usersAndGroupsEditor(t *testing.T) {
	key := "GRAFANA_SSO_USER_ID"
	userID := os.Getenv(key)
	if userID == "" {
		t.Skipf("Environment variable %s is not set", key)
	}
	key = "GRAFANA_SSO_GROUP_ID"
	groupID := os.Getenv(key)
	if groupID == "" {
		t.Skipf("Environment variable %s is not set", key)
	}

	role := managedgrafana.RoleEditor
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_grafana_role_association.test"
	workspaceResourceName := "aws_grafana_workspace.test"

	resource.Test(t,
		resource.TestCase{
			PreCheck: func() {
				acctest.PreCheck(t)
				acctest.PreCheckPartitionHasService(managedgrafana.EndpointsID, t)
				acctest.PreCheckSSOAdminInstances(t)
			},
			ErrorCheck:   acctest.ErrorCheck(t, managedgrafana.EndpointsID),
			CheckDestroy: testAccCheckRoleAssociationDestroy,
			Providers:    acctest.Providers,
			Steps: []resource.TestStep{
				{
					Config: testAccWorkspaceGrafanaRoleAssociationUsersAndGroups(rName, role, userID, groupID),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccCheckRoleAssociationExists(resourceName),
						resource.TestCheckResourceAttr(resourceName, "group_ids.#", "1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "group_ids.*", groupID),
						resource.TestCheckResourceAttr(resourceName, "role", role),
						resource.TestCheckResourceAttr(resourceName, "user_ids.#", "1"),
						resource.TestCheckTypeSetElemAttr(resourceName, "user_ids.*", userID),
						resource.TestCheckResourceAttrPair(resourceName, "workspace_id", workspaceResourceName, "id"),
					),
				},
			},
		})
}

func testAccWorkspaceGrafanaRoleAssociationUsers(rName, role, userID string) string {
	return acctest.ConfigCompose(testAccWorkspaceConfigAuthenticationProvider(rName, "AWS_SSO"), fmt.Sprintf(`
resource "aws_grafana_role_association" "test" {
  role         = %[1]q
  user_ids     = [%[2]q]
  workspace_id = aws_grafana_workspace.test.id
}
`, role, userID))
}

func testAccWorkspaceGrafanaRoleAssociationGroups(rName, role, groupID string) string {
	return acctest.ConfigCompose(testAccWorkspaceConfigAuthenticationProvider(rName, "AWS_SSO"), fmt.Sprintf(`
resource "aws_grafana_role_association" "test" {
  role         = %[1]q
  group_ids    = [%[2]q]
  workspace_id = aws_grafana_workspace.test.id
}
`, role, groupID))
}

func testAccWorkspaceGrafanaRoleAssociationUsersAndGroups(rName, role, userID, groupID string) string {
	return acctest.ConfigCompose(testAccWorkspaceConfigAuthenticationProvider(rName, "AWS_SSO"), fmt.Sprintf(`
resource "aws_grafana_role_association" "test" {
  role         = %[1]q
  user_ids     = [%[2]q]
  group_ids    = [%[3]q]
  workspace_id = aws_grafana_workspace.test.id
}
`, role, userID, groupID))
}

func testAccCheckRoleAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).GrafanaConn

		_, err := tfgrafana.FindRoleAssociationsByRoleAndWorkspaceID(conn, rs.Primary.Attributes["role"], rs.Primary.Attributes["workspace_id"])

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckRoleAssociationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).GrafanaConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_grafana_role_association" {
			continue
		}

		_, err := tfgrafana.FindRoleAssociationsByRoleAndWorkspaceID(conn, rs.Primary.Attributes["role"], rs.Primary.Attributes["workspace_id"])

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Grafana Workspace Role Association %s still exists", rs.Primary.ID)
	}
	return nil
}
