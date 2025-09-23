// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadEvidenceError_JSON(t *testing.T) {
	tests := []struct {
		name     string
		err      BadEvidenceError
		expected string
	}{
		{
			name: "string error",
			err:  BadEvidenceError{"test error"},
			expected: `{
				"error": "bad evidence",
				"detail-type": "string",
				"detail": "test error"
			}`,
		},
		{
			name: "wrapped error",
			err:  BadEvidenceError{fmt.Errorf("wrapped: %w", errors.New("inner"))},
			expected: `{
				"error": "bad evidence",
				"detail-type": "error",
				"detail": ["wrapped: inner", "inner"]
			}`,
		},
		{
			name: "non-error type",
			err:  BadEvidenceError{42},
			expected: `{
				"error": "bad evidence",
				"detail-type": "other",
				"detail": 42
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.err)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))
		})
	}
}

func TestBadEvidenceError_Unwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	wrapped := fmt.Errorf("wrapped: %w", innerErr)

	tests := []struct {
		name     string
		err      BadEvidenceError
		wantNil  bool
		wantText string
	}{
		{
			name:     "unwrappable error",
			err:      BadEvidenceError{wrapped},
			wantNil:  false,
			wantText: "wrapped: inner error",
		},
		{
			name:     "non-wrapped error",
			err:      BadEvidenceError{errors.New("simple")},
			wantNil:  false,
			wantText: "simple",
		},
		{
			name:    "non-error value",
			err:     BadEvidenceError{"string"},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unwrapped := tt.err.Unwrap()
			if tt.wantNil {
				assert.Nil(t, unwrapped)
			} else {
				assert.NotNil(t, unwrapped)
				assert.Equal(t, tt.wantText, unwrapped.Error())
			}
		})
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name:     "matching errors",
			err:      BadEvidenceError{"test"},
			target:   BadEvidenceError{"other"},
			expected: true,
		},
		{
			name:     "different error types",
			err:      BadEvidenceError{"test"},
			target:   errors.New("test"),
			expected: false,
		},
		{
			name:     "wrapped error matching",
			err:      fmt.Errorf("wrapped: %w", BadEvidenceError{"test"}),
			target:   BadEvidenceError{"other"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, errors.Is(tt.err, tt.target))
		})
	}
}