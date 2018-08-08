package nitro

import (
	"fmt"
	"strings"
)

func validateFieldType(fieldType string, resources map[string]*SpecFile) bool {
	ok := false

	if strings.HasSuffix(fieldType, "[]") {
		fieldType = strings.TrimSuffix(fieldType, "[]")
	}

	if fieldType == "ip" || fieldType == "ip_mask" || fieldType == "int" || fieldType == "string" || fieldType == "bool" || fieldType == "double" {
		ok = true
	} else if strings.HasPrefix(fieldType, "(") && strings.HasSuffix(fieldType, ")") {
		ok = true
	} else if strings.Contains(fieldType, ".") {
		parts := strings.Split(fieldType, ".")

		if len(parts) == 2 {
			target, found := resources[parts[0]]

			if found {
				_, found = target.Fields[parts[1]]
			}

			if found {
				ok = true
			}
		}
	}

	return ok
}

func validateResources(resources map[string]*SpecFile) error {
	for key, resource := range resources {
		if resource.Scope == "" {
			return fmt.Errorf("Invalid resource spec, no scope defined : %v", key)
		}

		if resource.Key == nil {
			return fmt.Errorf("Invalid resource spec, no key name defined : %v", key)
		}

		for _, field := range resource.Key {
			_, ok := resource.Fields[field]
			if !ok {
				return fmt.Errorf("Invalid resource spec, key field unknown : %v.%v", key, field)
			}
		}

		if resource.Update != nil {
			for _, field := range resource.Update {
				_, ok := resource.Fields[field]
				if !ok {
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

		for _, operation := range resource.Operations {
			if operation == "rename" {
				if len(resource.Key) != 1 {
					return fmt.Errorf("Invalid resource spec, rename not supported when key is greater than one : %v", key)
				}
			} else if operation == "update" || operation == "unset" {
				if len(resource.Update) < 1 {
					return fmt.Errorf("Invalid resource spec, update/unset not supported when no fields are updatable : %v", key)
				}
			}
		}
	}

	return nil
}

func validateBindings(resources map[string]*SpecFile, bindings map[string]*SpecFile) error {
	for key, binding := range bindings {
		if binding.Scope == "" {
			return fmt.Errorf("Invalid binding spec, no scope defined : %v", key)
		}

		if binding.Key == nil {
			return fmt.Errorf("Invalid binding spec, no key defined : %v", key)
		}

		if binding.Update != nil {
			return fmt.Errorf("Invalid binding spec, update not supported : %v", key)
		}

		if binding.Operations != nil {
			return fmt.Errorf("Invalid binding spec, operations not supported : %v", key)
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

func validateSpec(resources map[string]*SpecFile, bindings map[string]*SpecFile) error {
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
