package cdn_logs_uploader_target_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCDNLogsUploaderTarget_s3Other(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderTargetDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create S3 target with all config fields
			{
				Config: testAccCDNLogsUploaderTargetS3OtherConfig(rName, "Initial description"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("description"), knownvalue.StringExact("Initial description")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("storage_type"), knownvalue.StringExact("s3_other")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("client_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("bucket_name"), knownvalue.StringExact("acctest-bucket")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("region"), knownvalue.StringExact("us-east-1")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("endpoint"), knownvalue.StringExact("https://s3.example.com")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("use_path_style"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("directory"), knownvalue.StringExact("/logs/cdn")),
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 2: No drift
			{
				Config: testAccCDNLogsUploaderTargetS3OtherConfig(rName, "Initial description"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 3: Update config fields within same storage type
			{
				Config: testAccCDNLogsUploaderTargetS3OtherUpdatedConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("description"), knownvalue.StringExact("Updated description")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("bucket_name"), knownvalue.StringExact("acctest-bucket-v2")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("region"), knownvalue.StringExact("eu-west-1")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("endpoint"), knownvalue.StringExact("https://s3-v2.example.com")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("use_path_style"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("directory"), knownvalue.StringExact("/logs/cdn/v2")),
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 4: No drift after update
			{
				Config: testAccCDNLogsUploaderTargetS3OtherUpdatedConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 5: Remove optional fields (description, directory)
			{
				Config: testAccCDNLogsUploaderTargetS3OtherMinimalConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("description"), knownvalue.StringExact("")),
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 6: No drift after removing fields
			{
				Config: testAccCDNLogsUploaderTargetS3OtherMinimalConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 7: Import
			{
				ResourceName:            "gcore_cdn_logs_uploader_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCDNLogsUploaderTargetImportStateIDFunc("gcore_cdn_logs_uploader_target.test"),
				ImportStateVerifyIgnore: []string{"config.secret_access_key"},
			},
		},
	})
}

func TestAccCDNLogsUploaderTarget_sftp(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderTargetDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create SFTP target with password
			{
				Config: testAccCDNLogsUploaderTargetSFTPPasswordConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("storage_type"), knownvalue.StringExact("sftp")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("hostname"), knownvalue.StringExact("sftp.example.com")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("user"), knownvalue.StringExact("loguser")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("timeout_seconds"), knownvalue.Int64Exact(30)),
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 2: No drift (password no_refresh)
			{
				Config: testAccCDNLogsUploaderTargetSFTPPasswordConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 3: Update to private key auth
			{
				Config: testAccCDNLogsUploaderTargetSFTPKeyConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("hostname"), knownvalue.StringExact("sftp-key.example.com")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("user"), knownvalue.StringExact("keyuser")),
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 4: No drift (private_key + key_passphrase no_refresh)
			{
				Config: testAccCDNLogsUploaderTargetSFTPKeyConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccCDNLogsUploaderTarget_ftp(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNLogsUploaderTargetFTPConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("storage_type"), knownvalue.StringExact("ftp")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("hostname"), knownvalue.StringExact("ftp.example.com")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("user"), knownvalue.StringExact("ftpuser")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("timeout_seconds"), knownvalue.Int64Exact(15)),
				},
			},
			// No drift (password no_refresh)
			{
				Config: testAccCDNLogsUploaderTargetFTPConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccCDNLogsUploaderTarget_http(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderTargetDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create HTTP target with upload, append, retry, auth
			{
				Config: testAccCDNLogsUploaderTargetHTTPFullConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("storage_type"), knownvalue.StringExact("http")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("content_type"), knownvalue.StringExact("json")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("upload").AtMapKey("url"), knownvalue.StringExact("https://logs.example.com/upload")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("upload").AtMapKey("method"), knownvalue.StringExact("POST")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("upload").AtMapKey("timeout_seconds"), knownvalue.Int64Exact(60)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("upload").AtMapKey("use_compression"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("auth").AtMapKey("type"), knownvalue.StringExact("token")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("auth").AtMapKey("config").AtMapKey("header_name"), knownvalue.StringExact("Authorization")),
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 2: No drift (token no_refresh, response_actions defaults)
			{
				Config: testAccCDNLogsUploaderTargetHTTPFullConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 3: Update — change URLs, remove append/retry, update auth
			{
				Config: testAccCDNLogsUploaderTargetHTTPUpdatedConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("content_type"), knownvalue.StringExact("text")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("upload").AtMapKey("url"), knownvalue.StringExact("https://logs-v2.example.com/upload")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("upload").AtMapKey("method"), knownvalue.StringExact("PUT")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("auth").AtMapKey("config").AtMapKey("header_name"), knownvalue.StringExact("X-Auth-Token")),
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 4: No drift after update
			{
				Config: testAccCDNLogsUploaderTargetHTTPUpdatedConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 5: Add back append/retry blocks
			{
				Config: testAccCDNLogsUploaderTargetHTTPFullConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("append").AtMapKey("url"), knownvalue.StringExact("https://logs.example.com/append")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("retry").AtMapKey("url"), knownvalue.StringExact("https://logs.example.com/retry")),
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 6: No drift after adding blocks back
			{
				Config: testAccCDNLogsUploaderTargetHTTPFullConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccCDNLogsUploaderTarget_slsAkSk(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderTargetDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create SLS target with ak_sk auth
			{
				Config: testAccCDNLogsUploaderTargetSLSAkSkConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("storage_type"), knownvalue.StringExact("sls")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("region"), knownvalue.StringExact("eu-central-1")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("project"), knownvalue.StringExact("acctest-project")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("log_store"), knownvalue.StringExact("acctest-logstore")),
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("auth").AtMapKey("type"), knownvalue.StringExact("ak_sk")),
					// access_key_id is not masked by the API and must round-trip
					statecheck.ExpectKnownValue("gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("auth").AtMapKey("config").AtMapKey("access_key_id"),
						knownvalue.StringExact("acctest-alibaba-ak")),
					compareIDSame.AddStateValue(
						"gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 2: No drift (secret_access_key is masked on read; no_refresh must prevent a diff)
			{
				Config: testAccCDNLogsUploaderTargetSLSAkSkConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 3: Import — the masked secret cannot be verified against config
			{
				ResourceName:            "gcore_cdn_logs_uploader_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCDNLogsUploaderTargetImportStateIDFunc("gcore_cdn_logs_uploader_target.test"),
				ImportStateVerifyIgnore: []string{"config.auth.config.secret_access_key"},
			},
		},
	})
}

// Helpers

func testAccCDNLogsUploaderTargetImportStateIDFunc(resourceName string) func(*terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return rs.Primary.Attributes["id"], nil
	}
}

func testAccCheckCDNLogsUploaderTargetDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cdn_logs_uploader_target" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.Attributes["id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing id: %w", err)
		}

		_, err = client.CDN.LogsUploader.Targets.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("CDN logs uploader target %d still exists", id)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking CDN logs uploader target deletion: %w", err)
		}
	}
	return nil
}

// Config helpers

func testAccCDNLogsUploaderTargetS3OtherConfig(name, description string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_target" "test" {
  name         = %[1]q
  description  = %[2]q
  storage_type = "s3_other"

  config = {
    access_key_id     = "acctest-access-key"
    secret_access_key = "acctest-secret-key-12345"
    bucket_name       = "acctest-bucket"
    region            = "us-east-1"
    endpoint          = "https://s3.example.com"
    use_path_style    = true
    directory         = "/logs/cdn"
  }
}`, name, description)
}

func testAccCDNLogsUploaderTargetS3OtherUpdatedConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_target" "test" {
  name         = %[1]q
  description  = "Updated description"
  storage_type = "s3_other"

  config = {
    access_key_id     = "acctest-access-key"
    secret_access_key = "acctest-secret-key-12345"
    bucket_name       = "acctest-bucket-v2"
    region            = "eu-west-1"
    endpoint          = "https://s3-v2.example.com"
    use_path_style    = false
    directory         = "/logs/cdn/v2"
  }
}`, name)
}

func testAccCDNLogsUploaderTargetS3OtherMinimalConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_target" "test" {
  name         = %[1]q
  storage_type = "s3_other"

  config = {
    access_key_id     = "acctest-access-key"
    secret_access_key = "acctest-secret-key-12345"
    bucket_name       = "acctest-bucket-v2"
    region            = "eu-west-1"
    endpoint          = "https://s3-v2.example.com"
  }
}`, name)
}

func testAccCDNLogsUploaderTargetSFTPPasswordConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_target" "test" {
  name         = %[1]q
  storage_type = "sftp"

  config = {
    hostname        = "sftp.example.com"
    user            = "loguser"
    password        = "acctest-sftp-password"
    directory       = "/uploads/logs"
    timeout_seconds = 30
  }
}`, name)
}

func testAccCDNLogsUploaderTargetSFTPKeyConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_target" "test" {
  name         = %[1]q
  storage_type = "sftp"

  config = {
    hostname        = "sftp-key.example.com"
    user            = "keyuser"
    private_key     = "-----BEGIN OPENSSH PRIVATE KEY-----\nfake-key-content\n-----END OPENSSH PRIVATE KEY-----"
    key_passphrase  = "acctest-passphrase"
    directory       = "/secure-logs"
    timeout_seconds = 30
  }
}`, name)
}

func testAccCDNLogsUploaderTargetFTPConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_target" "test" {
  name         = %[1]q
  storage_type = "ftp"

  config = {
    hostname        = "ftp.example.com"
    user            = "ftpuser"
    password        = "acctest-ftp-password"
    directory       = "/logs"
    timeout_seconds = 15
  }
}`, name)
}

func testAccCDNLogsUploaderTargetHTTPFullConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_target" "test" {
  name         = %[1]q
  description  = "HTTP target acctest"
  storage_type = "http"

  config = {
    content_type = "json"

    upload = {
      url             = "https://logs.example.com/upload"
      method          = "POST"
      timeout_seconds = 60
      use_compression = true
      headers = {
        "X-Custom" = "test-value"
      }
      response_actions = [
        {
          action            = "retry"
          match_status_code = 500
          match_payload     = "error"
          description       = "Retry on server error"
        }
      ]
    }

    append = {
      url    = "https://logs.example.com/append"
      method = "PUT"
    }

    retry = {
      url             = "https://logs.example.com/retry"
      timeout_seconds = 10
    }

    auth = {
      type = "token"
      config = {
        token       = "acctest-auth-token-xyz"
        header_name = "Authorization"
      }
    }
  }
}`, name)
}

func testAccCDNLogsUploaderTargetSLSAkSkConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_target" "test" {
  name         = %[1]q
  storage_type = "sls"

  config = {
    region    = "eu-central-1"
    endpoint  = "eu-central-1.log.aliyuncs.com"
    project   = "acctest-project"
    log_store = "acctest-logstore"
    topic     = "acctest-topic"

    auth = {
      type = "ak_sk"
      config = {
        access_key_id     = "acctest-alibaba-ak"
        secret_access_key = "acctest-alibaba-sk-12345"
      }
    }
  }
}`, name)
}

func testAccCDNLogsUploaderTargetHTTPUpdatedConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_target" "test" {
  name         = %[1]q
  description  = "HTTP target acctest - updated"
  storage_type = "http"

  config = {
    content_type = "text"

    upload = {
      url             = "https://logs-v2.example.com/upload"
      method          = "PUT"
      timeout_seconds = 120
      use_compression = false
      headers = {
        "X-Custom" = "updated-value"
      }
    }

    auth = {
      type = "token"
      config = {
        token       = "acctest-rotated-token-456"
        header_name = "X-Auth-Token"
      }
    }
  }
}`, name)
}
