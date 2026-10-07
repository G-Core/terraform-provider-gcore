package sweep

import (
	"log"
	"os"
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

// ShouldSweep reports whether a resource may be deleted.
//
// A prefix match is not sufficient. Deletion also requires the resource to
// carry the run token named by EnvSweepRunID, so a sweep can only reclaim a run
// whose completion the caller has established — a scheduled job cannot know
// that, which is why no tokenised resource is swept when the variable is
// unset. The exception is EnvSweepUntokenized, which sweeps names carrying no
// token regardless; see its doc for why it cannot be scoped to a run.
func ShouldSweep(resourceType, name string) bool {
	if !IsTestResource(name) {
		log.Printf("[DEBUG] Skipping %s: %s (not a test resource)", resourceType, name)
		return false
	}

	token, tokenized := TokenFromName(name)
	if !tokenized {
		if os.Getenv(EnvSweepUntokenized) == "1" {
			log.Printf("[INFO] Sweeping %s: %s (predates run tokens, %s=1)", resourceType, name, EnvSweepUntokenized)
			return true
		}
		log.Printf("[INFO] Skipping %s: %s (predates run tokens; set %s=1 to reclaim)", resourceType, name, EnvSweepUntokenized)
		return false
	}

	want, ok := os.LookupEnv(EnvSweepRunID)
	if !ok || want == "" {
		log.Printf("[WARN] Skipping %s: %s (%s is unset; a sweep must name the finished run it reclaims)", resourceType, name, EnvSweepRunID)
		return false
	}

	if wantToken := Token(want); token != wantToken {
		log.Printf("[DEBUG] Skipping %s: %s (token %q, sweeping %q)", resourceType, name, token, wantToken)
		return false
	}

	log.Printf("[INFO] Sweeping %s: %s", resourceType, name)
	return true
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
