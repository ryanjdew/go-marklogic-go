package util

import "net/url"

// Transform represents a transform for a request
type Transform struct {
	Name       string
	Parameters map[string]string
}

// ToParameters converts a transform to its url form
func (t *Transform) ToParameters() string {
	params := "&transform=" + url.QueryEscape(t.Name)
	for key, val := range t.Parameters {
		params = params + "&trans:" + key + "=" + url.QueryEscape(val)
	}
	return params
}
