package common

import (
	"strings"
)

func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")

	for i := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(parts[i])
		} else if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
		}
	}

	return strings.Join(parts, "")
}

func ToPascalCase(s string) string {
	parts := strings.Split(s, "_")
	var result strings.Builder

	for _, part := range parts {
		if part == "" {
			continue
		}
		part = strings.ToLower(part)
		result.WriteString(strings.ToUpper(part[:1]))
		result.WriteString(part[1:])
	}

	return result.String()
}
