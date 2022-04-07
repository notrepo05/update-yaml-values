package main

import (
	"reflect"
	"testing"
)

var parseTests = []struct {
	filename string
	expected string
}{
	{"./fixtures/pipeline-2.yml", "./fixtures/pipeline-2-updated.yml"},
}

func TestWalk(t *testing.T) {
	for _, test := range parseTests {
		got, err := UpdateSecrets(test.filename)
		if err != nil {
			t.Errorf("err %+v", err)
		}
		expected, err := loadYaml(test.expected)
		if err != nil {
			t.Error("failed to load test data")
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("err: %+v != %+v", got, expected)
		}
	}
}
