package domain

import (
	"slices"
)

type HeaderType string

type Headers map[string]Header

const (
	StaticHeaderType  = HeaderType("static")
	DynamicHeaderType = HeaderType("dynamic")
)

var AllowedHeaderTypes = []HeaderType{StaticHeaderType, DynamicHeaderType}

type Header struct {
	Type  HeaderType `json:"type"`
	Value *string    `json:"value,omitempty"`
}

func (h *Headers) Validate() bool {
	if h == nil {
		return true
	}

	if len(*h) > 0 {
		for key := range *h {
			header := (*h)[key]

			if !(header).Validate() {
				return false
			}
		}
	}

	return true
}

func (h *Header) Validate() bool {
	if h == nil {
		return false
	}

	if !slices.Contains(AllowedHeaderTypes, h.Type) {
		return false
	}

	if h.Type == StaticHeaderType && (h.Value == nil || *h.Value == "") {
		return false
	}

	return true
}
