package importadoption_test

import (
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/importadoption"
)

func TestKeepKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		body    string
		keys    []string
		want    string
		wantErr bool
	}{
		{name: "allowlist keeps only listed", body: `{"name":"n","tags":{"a":"b"},"type_name":"ssd_hiiops","image_id":"i"}`, keys: []string{"name", "tags"}, want: `{"name":"n","tags":{"a":"b"}}`},
		{name: "nothing listed present", body: `{"type_name":"ssd_hiiops"}`, keys: []string{"name", "tags"}, want: `{}`},
		{name: "empty object", body: `{}`, keys: []string{"name"}, want: `{}`},
		{name: "null passes through", body: `null`, keys: []string{"name"}, want: `null`},
		{name: "invalid json errors", body: `{`, keys: []string{"name"}, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := importadoption.KeepKeys([]byte(tc.body), tc.keys...)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && string(got) != tc.want {
				t.Errorf("KeepKeys = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestIsEmptyBody(t *testing.T) {
	t.Parallel()

	for body, want := range map[string]bool{
		"": true, "{}": true, " {} \n": true, "null": true,
		`{"a":1}`: false, `{"origin":"x"}`: false,
	} {
		if got := importadoption.IsEmptyBody([]byte(body)); got != want {
			t.Errorf("IsEmptyBody(%q) = %v, want %v", body, got, want)
		}
	}
}
