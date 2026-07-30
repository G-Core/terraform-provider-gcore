package cloud_instance_test

import (
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// latestUbuntuImageID returns the ID of the latest public Ubuntu x86_64 image.
// The lookup is cached and skips the test when TF_ACC is unset; see
// acctest.LatestUbuntuImage.
func latestUbuntuImageID(t *testing.T) string {
	t.Helper()
	return acctest.LatestUbuntuImageID(t)
}
