package testutil

import (
	"reflect"
	"strings"
	"testing"
)

func AssertTrue(tb testing.TB, actual bool, name string) {
	AssertEquals(tb, true, actual, name)
}

func AssertFalse(tb testing.TB, actual bool, name string) {
	AssertEquals(tb, false, actual, name)
}

func AssertNil(tb testing.TB, actual interface{}, name string) {
	if actual == nil {
		return
	}

	rv := reflect.ValueOf(actual)
	if rv.Kind() == reflect.Pointer && rv.IsNil() {
		return
	}

	tb.Errorf("Incorrect result, expected %s to be '<nil>', got '%v'", name, actual)
}

func AssertNotNil(tb testing.TB, actual interface{}, name string) {
	if actual == nil {
		tb.Errorf("Incorrect result, expected %s to be not '<nil>', got '%v'", name, actual)
	}

	rv := reflect.ValueOf(actual)
	if rv.Kind() == reflect.Pointer && rv.IsNil() {
		tb.Errorf("Incorrect result, expected %s to be not '<nil>', got '%v'", name, actual)
	}
}

func AssertEquals(tb testing.TB, expected, actual interface{}, name string) {
	if expected != actual {
		tb.Errorf("Incorrect result, expected %s to be '%v', got '%v'", name, expected, actual)
	}
}

func AssertContains(tb testing.TB, s string, substr string) {
	if !strings.Contains(s, substr) {
		tb.Errorf("Incorrect result, expected '%s' to contains '%s'", s, substr)
	}
}
