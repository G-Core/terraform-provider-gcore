package storage_sftp_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/storage"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

const (
	// "lux" is an S3 location and the API rejects it here; SFTP is served from
	// mia / frn / ams.
	testAccStorageSftpLocation = "frn"
	testAccStorageSftpPassword = "Acc-Test-Pw-8Kd2xQ"
	testAccStorageSftpRotated  = "Acc-Test-Pw-9Lm4zR"
)

func TestAccStorageSftp_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpPassword, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("address"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("full_name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("provisioning_status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("location_name"), knownvalue.StringExact(testAccStorageSftpLocation)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo_version"), knownvalue.Int64Exact(1)),
				},
			},
			{
				Config: testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpPassword, 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo"), knownvalue.Null()),
				},
			},
		},
	})
}

// This test cannot prove that the PATCH contains the new password without an
// SFTP login. TestStorageSftpUpdateSendsPasswordDecision checks the body.
func TestAccStorageSftp_passwordRotation(t *testing.T) {
	rName := acctest.RandomName()
	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpPassword, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDSame.AddStateValue("gcore_storage_sftp.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpRotated, 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpRotated, 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_storage_sftp.test",
							plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("gcore_storage_sftp.test",
							tfjsonpath.New("has_password"), knownvalue.Bool(true)),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo_version"), knownvalue.Int64Exact(2)),
					compareIDSame.AddStateValue("gcore_storage_sftp.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpRotated, 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config:  testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpRotated, 2),
				Destroy: true,
			},
		},
	})
}

// the is_http_disabled steps double as the plain sibling-update coverage, so
// no separate update test exists
func TestAccStorageSftp_passwordRemoval(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpConfigPasswordAndHTTP(rName, testAccStorageSftpPassword, 1, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(true)),
				},
			},
			{
				Config: testAccStorageSftpConfigPasswordAndHTTP(rName, testAccStorageSftpPassword, 1, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_storage_sftp.test",
							plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("is_http_disabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo_version"), knownvalue.Int64Exact(1)),
				},
			},
			{
				Config: testAccStorageSftpConfigNoPasswordWithHTTP(rName, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_storage_sftp.test",
							plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("gcore_storage_sftp.test",
							tfjsonpath.New("has_password"), knownvalue.Bool(false)),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo_version"), knownvalue.Null()),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("is_http_disabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				Config: testAccStorageSftpConfigNoPasswordWithHTTP(rName, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccStorageSftp_passwordSetWithSiblingChange(t *testing.T) {
	rName := acctest.RandomName()
	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpConfigPasswordAndHTTP(rName, testAccStorageSftpPassword, 1, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("is_http_disabled"), knownvalue.Bool(false)),
					compareIDSame.AddStateValue("gcore_storage_sftp.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccStorageSftpConfigPasswordAndHTTP(rName, testAccStorageSftpRotated, 2, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_storage_sftp.test",
							plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("gcore_storage_sftp.test",
							tfjsonpath.New("has_password"), knownvalue.Bool(true)),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("is_http_disabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo_version"), knownvalue.Int64Exact(2)),
					compareIDSame.AddStateValue("gcore_storage_sftp.test", tfjsonpath.New("id")),
				},
			},
		},
	})
}

func TestAccStorageSftp_passwordAddedToPasswordlessStorage(t *testing.T) {
	rName := acctest.RandomName()
	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpConfigNoPassword(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo_version"), knownvalue.Null()),
					compareIDSame.AddStateValue("gcore_storage_sftp.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpPassword, 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_storage_sftp.test",
							plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("gcore_storage_sftp.test",
							tfjsonpath.New("has_password"), knownvalue.Bool(true)),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo_version"), knownvalue.Int64Exact(1)),
					compareIDSame.AddStateValue("gcore_storage_sftp.test", tfjsonpath.New("id")),
				},
			},
		},
	})
}

func TestAccStorageSftp_importedStorageTakesPassword(t *testing.T) {
	cases := map[string]struct {
		initialPassword   string
		hasPasswordBefore bool
	}{
		"with a password":    {initialPassword: testAccStorageSftpPassword, hasPasswordBefore: true},
		"without a password": {initialPassword: "", hasPasswordBefore: false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			adopted := &testAccAdoptedStorageSftp{t: t, name: acctest.RandomName()}
			compareIDSame := statecheck.CompareValue(compare.ValuesSame())

			resource.ParallelTest(t, resource.TestCase{
				PreCheck:                 func() { acctest.PreCheck(t) },
				ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
				CheckDestroy:             testAccCheckStorageSftpDestroy,
				Steps: []resource.TestStep{
					{
						PreConfig:  adopted.create(tc.initialPassword),
						ConfigFile: adopted.configFile(testAccStorageSftpConfigNoPassword(adopted.name)),
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("has_password"), knownvalue.Bool(tc.hasPasswordBefore)),
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("password_wo_version"), knownvalue.Null()),
							compareIDSame.AddStateValue("gcore_storage_sftp.test", tfjsonpath.New("id")),
						},
					},
					{
						ConfigFile: adopted.configFile(testAccStorageSftpConfigWithPassword(adopted.name, testAccStorageSftpRotated, 1)),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction("gcore_storage_sftp.test",
									plancheck.ResourceActionUpdate),
								plancheck.ExpectKnownValue("gcore_storage_sftp.test",
									tfjsonpath.New("has_password"), knownvalue.Bool(true)),
							},
						},
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("has_password"), knownvalue.Bool(true)),
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("password_wo"), knownvalue.Null()),
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("password_wo_version"), knownvalue.Int64Exact(1)),
							compareIDSame.AddStateValue("gcore_storage_sftp.test", tfjsonpath.New("id")),
						},
					},
					{
						ResourceName:       "gcore_storage_sftp.test",
						ImportState:        true,
						ImportStateKind:    resource.ImportBlockWithID,
						Config:             testAccStorageSftpConfigWithPassword(adopted.name, testAccStorageSftpRotated, 1),
						ExpectNonEmptyPlan: true,
						ImportPlanChecks: resource.ImportPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction("gcore_storage_sftp.test",
									plancheck.ResourceActionUpdate),
							},
						},
					},
				},
			})
		})
	}
}

func TestAccStorageSftp_adoptedPasswordLeftAlone(t *testing.T) {
	adopted := &testAccAdoptedStorageSftp{t: t, name: acctest.RandomName()}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:  adopted.create(testAccStorageSftpPassword),
				ConfigFile: adopted.configFile(testAccStorageSftpConfigNoPassword(adopted.name)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo_version"), knownvalue.Null()),
				},
			},
			{
				ConfigFile: adopted.configFile(testAccStorageSftpConfigNoPassword(adopted.name)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(true)),
				},
			},
		},
	})
}

func TestAccStorageSftp_adoptedStorageWithoutPasswordPlansCleanly(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpConfigNoPassword(rName),
			},
			{
				ResourceName:    "gcore_storage_sftp.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				Config:          testAccStorageSftpConfigNoPassword(rName),
			},
		},
	})
}

func TestAccStorageSftp_outOfBandPasswordRemovalIsRepaired(t *testing.T) {
	rName := acctest.RandomName()
	var id int64

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpPassword, 1),
				Check:  resource.TestCheckResourceAttrWith("gcore_storage_sftp.test", "id", testAccStorageSftpCaptureID(&id)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(true)),
				},
			},
			{
				PreConfig: testAccStorageSftpPatchPasswordMode(t, &id, "none"),
				Config:    testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpPassword, 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_storage_sftp.test",
							plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("gcore_storage_sftp.test",
							tfjsonpath.New("has_password"), knownvalue.Bool(true)),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password_wo_version"), knownvalue.Int64Exact(1)),
				},
			},
			{
				Config: testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpPassword, 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccStorageSftp_outOfBandPasswordRotationIsInvisible(t *testing.T) {
	rName := acctest.RandomName()
	var id int64

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpPassword, 1),
				Check:  resource.TestCheckResourceAttrWith("gcore_storage_sftp.test", "id", testAccStorageSftpCaptureID(&id)),
			},
			{
				PreConfig: testAccStorageSftpPatchPasswordMode(t, &id, "auto"),
				Config:    testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpPassword, 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("has_password"), knownvalue.Bool(true)),
				},
			},
		},
	})
}

func TestAccStorageSftp_invalidConfigs(t *testing.T) {
	cases := map[string]struct {
		attrs       string
		extra       string
		expectError *regexp.Regexp
	}{
		"password_wo without version": {
			attrs:       `password_wo = "Acc-Test-Pw-8Kd2xQ"`,
			expectError: regexp.MustCompile(`must be configured together`),
		},
		"password_wo_version without password_wo": {
			attrs:       `password_wo_version = 1`,
			expectError: regexp.MustCompile(`must be configured together`),
		},
		"unknown password_wo without version": {
			attrs: `password_wo = gcore_storage_sftp.unknown_source.full_name`,
			extra: fmt.Sprintf(`
resource "gcore_storage_sftp" "unknown_source" {
  name          = %[1]q
  location_name = %[2]q
}`, acctest.RandomName(), testAccStorageSftpLocation),
			expectError: regexp.MustCompile(`Missing password_wo_version`),
		},
		"password_wo too short": {
			attrs: `password_wo         = "short"
  password_wo_version = 1`,
			expectError: regexp.MustCompile(`must be between 8 and 63`),
		},
		"password_mode is not an argument": {
			attrs:       `password_mode = "set"`,
			expectError: regexp.MustCompile(`Unsupported argument|not expected here`),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:                 func() { acctest.PreCheck(t) },
				ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
resource "gcore_storage_sftp" "test" {
  name          = %[1]q
  location_name = %[2]q
  %[3]s
}
%[4]s`, acctest.RandomName(), testAccStorageSftpLocation, tc.attrs, tc.extra),
						PlanOnly:    true,
						ExpectError: tc.expectError,
					},
				},
			})
		})
	}
}

// This is the last release with the old SFTP password schema.
const testAccStorageSftpLegacyVersion = "2.0.0-rc.1"

// Run this test without a CLI dev override. An override would replace rc.1 with
// the local binary. The namespace must also match the provider address stored
// by rc.1 in state and the lock file.
func TestAccStorageSftp_upgradeFromV2RC1(t *testing.T) {
	t.Setenv(resource.EnvTfAccProviderNamespace, "g-core")

	cases := map[string]struct {
		legacyArgs  string
		hasPassword bool
	}{
		"auto": {legacyArgs: `password_mode = "auto"`, hasPassword: true},
		"set": {
			legacyArgs:  fmt.Sprintf("password_mode = \"set\"\n  sftp_password = %q", testAccStorageSftpPassword),
			hasPassword: true,
		},
		"none": {legacyArgs: `password_mode = "none"`, hasPassword: false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if name == "set" {
				testAccStorageSftpSkipUnlessLegacySetAccepted(t)
			}

			rName := acctest.RandomName()
			compareIDSame := statecheck.CompareValue(compare.ValuesSame())

			resource.ParallelTest(t, resource.TestCase{
				PreCheck:     func() { acctest.PreCheck(t) },
				CheckDestroy: testAccCheckStorageSftpDestroy,
				Steps: []resource.TestStep{
					{
						ExternalProviders: map[string]resource.ExternalProvider{
							"gcore": {Source: "G-Core/gcore", VersionConstraint: testAccStorageSftpLegacyVersion},
						},
						Config: testAccStorageSftpConfigLegacy(rName, tc.legacyArgs),
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("has_password"), knownvalue.Bool(tc.hasPassword)),
							compareIDSame.AddStateValue("gcore_storage_sftp.test", tfjsonpath.New("id")),
						},
					},
					{
						ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
						Config:                   testAccStorageSftpConfigNoPassword(rName),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectEmptyPlan(),
							},
						},
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckNoResourceAttr("gcore_storage_sftp.test", "password"),
							resource.TestCheckNoResourceAttr("gcore_storage_sftp.test", "password_mode"),
							resource.TestCheckNoResourceAttr("gcore_storage_sftp.test", "sftp_password"),
						),
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("has_password"), knownvalue.Bool(tc.hasPassword)),
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("password_wo_version"), knownvalue.Null()),
							compareIDSame.AddStateValue("gcore_storage_sftp.test", tfjsonpath.New("id")),
						},
					},
					{
						ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
						Config:                   testAccStorageSftpConfigWithPassword(rName, testAccStorageSftpRotated, 1),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction("gcore_storage_sftp.test",
									plancheck.ResourceActionUpdate),
								plancheck.ExpectKnownValue("gcore_storage_sftp.test",
									tfjsonpath.New("has_password"), knownvalue.Bool(true)),
							},
						},
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("has_password"), knownvalue.Bool(true)),
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("password_wo"), knownvalue.Null()),
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("password_wo_version"), knownvalue.Int64Exact(1)),
							compareIDSame.AddStateValue("gcore_storage_sftp.test", tfjsonpath.New("id")),
						},
					},
					{
						ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
						Config:                   testAccStorageSftpConfigNoPassword(rName),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction("gcore_storage_sftp.test",
									plancheck.ResourceActionUpdate),
							},
						},
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("has_password"), knownvalue.Bool(false)),
							statecheck.ExpectKnownValue("gcore_storage_sftp.test",
								tfjsonpath.New("password_wo_version"), knownvalue.Null()),
							compareIDSame.AddStateValue("gcore_storage_sftp.test", tfjsonpath.New("id")),
						},
					},
				},
			})
		})
	}
}

func testAccCheckStorageSftpDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_storage_sftp" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse id %q: %w", rs.Primary.ID, err)
		}

		_, err = client.Storage.SftpStorages.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("storage sftp %d still exists", id)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking sftp storage deletion: %w", err)
		}
	}

	return nil
}

func testAccStorageSftpCaptureID(id *int64) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("parsing id %q: %w", value, err)
		}
		*id = parsed
		return nil
	}
}

func testAccStorageSftpPatchPasswordMode(t *testing.T, id *int64, mode string) func() {
	t.Helper()
	return func() {
		client, err := acctest.NewGcoreClient()
		if err != nil {
			t.Fatalf("creating client: %s", err)
		}
		_, err = client.Storage.SftpStorages.UpdateAndPoll(
			context.Background(),
			*id,
			storage.SftpStorageUpdateParams{},
			option.WithRequestBody("application/json", []byte(fmt.Sprintf(`{"password_mode":%q}`, mode))),
		)
		if err != nil {
			t.Fatalf("out-of-band PATCH password_mode %q on storage %d: %s", mode, *id, err)
		}
	}
}

// An applied import block persists the adopted state. The testing library's
// import block is plan-only, while its persistent import conflicts with the
// storage created by PreConfig.
type testAccAdoptedStorageSftp struct {
	t    *testing.T
	name string
	id   int64
}

func (a *testAccAdoptedStorageSftp) create(password string) func() {
	a.t.Helper()
	return func() {
		client, err := acctest.NewGcoreClient()
		if err != nil {
			a.t.Fatalf("creating client: %s", err)
		}

		body := map[string]any{
			"name":          a.name,
			"location_name": testAccStorageSftpLocation,
			"password_mode": "none",
		}
		if password != "" {
			body["password_mode"] = "set"
			body["password"] = password
		}
		raw, err := json.Marshal(body)
		if err != nil {
			a.t.Fatalf("marshalling create body: %s", err)
		}

		created, err := client.Storage.SftpStorages.NewAndPoll(
			context.Background(),
			storage.SftpStorageNewParams{},
			option.WithRequestBody("application/json", raw),
		)
		if err != nil {
			a.t.Fatalf("creating storage %q outside Terraform: %s", a.name, err)
		}
		a.id = created.ID
	}
}

// ConfigFile runs before and after PreConfig, so its first call may use ID 0.
func (a *testAccAdoptedStorageSftp) configFile(resourceConfig string) config.TestStepConfigFunc {
	a.t.Helper()
	dir := a.t.TempDir()
	return func(req config.TestStepConfigRequest) string {
		file := filepath.Join(dir, fmt.Sprintf("step_%d.tf", req.StepNumber))
		content := fmt.Sprintf(`
import {
  to = gcore_storage_sftp.test
  id = %q
}
%s`, strconv.FormatInt(a.id, 10), resourceConfig)
		if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
			a.t.Fatalf("writing %s: %s", file, err)
		}
		return file
	}
}

// rc.1 sends sftp_password. Skip this case if the API no longer accepts that
// legacy field; the auto case still covers removal of a stored password.
func testAccStorageSftpSkipUnlessLegacySetAccepted(t *testing.T) {
	t.Helper()

	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}
	acctest.PreCheck(t)

	client, err := acctest.NewGcoreClient()
	if err != nil {
		t.Fatalf("creating client: %s", err)
	}

	ctx := context.Background()
	body := fmt.Sprintf(`{"name":%q,"location_name":%q,"password_mode":"set","sftp_password":%q}`,
		acctest.RandomName(), testAccStorageSftpLocation, testAccStorageSftpPassword)
	created, err := client.Storage.SftpStorages.NewAndPoll(
		ctx,
		storage.SftpStorageNewParams{},
		option.WithRequestBody("application/json", []byte(body)),
	)
	if err != nil {
		var apiErr *gcore.Error
		if errors.As(err, &apiErr) && apiErr.StatusCode == 400 {
			t.Skipf("the API rejects the rc.1 create body (sftp_password) with 400, so the %q starting state cannot be created by the rc.1 provider: %s", "set", err)
		}
		t.Fatalf("probing the legacy create body: %s", err)
	}
	if err := client.Storage.SftpStorages.DeleteAndPoll(ctx, created.ID); err != nil {
		t.Logf("deleting probe storage %d: %s (the sweeper will collect it)", created.ID, err)
	}
}

func testAccStorageSftpConfigWithPassword(name, password string, version int) string {
	return fmt.Sprintf(`
resource "gcore_storage_sftp" "test" {
  name                = %[1]q
  location_name       = %[2]q
  password_wo         = %[3]q
  password_wo_version = %[4]d
}`, name, testAccStorageSftpLocation, password, version)
}

func testAccStorageSftpConfigNoPassword(name string) string {
	return fmt.Sprintf(`
resource "gcore_storage_sftp" "test" {
  name          = %[1]q
  location_name = %[2]q
}`, name, testAccStorageSftpLocation)
}

func testAccStorageSftpConfigPasswordAndHTTP(name, password string, version int, httpDisabled bool) string {
	return fmt.Sprintf(`
resource "gcore_storage_sftp" "test" {
  name                = %[1]q
  location_name       = %[2]q
  password_wo         = %[3]q
  password_wo_version = %[4]d
  is_http_disabled    = %[5]t
}`, name, testAccStorageSftpLocation, password, version, httpDisabled)
}

func testAccStorageSftpConfigNoPasswordWithHTTP(name string, httpDisabled bool) string {
	return fmt.Sprintf(`
resource "gcore_storage_sftp" "test" {
  name             = %[1]q
  location_name    = %[2]q
  is_http_disabled = %[3]t
}`, name, testAccStorageSftpLocation, httpDisabled)
}

func testAccStorageSftpConfigLegacy(name, legacyArgs string) string {
	return fmt.Sprintf(`
resource "gcore_storage_sftp" "test" {
  name          = %[1]q
  location_name = %[2]q
  %[3]s
}`, name, testAccStorageSftpLocation, legacyArgs)
}
