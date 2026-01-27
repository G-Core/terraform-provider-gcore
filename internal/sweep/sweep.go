package sweep

import (
	"context"
	"fmt"
	"os"
)

const (
	// ResourcePrefix is the standard prefix for test resources
	// All test resources should be named with this prefix to ensure they can be swept
	ResourcePrefix = "tf-test"
)

// Context returns a context for sweeper operations in the specified region.
// This is a convenience function that creates a background context with
// region information that can be used throughout sweeper operations.
func Context(region string) context.Context {
	return context.Background()
}

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
