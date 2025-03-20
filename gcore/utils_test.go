package gcore

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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

type contextKey string

const (
	contextKeyExpected contextKey = "expected"
	contextKeyTesting  contextKey = "testing"
	contextKeyReturn   contextKey = "return"
)

func contextWithExpected[T any](ctx context.Context, expected T) context.Context {
	return context.WithValue(ctx, contextKeyExpected, expected)
}
func contextWithTesting(ctx context.Context, t *testing.T) context.Context {
	return context.WithValue(ctx, contextKeyTesting, t)
}
func contextWithReturn(ctx context.Context, ret *http.Response) context.Context {
	return context.WithValue(ctx, contextKeyReturn, ret)
}

func contextExpected[T any](ctx context.Context) T {
	return ctx.Value(contextKeyExpected).(T)
}
func contextTesting(ctx context.Context) *testing.T {
	return ctx.Value(contextKeyTesting).(*testing.T)
}
func contextReturn(ctx context.Context) *http.Response {
	return ctx.Value(contextKeyReturn).(*http.Response)
}

func compare[T any](t *testing.T, actual, expected T) {
	t.Helper()
	if !cmp.Equal(actual, expected) {
		t.Errorf("unexpected value, got: %v, expected: %v", actual, expected)
	}
}
