package gcore

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

var baseSecretJson string = `
	"name": "top_secret",
	"comment": "Super-duper secret",
	"secret_slots": [
`

func slotJsonRsp(id int64, checksum string) string {
	return fmt.Sprintf(`{"slot": %d, "checksum": "%s"}`, id, checksum)
}

var baseSecret = sdk.Secret{
	Name:    ptr("top_secret"),
	Comment: ptr("Super-duper secret"),
}

func secretSlot(id int64, val string) sdk.SecretSlot {
	return sdk.SecretSlot{Slot: id, Value: &val}
}

var baseTfFastEdgeSecretConfig string = `resource "gcore_fastedge_secret" "test" {
	name = "top_secret"
	comment = "Super-duper secret"
`

func slotTf(id int64, val string) string {
	return fmt.Sprintf(`
	slot {
		id = %d
		value = "%s"
	}`, id, val)
}

var (
	secret0   = "secret_value0"
	secret1   = "secret_value1"
	checksum0 = "669b7529549384204b205a264c0b1df025905092478cbc2e2c421ab85b9d1c12"
	checksum1 = "e5727128c402caf486633d235da4728e4f9fdeed402adf8f1794887faa0680ce"
)

func TestFastEdgeSecret_basic(t *testing.T) {
	updatedSecret := baseSecret
	baseSecret.SecretSlots = &[]sdk.SecretSlot{
		secretSlot(0, secret0),
	}
	updatedSecret.SecretSlots = &[]sdk.SecretSlot{
		secretSlot(0, secret1),
	}
	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetSecret": {
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, ` + baseSecretJson + slotJsonRsp(0, checksum0) + `]}`,
					},
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, ` + baseSecretJson + slotJsonRsp(0, checksum0) + `]}`,
					},
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, ` + baseSecretJson + slotJsonRsp(0, checksum1) + `]}`,
					},
				},
			},
			"AddSecret": {
				params: []mockParams{
					{
						expectPayload: sdk.AddSecretJSONRequestBody{Secret: baseSecret},
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, ` + baseSecretJson + slotJsonRsp(0, checksum0) + `]}`,
					},
				},
			},
			"UpdateSecret": {
				params: []mockParams{
					{
						expectPayload: sdk.UpdateSecretJSONRequestBody{Secret: updatedSecret},
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, ` + baseSecretJson + slotJsonRsp(0, checksum1) + `]}`,
					},
				},
			},
			"DeleteSecret": {
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
			{ // create resource
				Config: baseTfFastEdgeSecretConfig + slotTf(0, secret0) + `
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_secret.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_secret.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_secret.test", "comment", "Super-duper secret"),
					resource.TestCheckResourceAttr("gcore_fastedge_secret.test", "slot.0.id", "0"),
					resource.TestCheckResourceAttr("gcore_fastedge_secret.test", "slot.0.value", ""),
					resource.TestCheckResourceAttr("gcore_fastedge_secret.test", "slot.0.checksum", "669b7529549384204b205a264c0b1df025905092478cbc2e2c421ab85b9d1c12"),
				),
			},
			{ // update resource
				Config: baseTfFastEdgeSecretConfig + slotTf(0, secret1) + `
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_secret.test"),
					resource.TestCheckResourceAttr("gcore_fastedge_secret.test", "id", "42"),
					resource.TestCheckResourceAttr("gcore_fastedge_secret.test", "comment", "Super-duper secret"),
					resource.TestCheckResourceAttr("gcore_fastedge_secret.test", "slot.0.id", "0"),
//					resource.TestCheckResourceAttr("gcore_fastedge_secret.test", "slot.0.value", ""),
					resource.TestCheckResourceAttr("gcore_fastedge_secret.test", "slot.0.checksum", "e5727128c402caf486633d235da4728e4f9fdeed402adf8f1794887faa0680ce"),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}

func TestFastEdgeSecret_disappear(t *testing.T) {
	updatedSecret := baseSecret
	baseSecret.SecretSlots = &[]sdk.SecretSlot{
		secretSlot(0, secret0),
	}
	updatedSecret.SecretSlots = &[]sdk.SecretSlot{
		secretSlot(0, secret1),
	}
	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetSecret": {
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, ` + baseSecretJson + slotJsonRsp(0, checksum0) + `]}`,
					},
					{
						expectId:  42,
						retStatus: http.StatusNotFound, // resource disappeared from the backend
					},
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, ` + baseSecretJson + slotJsonRsp(0, checksum1) + `]}`,
					},
				},
			},
			"AddSecret": {
				params: []mockParams{
					{
						expectPayload: sdk.AddSecretJSONRequestBody{Secret: baseSecret},
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, ` + baseSecretJson + slotJsonRsp(0, checksum0) + `]}`,
					},
					{
						expectPayload: sdk.AddSecretJSONRequestBody{Secret: updatedSecret},
						retStatus:     http.StatusOK,
						retBody:       `{"id": 42, ` + baseSecretJson + slotJsonRsp(0, checksum1) + `]}`,
					},
				},
			},
			"DeleteSecret": {
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
			{ // create resource
				Config: baseTfFastEdgeSecretConfig + slotTf(0, secret0) + `
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_secret.test"),
				),
			},
			{ // update resource
				Config: baseTfFastEdgeSecretConfig + slotTf(0, secret1) + `
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_secret.test"),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}

func TestFastEdgeSecret_import(t *testing.T) {
	mock := &mockSDK{
		t: t,
		mocks: map[string]*funcMock{
			"GetSecret": {
				params: []mockParams{
					{
						expectId:  42,
						retStatus: http.StatusOK,
						retBody:   `{"id": 42, ` + baseSecretJson + slotJsonRsp(0, checksum0) + `]}`,
					},
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: fastedgeMockProvider(mock),
		IsUnitTest:        true,
		Steps: []resource.TestStep{
			{ // import resource
				Config: baseTfFastEdgeSecretConfig + slotTf(0, secret0) + `
				}`,
				ImportState:   true,
				ImportStateId: "42",
				ResourceName:  "gcore_fastedge_secret.test",
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("gcore_fastedge_secret.test"),
				),
			},
		},
	})

	mock.ExpectationsWereMet(t)
}
