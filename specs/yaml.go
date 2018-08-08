package nitro

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type SpecFile struct {
	Scope      string
	Fields     map[string]string
	Key        []string
	Update     []string
	Operations []string
}

func readFile(file string) ([]byte, error) {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		log.Println("Failed to read the input file : ", err)
	}

	return b, err
}

func parseSpecFile(file string) (*SpecFile, error) {
	data, err := readFile(file)

	if err != nil {
		return nil, err
	}

	resource := &SpecFile{}

	err = yaml.Unmarshal([]byte(data), &resource)

	if err != nil {
		log.Println("Failed to parse specfile : ", err)

		return nil, err
	}

	return resource, err
}
