package gcore

import "testing"

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
