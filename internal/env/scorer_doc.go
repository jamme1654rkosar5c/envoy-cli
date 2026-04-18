// Package env provides utilities for managing environment variable entries.
//
// # Scorer
//
// The scorer module evaluates the quality of env entries and assigns
// a numeric score from 0 to 100. Penalties are applied based on
// configurable options:
//
//   - PenalizeEmpty: deducts 30 points for entries with empty values.
//   - PenalizeUnquoted: deducts 20 points for values containing spaces
//     that are not wrapped in quotes.
//   - PenalizeNoComment: deducts 10 points for entries without a comment.
//
// Example usage:
//
//	opts := env.DefaultScoreOptions()
//	scores := env.Score(entries, opts)
//	fmt.Print(env.FormatScores(scores))
//	fmt.Printf("Average: %.1f\n", env.AverageScore(scores))
package env
