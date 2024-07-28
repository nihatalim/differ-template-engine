package service

import (
	"differ-template-engine/application/domain"
	"reflect"
	"testing"
)

func TestNormalizeUrl(t *testing.T) {
	tests := []struct {
		name                string
		url                 string
		executionParameters map[string]string
		expected            string
	}{
		{
			name: "Single parameter replacement",
			url:  "http://example.com/${param1}",
			executionParameters: map[string]string{
				"param1": "value1",
			},
			expected: "http://example.com/value1",
		},
		{
			name: "Multiple parameter replacement",
			url:  "http://example.com/${param1}/test/${param2}",
			executionParameters: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
			expected: "http://example.com/value1/test/value2",
		},
		{
			name:                "No parameter replacement",
			url:                 "http://example.com/test",
			executionParameters: map[string]string{},
			expected:            "http://example.com/test",
		},
		{
			name: "Parameter not in executionParameters",
			url:  "http://example.com/${param1}",
			executionParameters: map[string]string{
				"param2": "value2",
			},
			expected: "http://example.com/${param1}",
		},
		{
			name: "Parameter with special characters",
			url:  "http://example.com/${param1}",
			executionParameters: map[string]string{
				"param1": "value@1!",
			},
			expected: "http://example.com/value@1!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeUrl(tt.url, tt.executionParameters)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNormalizeHeader(t *testing.T) {
	staticValue := "StaticValue"
	tests := []struct {
		name                string
		templateHeaders     domain.Headers
		executionParameters map[string]string
		expected            map[string]string
	}{
		{
			name: "Static and Dynamic Headers",
			templateHeaders: domain.Headers{
				"Header1": {Type: domain.StaticHeaderType, Value: &staticValue},
				"Header2": {Type: domain.DynamicHeaderType},
			},
			executionParameters: map[string]string{
				"Header2": "DynamicValue",
			},
			expected: map[string]string{
				"Header1": "StaticValue",
				"Header2": "DynamicValue",
			},
		},
		{
			name: "Only Static Header",
			templateHeaders: domain.Headers{
				"Header1": {Type: domain.StaticHeaderType, Value: &staticValue},
			},
			executionParameters: map[string]string{},
			expected: map[string]string{
				"Header1": "StaticValue",
			},
		},
		{
			name: "Only Dynamic Header",
			templateHeaders: domain.Headers{
				"Header1": {Type: domain.DynamicHeaderType},
			},
			executionParameters: map[string]string{
				"Header1": "DynamicValue",
			},
			expected: map[string]string{
				"Header1": "DynamicValue",
			},
		},
		{
			name: "Dynamic Header Without Execution Parameter",
			templateHeaders: domain.Headers{
				"Header1": {Type: domain.DynamicHeaderType},
			},
			executionParameters: map[string]string{},
			expected:            map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeHeader(tt.templateHeaders, tt.executionParameters)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}
