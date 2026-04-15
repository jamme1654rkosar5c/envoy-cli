package exporter

import (
	"encoding/json"
	"fmt"
	"os"\n	"strings"

	"github.com/envoy-cli/internal/parser"
)

// Format represents the output format for exporting env files.
type Format string

const (
	FormatDotEnv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatExport Format = "export"
)

// Export writes the entries of an EnvFile to the given file path in the specified format.
func Export(env *parser.EnvFile, destPath string, format Format) error {
	var content string
	var err error

	switch format {
	case FormatDotEnv:
		content = toDotEnv(env)
	case FormatJSON:
		content, err = toJSON(env)
		if err != nil {
			return fmt.Errorf("exporter: json marshal failed: %w", err)
		}
	case FormatExport:
		content = toExportScript(env)
	default:
		return fmt.Errorf("exporter: unsupported format %q", format)
	}

	if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("exporter: failed to write file: %w", err)
	}
	return nil
}

func toDotEnv(env *parser.EnvFile) string {
	var sb strings.Builder
	for _, entry := range env.Entries {
		sb.WriteString(fmt.Sprintf("%s=%s\n", entry.Key, entry.Value))
	}
	return sb.String()
}

func toJSON(env *parser.EnvFile) (string, error) {
	m := make(map[string]string, len(env.Entries))
	for _, entry := range env.Entries {
		m[entry.Key] = entry.Value
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

func toExportScript(env *parser.EnvFile) string {
	var sb strings.Builder
	for _, entry := range env.Entries {
		sb.WriteString(fmt.Sprintf("export %s=%q\n", entry.Key, entry.Value))
	}
	return sb.String()
}
