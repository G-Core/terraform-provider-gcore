package custom

import "github.com/hashicorp/terraform-plugin-framework/types"

// TagsEqual compares two maps of tags for equality.
func TagsEqual(a, b *map[string]types.String) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(*a) != len(*b) {
		return false
	}
	for k, v := range *a {
		if bv, ok := (*b)[k]; !ok || !v.Equal(bv) {
			return false
		}
	}
	return true
}
