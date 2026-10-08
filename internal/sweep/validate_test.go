package sweep_test

import (
	"os"
	"strings"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/sweep"
)

// TestValidateSweeperEnvironment pins the value checks. Before them, any
// non-empty project or region id passed, and the sweep then ran against
// whatever the API made of it.
func TestValidateSweeperEnvironment(t *testing.T) {
	const unset = "<unset>" // never set as a value; means os.Unsetenv

	for _, tc := range []struct {
		name      string
		apiKey    string
		projectID string
		regionID  string
		wantErr   string // substring; empty means valid
	}{
		{"all valid", "key", "893427", "4", ""},

		{"api key unset", unset, "893427", "4", "GCORE_API_KEY must be set"},
		{"api key empty", "", "893427", "4", "GCORE_API_KEY must be set"},

		{"project unset", "key", unset, "4", "GCORE_CLOUD_PROJECT_ID must be an integer"},
		{"project empty", "key", "", "4", "GCORE_CLOUD_PROJECT_ID must be an integer"},
		{"project non-numeric", "key", "abc", "4", "GCORE_CLOUD_PROJECT_ID must be an integer"},
		{"project fractional", "key", "1.5", "4", "GCORE_CLOUD_PROJECT_ID must be an integer"},
		{"project zero", "key", "0", "4", "GCORE_CLOUD_PROJECT_ID must be positive"},
		{"project negative", "key", "-5", "4", "GCORE_CLOUD_PROJECT_ID must be positive"},

		{"region unset", "key", "893427", unset, "GCORE_CLOUD_REGION_ID must be an integer"},
		{"region non-numeric", "key", "893427", "eu", "GCORE_CLOUD_REGION_ID must be an integer"},
		{"region zero", "key", "893427", "0", "GCORE_CLOUD_REGION_ID must be positive"},
		{"region negative", "key", "893427", "-1", "GCORE_CLOUD_REGION_ID must be positive"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for name, value := range map[string]string{
				"GCORE_API_KEY":          tc.apiKey,
				"GCORE_CLOUD_PROJECT_ID": tc.projectID,
				"GCORE_CLOUD_REGION_ID":  tc.regionID,
			} {
				if value == unset {
					// t.Setenv first so the original value is restored
					// afterwards; it has no unset counterpart.
					t.Setenv(name, "")
					os.Unsetenv(name)
					continue
				}
				t.Setenv(name, value)
			}

			err := sweep.ValidateSweeperEnvironment()
			switch {
			case tc.wantErr == "" && err != nil:
				t.Errorf("ValidateSweeperEnvironment() = %v, want nil", err)
			case tc.wantErr != "" && err == nil:
				t.Errorf("ValidateSweeperEnvironment() = nil, want error containing %q", tc.wantErr)
			case tc.wantErr != "" && !strings.Contains(err.Error(), tc.wantErr):
				t.Errorf("ValidateSweeperEnvironment() = %q, want error containing %q", err, tc.wantErr)
			}
		})
	}
}
