package cloud_gpu_baremetal_cluster_image

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// adoptingURL reports whether the plan records url from configuration over a
// null prior value. url is create-only and never returned by the API, so this
// only happens on the first apply after `terraform import`; a create has no
// prior state at all and is excluded by the caller.
func adoptingURL(prior, planned types.String) bool {
	return prior.IsNull() && !planned.IsNull() && !planned.IsUnknown()
}

// urlAdoptionWarning tells the user the one way adoption can go wrong: a
// configured url that is not the image's real source. Nothing is sent to the
// API for it and no refresh can correct it afterwards. The value itself is
// Sensitive, so the warning is attached to the resource and never quotes it.
//
// Attribute-level replacement decisions are not visible here, so the text
// also covers a plan where another attribute already forces a replacement:
// that path re-uploads from the configured url and adopts nothing.
func urlAdoptionWarning() diag.Diagnostic {
	return diag.NewWarningDiagnostic(
		"Adopting url into state after import",
		"Terraform is recording \"url\" from configuration because this image was imported "+
			"and the API does not return the source url.\n\n"+
			"Unless this plan replaces the image, nothing is sent to the API for it and the "+
			"configured value is trusted as-is. If it does not match the source the image was "+
			"uploaded from, state will be wrong and no future refresh can correct it.\n\n"+
			"Check the url against the image before applying. To re-upload from the configured "+
			"url instead, run `terraform apply -replace=` for this resource.",
	)
}
