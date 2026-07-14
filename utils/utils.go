// Package utils holds small, reusable helpers that don't belong to any
// single domain package.
package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ToCSV returns a [][]string from data — any value that can be JSON-marshaled
// into a list of objects, such as any go struct — ready for io.Writer interface
// when needed.
//
// keys must match the JSON key names produced by data's `json:"..."` tags,
// not the Go struct field names, since ToCSV operates on data's marshaled
// JSON representation rather than the struct itself.
func ToCSV(data any, keys []string, header bool) ([][]string, error) {
	// result variable
	var csvData [][]string

	rows, err := toRowMaps(data)
	if err != nil {
		return nil, fmt.Errorf("utils: converting data to rows: %w", err)
	}

	if len(rows) > 0 {
		unknown := unknownKeys(rows[0], keys)
		if len(unknown) > 0 {
			return nil, fmt.Errorf("utils: unknown CSV field(s): %s", strings.Join(unknown, ", "))
		}
	}

	// optional: add header to CSV file
	if header {
		csvData = append(csvData, keys)
	}

	for _, row := range rows {
		record := make([]string, len(keys))
		for j, key := range keys {
			record[j] = stringify(row[key])
		}
		csvData = append(csvData, record)
	}

	return csvData, nil
}

// toRowMaps normalizes arbitrary JSON-marshalable data into a slice of
// generic key/value rows, so ToCSV doesn't need to know the concrete type
// (e.g. []snmpmodules.ReachableDevice) it was called with.
func toRowMaps(data any) ([]map[string]any, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshaling to JSON: %w", err)
	}

	var rows []map[string]any
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, fmt.Errorf("unmarshaling into rows: %w", err)
	}

	return rows, nil
}

// unknownKeys returns any requested keys that don't exist in row. Only the
// first row is checked (see ToCSV), which assumes all rows share the same
// shape — true for a slice of a single struct type, which is the expected
// input here.
func unknownKeys(row map[string]any, keys []string) []string {
	var unknown []string
	for _, k := range keys {
		if _, ok := row[k]; !ok {
			unknown = append(unknown, k)
		}
	}
	return unknown
}

// stringify converts a single decoded JSON value into its CSV cell
// representation. Scalars are formatted directly; missing values, maps,
// and slices fall back to compact JSON so nested data (e.g. a credential
// map) isn't silently dropped, just condensed into one cell.
func stringify(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case float64:
		// encoding/json decodes every JSON number as float64; format
		// whole numbers without a trailing ".0".
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		// map[string]any, []any, or anything else JSON produced.
		raw, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(raw)
	}
}
