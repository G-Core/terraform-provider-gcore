package sweep

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"

	"github.com/G-Core/terraform-provider-gcore/internal/custom"
)

// EnvAccountID is the account a sweep expects its credentials to resolve to.
//
// Run tokens already mean that wrong credentials mostly make a sweep find
// nothing rather than delete the wrong thing. The combination this guards is
// wrong credentials plus EnvSweepUntokenized, which would delete untagged
// resources in an account nobody intended to touch.
const EnvAccountID = "GCORE_SWEEP_ACCOUNT_ID"

// Target is the account, project and region a sweep operates on, resolved once
// and checked against what the operator said to expect.
type Target struct {
	Client    *gcore.Client
	ProjectID int64
	RegionID  int64
}

var (
	targetOnce sync.Once
	target     *Target
	targetErr  error
)

// Resolve returns the client and scope every sweeper should use. The account
// assertion costs one API call for the whole process, not one per sweeper.
func Resolve(ctx context.Context) (*Target, error) {
	targetOnce.Do(func() { target, targetErr = resolve(ctx) })
	return target, targetErr
}

func resolve(ctx context.Context) (*Target, error) {
	if err := ValidateSweeperEnvironment(); err != nil {
		return nil, err
	}

	projectID, err := envInt64("GCORE_CLOUD_PROJECT_ID")
	if err != nil {
		return nil, err
	}
	regionID, err := envInt64("GCORE_CLOUD_REGION_ID")
	if err != nil {
		return nil, err
	}

	client := gcore.NewClient(
		option.WithAPIKey(os.Getenv("GCORE_API_KEY")),
		option.WithCloudProjectID(projectID),
		option.WithCloudRegionID(regionID),
		custom.UserAgentOption("sweeper"),
	)

	if err := assertAccount(ctx, &client); err != nil {
		return nil, err
	}

	return &Target{Client: &client, ProjectID: projectID, RegionID: regionID}, nil
}

// assertAccount refuses to sweep until the credentials are shown to resolve to
// the expected account. A failure here is never routed through SkipSweepError:
// not knowing which account this is has to stop the sweep, not skip it.
func assertAccount(ctx context.Context, client *gcore.Client) error {
	want, ok := os.LookupEnv(EnvAccountID)
	if !ok || want == "" {
		return fmt.Errorf("%s must be set: refusing to sweep without an explicit account assertion", EnvAccountID)
	}

	overview, err := client.Iam.GetAccountOverview(ctx)
	if err != nil {
		return fmt.Errorf("verifying the sweep target account: %w", err)
	}

	if got := strconv.FormatInt(overview.ID, 10); got != want {
		return fmt.Errorf("refusing to sweep: credentials resolve to account %s (%q) but %s is %s",
			got, overview.CompanyName, EnvAccountID, want)
	}

	log.Printf("[INFO] sweep target confirmed: account %d (%q) is_test=%t",
		overview.ID, overview.CompanyName, overview.IsTest)
	return nil
}

func envInt64(name string) (int64, error) {
	raw := os.Getenv(name)
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer, got %q: %w", name, raw, err)
	}
	if v <= 0 {
		return 0, fmt.Errorf("%s must be positive, got %d", name, v)
	}
	return v, nil
}
