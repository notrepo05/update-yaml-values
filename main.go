package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
)
import "gopkg.in/yaml.v2"

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

func UpdateSecrets(filename string) (map[string]interface{}, error) {
	pipeline, _ := loadYaml(filename)

	walkAndUpdate(reflect.ValueOf(pipeline))

	fmt.Printf("\n\n+%v\n\n", pipeline)
	return pipeline, nil
}

func walkAndUpdate(v reflect.Value) {
	if v.IsNil() {
		return
	}
	// dereference pointers
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
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
				matchString, _ := regexp.MatchString(`^\(\(.*\)\).*`, elem.String())
				if matchString {
					v.SetMapIndex(k, reflect.ValueOf("newSecret"))
				}
			}

			walkAndUpdate(v.MapIndex(k))
		}
	default:
	}
}

func main() {
	fmt.Println("hello")
}
