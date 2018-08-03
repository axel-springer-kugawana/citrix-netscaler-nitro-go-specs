package nitro

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Resource struct {
	Key    *Key
	State  string
	Fields map[string]string
	Update []string
}

type Key struct {
	Name   string
	Fields []string
}

type Binding struct {
	Key    []string
	Fields map[string]string
}

func readFile(file string) ([]byte, error) {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		log.Println("Failed to read the input file : ", err)
	}

	return b, err
}

func parseResource(file string) (*Resource, error) {
	data, err := readFile(file)

	if err != nil {
		return nil, err
	}

	resource := &Resource{}

	err = yaml.Unmarshal([]byte(data), &resource)

	if err != nil {
		log.Println("Failed to parse resource : ", err)

		return nil, err
	}

	return resource, err
}

func parseBinding(file string) (*Binding, error) {
	data, err := readFile(file)

	if err != nil {
		return nil, err
	}

	binding := &Binding{}

	err = yaml.Unmarshal([]byte(data), &binding)

	if err != nil {
		log.Println("Failed to parse binding : ", err)

		return nil, err
	}

	return binding, err
}
