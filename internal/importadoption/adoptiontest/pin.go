// Package adoptiontest ships the test helpers every resource uses to pin its
// import-adoption wiring, so the needles cannot drift between copies. It is a
// separate package from importadoption so the provider binary never links
// "testing" (the same split as internal/test_helpers vs the code it checks).
package adoptiontest

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Wiring names the literal, whitespace-stripped call sites a resource.go must
// carry, per method. Needles are matched against the method body only, sliced
// by its `func (r *<Receiver>) <Method>(` header up to the next top-level func.
type Wiring struct {
	Receiver string
	Needles  map[string][]string
}

// StandardWiring is the four-hook contract. planHook "" means the resource owns
// its ModifyPlan body and pins it with With().
func StandardWiring(receiver, marker, planHook string) Wiring {
	w := Wiring{Receiver: receiver, Needles: map[string][]string{
		"ImportState": {marker + ".Mark(ctx,resp.Private)"},
		"Read":        {marker + ".Advance(ctx,resp.Private)"},
		"Update":      {marker + ".BeginUpdate(ctx,req.Private,resp.Private,&resp.Diagnostics)", "defersettle()"},
	}}
	if planHook != "" {
		w.Needles["ModifyPlan"] = []string{planHook + ".ModifyPlan(ctx,req,resp)"}
	}
	return w
}

// With adds needles to a method (creating the entry when absent).
func (w Wiring) With(method string, needles ...string) Wiring {
	w.Needles[method] = append(w.Needles[method], needles...)
	return w
}

// AssertWired reads resource.go off disk and asserts every needle literally.
// resource.go is generated and these hooks are sealed custom hunks - small hunks
// in generated files are the class of change that silently evaporates when the
// seal is rebuilt against a regen. Behavioural tests call the helpers directly
// and would never notice the call sites vanishing.
func AssertWired(t testing.TB, resourceGoPath string, w Wiring) {
	t.Helper()
	raw, err := os.ReadFile(resourceGoPath)
	if err != nil {
		t.Fatalf("reading %s: %s", resourceGoPath, err)
	}
	AssertWiredSource(t, string(raw), w)
}

// AssertWiredSource is AssertWired over an in-memory source, so the helper can
// be proven to fail closed against a fixture.
func AssertWiredSource(t testing.TB, src string, w Wiring) {
	t.Helper()
	for method, needles := range w.Needles {
		body := StripSpace(MethodBody(t, src, w.Receiver, method))
		for _, n := range needles {
			if !strings.Contains(body, n) {
				t.Errorf("%s.%s no longer contains %q: a reseal dropped a sealed hunk; "+
					"either import plans a destroy+recreate or a day-two attribute set is "+
					"silently adopted into state", w.Receiver, method, n)
			}
		}
	}
}

// CountInMethod returns how many times a whitespace-stripped needle occurs in
// the named method of the source file - for hooks that must appear on every
// decode path of one method (the pool's Update has two).
func CountInMethod(t testing.TB, resourceGoPath, receiver, method, needle string) int {
	t.Helper()
	raw, err := os.ReadFile(resourceGoPath)
	if err != nil {
		t.Fatalf("reading %s: %s", resourceGoPath, err)
	}
	return strings.Count(StripSpace(MethodBody(t, string(raw), receiver, method)), needle)
}

// AssertOrderInMethod asserts `first` occurs before `second` (both
// whitespace-stripped) inside the named method, and that both are present.
func AssertOrderInMethod(t testing.TB, resourceGoPath, receiver, method, first, second string) {
	t.Helper()
	raw, err := os.ReadFile(resourceGoPath)
	if err != nil {
		t.Fatalf("reading %s: %s", resourceGoPath, err)
	}
	body := StripSpace(MethodBody(t, string(raw), receiver, method))
	i, j := strings.Index(body, first), strings.Index(body, second)
	switch {
	case i < 0:
		t.Errorf("%s.%s no longer contains %q", receiver, method, first)
	case j < 0:
		t.Errorf("%s.%s no longer contains %q", receiver, method, second)
	case i > j:
		t.Errorf("%s.%s: %q must come before %q", receiver, method, first, second)
	}
}

// AssertAttributeCallArgs reads schema.go off disk and asserts that inside the
// named top-level attribute's schema literal (`"<attribute>": schema.XAttribute{...}`
// directly under the schema's Attributes map) every call to exactly fn - the
// fully-qualified callee as printed by go/format, e.g.
// "planmodifiers.StringRequiresReplaceUnlessAdopting" - carries exactly
// wantArgs (whitespace-stripped, comma-joined; "" pins a no-argument call),
// and that at least one such call exists. %T pins the modifier type but not
// its constructor: a hand-resolved seal that passes the wrong private-state
// key, watches the wrong sibling attribute, swaps in a sibling constructor of
// the same type, or re-emits the framework's
// RequiresReplace() where RequiresReplaceIfConfigured() was (same %T) is
// invisible to it. The callee is matched whole, so a longer sibling name
// never satisfies the pin for the shorter one.
func AssertAttributeCallArgs(t testing.TB, schemaGoPath, attribute, fn, wantArgs string) {
	t.Helper()
	raw, err := os.ReadFile(schemaGoPath)
	if err != nil {
		t.Fatalf("reading %s: %s", schemaGoPath, err)
	}
	AssertAttributeCallArgsSource(t, string(raw), attribute, fn, wantArgs)
}

// AssertAttributeCallArgsSource is AssertAttributeCallArgs over an in-memory
// source.
func AssertAttributeCallArgsSource(t testing.TB, src, attribute, fn, wantArgs string) {
	t.Helper()
	got := AttributeCallArgs(t, src, attribute, fn)
	if len(got) == 0 {
		t.Errorf("attribute %q: no call to %s( in its schema literal", attribute, fn)
	}
	for _, args := range got {
		if args != StripSpace(wantArgs) {
			t.Errorf("attribute %q: %s(%s), want (%s)", attribute, fn, args, StripSpace(wantArgs))
		}
	}
}

// AttributeCallArgs returns the whitespace-stripped argument list of every
// call to exactly fn (the callee as printed by go/format, package-qualified)
// inside the named top-level attribute's schema literal - the
// `"<attribute>": ...` element of the `Attributes:` map of the first
// `schema.Schema{...}` literal in the file, so a same-named key under a nested
// attribute never stands in for it. Source that does not parse, or has no
// schema.Schema literal with an Attributes map, is fatal.
func AttributeCallArgs(t testing.TB, src, attribute, fn string) []string {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "schema.go", src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parsing source: %s", err)
		return nil
	}
	attrs := topLevelAttributes(fset, file)
	if attrs == nil {
		t.Fatalf("no schema.Schema{Attributes: map[string]schema.Attribute{...}} literal in source")
		return nil
	}
	var out []string
	for _, elt := range attrs.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.BasicLit)
		if !ok || key.Kind != token.STRING || key.Value != strconv.Quote(attribute) {
			continue
		}
		ast.Inspect(kv.Value, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			var name bytes.Buffer
			if err := format.Node(&name, fset, call.Fun); err != nil || name.String() != fn {
				return true
			}
			args := make([]string, 0, len(call.Args))
			for _, a := range call.Args {
				var buf bytes.Buffer
				if err := format.Node(&buf, fset, a); err != nil {
					t.Fatalf("printing argument of %s: %s", name.String(), err)
				}
				args = append(args, StripSpace(buf.String()))
			}
			out = append(out, strings.Join(args, ","))
			return true
		})
	}
	return out
}

// topLevelAttributes returns the composite literal of the Attributes map of
// the first schema.Schema{...} literal in the file, or nil.
func topLevelAttributes(fset *token.FileSet, file *ast.File) *ast.CompositeLit {
	var found *ast.CompositeLit
	ast.Inspect(file, func(n ast.Node) bool {
		if found != nil {
			return false
		}
		lit, ok := n.(*ast.CompositeLit)
		if !ok || lit.Type == nil {
			return true
		}
		var typ bytes.Buffer
		if err := format.Node(&typ, fset, lit.Type); err != nil || typ.String() != "schema.Schema" {
			return true
		}
		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			if id, ok := kv.Key.(*ast.Ident); ok && id.Name == "Attributes" {
				if m, ok := kv.Value.(*ast.CompositeLit); ok {
					found = m
				}
			}
		}
		return found == nil
	})
	return found
}

// AssertAdoptingModifiers pins, by %T, that each named attribute carries a
// planmodifiers.*RequiresReplaceUnlessAdopting modifier. Positive-only, fails
// closed: a reseal that re-emits the framework built-in necessarily drops ours.
func AssertAdoptingModifiers(t testing.TB, attrs map[string]schema.Attribute, names ...string) {
	t.Helper()
	for _, name := range names {
		got := ModifierTypeNames(t, attrs, name)
		found := false
		for _, g := range got {
			if strings.HasPrefix(g, "planmodifiers.") && strings.HasSuffix(g, "RequiresReplaceUnlessAdoptingModifier") {
				found = true
			}
		}
		if !found {
			t.Errorf("%s: missing planmodifiers.*RequiresReplaceUnlessAdopting(importAdoption.Key) - "+
				"`terraform import` plans a destroy+recreate or a day-two set is silently adopted; modifiers are %v", name, got)
		}
	}
}

// AssertModifierTypes pins an exact ordered %T list - used where the correct
// modifiers are framework built-ins and their ORDER matters (UseNonNull... first).
func AssertModifierTypes(t testing.TB, attrs map[string]schema.Attribute, name string, want ...string) {
	t.Helper()
	got := ModifierTypeNames(t, attrs, name)
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("%s: plan modifiers are %v, want %v", name, got, want)
	}
}

// AssertComputedOptional pins Computed && Optional on each named attribute.
func AssertComputedOptional(t testing.TB, attrs map[string]schema.Attribute, names ...string) {
	t.Helper()
	for _, name := range names {
		a, ok := attrs[name]
		if !ok {
			t.Fatalf("attribute %q not in schema", name)
		}
		if !a.IsComputed() || !a.IsOptional() {
			t.Errorf("%s: want Computed+Optional (removal from config must be a no-op), got computed=%v optional=%v",
				name, a.IsComputed(), a.IsOptional())
		}
	}
}

// ModifierTypeNames returns the %T of every plan modifier on a schema attribute.
func ModifierTypeNames(t testing.TB, attrs map[string]schema.Attribute, name string) []string {
	t.Helper()
	a, ok := attrs[name]
	if !ok {
		t.Fatalf("attribute %q not in schema", name)
	}
	var out []string
	switch v := a.(type) {
	case schema.StringAttribute:
		for _, m := range v.PlanModifiers {
			out = append(out, fmt.Sprintf("%T", m))
		}
	case schema.BoolAttribute:
		for _, m := range v.PlanModifiers {
			out = append(out, fmt.Sprintf("%T", m))
		}
	case schema.Int64Attribute:
		for _, m := range v.PlanModifiers {
			out = append(out, fmt.Sprintf("%T", m))
		}
	case schema.MapAttribute:
		for _, m := range v.PlanModifiers {
			out = append(out, fmt.Sprintf("%T", m))
		}
	case schema.ListAttribute:
		for _, m := range v.PlanModifiers {
			out = append(out, fmt.Sprintf("%T", m))
		}
	case schema.SingleNestedAttribute:
		for _, m := range v.PlanModifiers {
			out = append(out, fmt.Sprintf("%T", m))
		}
	default:
		t.Fatalf("attribute %q: unsupported type %T", name, a)
	}
	return out
}

// MethodBody returns the source of the named method's body, comments
// stripped: the file is parsed without comments and the *ast.FuncDecl with a
// pointer receiver of the named type is printed back. A needle inside a
// comment (a call commented out while resolving a seal conflict) therefore no
// longer satisfies a pin. A file that does not parse, or has no such method,
// is fatal - never silently green.
func MethodBody(t testing.TB, src, receiver, method string) string {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "resource.go", src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parsing source: %s", err)
		return "" // a testing.TB whose Fatalf does not stop the goroutine (the self-check recorder)
	}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != method || fn.Recv == nil || len(fn.Recv.List) != 1 || fn.Body == nil {
			continue
		}
		star, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			continue
		}
		if ident, ok := star.X.(*ast.Ident); !ok || ident.Name != receiver {
			continue
		}
		var buf bytes.Buffer
		if err := format.Node(&buf, fset, fn.Body); err != nil {
			t.Fatalf("printing %s.%s: %s", receiver, method, err)
			return ""
		}
		return buf.String()
	}
	t.Fatalf("could not find func (r *%s) %s(", receiver, method)
	return ""
}

// StripSpace removes all whitespace so the assertions survive gofmt re-wrapping
// without loosening what they match.
func StripSpace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}
