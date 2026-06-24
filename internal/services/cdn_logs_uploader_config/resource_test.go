package cdn_logs_uploader_config_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cdn"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// TestAccCDNLogsUploaderConfig_basic creates a config and asserts the apply is
// stable. The plugin-testing framework runs a refresh+plan after the step and
// fails on any non-empty plan, so this step is the regression guard for
// GCLOUD2-22619: the generated schema produced a perpetual update plan because
// the API echoes `resources` as `[]` for a null input. The custom-code fix
// (resources -> computed_optional + UseStateForUnknown) is what keeps this
// step's post-apply plan empty.
func TestAccCDNLogsUploaderConfig_basic(t *testing.T) {
	rName := acctest.RandomName()
	policyID, targetID := testAccLogsUploaderConfigDeps(t, rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNLogsUploaderConfigConfig(rName, policyID, targetID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("policy"), knownvalue.Int64Exact(policyID)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("target"), knownvalue.Int64Exact(targetID)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					// resources is null in config; API materializes it as []. The fix
					// preserves that empty list in state instead of perpetually
					// diffing [] -> null.
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("resources"), knownvalue.ListExact([]knownvalue.Check{})),
				},
			},
		},
	})
}

// TestAccCDNLogsUploaderConfig_update changes mutable fields and asserts the
// resource is updated in place (id stable) and that the immutable computed
// fields client_id/created do not flip to (known after apply) during update —
// the second half of the GCLOUD2-22619 fix.
func TestAccCDNLogsUploaderConfig_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()
	policyID, targetID := testAccLogsUploaderConfigDeps(t, rName)

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())
	compareClientIDSame := statecheck.CompareValue(compare.ValuesSame())
	compareCreatedSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNLogsUploaderConfigConfigEnabled(rName, policyID, targetID, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					compareIDSame.AddStateValue("gcore_cdn_logs_uploader_config.test", tfjsonpath.New("id")),
					compareClientIDSame.AddStateValue("gcore_cdn_logs_uploader_config.test", tfjsonpath.New("client_id")),
					compareCreatedSame.AddStateValue("gcore_cdn_logs_uploader_config.test", tfjsonpath.New("created")),
				},
			},
			{
				Config: testAccCDNLogsUploaderConfigConfigEnabled(rNameUpdated, policyID, targetID, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("enabled"), knownvalue.Bool(false)),
					// in-place update: id, client_id, created must be unchanged
					compareIDSame.AddStateValue("gcore_cdn_logs_uploader_config.test", tfjsonpath.New("id")),
					compareClientIDSame.AddStateValue("gcore_cdn_logs_uploader_config.test", tfjsonpath.New("client_id")),
					compareCreatedSame.AddStateValue("gcore_cdn_logs_uploader_config.test", tfjsonpath.New("created")),
				},
			},
		},
	})
}

// TestAccCDNLogsUploaderConfig_import verifies the config round-trips through
// import by its numeric id.
func TestAccCDNLogsUploaderConfig_import(t *testing.T) {
	rName := acctest.RandomName()
	policyID, targetID := testAccLogsUploaderConfigDeps(t, rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNLogsUploaderConfigConfig(rName, policyID, targetID),
			},
			{
				ResourceName:      "gcore_cdn_logs_uploader_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccLogsUploaderConfigDeps creates a logs-uploader policy and HTTP target
// directly via the SDK (those resources are not part of this provider build)
// and registers cleanup. It returns the IDs to wire into the config resource.
func testAccLogsUploaderConfigDeps(t *testing.T, name string) (policyID, targetID int64) {
	t.Helper()

	// Mirror resource.ParallelTest's TF_ACC gate so deps aren't created during
	// a plain `go test` run.
	if os.Getenv("TF_ACC") == "" {
		return 0, 0
	}
	acctest.PreCheck(t)

	client, err := acctest.NewGcoreClient()
	if err != nil {
		t.Fatalf("failed to create client: %s", err)
	}
	ctx := context.Background()

	policy, err := client.CDN.LogsUploader.Policies.New(ctx, cdn.LogsUploaderPolicyNewParams{
		Name: gcore.String(name + "-policy"),
	})
	if err != nil {
		t.Fatalf("failed to create logs uploader policy: %s", err)
	}
	t.Cleanup(func() {
		if err := client.CDN.LogsUploader.Policies.Delete(context.Background(), policy.ID); err != nil {
			t.Logf("warning: failed to clean up logs uploader policy %d: %s", policy.ID, err)
		}
	})

	target, err := client.CDN.LogsUploader.Targets.New(ctx, cdn.LogsUploaderTargetNewParams{
		Name:        gcore.String(name + "-target"),
		StorageType: cdn.LogsUploaderTargetNewParamsStorageTypeHTTP,
		Config: cdn.LogsUploaderTargetNewParamsConfigUnion{
			OfHTTPConfig: &cdn.LogsUploaderTargetNewParamsConfigHTTPConfig{
				Upload: cdn.LogsUploaderTargetNewParamsConfigHTTPConfigUpload{
					URL: "https://example.com/logs",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create logs uploader target: %s", err)
	}
	t.Cleanup(func() {
		if err := client.CDN.LogsUploader.Targets.Delete(context.Background(), target.ID); err != nil {
			t.Logf("warning: failed to clean up logs uploader target %d: %s", target.ID, err)
		}
	})

	return policy.ID, target.ID
}

func testAccCheckCDNLogsUploaderConfigDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cdn_logs_uploader_config" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.Attributes["id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing id: %w", err)
		}

		_, err = client.CDN.LogsUploader.Configs.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("CDN logs uploader config %d still exists", id)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking CDN logs uploader config deletion: %w", err)
		}
	}

	return nil
}

func testAccCDNLogsUploaderConfigConfig(name string, policyID, targetID int64) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_config" "test" {
  name   = %[1]q
  policy = %[2]d
  target = %[3]d
}`, name, policyID, targetID)
}

func testAccCDNLogsUploaderConfigConfigEnabled(name string, policyID, targetID int64, enabled bool) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_config" "test" {
  name    = %[1]q
  policy  = %[2]d
  target  = %[3]d
  enabled = %[4]t
}`, name, policyID, targetID, enabled)
}
