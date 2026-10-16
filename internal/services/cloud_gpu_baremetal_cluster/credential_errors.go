package cloud_gpu_baremetal_cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const (
	unsupportedCredentialUpdateSummary = "Unsupported credential update"
	unsupportedCredentialUpdateDetail  = "Updating username or password is not supported. These credentials can only be set during resource creation. To change them, the resource must be recreated."

	importedCredentialsHint = "\n\nThis cluster was imported: the API never returns the password or " +
		"password_wo_version, so the configured credentials cannot be matched against the cluster. Keep them as they " +
		"were at creation and add\n\n" +
		"  lifecycle {\n" +
		"    ignore_changes = [servers_settings.credentials]\n" +
		"  }\n\n" +
		"to the resource for as long as the cluster exists. See the Import section of the " +
		"gcore_cloud_gpu_baremetal_cluster documentation."

	// importedKey marks, in private state, a cluster that came in through
	// ImportState. Only an import can leave state that does not hold the
	// credentials the cluster was created with: Create stores the configured
	// ones, and the API never returns the password (nor the username when the
	// password went to the image's default user), so state alone cannot tell
	// an imported password cluster from one created without credentials.
	//
	// Retirement: none needed while the cluster lives, the credentials stay
	// unverifiable for its whole life. A replacement runs Create, which starts
	// with empty private state, so the new cluster is unmarked.
	importedKey = "gpu_baremetal_cluster_imported"
)

type privateReader interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
}

type privateWriter interface {
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

func markImported(ctx context.Context, w privateWriter) diag.Diagnostics {
	return w.SetKey(ctx, importedKey, []byte("true"))
}

// wasImported reports the import mark. A nil private state (unit tests, some
// framework paths) reads as not imported.
func wasImported(ctx context.Context, r privateReader) bool {
	if r == nil {
		return false
	}
	v, _ := r.GetKey(ctx, importedKey)
	return len(v) > 0
}

// unsupportedCredentialUpdateError is the error for a username or password
// change. It is resource-level on purpose: an attribute diagnostic would
// print the credentials block, including the password_wo line.
func unsupportedCredentialUpdateError(imported bool) diag.Diagnostic {
	detail := unsupportedCredentialUpdateDetail
	if imported {
		detail += importedCredentialsHint
	}
	return diag.NewErrorDiagnostic(unsupportedCredentialUpdateSummary, detail)
}
