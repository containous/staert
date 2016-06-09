package staert

import (
	"fmt"
	"strings"
)

func generateMapstructure(input map[string]string, prefix string) (map[string]interface{}, error) {
	raw := make(map[string]interface{})
	for k, v := range input {
		// Trim the prefix off our key first
		key := strings.TrimPrefix(k, prefix)

		// Determine what map we're writing the value to. We split by '/'
		// to determine any sub-maps that need to be created.
		m := raw
		children := strings.Split(key, "/")
		if len(children) > 0 {
			key = children[len(children)-1]
			children = children[:len(children)-1]
			for _, child := range children {
				if m[child] == nil {
					m[child] = make(map[string]interface{})
				}
				subm, ok := m[child].(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("child is both a data item and dir: %s", child)
				}
				m = subm
			}

		}
		m[key] = string(v)
	}
	return raw, nil
}
