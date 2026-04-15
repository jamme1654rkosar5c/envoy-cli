package secrets

import (
	"regexp"
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// SensitivePatterns holds regex patterns that identify secret-like keys.
var SensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|pwd)`),
	regexp.MustCompile(`(?i)(secret|token|apikey|api_key)`),
	regexp.MustCompile(`(?i)(private_key|priv_key)`),
	regexp.MustCompile(`(?i)(auth|credential|credentials)`),
	regexp.MustCompile(`(?i)(access_key|access_token)`),
}

const redactedValue = "[REDACTED]"

// IsSensitive returns true if the given key matches any sensitive pattern.
func IsSensitive(key string) bool {
	for _, pattern := range SensitivePatterns {
		if pattern.MatchString(key) {
			return true
		}
	}
	return false
}

// Redact returns a copy of the EnvFile with sensitive values replaced.
func Redact(file parser.EnvFile) parser.EnvFile {
	redacted := make([]parser.EnvEntry, len(file.Entries))
	for i, entry := range file.Entries {
		if IsSensitive(entry.Key) {
			entry.Value = redactedValue
		}
		redacted[i] = entry
	}
	return parser.EnvFile{
		Path:    file.Path,
		Entries: redacted,
	}
}

// RedactValue masks a secret, revealing only the first 3 characters.
func RedactValue(value string) string {
	if len(value) <= 3 {
		return strings.Repeat("*", len(value))
	}
	return value[:3] + strings.Repeat("*", len(value)-3)
}
