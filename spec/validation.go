package nitro

import (
	"fmt"
	"strings"
)

func validateFieldType(fieldType string, resources map[string]*Resource) bool {
	ok := false

	if strings.HasSuffix(fieldType, "[]") {
		fieldType = strings.TrimSuffix(fieldType, "[]")
	}

	if fieldType == "ip" || fieldType == "ip_mask"  || fieldType == "int"  || fieldType == "string"  || fieldType == "bool"  || fieldType == "double" {
		ok = true
	} else if strings.HasPrefix(fieldType, "(") && strings.HasSuffix(fieldType, ")")  {
		ok = true
	} else if strings.Contains(fieldType, ".") {
		parts := strings.Split(fieldType, ".")

		if len(parts) == 2 {
			target, found := resources[parts[0]]

			if found && parts[1] != target.Key.Name && parts[1] != target.State  {
				_, found = target.Fields[parts[1]]
			}

			if found {
				ok = true
			}
		}
	}

	return ok
}

func validateResources(resources map[string]*Resource) error {
	for key, resource := range resources {
		if resource.Key == nil || resource.Key.Name == "" {
			return fmt.Errorf("Invalid resource spec, no key name defined : %v", key)
		}

		if resource.Key.Fields != nil {
			for _, field := range resource.Key.Fields {
				_, ok := resource.Fields[field]
				if !ok {
					return fmt.Errorf("Invalid resource spec, key field unknown : %v.%v", key, field)
				}
			}
		}

		if resource.Update != nil {
			for _, field := range resource.Update {
				_, ok := resource.Fields[field]
				if !ok && field != resource.State {
					return fmt.Errorf("Invalid resource spec, update field unknown : %v.%v", key, field)
				}
			}
		}

		for field, fieldType := range resource.Fields {
			ok := validateFieldType(fieldType, resources)

			if !ok {
				return fmt.Errorf("Invalid resource spec, invalid field type : %v.%v (%v)", key, field, fieldType)
			}
		}
	}

	return nil
}

func validateBindings(resources map[string]*Resource, bindings map[string]*Binding) error {
	for key, binding := range bindings {
		if binding.Key == nil {
			return fmt.Errorf("Invalid binding spec, no key defined : %v", key)
		}

		for _, field := range binding.Key {
			_, ok := binding.Fields[field]
			if !ok {
				return fmt.Errorf("Invalid binding spec, key field unknown : %v.%v", key, field)
			}
		}

		for field, fieldType := range binding.Fields {
			ok := validateFieldType(fieldType, resources)

			if !ok {
				return fmt.Errorf("Invalid binding spec, invalid field type : %v.%v (%v)", key, field, fieldType)
			}
		}
	}

	return nil
}

func validateSpec(resources map[string]*Resource, bindings map[string]*Binding) error {
	err := validateResources(resources)

	if err != nil {
		return err
	}

	err = validateBindings(resources, bindings)

	if err != nil {
		return err
	}

	return nil
}
