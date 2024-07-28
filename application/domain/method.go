package domain

import "slices"

type Method string

const MethodGet = Method("GET")

var AllowedMethods = []Method{MethodGet}

func (m *Method) Validate() bool {
	if m == nil {
		return false
	}

	return slices.Contains(AllowedMethods, *m)
}
