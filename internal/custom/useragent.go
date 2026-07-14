package custom

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/G-Core/gcore-go/option"
)

// sdkModulePath is the Go module path of the underlying Gcore SDK. Its version
// is resolved at runtime from the build info so the User-Agent always carries a
// concrete SDK version rather than a placeholder.
const sdkModulePath = "github.com/G-Core/gcore-go"

// UserAgentOption returns a request option that sets the User-Agent header to
// identify the Terraform provider and its version, following the Gcore
// convention: "terraform-provider-gcore/<providerVersion> Gcore/Go <sdkVersion>".
//
// It uses WithHeader (overwrite) rather than WithHeaderAdd because it composes
// the full single-value header itself; the SDK applies its default User-Agent
// first, and this cleanly replaces it.
func UserAgentOption(providerVersion string) option.RequestOption {
	return option.WithHeader(
		"User-Agent",
		fmt.Sprintf("terraform-provider-gcore/%s Gcore/Go %s", providerVersion, sdkVersion()),
	)
}

// sdkVersion resolves the version of the underlying Gcore SDK module from the
// binary's build info. It returns "unknown" when build info is unavailable or
// the dependency cannot be found.
func sdkVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	for _, dep := range info.Deps {
		if dep.Path == sdkModulePath {
			return strings.TrimPrefix(dep.Version, "v")
		}
	}
	return "unknown"
}
