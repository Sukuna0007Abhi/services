// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package handler

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/veraison/services/plugin"
)

// MockEndorsementHandler implements IEndorsementHandler for testing
type MockEndorsementHandler struct {
	plugin.Plugin[plugin.IPluggable]
	InitCalled    bool
	InitError     error
	CloseCalled   bool
	CloseError    error
	DecodeCalled  bool
	DecodeError   error
	MockResponse  *EndorsementHandlerResponse
}

func (m *MockEndorsementHandler) Init(params EndorsementHandlerParams) error {
	m.InitCalled = true
	return m.InitError
}

func (m *MockEndorsementHandler) Close() error {
	m.CloseCalled = true
	return m.CloseError
}

func (m *MockEndorsementHandler) Decode(data []byte, mediaType string, caCertPool []byte) (*EndorsementHandlerResponse, error) {
	m.DecodeCalled = true
	if m.DecodeError != nil {
		return nil, m.DecodeError
	}
	return m.MockResponse, nil
}

func TestEndorsementHandlerResponse(t *testing.T) {
	tests := []struct {
		name        string
		response    EndorsementHandlerResponse
		wantJSON    string
		endorsement Endorsement
	}{
		{
			name: "empty response",
			response: EndorsementHandlerResponse{
				ReferenceValues: []Endorsement{},
				TrustAnchors:   []Endorsement{},
				SignerInfo:     map[string]string{},
			},
			wantJSON: `{"ReferenceValues":[],"TrustAnchors":[],"SignerInfo":{}}`,
		},
		{
			name: "with reference values",
			endorsement: Endorsement{
				Scheme:   "test-scheme",
				Type:     EndorsementType_REFERENCE_VALUE,
				SubType:  "test-subtype",
				Attributes: json.RawMessage(`{"key":"value"}`),
			},
			response: EndorsementHandlerResponse{
				ReferenceValues: []Endorsement{
					{
						Scheme:   "test-scheme",
						Type:     EndorsementType_REFERENCE_VALUE,
						SubType:  "test-subtype",
						Attributes: json.RawMessage(`{"key":"value"}`),
					},
				},
				TrustAnchors: []Endorsement{},
				SignerInfo:   map[string]string{"signer": "info"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantJSON != "" {
				data, err := json.Marshal(tt.response)
				require.NoError(t, err)
				assert.JSONEq(t, tt.wantJSON, string(data))
			}

			if tt.endorsement.Scheme != "" {
				assert.Equal(t, tt.endorsement.Scheme, tt.response.ReferenceValues[0].Scheme)
				assert.Equal(t, tt.endorsement.Type, tt.response.ReferenceValues[0].Type)
				assert.Equal(t, tt.endorsement.SubType, tt.response.ReferenceValues[0].SubType)
				assert.JSONEq(t, string(tt.endorsement.Attributes), string(tt.response.ReferenceValues[0].Attributes))
			}
		})
	}
}

func TestEndorsementTypes(t *testing.T) {
	assert.Equal(t, "unspecified", EndorsementType_UNSPECIFIED)
	assert.Equal(t, "reference value", EndorsementType_REFERENCE_VALUE)
	assert.Equal(t, "trust anchor", EndorsementType_VERIFICATION_KEY)
}

func TestMockEndorsementHandler(t *testing.T) {
	handler := &MockEndorsementHandler{
		MockResponse: &EndorsementHandlerResponse{
			ReferenceValues: []Endorsement{
				{
					Scheme:   "test-scheme",
					Type:     EndorsementType_REFERENCE_VALUE,
					SubType:  "test-subtype",
					Attributes: json.RawMessage(`{"key":"value"}`),
				},
			},
		},
	}

	// Test Init
	err := handler.Init(EndorsementHandlerParams{"param": "value"})
	assert.NoError(t, err)
	assert.True(t, handler.InitCalled)

	// Test Decode
	resp, err := handler.Decode([]byte("test"), "application/json", nil)
	assert.NoError(t, err)
	assert.True(t, handler.DecodeCalled)
	assert.NotNil(t, resp)
	assert.Len(t, resp.ReferenceValues, 1)

	// Test Close
	err = handler.Close()
	assert.NoError(t, err)
	assert.True(t, handler.CloseCalled)
}