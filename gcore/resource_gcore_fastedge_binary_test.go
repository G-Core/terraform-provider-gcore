package gcore

import (
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestFastEdgeBinary_basic(t *testing.T) {
	checksum, _ := fileChecksum(os.Args[0])
	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetBinary": &funcMock{
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "checksum": "` + checksum + `"}`,
					},
				},
			},
			"StoreBinary": &funcMock{
				params: []mockParams{
					{
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "checksum": "` + checksum + `"}`,
					},
				},
			},
			"DelBinary": &funcMock{
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusNoContent,
					},
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: fastedgeMockProvider(mock),
		IsUnitTest:        true,
		Steps: []resource.TestStep{
			{
				Config: `resource "gcore_fastedge_binary" "test" {	filename = "` + os.Args[0] + `"}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gcore_fastedge_binary.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_binary.test", "checksum", checksum),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}

func TestFastEdgeBinary_corrupted(t *testing.T) {
	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"StoreBinary": &funcMock{
				params: []mockParams{
					{
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "checksum": "xyz"}`, // checksum cannot possibly match the real one
					},
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: fastedgeMockProvider(mock),
		IsUnitTest:        true,
		Steps: []resource.TestStep{
			{
				Config:      `resource "gcore_fastedge_binary" "test" {	filename = "` + os.Args[0] + `"}`,
				ExpectError: regexp.MustCompile(`uploaded binary checksum \(xyz\) does not match expected \(.*\), please retry`),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}

func TestFastEdgeBinary_disappear(t *testing.T) {
	checksum, _ := fileChecksum(os.Args[0])
	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetBinary": &funcMock{
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "checksum": "` + checksum + `"}`,
					},
					{
						expectId:  42,
						retStatus: http.StatusNotFound, // resource disappeared from the backend
					},
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "checksum": "` + checksum + `"}`,
					},
				},
			},
			"StoreBinary": &funcMock{
				params: []mockParams{
					{
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "checksum": "` + checksum + `"}`,
					},
					{
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "checksum": "` + checksum + `"}`,
					},
				},
			},
			"DelBinary": &funcMock{
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusNoContent,
					},
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: fastedgeMockProvider(mock),
		IsUnitTest:        true,
		Steps: []resource.TestStep{
			{
				Config: `resource "gcore_fastedge_binary" "test" {	filename = "` + os.Args[0] + `"}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gcore_fastedge_binary.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_binary.test", "checksum", checksum),
				),
			},
			{
				Config: `resource "gcore_fastedge_binary" "test" {	filename = "` + os.Args[0] + `"}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gcore_fastedge_binary.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_binary.test", "checksum", checksum),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}

func TestFastEdgeBinary_import(t *testing.T) {
	checksum, _ := fileChecksum(os.Args[0])
	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetBinary": &funcMock{
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, "checksum": "` + checksum + `"}`,
					},
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: fastedgeMockProvider(mock),
		IsUnitTest:        true,
		Steps: []resource.TestStep{
			{
				Config:        `resource "gcore_fastedge_binary" "test" {	filename = "` + os.Args[0] + `"}`,
				ImportState:   true,
				ImportStateId: "42",
				ResourceName:  "gcore_fastedge_binary.test",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gcore_fastedge_binary.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_binary.test", "checksum", checksum),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}
