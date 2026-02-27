package sweep_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	// Import sweeper registrations
	_ "github.com/stainless-sdks/gcore-terraform/internal/services/cdn_certificate"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
