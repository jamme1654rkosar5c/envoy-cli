// Package env provides utilities for managing, transforming, and inspecting
// environment variable entries.
//
// # Classifier
//
// The Classifier assigns a semantic Category to each env entry based on its
// key name, making it easy to audit and group variables by concern.
//
// Built-in categories:
//
//   - secret   – keys containing SECRET, PASSWORD, TOKEN, API_KEY, etc.
//   - database – keys prefixed with DB_, DATABASE_, or named after databases
//   - network  – keys related to HOST, PORT, URL, ENDPOINT, etc.
//   - feature  – keys prefixed with FEATURE_, FLAG_, ENABLE_, DISABLE_
//   - general  – everything else
//
// Custom rules can be supplied via ClassifyOptions.CustomRules, which maps a
// category name to a list of key prefixes that take precedence over built-ins.
//
// Example:
//
//	opts := env.DefaultClassifyOptions()
//	opts.CustomRules["billing"] = []string{"STRIPE_", "PAYMENT_"}
//	results := env.Classify(entries, opts)
//	summaries := env.BuildClassifySummaries(results)
//	fmt.Print(env.FormatClassify(summaries))
package env
