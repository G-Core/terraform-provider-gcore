package acctest

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// ProjectID returns the project ID from environment variable
// Used for constructing test configurations
func ProjectID() string {
	return os.Getenv("GCORE_CLOUD_PROJECT_ID")
}

// RegionID returns the region ID from environment variable
// Used for constructing test configurations
func RegionID() string {
	return os.Getenv("GCORE_CLOUD_REGION_ID")
}

// RandomName generates a random resource name for testing
// ensures test resources don't conflict when running concurrently
// uses RandString for short random suffix (10 alphanumeric chars)
// example: RandomName() -> "tf-test-a1b2c3d4e5"
func RandomName() string {
	return fmt.Sprintf("tf-test-%s", acctest.RandString(10))
}
