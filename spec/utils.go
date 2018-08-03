package nitro

import (
	"strings"
)

func IsArray(fieldType string) bool {
	return strings.HasSuffix(fieldType, "[]")
}

func IsEnum(fieldType string) bool {
	return strings.HasPrefix(fieldType, "(") && strings.HasSuffix(fieldType, ")")
}

func Name(name string) string {
	parts := strings.Split(name, "_")

	for i, part := range parts {
		parts[i] = strings.Title(part)
	}

	return strings.Join(parts, "")
}

func GoBaseType(fieldType string) string {
	fieldType = strings.TrimSuffix(fieldType, "[]")

	isEnum := IsEnum(fieldType)

	fieldType = strings.TrimPrefix(fieldType, "(")
	fieldType = strings.TrimSuffix(fieldType, ")")

	base := "string"

	switch fieldType {
	case "ip", "ip_mask", "string":
		base = "string"
	case "double", "int":
		base = "int"
	case "bool":
		base = "bool"
	}

	if isEnum {
		base = "string"
	}

	return base
}

func GoType(fieldType string) string {
	isArray := IsArray(fieldType)

	base := GoBaseType(fieldType)

	if isArray {
		base = "[]" + base
	}

	return base
}
