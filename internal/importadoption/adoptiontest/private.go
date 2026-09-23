package adoptiontest

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// PrivateState is the framework's private state as a test uses it. It matches
// importadoption.ReadWriter, so a Marker can Mark it directly.
type PrivateState interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

// NewPrivate installs an empty framework private state on req and returns it.
// req must point to a framework request struct with a Private field, such as
// *resource.ModifyPlanRequest or *planmodifier.StringRequest.
//
// The field's type, *privatestate.ProviderData, lives in an internal framework
// package, so a test cannot name it. reflect can still allocate one, and its
// zero value is ready to use: SetKey creates the map on first write. That
// lets a test drive the adopting path through the real request type instead
// of a pure-function stand-in.
func NewPrivate(t testing.TB, req any) PrivateState {
	t.Helper()
	v := reflect.ValueOf(req)
	if v.Kind() != reflect.Pointer || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		t.Fatalf("NewPrivate: want a non-nil pointer to a request struct, got %T", req)
		return nil
	}
	f := v.Elem().FieldByName("Private")
	if !f.IsValid() || f.Kind() != reflect.Pointer || !f.CanSet() {
		t.Fatalf("NewPrivate: %T has no settable Private pointer field", req)
		return nil
	}
	f.Set(reflect.New(f.Type().Elem()))
	p, ok := f.Interface().(PrivateState)
	if !ok {
		t.Fatalf("NewPrivate: %s does not implement GetKey and SetKey", f.Type())
		return nil
	}
	return p
}
