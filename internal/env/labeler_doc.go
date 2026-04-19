// Package env provides utilities for manipulating .env file entries.
//
// # Labeler
//
// The labeler module allows attaching short named labels to individual
// environment variable entries via their comment field.
//
// Labels are stored inline using the "label:<value>" convention and can
// coexist with other comment text.
//
// Example usage:
//
//	entry before: APP_ENV=production  # some note
//	after Label:  APP_ENV=production  # some note label:infra
//
// Functions:
//   - Label(entries, key, label, opts) – attach a label to a key
//   - Unlabel(entries, key)            – remove the label from a key
//   - GetLabel(entries, key)           – retrieve the label value for a key
package env
