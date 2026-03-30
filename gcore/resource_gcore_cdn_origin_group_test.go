//go:build !cloud
// +build !cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOriginGroup(t *testing.T) {
	fullName := "gcore_cdn_origingroup.acctest"

	type Params struct {
		Source  string
		Enabled string
	}

	create := Params{"google.com", "true"}
	update := Params{"tut.by", "false"}

	template := func(params *Params) string {
		return fmt.Sprintf(`
            resource "gcore_cdn_origingroup" "acctest" {
			  name = "terraform_acctest_group"
			  use_next = true

			  origin {
			    source = "%s"
				enabled = %s
			  }

			  origin {
			    source = "yandex.ru"
			    enabled = true
			  }
			}
		`, params.Source, params.Enabled)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_CDN_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: template(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "terraform_acctest_group"),
					resource.TestCheckResourceAttr(fullName, "origin.0.source", create.Source),
					resource.TestCheckResourceAttr(fullName, "origin.0.enabled", create.Enabled),
					resource.TestCheckResourceAttr(fullName, "origin.1.source", "yandex.ru"),
					resource.TestCheckResourceAttr(fullName, "origin.1.enabled", "true"),
				),
			},
			{
				Config: template(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "terraform_acctest_group"),
					resource.TestCheckResourceAttr(fullName, "origin.0.source", update.Source),
					resource.TestCheckResourceAttr(fullName, "origin.0.enabled", update.Enabled),
					resource.TestCheckResourceAttr(fullName, "origin.1.source", "yandex.ru"),
					resource.TestCheckResourceAttr(fullName, "origin.1.enabled", "true"),
				),
			},
		},
	})
}

func TestAccOriginGroupMixed(t *testing.T) {
	fullName := "gcore_cdn_origingroup.acctest_mixed"

	config := `
		resource "gcore_cdn_origingroup" "acctest_mixed" {
		  name     = "terraform_acctest_mixed_group"
		  use_next = true

		  origin {
		    source  = "cdn.example.com"
		    enabled = true
		  }

		  origin {
		    origin_type = "s3"
		    enabled     = true
		    backup      = true
		    config {
		      s3_type              = "amazon"
		      s3_bucket_name       = "test-bucket"
		      s3_region            = "eu-west-1"
		      s3_access_key_id     = "AKIAIOSFODNN7EXAMPLE"
		      s3_secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
		    }
		  }
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_CDN_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "terraform_acctest_mixed_group"),
					resource.TestCheckResourceAttr(fullName, "origin.#", "2"),
				),
			},
		},
	})
}

func TestAccOriginGroupS3Origin(t *testing.T) {
	fullName := "gcore_cdn_origingroup.acctest_s3"

	config := `
		resource "gcore_cdn_origingroup" "acctest_s3" {
		  name     = "terraform_acctest_s3_group"
		  use_next = true

		  origin {
		    origin_type = "s3"
		    enabled     = true
		    config {
		      s3_type              = "amazon"
		      s3_bucket_name       = "test-bucket"
		      s3_region            = "eu-west-1"
		      s3_access_key_id     = "AKIAIOSFODNN7EXAMPLE"
		      s3_secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
		    }
		  }
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_CDN_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "terraform_acctest_s3_group"),
					resource.TestCheckResourceAttr(fullName, "origin.#", "1"),
				),
			},
		},
	})
}

