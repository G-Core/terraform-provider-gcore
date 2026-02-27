package sweep_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	// Import sweeper registrations
	// CDN
	_ "github.com/stainless-sdks/gcore-terraform/internal/services/cdn_certificate"

	// Cloud
	_ "github.com/stainless-sdks/gcore-terraform/internal/services/cloud_ssh_key"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
