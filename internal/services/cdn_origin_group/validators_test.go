package cdn_origin_group

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
)

func gcoreStorageConfig() *CDNOriginGroupSourcesConfigModel {
	return &CDNOriginGroupSourcesConfigModel{
		S3Type:       types.StringValue("gcore"),
		S3BucketName: types.StringValue("bucket"),
		StorageID:    types.Int64Value(123),
	}
}

func amazonConfig() *CDNOriginGroupSourcesConfigModel {
	return &CDNOriginGroupSourcesConfigModel{
		S3Type:            types.StringValue("amazon"),
		S3BucketName:      types.StringValue("bucket"),
		S3AccessKeyID:     types.StringValue("access-key"),
		S3SecretAccessKey: types.StringValue("secret-access-key"),
		S3Region:          types.StringValue("eu-west-1"),
	}
}

func s3Source(cfg *CDNOriginGroupSourcesConfigModel) CDNOriginGroupSourcesModel {
	return CDNOriginGroupSourcesModel{
		OriginType: types.StringValue("s3"),
		Config:     customfield.NewObjectMust(context.Background(), cfg),
	}
}

func originGroup(credentialsVersion types.Int64, sources ...CDNOriginGroupSourcesModel) *CDNOriginGroupModel {
	ctx := context.Background()
	return &CDNOriginGroupModel{
		Name:                 types.StringValue("group"),
		S3CredentialsVersion: credentialsVersion,
		ProxyNextUpstream:    customfield.NullList[types.String](ctx),
		Sources:              customfield.NewObjectListMust(ctx, sources),
	}
}

// unknownConfigAt marks the attribute of sources[0].config with the given name as unknown.
func unknownConfigAt(name string) func(tftypes.Value) tftypes.Value {
	return unknownAt(tftypes.AttributeName("sources"), tftypes.ElementKeyInt(0),
		tftypes.AttributeName("config"), tftypes.AttributeName(name))
}

// unknownAt marks the value at the given path as unknown.
func unknownAt(steps ...tftypes.AttributePathStep) func(tftypes.Value) tftypes.Value {
	want := tftypes.NewAttributePathWithSteps(steps)
	return func(v tftypes.Value) tftypes.Value {
		out, err := tftypes.Transform(v, func(p *tftypes.AttributePath, v tftypes.Value) (tftypes.Value, error) {
			if p.Equal(want) {
				return tftypes.NewValue(v.Type(), tftypes.UnknownValue), nil
			}
			return v, nil
		})
		if err != nil {
			panic(err)
		}
		return out
	}
}

func TestValidateConfigS3Type(t *testing.T) {
	noVersion := types.Int64Null()
	version := types.Int64Value(1)

	tests := []struct {
		name    string
		model   func() *CDNOriginGroupModel
		raw     func(tftypes.Value) tftypes.Value
		wantErr string
	}{
		{name: "valid gcore without credentials version", model: func() *CDNOriginGroupModel {
			return originGroup(noVersion, s3Source(gcoreStorageConfig()))
		}},
		{name: "valid gcore in upper case", model: func() *CDNOriginGroupModel {
			cfg := gcoreStorageConfig()
			cfg.S3Type = types.StringValue("GCORE")
			return originGroup(noVersion, s3Source(cfg))
		}},
		{name: "valid amazon", model: func() *CDNOriginGroupModel {
			return originGroup(version, s3Source(amazonConfig()))
		}},
		{name: "valid other", model: func() *CDNOriginGroupModel {
			cfg := amazonConfig()
			cfg.S3Type = types.StringValue("other")
			cfg.S3Region = types.StringNull()
			cfg.S3StorageHostname = types.StringValue("s3.example.com")
			return originGroup(version, s3Source(cfg))
		}},
		{name: "valid gcore next to host origin", model: func() *CDNOriginGroupModel {
			return originGroup(noVersion,
				CDNOriginGroupSourcesModel{Source: types.StringValue("example.com")},
				s3Source(gcoreStorageConfig()),
			)
		}},
		{name: "gcore without storage_id", wantErr: `storage_id is required when s3_type is "gcore"`, model: func() *CDNOriginGroupModel {
			cfg := gcoreStorageConfig()
			cfg.StorageID = types.Int64Null()
			return originGroup(noVersion, s3Source(cfg))
		}},
		{name: "gcore without bucket", wantErr: "'s3_bucket_name' is required", model: func() *CDNOriginGroupModel {
			cfg := gcoreStorageConfig()
			cfg.S3BucketName = types.StringNull()
			return originGroup(noVersion, s3Source(cfg))
		}},
		{name: "gcore with access key", wantErr: `s3_access_key_id is not allowed when s3_type is "gcore"`, model: func() *CDNOriginGroupModel {
			cfg := gcoreStorageConfig()
			cfg.S3AccessKeyID = types.StringValue("access-key")
			return originGroup(noVersion, s3Source(cfg))
		}},
		{name: "gcore with secret key", wantErr: `s3_secret_access_key is not allowed when s3_type is "gcore"`, model: func() *CDNOriginGroupModel {
			cfg := gcoreStorageConfig()
			cfg.S3SecretAccessKey = types.StringValue("secret-access-key")
			return originGroup(noVersion, s3Source(cfg))
		}},
		{name: "gcore with region", wantErr: `s3_region is not allowed when s3_type is "gcore"`, model: func() *CDNOriginGroupModel {
			cfg := gcoreStorageConfig()
			cfg.S3Region = types.StringValue("eu-west-1")
			return originGroup(noVersion, s3Source(cfg))
		}},
		{name: "gcore with storage hostname", wantErr: `s3_storage_hostname is not allowed when s3_type is "gcore"`, model: func() *CDNOriginGroupModel {
			cfg := gcoreStorageConfig()
			cfg.S3StorageHostname = types.StringValue("s3.example.com")
			return originGroup(noVersion, s3Source(cfg))
		}},
		{name: "gcore with empty region", wantErr: `s3_region is not allowed when s3_type is "gcore"`, model: func() *CDNOriginGroupModel {
			cfg := gcoreStorageConfig()
			cfg.S3Region = types.StringValue("")
			return originGroup(noVersion, s3Source(cfg))
		}},
		{name: "gcore with empty access key", wantErr: `s3_access_key_id is not allowed when s3_type is "gcore"`, model: func() *CDNOriginGroupModel {
			cfg := gcoreStorageConfig()
			cfg.S3AccessKeyID = types.StringValue("")
			return originGroup(noVersion, s3Source(cfg))
		}},
		{name: "gcore with unknown access key", raw: unknownConfigAt("s3_access_key_id"), model: func() *CDNOriginGroupModel {
			return originGroup(noVersion, s3Source(gcoreStorageConfig()))
		}},
		{name: "gcore with host_header_override", wantErr: `host_header_override is not allowed when s3_type is "gcore"`, model: func() *CDNOriginGroupModel {
			src := s3Source(gcoreStorageConfig())
			src.HostHeaderOverride = types.StringValue("x.example.com")
			return originGroup(noVersion, src)
		}},
		{name: "gcore with empty host_header_override", wantErr: `host_header_override is not allowed when s3_type is "gcore"`, model: func() *CDNOriginGroupModel {
			src := s3Source(gcoreStorageConfig())
			src.HostHeaderOverride = types.StringValue("")
			return originGroup(noVersion, src)
		}},
		{name: "gcore with unknown host_header_override", raw: unknownAt(tftypes.AttributeName("sources"), tftypes.ElementKeyInt(0), tftypes.AttributeName("host_header_override")), model: func() *CDNOriginGroupModel {
			src := s3Source(gcoreStorageConfig())
			src.HostHeaderOverride = types.StringValue("x.example.com")
			return originGroup(noVersion, src)
		}},
		{name: "amazon with host_header_override", model: func() *CDNOriginGroupModel {
			src := s3Source(amazonConfig())
			src.HostHeaderOverride = types.StringValue("x.example.com")
			return originGroup(version, src)
		}},
		{name: "amazon with storage_id", wantErr: `storage_id is only allowed when s3_type is "gcore"`, model: func() *CDNOriginGroupModel {
			cfg := amazonConfig()
			cfg.StorageID = types.Int64Value(123)
			return originGroup(version, s3Source(cfg))
		}},
		{name: "amazon without access key", wantErr: "'s3_access_key_id' is required", model: func() *CDNOriginGroupModel {
			cfg := amazonConfig()
			cfg.S3AccessKeyID = types.StringNull()
			return originGroup(version, s3Source(cfg))
		}},
		{name: "amazon without credentials version", wantErr: "s3_credentials_version is required", model: func() *CDNOriginGroupModel {
			return originGroup(noVersion, s3Source(amazonConfig()))
		}},
		{name: "amazon next to gcore without credentials version", wantErr: "s3_credentials_version is required", model: func() *CDNOriginGroupModel {
			return originGroup(noVersion, s3Source(gcoreStorageConfig()), s3Source(amazonConfig()))
		}},
		{name: "storage_id on host origin", wantErr: `'origin_type' must be "s3"`, model: func() *CDNOriginGroupModel {
			return originGroup(noVersion, CDNOriginGroupSourcesModel{
				Source: types.StringValue("example.com"),
				Config: customfield.NewObjectMust(context.Background(), &CDNOriginGroupSourcesConfigModel{
					StorageID: types.Int64Value(123),
				}),
			})
		}},
		// Unknown values must never fail validation.
		{name: "unknown s3_type without keys", raw: unknownConfigAt("s3_type"), model: func() *CDNOriginGroupModel {
			return originGroup(noVersion, s3Source(gcoreStorageConfig()))
		}},
		{name: "unknown storage_id", raw: unknownConfigAt("storage_id"), model: func() *CDNOriginGroupModel {
			return originGroup(noVersion, s3Source(gcoreStorageConfig()))
		}},
		{name: "unknown s3_type with keys without credentials version", raw: unknownConfigAt("s3_type"), wantErr: "s3_credentials_version is required", model: func() *CDNOriginGroupModel {
			return originGroup(noVersion, s3Source(amazonConfig()))
		}},
		{name: "amazon with unknown region", raw: unknownConfigAt("s3_region"), model: func() *CDNOriginGroupModel {
			return originGroup(version, s3Source(amazonConfig()))
		}},
		{name: "other with unknown storage hostname", raw: unknownConfigAt("s3_storage_hostname"), model: func() *CDNOriginGroupModel {
			cfg := amazonConfig()
			cfg.S3Type = types.StringValue("other")
			cfg.S3Region = types.StringNull()
			cfg.S3StorageHostname = types.StringValue("s3.example.com")
			return originGroup(version, s3Source(cfg))
		}},
	}

	ctx := context.Background()
	schema := ResourceSchema(ctx)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// A plan, unlike a state, can hold write-only values such as the access keys
			plan := tfsdk.Plan{Schema: schema, Raw: tftypes.NewValue(schema.Type().TerraformType(ctx), nil)}
			if diags := plan.Set(ctx, tt.model()); diags.HasError() {
				t.Fatal(diags)
			}
			raw := plan.Raw
			if tt.raw != nil {
				raw = tt.raw(raw)
			}

			resp := &resource.ValidateConfigResponse{}
			(&CDNOriginGroupResource{}).ValidateConfig(ctx, resource.ValidateConfigRequest{
				Config: tfsdk.Config{Schema: schema, Raw: raw},
			}, resp)

			// Terraform prints the config line of the attribute a diagnostic points at
			secret := path.Root("sources").AtListIndex(0).AtName("config").AtName("s3_secret_access_key")
			for _, d := range resp.Diagnostics {
				if withPath, ok := d.(diag.DiagnosticWithPath); ok && withPath.Path().Equal(secret) {
					t.Fatalf("diagnostic %q points at s3_secret_access_key, whose value Terraform would print", d.Summary())
				}
			}

			if tt.wantErr == "" {
				if resp.Diagnostics.HasError() {
					t.Fatalf("unexpected error: %v", resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.ErrorsCount() != 1 {
				t.Fatalf("want one error containing %q, got %v", tt.wantErr, resp.Diagnostics)
			}
			if d := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(d, tt.wantErr) {
				t.Fatalf("error %q does not contain %q", d, tt.wantErr)
			}
		})
	}
}
