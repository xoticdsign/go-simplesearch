package utils

import (
	"bytes"
	"encoding/json"
	"io"
)

// JSONEncode() serializes an input value into JSON format and writes it to a buffer.
//
// This function takes an input value 'v' of any type, serializes it into JSON format,
// and writes the result to a bytes buffer. If successful, it returns the buffer and a nil error.
// If an error occurs during encoding, the buffer and the error are returned.
func JSONEncode(v any) (bytes.Buffer, error) {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

// JSONDecode() reads JSON data from an io.Reader and deserializes it into an arbitrary value.
//
// This function reads JSON data from the provided reader 'r' and decodes it into a variable
// of type 'any' (interface{}). If successful, it returns the decoded value and a nil error.
// If an error occurs during decoding, it returns the value (which may be partially decoded) and the error.
func JSONDecode(r io.Reader) (any, error) {
	var out any

	err := json.NewDecoder(r).Decode(&out)
	if err != nil {
		return out, err
	}
	return out, nil
}
