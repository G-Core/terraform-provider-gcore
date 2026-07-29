package sweep_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	// Import sweeper registrations.
	//
	// Each sweeper registers itself from an init() in its own service package, so a
	// sweeper that is not imported here is never registered and can never run: the
	// sweep exits reporting success while leaving the resources behind. Every package
	// containing a sweep.go must be listed, and nothing in the build or lint step
	// catches an omission - keep this in sync when adding a sweeper.

	// CDN
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cdn_certificate"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cdn_logs_uploader_policy"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cdn_logs_uploader_target"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cdn_origin_group"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cdn_resource_rule"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cdn_trusted_ca_certificate"

	// Cloud
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_file_share"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_floating_ip"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_gpu_baremetal_cluster"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_gpu_baremetal_cluster_image"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_gpu_virtual_cluster"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_inference_registry_credential"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_inference_secret"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_instance"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_instance_image"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_k8s_cluster"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_load_balancer"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_network"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_network_router"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_network_subnet"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_placement_group"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_reserved_fixed_ip"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_secret"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_security_group"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_security_group_rule"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_ssh_key"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_volume"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/cloud_volume_snapshot"

	// DNS
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/dns_network_mapping"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/dns_zone"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/dns_zone_rrset"

	// FastEdge
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/fastedge_app"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/fastedge_binary"

	// Storage
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/storage_sftp"
	_ "github.com/G-Core/terraform-provider-gcore/internal/services/storage_ssh_key"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
