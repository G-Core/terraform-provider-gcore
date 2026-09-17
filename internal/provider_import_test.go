package internal

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var importSupportSkipReasons = map[string]string{
	"gcore_cloud_file_share_access_rule":  "GCLOUD2-20036: API has no single-resource read endpoint",
	"gcore_storage_object_storage_bucket": "GCLOUD2-29373: blocked by STO-529 bucket create API behavior",
}

func TestRegisteredResourcesSupportImport(t *testing.T) {
	t.Parallel()

	repositoryRoot := repositoryRoot(t)
	registered := make(map[string]struct{})
	provider := &GcoreProvider{}

	for _, factory := range provider.Resources(context.Background()) {
		managedResource := factory()
		response := resource.MetadataResponse{}
		managedResource.Metadata(
			context.Background(),
			resource.MetadataRequest{ProviderTypeName: "gcore"},
			&response,
		)

		resourceName := response.TypeName
		if resourceName == "" {
			t.Error("registered resource returned an empty type name")
			continue
		}
		if _, exists := registered[resourceName]; exists {
			t.Errorf("resource %s is registered more than once", resourceName)
			continue
		}
		registered[resourceName] = struct{}{}

		_, implementsImport := managedResource.(resource.ResourceWithImportState)
		examplePath := filepath.Join(repositoryRoot, "examples", "resources", resourceName, "import.sh")
		hasExample := regularFileExists(t, examplePath)
		skipReason, skipped := importSupportSkipReasons[resourceName]

		switch {
		case skipped && strings.TrimSpace(skipReason) == "":
			t.Errorf("%s has an empty import skip reason", resourceName)
		case skipped && (implementsImport || hasExample):
			t.Errorf("%s now has partial or complete import support; remove it from importSupportSkipReasons", resourceName)
		case skipped:
			continue
		case implementsImport && !hasExample:
			t.Errorf("%s implements ResourceWithImportState but has no import.sh", resourceName)
		case !implementsImport && hasExample:
			t.Errorf("%s has import.sh but does not implement ResourceWithImportState", resourceName)
		case !implementsImport:
			t.Errorf("%s does not implement ResourceWithImportState and has no import.sh", resourceName)
		}
	}

	for resourceName := range importSupportSkipReasons {
		if _, exists := registered[resourceName]; !exists {
			t.Errorf("%s is in importSupportSkipReasons but is not a registered resource", resourceName)
		}
	}

	examplesRoot := filepath.Join(repositoryRoot, "examples", "resources")
	entries, err := os.ReadDir(examplesRoot)
	if err != nil {
		t.Fatalf("read resource examples: %v", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if !regularFileExists(t, filepath.Join(examplesRoot, entry.Name(), "import.sh")) {
			continue
		}
		if _, exists := registered[entry.Name()]; !exists {
			t.Errorf("%s has import.sh but is not a registered resource", entry.Name())
		}
	}
}

func repositoryRoot(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate provider import test")
	}
	return filepath.Dir(filepath.Dir(filename))
}

func regularFileExists(t *testing.T, path string) bool {
	t.Helper()

	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if !info.Mode().IsRegular() {
		t.Fatalf("%s exists but is not a regular file", path)
	}
	return true
}
