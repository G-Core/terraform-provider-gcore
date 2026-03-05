package fastedge_binary

import (
	"log"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/stainless-sdks/gcore-terraform/internal/sweep"
)

func init() {
	resource.AddTestSweepers("gcore_fastedge_binary", &resource.Sweeper{
		Name: "gcore_fastedge_binary",
		F:    sweepFastedgeBinaries,
	})
}

func sweepFastedgeBinaries(_ string) error {
	if err := sweep.ValidateSweeperEnvironment(); err != nil {
		return err
	}

	// FastEdge binaries don't have names - they only have numeric IDs and checksums.
	// There's no way to identify which binaries are test resources vs production.
	// Test binaries are cleaned up by terraform destroy.
	// Unreferenced binaries are automatically cleaned up by the FastEdge platform.
	log.Printf("[INFO] FastEdge binary sweeper: binaries have no names for identification. " +
		"Test binaries are cleaned up by terraform destroy. " +
		"Unreferenced binaries are automatically cleaned up by the FastEdge platform.")

	return nil
}
