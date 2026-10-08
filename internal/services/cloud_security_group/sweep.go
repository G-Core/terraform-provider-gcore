package cloud_security_group

import (
	"context"
	"fmt"
	"log"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_cloud_security_group", &resource.Sweeper{
		Name: "gcore_cloud_security_group",
		F:    sweepCloudSecurityGroups,
	})
}

func sweepCloudSecurityGroups(_ string) error {
	const resourceType = "gcore_cloud_security_group"
	defer sweep.Finish(resourceType)

	ctx := context.Background()
	target, err := sweep.Resolve(ctx)
	if err != nil {
		return err
	}

	// Collect first, delete second. Deleting while iterating an
	// offset-paginated collection shifts it, so the next offset skips exactly
	// as many items as were removed.
	var doomed []cloud.SecurityGroup
	iter := target.Client.Cloud.SecurityGroups.ListAutoPaging(ctx, cloud.SecurityGroupListParams{
		ProjectID: param.NewOpt(target.ProjectID),
		RegionID:  param.NewOpt(target.RegionID),
	})
	for iter.Next() {
		group := iter.Current()
		sweep.RecordListed(resourceType, 1)
		if sweep.ShouldSweep(resourceType, group.Name) {
			doomed = append(doomed, group)
		}
	}
	if err := iter.Err(); err != nil {
		if sweep.SkipSweepError(err) {
			log.Printf("[WARN] Skipping security group sweep: %s", err)
			return nil
		}
		return fmt.Errorf("error listing security groups: %w", err)
	}

	for _, group := range doomed {
		log.Printf("[INFO] Deleting security group: %s (%s)", group.Name, group.ID)
		if err := target.Client.Cloud.SecurityGroups.Delete(ctx, group.ID, cloud.SecurityGroupDeleteParams{
			ProjectID: param.NewOpt(target.ProjectID),
			RegionID:  param.NewOpt(target.RegionID),
		}); err != nil {
			sweep.RecordFailure(resourceType, group.Name, err)
			continue
		}
		sweep.RecordDeleted(resourceType, group.Name)
	}

	return nil
}
