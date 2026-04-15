package audit

import "time"

// FilterOptions controls which entries are returned by Filter.
type FilterOptions struct {
	Event   EventType // if non-empty, only entries with this event type
	File    string    // if non-empty, only entries for this file
	Since   time.Time // if non-zero, only entries at or after this time
	MaxRows int       // if > 0, limit results to this many entries
}

// Filter returns a subset of log entries matching the given options.
func Filter(log *Log, opts FilterOptions) []Entry {
	var result []Entry

	for _, e := range log.Entries {
		if opts.Event != "" && e.Event != opts.Event {
			continue
		}
		if opts.File != "" && e.File != opts.File {
			continue
		}
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		result = append(result, e)
	}

	if opts.MaxRows > 0 && len(result) > opts.MaxRows {
		result = result[len(result)-opts.MaxRows:]
	}

	return result
}
