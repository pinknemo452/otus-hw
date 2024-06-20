package hw04lrucache

import (
	"reflect"
	"testing"
)

// require is banned by linter(

func requireFalse(t *testing.T, val bool) {
	t.Helper()
	if val {
		t.Fatal("Should be false")
	}
}

func requireTrue(t *testing.T, val bool) {
	t.Helper()
	if !val {
		t.Fatal("Should be true")
	}
}

func requireEqual(t *testing.T, expected interface{}, actual interface{}, _ ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatal("Should be equal")
	}
}

func requireNil(t *testing.T, object interface{}, _ ...interface{}) {
	t.Helper()
	if object == nil {
		return
	}
	if !reflect.ValueOf(object).IsNil() {
		t.Fatal("Should be nil")
	}
}
