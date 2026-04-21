package env

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// DigestAlgorithm specifies the hashing algorithm to use.
type DigestAlgorithm string

const (
	DigestMD5    DigestAlgorithm = "md5"
	DigestSHA256 DigestAlgorithm = "sha256"
)

// DigestOptions controls how the digest is computed.
type DigestOptions struct {
	Algorithm    DigestAlgorithm
	SortKeys     bool
	IncludeKeys  []string
	ExcludeKeys  []string
}

// DefaultDigestOptions returns sensible defaults.
func DefaultDigestOptions() DigestOptions {
	return DigestOptions{
		Algorithm: DigestSHA256,
		SortKeys:  true,
	}
}

// Digest computes a deterministic hash over the entries.
func Digest(entries []parser.Entry, opts DigestOptions) (string, error) {
	filtered := filterForDigest(entries, opts)

	if opts.SortKeys {
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Key < filtered[j].Key
		})
	}

	var sb strings.Builder
	for _, e := range filtered {
		sb.WriteString(e.Key)
		sb.WriteByte('=')
		sb.WriteString(e.Value)
		sb.WriteByte('\n')
	}

	raw := sb.String()
	switch opts.Algorithm {
	case DigestMD5:
		sum := md5.Sum([]byte(raw))
		return fmt.Sprintf("%x", sum), nil
	case DigestSHA256:
		sum := sha256.Sum256([]byte(raw))
		return fmt.Sprintf("%x", sum), nil
	default:
		return "", fmt.Errorf("unsupported digest algorithm: %s", opts.Algorithm)
	}
}

func filterForDigest(entries []parser.Entry, opts DigestOptions) []parser.Entry {
	include := toSet(opts.IncludeKeys)
	exclude := toSet(opts.ExcludeKeys)

	var out []parser.Entry
	for _, e := range entries {
		if len(include) > 0 && !include[e.Key] {
			continue
		}
		if exclude[e.Key] {
			continue
		}
		out = append(out, e)
	}
	return out
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
