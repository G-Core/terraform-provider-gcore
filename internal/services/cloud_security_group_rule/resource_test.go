package cloud_security_group_rule_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/stainless-sdks/gcore-terraform/internal/acctest"
)

func TestAccCloudSecurityGroupRule_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityGroupRuleConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_security_group_rule.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_security_group_rule.test",
						tfjsonpath.New("direction"), knownvalue.StringExact("ingress")),
					statecheck.ExpectKnownValue("gcore_cloud_security_group_rule.test",
						tfjsonpath.New("protocol"), knownvalue.StringExact("tcp")),
					statecheck.ExpectKnownValue("gcore_cloud_security_group_rule.test",
						tfjsonpath.New("ethertype"), knownvalue.StringExact("IPv4")),
					statecheck.ExpectKnownValue("gcore_cloud_security_group_rule.test",
						tfjsonpath.New("port_range_min"), knownvalue.Int64Exact(80)),
					statecheck.ExpectKnownValue("gcore_cloud_security_group_rule.test",
						tfjsonpath.New("port_range_max"), knownvalue.Int64Exact(80)),
					statecheck.ExpectKnownValue("gcore_cloud_security_group_rule.test",
						tfjsonpath.New("remote_ip_prefix"), knownvalue.StringExact("0.0.0.0/0")),
				},
			},
		},
	})
}

func TestAccCloudSecurityGroupRule_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityGroupRuleConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_security_group_rule.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:      "gcore_cloud_security_group_rule.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_security_group_rule.test", "project_id", "region_id", "group_id", "id"),
			},
		},
	})
}

func testAccCheckCloudSecurityGroupRuleDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_security_group_rule" {
			continue
		}

		projectID, err := strconv.ParseInt(rs.Primary.Attributes["project_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing project_id: %w", err)
		}
		regionID, err := strconv.ParseInt(rs.Primary.Attributes["region_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing region_id: %w", err)
		}
		groupID := rs.Primary.Attributes["group_id"]
		ruleID := rs.Primary.ID

		// Get the parent security group to check if rule still exists
		sg, err := client.Cloud.SecurityGroups.Get(context.Background(), groupID, cloud.SecurityGroupGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err != nil {
			// If security group is not found, the rule is also gone
			if acctest.IsNotFoundError(err) {
				return nil
			}
			return fmt.Errorf("error checking security group: %w", err)
		}

		// Check if the rule still exists in the security group
		for _, rule := range sg.SecurityGroupRules {
			if rule.ID == ruleID {
				return fmt.Errorf("security group rule %s still exists", ruleID)
			}
		}
	}

	return nil
}

func testAccCloudSecurityGroupRuleConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_security_group" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
}

resource "gcore_cloud_security_group_rule" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  group_id         = gcore_cloud_security_group.test.id
  direction        = "ingress"
  protocol         = "tcp"
  ethertype        = "IPv4"
  port_range_min   = 80
  port_range_max   = 80
  remote_ip_prefix = "0.0.0.0/0"
  description      = ""
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
