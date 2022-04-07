package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
)
import "gopkg.in/yaml.v2"

func unmarshalYaml(filename string) (map[string]interface{}, error) {
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

func marshalYaml(unmarshalledYaml map[string]interface{}, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)
	err := encoder.Encode(unmarshalledYaml)
	return err
}

func UpdateSecrets(filename string) (map[string]interface{}, error) {
	pipeline, err := unmarshalYaml(filename)
	if err != nil {
		return nil, err
	}
	walkAndUpdate(reflect.ValueOf(pipeline))

	return pipeline, nil
}

func walkAndUpdate(v reflect.Value) {
	if v.IsNil() {
		return
	}
	// dereference pointers (to pointers of pointers of pointers...)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			walkAndUpdate(v.Index(i))
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			elem := v.MapIndex(k).Elem()
			if elem.Kind() == reflect.String {
				secretMatcher := regexp.MustCompile(`^\(\((.*)\)\).*`)
				secretGroups := secretMatcher.FindStringSubmatch(elem.String())
				if len(secretGroups) > 1 {
					splitSections := strings.Split(secretGroups[1], ".")
					newSecret := fmt.Sprintf("((%s.%s))", splitSections[0], strings.Join(splitSections, "."))
					v.SetMapIndex(k, reflect.ValueOf(newSecret))
				}
			}

			walkAndUpdate(v.MapIndex(k))
		}
	default:
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE: update-secrets filename")
		return
	}

	filename := os.Args[1]
	updatedFile, err := UpdateSecrets(filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "update-secrets failed: %s\n", err)
		os.Exit(1)
	}
	err = marshalYaml(updatedFile, os.Stdout)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "update-secrets failed: %s\n", err)
		os.Exit(1)
	}
}
