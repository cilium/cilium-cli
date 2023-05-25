package utils

import "fmt"

// Merge maps recursively merges the values of b into a copy of a, preferring the values from b
func MergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = MergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

// GetPath gets a value by a path in a nested map[string]interface{}.
// An error is returned if:
//   - At any point along the path, the map is not a map[string]interface{}
//   - At any point along the path, the path segment does not appear in the map
func GetPath(m map[string]interface{}, path []string) (interface{}, error) {
	for i, segment := range path {
		if i == len(path)-1 {
			return m[segment], nil
		}
		if _, ok := m[segment]; !ok {
			return nil, fmt.Errorf("map %+v does not contain path %v", m, path)
		}
		if m, ok := m[segment].(map[string]interface{}); !ok {
			return nil, fmt.Errorf("cannot traverse string path %v in map %+v", path, m)
		}
	}
	return nil, fmt.Errorf("map %+v does not contain path %v", m, path)
}
