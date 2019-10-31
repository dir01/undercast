package main

import (
	"reflect"
	"testing"
)

func assertDeepEquals(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Deep equality test failed.\nGot %#v\nWant %#v", actual, expected)
	}
}
