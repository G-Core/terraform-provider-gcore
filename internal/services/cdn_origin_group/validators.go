package cdn_origin_group

import (
	"context"
	"fmt"
	"strings"

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
//   - Host origins require source (hostname/IP). S3 and FastEdge origins must not set source
//     (the API derives the endpoint from config).
//   - origin_type = "s3" and a config block must appear together; same for "fastedge".
//   - For origin_type = "s3": s3_type and s3_bucket_name are required in config.
//   - For origin_type = "fastedge": app_id is required in config.
//   - s3_type = "amazon" or "other" requires s3_access_key_id and s3_secret_access_key, and
//     does not allow storage_id.
//   - s3_type = "amazon" requires s3_region (AWS needs a region to construct the endpoint).
//   - s3_type = "other" requires s3_storage_hostname (non-AWS S3 needs an explicit hostname).
//   - s3_type = "gcore" requires storage_id and does not allow the access keys, s3_region,
//     s3_storage_hostname or the source's host_header_override (the CDN configures the
//     connection to the Gcore Object Storage).
//   - s3_credentials_version is required when any source has s3_type = "amazon" or "other"
//     (because S3 credentials are write-only and this field is the only way to trigger
//     re-sends).
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

	hasS3Credentials := false

	for i, src := range sourceModels {
		if src.OriginType.IsUnknown() {
			continue
		}
		originType := strings.ToLower(src.OriginType.ValueString())
		isS3 := originType == "s3"
		isFastedge := originType == "fastedge"
		isHost := !isS3 && !isFastedge // host or unset

		srcPath := sourcesPath.AtListIndex(i)
		configPath := srcPath.AtName("config")
		sourcePath := srcPath.AtName("source")

		hasConfig := !src.Config.IsNull() && !src.Config.IsUnknown()
		hasSource := !src.Source.IsNull() && !src.Source.IsUnknown() && src.Source.ValueString() != ""

		// 1. Host origins require source
		if isHost && !hasSource && !src.Source.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				sourcePath,
				"Missing source",
				fmt.Sprintf("sources[%d]: source is required for host origins.", i),
			)
		}

		// 2. s3 / fastedge origins require a config block
		if (isS3 || isFastedge) && !hasConfig && !src.Config.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				configPath,
				"Missing config block",
				fmt.Sprintf("sources[%d]: config block is required when origin_type is %q.", i, originType),
			)
			continue
		}

		if !hasConfig {
			continue
		}

		cfg, diags := src.Config.Value(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() || cfg == nil {
			continue
		}

		appIDSet := isStringSet(cfg.AppID)
		s3FieldsSet := isStringSet(cfg.S3Type) ||
			isStringSet(cfg.S3BucketName) ||
			isStringSet(cfg.S3AccessKeyID) ||
			isStringSet(cfg.S3SecretAccessKey) ||
			isInt64Set(cfg.StorageID)
		originTypePath := srcPath.AtName("origin_type")

		// 3a. app_id is only valid for fastedge sources
		if appIDSet && !isFastedge {
			resp.Diagnostics.AddAttributeError(
				originTypePath,
				"Origin type does not match config",
				fmt.Sprintf("sources[%d]: 'origin_type' must be \"fastedge\" when 'config.app_id' is set.", i),
			)
			continue
		}

		// 3b. S3 fields are only valid for s3 sources
		if s3FieldsSet && !isS3 {
			resp.Diagnostics.AddAttributeError(
				originTypePath,
				"Origin type does not match config",
				fmt.Sprintf("sources[%d]: 'origin_type' must be \"s3\" when S3 fields are set in config.", i),
			)
			continue
		}

		// 3c. host origins must not set a config block at all
		if isHost {
			resp.Diagnostics.AddAttributeError(
				configPath,
				"Unexpected config block",
				fmt.Sprintf("sources[%d]: config block can only be specified when origin_type is \"s3\" or \"fastedge\".", i),
			)
			continue
		}

		switch {
		case isS3:
			if validateS3Config(i, srcPath, originType, src, cfg, resp) {
				hasS3Credentials = true
			}

		case isFastedge:
			requireField(configPath, "app_id", originType, cfg.AppID, resp)
		}
	}

	// s3_credentials_version is required when any source sends S3 credentials
	if hasS3Credentials {
		var credVersion types.Int64
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("s3_credentials_version"), &credVersion)...)
		if !resp.Diagnostics.HasError() && credVersion.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("s3_credentials_version"),
				"Missing S3 credentials version",
				"s3_credentials_version is required when any source has s3_type \"amazon\" or \"other\". "+
					"S3 credentials are write-only (not stored in state), so this field must be "+
					"incremented to force Terraform to re-send credentials to the API.",
			)
		}
	}
}

// validateS3Config checks the config of an S3 source against the rules for its s3_type, and
// reports whether the source sends S3 credentials.
func validateS3Config(i int, srcPath path.Path, originType string, src CDNOriginGroupSourcesModel, cfg *CDNOriginGroupSourcesConfigModel, resp *resource.ValidateConfigResponse) bool {
	configPath := srcPath.AtName("config")
	requireField(configPath, "s3_type", originType, cfg.S3Type, resp)
	requireField(configPath, "s3_bucket_name", originType, cfg.S3BucketName, resp)

	// The rules below depend on s3_type, so skip them until it is known. Known keys
	// are still sent, so they still need s3_credentials_version.
	if cfg.S3Type.IsUnknown() {
		return isKnown(cfg.S3AccessKeyID) || isKnown(cfg.S3SecretAccessKey)
	}

	s3Type := strings.ToLower(cfg.S3Type.ValueString())
	if s3Type == "gcore" {
		if cfg.StorageID.IsNull() {
			resp.Diagnostics.AddAttributeError(
				configPath.AtName("storage_id"),
				"Missing storage ID",
				fmt.Sprintf("sources[%d].config: storage_id is required when s3_type is \"gcore\".", i),
			)
		}
		for name, val := range map[string]types.String{
			"s3_access_key_id":     cfg.S3AccessKeyID,
			"s3_secret_access_key": cfg.S3SecretAccessKey,
			"s3_region":            cfg.S3Region,
			"s3_storage_hostname":  cfg.S3StorageHostname,
		} {
			// An empty string is still sent to the API, so reject any known value. Report
			// on s3_type: Terraform prints the config line of the attribute an error
			// points at, which would expose a literal secret key.
			if isKnown(val) {
				resp.Diagnostics.AddAttributeError(
					configPath.AtName("s3_type"),
					"Unexpected S3 field",
					fmt.Sprintf("sources[%d].config: %s is not allowed when s3_type is \"gcore\". The CDN configures the connection to the storage.", i, name),
				)
			}
		}
		// The CDN sets the Host header from the storage, so the API rejects any override
		if isKnown(src.HostHeaderOverride) {
			resp.Diagnostics.AddAttributeError(
				srcPath.AtName("host_header_override"),
				"Unexpected host header override",
				fmt.Sprintf("sources[%d]: host_header_override is not allowed when s3_type is \"gcore\". The CDN sets the Host header from the storage.", i),
			)
		}
		return false
	}

	if isInt64Set(cfg.StorageID) {
		resp.Diagnostics.AddAttributeError(
			configPath.AtName("storage_id"),
			"Unexpected storage ID",
			fmt.Sprintf("sources[%d].config: storage_id is only allowed when s3_type is \"gcore\".", i),
		)
	}

	requireField(configPath, "s3_access_key_id", originType, cfg.S3AccessKeyID, resp)
	requireField(configPath, "s3_secret_access_key", originType, cfg.S3SecretAccessKey, resp)

	switch s3Type {
	case "amazon":
		if !cfg.S3Region.IsUnknown() && cfg.S3Region.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				configPath.AtName("s3_region"),
				"Missing S3 region",
				fmt.Sprintf("sources[%d].config: s3_region is required when s3_type is \"amazon\".", i),
			)
		}
	case "other":
		if !cfg.S3StorageHostname.IsUnknown() && cfg.S3StorageHostname.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				configPath.AtName("s3_storage_hostname"),
				"Missing S3 storage hostname",
				fmt.Sprintf("sources[%d].config: s3_storage_hostname is required when s3_type is \"other\".", i),
			)
		}
	}
	return true
}

func requireField(parent path.Path, name, originType string, val types.String, resp *resource.ValidateConfigResponse) {
	if val.IsUnknown() || (!val.IsNull() && val.ValueString() != "") {
		return
	}
	resp.Diagnostics.AddAttributeError(
		parent.AtName(name),
		"Missing required field",
		fmt.Sprintf("'%s' is required when 'origin_type' is %q.", name, originType),
	)
}

func isStringSet(v types.String) bool {
	return !v.IsNull() && !v.IsUnknown() && v.ValueString() != ""
}

func isInt64Set(v types.Int64) bool {
	return !v.IsNull() && !v.IsUnknown()
}

func isKnown(v types.String) bool {
	return !v.IsNull() && !v.IsUnknown()
}
