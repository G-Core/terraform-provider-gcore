package sweep

import (
	"fmt"
	"os"
)

// ValidateSweeperEnvironment checks the environment a sweeper needs.
//
// It validates values rather than mere presence: a project or region id that is
// not a positive integer used to pass this check and then silently sweep
// whatever the API made of it.
func ValidateSweeperEnvironment() error {
	if os.Getenv("GCORE_API_KEY") == "" {
		return fmt.Errorf("GCORE_API_KEY must be set for sweepers")
	}
	if _, err := envInt64("GCORE_CLOUD_PROJECT_ID"); err != nil {
		return err
	}
	if _, err := envInt64("GCORE_CLOUD_REGION_ID"); err != nil {
		return err
	}
	return nil
}
