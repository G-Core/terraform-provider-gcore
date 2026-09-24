package sweep

import (
	"fmt"
	"os"
)

const (
	// ResourcePrefix is the standard prefix for test resources.
	// The trailing separator matters; see IsTestResource in framework.go.
	ResourcePrefix = "tf-test-"
)

// ValidateSweeperEnvironment checks that all required environment variables are set.
// Call this at the beginning of each sweeper function.
//
// Returns nil if environment is valid, error otherwise.
func ValidateSweeperEnvironment() error {
	apiKey := os.Getenv("GCORE_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("GCORE_API_KEY must be set for sweepers")
	}

	projectID := os.Getenv("GCORE_CLOUD_PROJECT_ID")
	if projectID == "" {
		return fmt.Errorf("GCORE_CLOUD_PROJECT_ID must be set for sweepers")
	}

	regionID := os.Getenv("GCORE_CLOUD_REGION_ID")
	if regionID == "" {
		return fmt.Errorf("GCORE_CLOUD_REGION_ID must be set for sweepers")
	}

	return nil
}
