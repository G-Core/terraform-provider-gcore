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

func sweepCloudReservedFixedIPs(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
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
		// Only sweep external reserved IPs: internal IPs are always attached to a port
		// and cannot be deleted independently.
		if !ip.IsExternal {
			continue
		}

		// Only sweep unattached IPs (status DOWN) to avoid breaking in-use resources.
		if ip.Reservation.Status != "DOWN" && ip.Reservation.Status != "" {
			log.Printf("[INFO] Skipping attached reserved fixed IP: %s (status: %s)", ip.PortID, ip.Reservation.Status)
			continue
		}

		// External IPs have auto-generated names; use ShouldSweep to filter by the
		// tf-test prefix when one is present, and allow unnamed IPs through as a
		// best-effort cleanup (they have no production significance).
		if ip.Name != "" && !sweep.ShouldSweep("gcore_cloud_reserved_fixed_ip", ip.Name) {
			continue
		}

		log.Printf("[INFO] Deleting reserved fixed IP: %s (external: %t, status: %s)", ip.PortID, ip.IsExternal, ip.Reservation.Status)
		_, err := client.Cloud.ReservedFixedIPs.Delete(ctx, ip.PortID, cloud.ReservedFixedIPDeleteParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete reserved fixed IP %s: %s", ip.PortID, err)
		}
	}

	return nil
}
