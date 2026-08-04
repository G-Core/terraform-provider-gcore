package cloud_instance_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_instance"
)

// Regression guard, not a schema policy: task_state was deliberately removed
// from the resource schema. It is transient background-task bookkeeping (null
// whenever no task is in flight) and is guaranteed to change during an apply,
// so it can neither be pinned with UseStateForUnknown nor left bare without
// failing Terraform's post-apply validation. It remains available on the
// gcore_cloud_instance data sources. Deliberate deletions live in custom code
// re-applied on every regeneration and have silently reappeared before; this
// test keeps that from shipping.
func TestCloudInstanceResourceForbiddenAttributes(t *testing.T) {
	t.Parallel()
	schema := cloud_instance.ResourceSchema(context.TODO())
	for _, name := range []string{"task_state"} {
		if _, ok := schema.Attributes[name]; ok {
			t.Errorf("attribute %q must not exist on the gcore_cloud_instance resource schema; it was deliberately removed and belongs on the data sources only", name)
		}
	}
}
