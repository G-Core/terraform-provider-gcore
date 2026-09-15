package storage_sftp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/storage"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
)

// option.WithRequestBody replaces the whole request body, discarding any
// WithJSONSet applied before it. Update relies on that contract by appending
// passwordRequestOptions after the body; this test guards the contract, so a
// silent SDK change cannot drop the password fields.
func TestStorageSftpPasswordOptionsBeforeBodyAreLost(t *testing.T) {
	t.Parallel()

	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Errorf("request body is not a JSON object: %s: %q", err, raw)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":1}`))
	}))
	defer srv.Close()

	client := gcore.NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("unit-test"),
		option.WithMaxRetries(0),
	)

	state := StorageSftpModel{
		ID:             types.Int64Value(1),
		Name:           types.StringValue("keep"),
		LocationName:   types.StringValue("frn"),
		IsHTTPDisabled: types.BoolValue(false),
	}
	plan := state
	plan.IsHTTPDisabled = types.BoolValue(true)

	body, err := plan.MarshalJSONForUpdate(state)
	if err != nil {
		t.Fatalf("marshal: %s", err)
	}

	// the wrong order: password options first, body second
	opts := append(
		passwordRequestOptions(passwordClear, types.StringNull()),
		option.WithRequestBody("application/json", body),
	)
	if _, err := client.Storage.SftpStorages.Update(context.Background(), 1, storage.SftpStorageUpdateParams{}, opts...); err != nil {
		t.Fatalf("update: %s", err)
	}

	want := map[string]any{"is_http_disabled": true}
	if !reflect.DeepEqual(gotBody, want) {
		t.Fatalf("PATCH body:\n got %v\nwant %v\n(password_mode surviving means the SDK contract changed and the ordering in Update must be revisited)", gotBody, want)
	}
}

func TestStorageSftpPasswordEchoIsNeverDecoded(t *testing.T) {
	t.Parallel()

	echo := []byte(`{"id":1,"has_password":true,"password":"Echoed-Pw-8Kd2xQ"}`)

	created := &StorageSftpModel{}
	if err := apijson.UnmarshalComputed(echo, &created); err != nil {
		t.Fatalf("UnmarshalComputed: %s", err)
	}
	if !created.PasswordWo.IsNull() {
		t.Fatalf("UnmarshalComputed decoded the echoed password into password_wo: %v", created.PasswordWo)
	}
	if !created.HasPassword.Equal(types.BoolValue(true)) {
		t.Fatalf("UnmarshalComputed did not decode has_password: %v", created.HasPassword)
	}

	read := &StorageSftpModel{}
	if err := apijson.Unmarshal(echo, &read); err != nil {
		t.Fatalf("Unmarshal: %s", err)
	}
	if !read.PasswordWo.IsNull() {
		t.Fatalf("Unmarshal decoded the echoed password into password_wo: %v", read.PasswordWo)
	}
	if !read.HasPassword.Equal(types.BoolValue(true)) {
		t.Fatalf("Unmarshal did not decode has_password: %v", read.HasPassword)
	}
}
