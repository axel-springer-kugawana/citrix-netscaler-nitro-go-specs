package nitro

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type Spec struct {
	Resources map[string]*Resource
	Bindings  map[string]*Binding
}

func ReadSpec(folder string) (*Spec, error) {
	resourcesFolder := filepath.Join(folder, "resources")
	bindingsFolder := filepath.Join(folder, "bindings")

	resourceFiles, err := ioutil.ReadDir(resourcesFolder)

	if err != nil {
		log.Println("Failed to enumerate resources folder : ", err)

		return nil, err
	}

	bindingFiles, err := ioutil.ReadDir(bindingsFolder)

	if err != nil {
		log.Println("Failed to enumerate bindings folder : ", err)

		return nil, err
	}

	resources := map[string]*Resource{}

	for _, file := range resourceFiles {
		fileName := filepath.Base(file.Name())
		resourceName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		resource, err := parseResource(filepath.Join(resourcesFolder, file.Name()))

		if err != nil {
			log.Println("Failed to parse resource : ", err)

			return nil, err
		}

		resources[resourceName] = resource
	}

	bindings := map[string]*Binding{}

	for _, file := range bindingFiles {
		fileName := filepath.Base(file.Name())
		bindingName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		binding, err := parseBinding(filepath.Join(bindingsFolder, file.Name()))

		if err != nil {
			log.Println("Failed to parse binding : ", err)

			return nil, err
		}

		bindings[bindingName] = binding
	}

	err = validateSpec(resources, bindings)

	if err != nil {
		log.Println("Failed to validate spec : ", err)

		return nil, err
	}

	spec := Spec{
		Resources: resources,
		Bindings: bindings,
	}

	return &spec, nil
}
