package toygo

import (
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, except interface{}, actual interface{}) {
	if except != actual {
		t.Fatalf("except: \"%v\" (%s), actual: \"%v\" (%s)", except, reflect.TypeOf(except), actual, reflect.TypeOf(actual))
	}
}
