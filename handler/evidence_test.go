// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package handler

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/services/plugin"
	"github.com/veraison/services/proto"
)

// Common errors for evidence handling
var (
	ErrUnsupportedFormat = errors.New("unsupported evidence format")
	ErrValidationFailed  = errors.New("evidence validation failed")
)

// MockEvidenceHandler implements IEvidenceHandler for testing
type MockEvidenceHandler struct {
	plugin.Plugin[plugin.IPluggable]
	ExtractClaimsCalled bool
	ExtractClaimsError  error
	MockClaims         map[string]interface{}
	ValidateEvidenceIntegrityCalled bool
	ValidateEvidenceIntegrityError  error
}

func (m *MockEvidenceHandler) ExtractClaims(token *proto.AttestationToken, trustAnchors []string) (map[string]interface{}, error) {
	m.ExtractClaimsCalled = true
	if m.ExtractClaimsError != nil {
		return nil, m.ExtractClaimsError
	}
	return m.MockClaims, nil
}

func (m *MockEvidenceHandler) ValidateEvidenceIntegrity(token *proto.AttestationToken, trustAnchors []string) error {
	m.ValidateEvidenceIntegrityCalled = true
	return m.ValidateEvidenceIntegrityError
}

func TestMockEvidenceHandler_ExtractClaims(t *testing.T) {
	mockClaims := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	tests := []struct {
		name         string
		handler      *MockEvidenceHandler
		token        *proto.AttestationToken
		trustAnchors []string
		wantClaims   map[string]interface{}
		wantErr      bool
	}{
		{
			name: "successful extraction",
			handler: &MockEvidenceHandler{
				MockClaims: mockClaims,
			},
			token: &proto.AttestationToken{
				TenantId: "test-tenant",
				Data:     []byte("test-data"),
			},
			trustAnchors: []string{"anchor1", "anchor2"},
			wantClaims:   mockClaims,
			wantErr:      false,
		},
		{
			name: "extraction error",
			handler: &MockEvidenceHandler{
				ExtractClaimsError: ErrUnsupportedFormat,
			},
			token: &proto.AttestationToken{
				TenantId: "test-tenant",
				Data:     []byte("test-data"),
			},
			trustAnchors: []string{"anchor1"},
			wantClaims:   nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := tt.handler.ExtractClaims(tt.token, tt.trustAnchors)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantClaims, claims)
			}
			assert.True(t, tt.handler.ExtractClaimsCalled)
		})
	}
}

func TestMockEvidenceHandler_ValidateEvidenceIntegrity(t *testing.T) {
	tests := []struct {
		name         string
		handler      *MockEvidenceHandler
		token        *proto.AttestationToken
		trustAnchors []string
		wantErr      bool
	}{
		{
			name:    "successful validation",
			handler: &MockEvidenceHandler{},
			token: &proto.AttestationToken{
				TenantId: "test-tenant",
				Data:     []byte("test-data"),
			},
			trustAnchors: []string{"anchor1", "anchor2"},
			wantErr:      false,
		},
		{
			name: "validation error",
			handler: &MockEvidenceHandler{
				ValidateEvidenceIntegrityError: ErrValidationFailed,
			},
			token: &proto.AttestationToken{
				TenantId: "test-tenant",
				Data:     []byte("invalid-data"),
			},
			trustAnchors: []string{"anchor1"},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.handler.ValidateEvidenceIntegrity(tt.token, tt.trustAnchors)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.True(t, tt.handler.ValidateEvidenceIntegrityCalled)
		})
	}
}