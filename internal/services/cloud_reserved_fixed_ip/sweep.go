package cloud_reserved_fixed_ip

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_cloud_reserved_fixed_ip", &resource.Sweeper{
		Name: "gcore_cloud_reserved_fixed_ip",
		F:    sweepCloudReservedFixedIPs,
	})
}

// sweepCloudReservedFixedIPs reclaims external reserved fixed IPs, but only
// when sweep.EnvSweepAllReservedFixedIPs is set.
//
// Nothing can tie one of these to the run that created it. The name is computed
// by the API, no create or update variant accepts a name or tags, and an
// external address sits on the region's shared public network rather than a
// test one, so not even the parent carries a token. Every ownership mechanism
// the other sweepers rely on is unavailable here.
//
// So this is deliberately an unsafe operation behind an explicit flag rather
// than a clever filter pretending to establish ownership. Unset, it deletes
// nothing; set, it deletes every external address in the DOWN state in this
// project and region. What bounds it is scope rather than precision: the list
// call is project-scoped, so it reaches only a project dedicated to acceptance
// tests.
//
// Internal addresses are never considered. They are attached to a port on a
// test network and are reclaimed with it, so they need nothing from here.
func sweepCloudReservedFixedIPs(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	if os.Getenv(sweep.EnvSweepAllReservedFixedIPs) != "1" {
		log.Printf("[INFO] Skipping reserved fixed IP sweep: the API exposes no name or tag to carry a "+
			"run token, so ownership cannot be established. Set %s=1 to delete every external address "+
			"in the DOWN state in this project.", sweep.EnvSweepAllReservedFixedIPs)
		return nil
	}

	apiKey := os.Getenv("GCORE_API_KEY")
	projectID, err := strconv.ParseInt(os.Getenv("GCORE_CLOUD_PROJECT_ID"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse GCORE_CLOUD_PROJECT_ID: %w", err)
	}
	regionID, err := strconv.ParseInt(os.Getenv("GCORE_CLOUD_REGION_ID"), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse GCORE_CLOUD_REGION_ID: %w", err)
	}

	client := gcore.NewClient(
		option.WithAPIKey(apiKey),
	)

	ctx := context.Background()

	log.Printf("[WARN] Reserved fixed IP sweep running with %s=1: every external address in the DOWN "+
		"state in project %d region %d will be deleted, including any a concurrent run is holding.",
		sweep.EnvSweepAllReservedFixedIPs, projectID, regionID)

	page, err := client.Cloud.ReservedFixedIPs.List(ctx, cloud.ReservedFixedIPListParams{
		ProjectID: param.NewOpt(projectID),
		RegionID:  param.NewOpt(regionID),
	})
	if err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping reserved fixed IP sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing reserved fixed IPs: %w", err)
	}

	for _, ip := range page.Results {
		// Internal addresses cannot be deleted independently of the port they
		// are attached to, and go with the test network that holds them.
		if !ip.IsExternal {
			continue
		}

		// Anything in another state is attached to something still in use,
		// whether or not a test owns it.
		if ip.Reservation.Status != "DOWN" && ip.Reservation.Status != "" {
			log.Printf("[DEBUG] Skipping reserved fixed IP %s: status %s", ip.PortID, ip.Reservation.Status)
			continue
		}

		// The flag is named in the line deliberately. With no ownership check,
		// a developer whose local run breaks because of this needs the log to
		// say why rather than being left to file a provider bug.
		log.Printf("[INFO] Deleting reserved fixed IP %s (external, status %q, created %s) under %s=1",
			ip.PortID, ip.Reservation.Status, ip.CreatedAt.Format("2006-01-02T15:04:05Z"),
			sweep.EnvSweepAllReservedFixedIPs)

		if _, err := client.Cloud.ReservedFixedIPs.Delete(ctx, ip.PortID, cloud.ReservedFixedIPDeleteParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		}); err != nil {
			log.Printf("[ERROR] Failed to delete reserved fixed IP %s: %s", ip.PortID, err)
		}
	}

	return nil
}
