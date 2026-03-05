package cloud_security_group_rule

import (
	"log"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stainless-sdks/gcore-terraform/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_cloud_security_group_rule", &resource.Sweeper{
		Name:         "gcore_cloud_security_group_rule",
		F:            sweepCloudSecurityGroupRules,
		Dependencies: []string{"gcore_cloud_security_group"},
	})
}

func sweepCloudSecurityGroupRules(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}
	log.Printf("[INFO] Security group rule sweeper: rules are deleted with their parent security groups. " +
		"The gcore_cloud_security_group sweeper handles cleanup.")
	return nil
}
