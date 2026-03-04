package sweep_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	// Import sweeper registrations
	// CDN
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cdn_certificate"

	// Cloud
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_placement_group"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_ssh_key"

	// DNS
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/dns_zone"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/dns_zone_rrset"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
