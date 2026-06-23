package cdn_logs_uploader_policy_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCDNLogsUploaderPolicy_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNLogsUploaderPolicyMinimalConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("client_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("created"), knownvalue.NotNull()),
					// defaults
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("escape_special_characters"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("include_empty_logs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("include_shield_logs"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("log_sample_rate"), knownvalue.Float64Exact(1)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("retry_interval_minutes"), knownvalue.Int64Exact(60)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("rotate_interval_minutes"), knownvalue.Int64Exact(5)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("rotate_threshold_lines"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("date_format"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("format_type"), knownvalue.StringExact("")),
				},
			},
			// re-apply same config to verify no drift
			{
				Config: testAccCDNLogsUploaderPolicyMinimalConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccCDNLogsUploaderPolicy_full(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderPolicyDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with all params
			{
				Config: testAccCDNLogsUploaderPolicyFullConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("description"), knownvalue.StringExact("acctest policy")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("format_type"), knownvalue.StringExact("json")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("field_delimiter"), knownvalue.StringExact(",")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("field_separator"), knownvalue.StringExact("|")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("include_shield_logs"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("include_empty_logs"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("escape_special_characters"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("log_sample_rate"), knownvalue.Float64Exact(0.5)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("retry_interval_minutes"), knownvalue.Int64Exact(30)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("rotate_interval_minutes"), knownvalue.Int64Exact(10)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("rotate_threshold_mb"), knownvalue.Int64Exact(200)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("rotate_threshold_lines"), knownvalue.Int64Exact(5000)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("fields"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("remote_addr"),
							knownvalue.StringExact("status"),
						})),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"env": knownvalue.StringExact("test"),
						})),
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 2: no drift
			{
				Config: testAccCDNLogsUploaderPolicyFullConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 3: Import
			{
				ResourceName:      "gcore_cdn_logs_uploader_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCDNLogsUploaderPolicyImportStateIDFunc("gcore_cdn_logs_uploader_policy.test"),
			},
		},
	})
}

func TestAccCDNLogsUploaderPolicy_update(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderPolicyDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with full params
			{
				Config: testAccCDNLogsUploaderPolicyFullConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("description"), knownvalue.StringExact("acctest policy")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("rotate_threshold_mb"), knownvalue.Int64Exact(200)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"env": knownvalue.StringExact("test"),
						})),
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 2: Update to different values
			{
				Config: testAccCDNLogsUploaderPolicyUpdatedConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("description"), knownvalue.StringExact("updated description")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("format_type"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("field_delimiter"), knownvalue.StringExact(";")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("log_sample_rate"), knownvalue.Float64Exact(0.75)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("retry_interval_minutes"), knownvalue.Int64Exact(15)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"env":  knownvalue.StringExact("staging"),
							"team": knownvalue.StringExact("platform"),
						})),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("fields"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("remote_addr"),
							knownvalue.StringExact("status"),
							knownvalue.StringExact("request_uri"),
						})),
					// ID should not change — in-place update
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 3: no drift after update
			{
				Config: testAccCDNLogsUploaderPolicyUpdatedConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 4: Strip back to minimal — removes tags, description, fields, rotate_threshold_mb, etc.
			{
				Config: testAccCDNLogsUploaderPolicyMinimalConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{})),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("escape_special_characters"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("log_sample_rate"), knownvalue.Float64Exact(1)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("rotate_threshold_mb"), knownvalue.Null()),
					// ID should not change
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 5: no drift after stripping back
			{
				Config: testAccCDNLogsUploaderPolicyMinimalConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// testAccCDNLogsUploaderPolicyImportStateIDFunc returns the numeric id as string for import.
func testAccCDNLogsUploaderPolicyImportStateIDFunc(resourceName string) func(*terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return rs.Primary.Attributes["id"], nil
	}
}

func testAccCheckCDNLogsUploaderPolicyDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cdn_logs_uploader_policy" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.Attributes["id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing id: %w", err)
		}

		_, err = client.CDN.LogsUploader.Policies.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("CDN logs uploader policy %d still exists", id)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking CDN logs uploader policy deletion: %w", err)
		}
	}
	return nil
}

func testAccCDNLogsUploaderPolicyMinimalConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_policy" "test" {
  name = %[1]q
}`, name)
}

func testAccCDNLogsUploaderPolicyFullConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_policy" "test" {
  name                     = %[1]q
  description              = "acctest policy"
  format_type              = "json"
  fields                   = ["remote_addr", "status"]
  field_delimiter          = ","
  field_separator          = "|"
  include_shield_logs      = true
  include_empty_logs       = true
  escape_special_characters = true
  log_sample_rate          = 0.5
  retry_interval_minutes   = 30
  rotate_interval_minutes  = 10
  rotate_threshold_mb      = 200
  rotate_threshold_lines   = 5000
  date_format              = "[02/Jan/2006:15:04:05 -0700]"
  file_name_template       = "{{YYYY}}_{{MM}}_{{DD}}_{{HH}}_{{mm}}_{{ss}}_{{HOST}}_{{CNAME}}_access.log.gz"
  tags = {
    env = "test"
  }
}`, name)
}

func testAccCDNLogsUploaderPolicyUpdatedConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_policy" "test" {
  name                     = %[1]q
  description              = "updated description"
  fields                   = ["remote_addr", "status", "request_uri"]
  field_delimiter          = ";"
  field_separator          = "|"
  include_shield_logs      = true
  include_empty_logs       = false
  escape_special_characters = false
  log_sample_rate          = 0.75
  retry_interval_minutes   = 15
  rotate_interval_minutes  = 10
  rotate_threshold_mb      = 200
  rotate_threshold_lines   = 5000
  date_format              = "[02/Jan/2006:15:04:05 -0700]"
  file_name_template       = "{{YYYY}}_{{MM}}_{{DD}}_{{HH}}_{{mm}}_{{ss}}_{{HOST}}_{{CNAME}}_access.log.gz"
  tags = {
    env  = "staging"
    team = "platform"
  }
}`, name)
}
