package cdn_origin_group

import (
	"context"
	"fmt"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.ResourceWithValidateConfig = (*CDNOriginGroupResource)(nil)

	sourcesPath = path.Root("sources")
)

// ValidateConfig performs cross-field validation on source entries within the sources list.
//
// The validations enforce:
//   - Host origins require source (hostname/IP). S3 origins must not set source (the API
//     ignores it and derives the endpoint from the config).
//   - origin_type = "s3" and config must appear together (S3 origins require credentials;
//     host origins must not have an S3 config block).
//   - s3_type = "amazon" requires s3_region (AWS needs a region to construct the endpoint).
//   - s3_type = "other" requires s3_storage_hostname (non-AWS S3 needs an explicit hostname).
//   - s3_credentials_version is required when any source has origin_type = "s3" (because
//     S3 credentials are write-only and this field is the only way to trigger re-sends).
func (r *CDNOriginGroupResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var sources customfield.NestedObjectList[CDNOriginGroupSourcesModel]
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, sourcesPath, &sources)...)
	if resp.Diagnostics.HasError() || sources.IsNull() || sources.IsUnknown() {
		return
	}

	sourceModels, diags := sources.AsStructSliceT(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasS3Source := false

	for i, src := range sourceModels {
		originType := src.OriginType.ValueString() // "" if null
		isS3 := originType == "s3"
		hasConfig := !src.Config.IsNull() && !src.Config.IsUnknown()
		hasSource := !src.Source.IsNull() && src.Source.ValueString() != ""

		// Host origins require source
		if !isS3 && !hasSource {
			resp.Diagnostics.AddAttributeError(
				sourcesPath.AtListIndex(i).AtName("source"),
				"Missing source",
				fmt.Sprintf("sources[%d]: source is required for host origins.", i),
			)
		}

		// origin_type = "s3" requires config
		if isS3 && !hasConfig {
			resp.Diagnostics.AddAttributeError(
				sourcesPath.AtListIndex(i).AtName("config"),
				"Missing S3 configuration",
				fmt.Sprintf("sources[%d]: config block is required when origin_type is \"s3\".", i),
			)
			continue
		}

		// config without origin_type = "s3" is invalid
		if hasConfig && !isS3 {
			resp.Diagnostics.AddAttributeError(
				sourcesPath.AtListIndex(i).AtName("config"),
				"Unexpected S3 configuration",
				fmt.Sprintf("sources[%d]: config block can only be specified when origin_type is \"s3\".", i),
			)
			continue
		}

		// Validate S3 config fields
		if isS3 && hasConfig {
			hasS3Source = true

			cfg, diags := src.Config.Value(ctx)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() || cfg == nil {
				continue
			}

			s3Type := cfg.S3Type.ValueString()
			switch s3Type {
			case "amazon":
				if cfg.S3Region.IsNull() || cfg.S3Region.ValueString() == "" {
					resp.Diagnostics.AddAttributeError(
						sourcesPath.AtListIndex(i).AtName("config").AtName("s3_region"),
						"Missing S3 region",
						fmt.Sprintf("sources[%d].config: s3_region is required when s3_type is \"amazon\".", i),
					)
				}
			case "other":
				if cfg.S3StorageHostname.IsNull() || cfg.S3StorageHostname.ValueString() == "" {
					resp.Diagnostics.AddAttributeError(
						sourcesPath.AtListIndex(i).AtName("config").AtName("s3_storage_hostname"),
						"Missing S3 storage hostname",
						fmt.Sprintf("sources[%d].config: s3_storage_hostname is required when s3_type is \"other\".", i),
					)
				}
			}
		}
	}

	// s3_credentials_version is required when any source is S3
	if hasS3Source {
		var credVersion types.Int64
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("s3_credentials_version"), &credVersion)...)
		if !resp.Diagnostics.HasError() && credVersion.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("s3_credentials_version"),
				"Missing S3 credentials version",
				"s3_credentials_version is required when any source has origin_type \"s3\". "+
					"S3 credentials are write-only (not stored in state), so this field must be "+
					"incremented to force Terraform to re-send credentials to the API.",
			)
		}
	}
}
