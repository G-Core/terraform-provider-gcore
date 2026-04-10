package custom

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
)

// apiTag represents a tag object from the Gcore API response.
// The API returns tags as [{"key":"k","value":"v","read_only":false}].
type apiTag struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	ReadOnly bool   `json:"read_only"`
}

// parseTagsFromJSON extracts the tags array from raw API JSON bytes.
// Tries "tags_v2" first (proper object format used by load_balancer and security_group),
// then falls back to "tags".
// Handles two API response formats:
//   - Object format: [{"key":"k","value":"v","read_only":false}] (most resources)
//   - String format: ["key=value"] (security_group "tags" field)
//
// Filters out read_only system tags. Returns nil if JSON parsing fails entirely.
func parseTagsFromJSON(jsonBytes []byte) *map[string]types.String {
	var raw struct {
		Tags   json.RawMessage `json:"tags"`
		TagsV2 json.RawMessage `json:"tags_v2"`
	}
	if err := json.Unmarshal(jsonBytes, &raw); err != nil {
		return nil
	}

	// Try tags_v2 first (always object format when present), then fall back to tags.
	// This ordering matters because some resources (security_group) have both:
	//   "tags_v2": [{"key":"k","value":"v","read_only":false}]  (structured)
	//   "tags": ["k=v"]                                          (legacy string format)
	sources := []json.RawMessage{raw.TagsV2, raw.Tags}
	for _, tagsJSON := range sources {
		if len(tagsJSON) == 0 || string(tagsJSON) == "null" {
			continue
		}

		// Try object format: [{"key":"k","value":"v","read_only":false}]
		var objectTags []apiTag
		if err := json.Unmarshal(tagsJSON, &objectTags); err == nil {
			result := make(map[string]types.String, len(objectTags))
			for _, tag := range objectTags {
				if tag.ReadOnly {
					continue
				}
				result[tag.Key] = types.StringValue(tag.Value)
			}
			return &result
		}

		// Try string format: ["key=value"]
		var stringTags []string
		if err := json.Unmarshal(tagsJSON, &stringTags); err == nil {
			result := make(map[string]types.String, len(stringTags))
			for _, tag := range stringTags {
				if idx := strings.Index(tag, "="); idx > 0 {
					result[tag[:idx]] = types.StringValue(tag[idx+1:])
				}
			}
			return &result
		}
	}

	result := make(map[string]types.String)
	return &result
}

// ConvertAPITagsToCustomfieldMap parses raw API JSON and extracts the tags array,
// returning a customfield.Map[types.String]. Filters out read_only system tags.
// Returns false only if JSON parsing fails entirely.
func ConvertAPITagsToCustomfieldMap(ctx context.Context, jsonBytes []byte) (customfield.Map[types.String], bool) {
	tags := parseTagsFromJSON(jsonBytes)
	if tags == nil {
		return customfield.Map[types.String]{}, false
	}
	return customfield.NewMapMust(ctx, *tags), true
}

// ConvertAPITagsToMap parses raw API JSON and extracts the tags array,
// returning a *map[string]types.String. Filters out read_only system tags.
// Returns nil only if JSON parsing fails entirely.
func ConvertAPITagsToMap(jsonBytes []byte) *map[string]types.String {
	return parseTagsFromJSON(jsonBytes)
}
