package importadoption

import (
	"encoding/json"
	"strings"
)

// KeepKeys reduces a JSON object body to an allowlist of top-level keys. Used
// where the PATCH endpoint accepts a fixed, small set (volume: name, tags) so
// an adopted create-only attribute never rides a coincident PATCH: apijson's
// MarshalForPatch emits any attribute whose prior state is null and plan is
// set, which is exactly the shape of an adoption.
func KeepKeys(body []byte, keys ...string) ([]byte, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, err
	}
	if obj == nil {
		return body, nil
	}
	kept := map[string]json.RawMessage{}
	for _, k := range keys {
		if v, ok := obj[k]; ok {
			kept[k] = v
		}
	}
	return json.Marshal(kept)
}

// StripKeys drops top-level keys from a JSON object body. Used where the PATCH
// endpoint accepts many keys but never this one (cdn: origin).
func StripKeys(body []byte, keys ...string) ([]byte, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, err
	}
	if obj == nil {
		return body, nil
	}
	for _, k := range keys {
		delete(obj, k)
	}
	return json.Marshal(obj)
}

// IsEmptyBody reports a body that carries nothing to send.
func IsEmptyBody(body []byte) bool {
	s := strings.TrimSpace(string(body))
	return s == "" || s == "{}" || s == "null"
}
