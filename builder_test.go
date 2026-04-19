package flexgo

import (
	"reflect"
	"testing"
)

func TestBuilderHasMethodForEveryExportedNodeField(t *testing.T) {
	nodeType := reflect.TypeOf(Node{})
	builderType := reflect.TypeOf(&NodeBuilder{})

	for i := 0; i < nodeType.NumField(); i++ {
		field := nodeType.Field(i)
		if !field.IsExported() {
			continue
		}
		if _, ok := builderType.MethodByName(field.Name); !ok {
			t.Fatalf("NodeBuilder missing method for Node field %q", field.Name)
		}
	}
}
