package sweep

import (
	"log"
	"strings"
)

// The only name prefixes a sweeper will delete. Unexported consts rather than
// a slice, so the allowlist cannot be widened at run time; an empty prefix
// would make every name match.
//
// The trailing separator is load-bearing: "tf-test" alone also matches
// unrelated names such as "tf-testing-prod".
const (
	testPrefixHyphen     = "tf-test-"
	testPrefixUnderscore = "tf_test_"
)

// IsTestResource reports whether name was created by this repo's acceptance
// tests.
//
// It is deliberately the narrowest predicate that still covers every name
// generator in the repo, and it must never guess: API keys are client-scoped
// and fourteen sweepers enumerate the whole account with no project boundary,
// so for those this is the only thing standing between a sweep and real
// infrastructure.
func IsTestResource(name string) bool {
	return strings.HasPrefix(name, testPrefixHyphen) ||
		strings.HasPrefix(name, testPrefixUnderscore)
}

// ShouldSweep logs and returns whether a resource should be swept.
//
// Parameters:
//   - resourceType: The Terraform resource type (e.g., "gcore_cloud_volume")
//   - name: The resource name to check
//
// Returns true if the resource should be deleted, false otherwise.
func ShouldSweep(resourceType, name string) bool {
	if IsTestResource(name) {
		log.Printf("[INFO] Sweeping %s: %s", resourceType, name)
		return true
	}

	log.Printf("[DEBUG] Skipping %s: %s (not a test resource)", resourceType, name)
	return false
}

// SkipSweepError determines if a sweeper error should be skipped.
// Some errors are expected and shouldn't cause the entire sweep to fail.
//
// Common skippable errors:
// - Access denied (sweeper doesn't have permission)
// - Resource not found (already deleted)
// - Region not available (region doesn't exist for this account)
//
// Returns true if the error should be ignored, false if it's a real failure.
func SkipSweepError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()

	// Skip access denied errors
	skipPatterns := []string{
		"Access Denied",
		"AccessDenied",
		"Forbidden",
		"403",
		"UnauthorizedOperation",
		"not authorized",
		"insufficient permissions",
		// Add more patterns as needed based on Gcore API errors
	}

	for _, pattern := range skipPatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}
