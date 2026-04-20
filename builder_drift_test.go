package flexgo

import (
	"reflect"
	"strings"
	"testing"
)

// TestBuilderCoversAllNodeFields fails when someone adds a field to
// Node without adding a matching setter to NodeBuilder. Keep the
// exclusion list tight — every new entry is a public API gap.
func TestBuilderCoversAllNodeFields(t *testing.T) {
	// Fields that intentionally have no direct builder setter.
	// Document the reason next to each.
	excluded := map[string]string{
		// Children has its own AddChildren/Children setters that take
		// varargs and are surfaced explicitly; reflection can't see the
		// pairing but both are present.
	}

	builderMethods := collectBuilderMethodNames()

	nodeType := reflect.TypeOf(Node{})
	for i := 0; i < nodeType.NumField(); i++ {
		field := nodeType.Field(i)
		if !field.IsExported() {
			continue
		}
		if _, ok := excluded[field.Name]; ok {
			continue
		}
		if !builderMethods[field.Name] {
			t.Errorf("Node field %q has no matching NodeBuilder setter", field.Name)
		}
	}
}

func collectBuilderMethodNames() map[string]bool {
	methods := map[string]bool{}
	t := reflect.TypeOf(&NodeBuilder{})
	for i := 0; i < t.NumMethod(); i++ {
		name := t.Method(i).Name
		// Treat exported methods as setters. AddChildren wraps Children.
		if name == "Build" || name == "AddChildren" {
			methods[strings.TrimPrefix(name, "Add")] = true
			continue
		}
		methods[name] = true
	}
	return methods
}
