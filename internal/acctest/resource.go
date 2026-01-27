package acctest

import (
	"fmt"
	"os"
	"strconv"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// NewGcoreClient creates a new Gcore API client for testing
// This is used by CheckDestroy and other test utilities that need to
// interact with the Gcore API directly
func NewGcoreClient() (*gcore.Client, error) {
	apiKey := os.Getenv("GCORE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GCORE_API_KEY must be set")
	}

	projectIDStr := os.Getenv("GCORE_CLOUD_PROJECT_ID")
	if projectIDStr == "" {
		return nil, fmt.Errorf("GCORE_CLOUD_PROJECT_ID must be set")
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GCORE_CLOUD_PROJECT_ID: %w", err)
	}

	client := gcore.NewClient(
		option.WithAPIKey(apiKey),
		option.WithCloudProjectID(projectID),
	)

	return &client, nil
}

// CheckResourceDestroyed is a helper for CheckDestroy functions
// It verifies that a specific resource no longer exists in the Gcore API
//
// Usage in CheckDestroy:
//
//	func testAccCheckExampleDestroy(s *terraform.State) error {
//	    return acctest.CheckResourceDestroyed(s, "gcore_example", func(client *gcore.Client, id string) error {
//	        _, err := client.Cloud.Example.Get(context.Background(), id)
//	        return err
//	    })
//	}
//
// The checkFunc should return nil if the resource exists, or an error if it doesn't
// A "not found" error is expected and indicates successful deletion
func CheckResourceDestroyed(s *terraform.State, resourceType string, checkFunc func(*gcore.Client, string) error) error {
	client, err := NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceType {
			continue
		}

		// Try to fetch the resource
		err := checkFunc(client, rs.Primary.ID)

		// If no error, resource still exists - that's a failure
		if err == nil {
			return fmt.Errorf("resource %s (%s) still exists", resourceType, rs.Primary.ID)
		}

		// If error is "not found", that's success
		// Note: Different resources may have different error messages for "not found"
		// This is a basic check - specific resources may need custom logic
		if !IsNotFoundError(err) {
			return fmt.Errorf("error checking resource deletion: %w", err)
		}
	}

	return nil
}

// IsNotFoundError checks if an error indicates a resource was not found
// This is used to distinguish between "resource doesn't exist" (expected after deletion)
// and other errors (unexpected)
//
// Common Gcore API "not found" patterns:
// - HTTP 404 status code
// - Error messages containing "not found", "does not exist", etc.
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()

	// Common "not found" patterns in Gcore API
	notFoundPatterns := []string{
		"not found",
		"Not Found",
		"NOT_FOUND",
		"does not exist",
		"doesn't exist",
		"404",
		"NotFound",
		"NoSuchEntity",
	}

	for _, pattern := range notFoundPatterns {
		if contains(errMsg, pattern) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (case-sensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// CheckResourceExists is a helper for test Check functions
// It verifies that a specific resource exists in the Gcore API
//
// Usage in test steps:
//
//	Check: resource.ComposeAggregateTestCheckFunc(
//	    acctest.CheckResourceExists("gcore_example.test", func(client *gcore.Client, id string) error {
//	        _, err := client.Cloud.Example.Get(context.Background(), id)
//	        return err
//	    }),
//	    // ... other checks
//	)
func CheckResourceExists(resourceName string, checkFunc func(*gcore.Client, string) error) func(*terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is not set")
		}

		client, err := NewGcoreClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		return checkFunc(client, rs.Primary.ID)
	}
}

// BuildImportID is a helper for constructing import IDs
// Many Gcore resources use composite IDs like "project_id/resource_name"
//
// Usage:
//
//	ImportStateIdFunc: acctest.BuildImportID("gcore_example.test", "project_id", "name"),
func BuildImportID(resourceName string, attrNames ...string) func(*terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		var parts []string
		for _, attrName := range attrNames {
			value := rs.Primary.Attributes[attrName]
			if value == "" {
				return "", fmt.Errorf("attribute %s is empty", attrName)
			}
			parts = append(parts, value)
		}

		// Join with "/"
		result := ""
		for i, part := range parts {
			if i > 0 {
				result += "/"
			}
			result += part
		}

		return result, nil
	}
}
