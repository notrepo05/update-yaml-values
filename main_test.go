package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"testing"
)

var parseTests = []struct {
	filename string
	expected []ResourceTypes
}{
	{"./fixtures/pipeline-1.yml", []ResourceTypes{{"gcs-resource"}, {"gcs-resource-2"}}},
	{"./fixtures/pipeline-2.yml", []ResourceTypes{{"gcs-resource"}}},
}

func loadYaml(filename string) (map[string]interface{}, error) {
	pipeline := map[string]interface{}{}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &pipeline)
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

func TestWalk(t *testing.T) {
	got, _ := Walk("./fixtures/pipeline-2.yml")
	expected, err := loadYaml("./fixtures/pipeline-2-updated.yml")
	if err != nil {
		t.Error("failed to load test data")
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got: %+v != expected: %+v", got, expected)
	}
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		got, err := Parse(test.filename)
		if err != nil {
			t.Errorf("err %+v", err)
		}
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("err: %+v != %+v", got, test.expected)
		}

	}
}
