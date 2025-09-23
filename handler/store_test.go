// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Common errors for store operations
var (
	ErrEmptyKey    = errors.New("empty key")
	ErrKeyNotFound = errors.New("key not found")
)

// MockStoreHandler implements IStoreHandler for testing
type MockStoreHandler struct {
	StoreCalled bool
	StoreError  error
	GetCalled   bool
	GetError    error
	MockData    []byte
}

func (m *MockStoreHandler) Store(ctx context.Context, key string, value []byte) error {
	m.StoreCalled = true
	return m.StoreError
}

func (m *MockStoreHandler) Get(ctx context.Context, key string) ([]byte, error) {
	m.GetCalled = true
	return m.MockData, m.GetError
}

func (m *MockStoreHandler) GetAttestationScheme() string {
	return "mock-scheme"
}

func TestMockStoreHandler(t *testing.T) {
	handler := &MockStoreHandler{
		MockData: []byte("test-data"),
	}

	// Test Store
	err := handler.Store(context.Background(), "test-key", []byte("test-value"))
	assert.NoError(t, err)
	assert.True(t, handler.StoreCalled)

	// Test Get
	data, err := handler.Get(context.Background(), "test-key")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test-data"), data)
	assert.True(t, handler.GetCalled)

	// Test GetAttestationScheme
	scheme := handler.GetAttestationScheme()
	assert.Equal(t, "mock-scheme", scheme)
}

func TestMockStoreHandler_Errors(t *testing.T) {
	tests := []struct {
		name      string
		handler   *MockStoreHandler
		operation string
		key       string
		value     []byte
		wantErr   error
	}{
		{
			name: "store - key not found",
			handler: &MockStoreHandler{
				StoreError: ErrKeyNotFound,
			},
			operation: "store",
			key:       "test-key",
			value:     []byte("test-value"),
			wantErr:   ErrKeyNotFound,
		},
		{
			name: "get - key not found",
			handler: &MockStoreHandler{
				GetError: ErrKeyNotFound,
			},
			operation: "get",
			key:       "test-key",
			wantErr:   ErrKeyNotFound,
		},
		{
			name: "store - empty key",
			handler: &MockStoreHandler{
				StoreError: ErrEmptyKey,
			},
			operation: "store",
			key:       "",
			value:     []byte("test-value"),
			wantErr:   ErrEmptyKey,
		},
		{
			name: "get - empty key",
			handler: &MockStoreHandler{
				GetError: ErrEmptyKey,
			},
			operation: "get",
			key:       "",
			wantErr:   ErrEmptyKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.operation == "store" {
				err = tt.handler.Store(context.Background(), tt.key, tt.value)
			} else {
				_, err = tt.handler.Get(context.Background(), tt.key)
			}
			assert.Equal(t, tt.wantErr, err)
		})
	}
}