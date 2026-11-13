package gcore

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTagsMergePatch(t *testing.T) {
	str := func(s string) *string { return &s }

	tests := []struct {
		name    string
		oldTags interface{}
		newTags interface{}
		want    map[string]*string
	}{
		{
			name:    "add tags",
			oldTags: map[string]interface{}{},
			newTags: map[string]interface{}{"key1": "val1"},
			want:    map[string]*string{"key1": str("val1")},
		},
		{
			name:    "change a value",
			oldTags: map[string]interface{}{"key1": "val1"},
			newTags: map[string]interface{}{"key1": "val2"},
			want:    map[string]*string{"key1": str("val2")},
		},
		{
			name:    "remove a key",
			oldTags: map[string]interface{}{"key1": "val1", "key2": "val2"},
			newTags: map[string]interface{}{"key1": "val1"},
			want:    map[string]*string{"key1": str("val1"), "key2": nil},
		},
		{
			name:    "replace all keys",
			oldTags: map[string]interface{}{"key1": "val1", "key2": "val2"},
			newTags: map[string]interface{}{"key3": "val3"},
			want:    map[string]*string{"key1": nil, "key2": nil, "key3": str("val3")},
		},
		{
			name:    "remove all keys",
			oldTags: map[string]interface{}{"key1": "val1"},
			newTags: map[string]interface{}{},
			want:    map[string]*string{"key1": nil},
		},
		{
			name:    "nil inputs",
			oldTags: nil,
			newTags: nil,
			want:    map[string]*string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tagsMergePatch(tt.oldTags, tt.newTags)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("tagsMergePatch() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExtractHosAndPath(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name     string
		args     args
		wantHost string
		wantPath string
		wantErr  bool
	}{
		{
			name: "long url success",
			args: args{
				uri: "https://test.url/with/path",
			},
			wantHost: "https://test.url",
			wantPath: "/with/path",
			wantErr:  false,
		},
		{
			name: "short url success",
			args: args{
				uri: "https://test.url",
			},
			wantHost: "https://test.url",
			wantPath: "",
			wantErr:  false,
		},
		{
			name: "error on empty",
			args: args{
				uri: "",
			},
			wantHost: "",
			wantPath: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotPath, err := ExtractHostAndPath(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractHostAndPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHost != tt.wantHost {
				t.Errorf("ExtractHostAndPath() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if gotPath != tt.wantPath {
				t.Errorf("ExtractHostAndPath() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}
