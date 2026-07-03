package cdn_origin_shielding_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// testAccCheckCDNOriginShieldingDestroy verifies shielding was disabled on destroy.
// After `terraform destroy` the resource's PUT sets shielding_pop to null; the parent
// CDN resource fixture still exists at this point (its t.Cleanup runs later), so a Get
// should report a null (0) shielding_pop. If the CDN resource is already gone, that also
// means shielding is gone.
func testAccCheckCDNOriginShieldingDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cdn_origin_shielding" {
			continue
		}
		resourceID, err := strconv.ParseInt(rs.Primary.Attributes["resource_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse resource_id: %w", err)
		}
		result, err := client.CDN.CDNResources.Shield.Get(ctx, resourceID)
		if err != nil {
			if acctest.IsNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking origin shielding for resource %d: %w", resourceID, err)
		}
		if result.ShieldingPop != 0 {
			return fmt.Errorf("origin shielding still enabled on resource %d (shielding_pop=%d)", resourceID, result.ShieldingPop)
		}
	}
	return nil
}

// TestAccCDNOriginShielding_basic covers create, no-drift re-apply (config drift check),
// and import.
func TestAccCDNOriginShielding_basic(t *testing.T) {
	fixture := setupCDNResource(t)
	requireOriginShielding(t, fixture.ResourceID)
	pop := shieldingLocationIDs(t)[0]

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNOriginShieldingDestroy,
		Steps: []resource.TestStep{
			// Step 1: create (enable shielding)
			{
				Config: testAccCDNOriginShieldingConfig(fixture.ResourceID, pop),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_origin_shielding.test",
						tfjsonpath.New("resource_id"), knownvalue.Int64Exact(fixture.ResourceID)),
					statecheck.ExpectKnownValue("gcore_cdn_origin_shielding.test",
						tfjsonpath.New("shielding_pop"), knownvalue.Int64Exact(pop)),
				},
			},
			// Step 2: re-apply the same config -> must be a no-op (no config drift)
			{
				Config: testAccCDNOriginShieldingConfig(fixture.ResourceID, pop),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			// Step 3: import by resource_id (the resource has no separate id)
			{
				ResourceName:                         "gcore_cdn_origin_shielding.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "resource_id",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return s.RootModule().Resources["gcore_cdn_origin_shielding.test"].Primary.Attributes["resource_id"], nil
				},
			},
		},
	})
}

// TestAccCDNOriginShielding_update changes the shielding location and asserts it is an
// in-place update (not a replace).
func TestAccCDNOriginShielding_update(t *testing.T) {
	fixture := setupCDNResource(t)
	requireOriginShielding(t, fixture.ResourceID)
	pops := shieldingLocationIDs(t)
	if len(pops) < 2 {
		t.Skip("need at least 2 origin shielding locations to test shielding_pop update")
	}
	pop1, pop2 := pops[0], pops[1]

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNOriginShieldingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNOriginShieldingConfig(fixture.ResourceID, pop1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_origin_shielding.test",
						tfjsonpath.New("shielding_pop"), knownvalue.Int64Exact(pop1)),
				},
			},
			{
				Config: testAccCDNOriginShieldingConfig(fixture.ResourceID, pop2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cdn_origin_shielding.test",
							plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_origin_shielding.test",
						tfjsonpath.New("shielding_pop"), knownvalue.Int64Exact(pop2)),
				},
			},
		},
	})
}

func testAccCDNOriginShieldingConfig(resourceID, pop int64) string {
	return fmt.Sprintf(`
resource "gcore_cdn_origin_shielding" "test" {
  resource_id   = %d
  shielding_pop = %d
}`, resourceID, pop)
}
