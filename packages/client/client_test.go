package client

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		apiKey   string
		wantErr  bool
	}{
		{
			name:     "valid URL and API key",
			baseURL:  "https://example.com",
			apiKey:   "test-key",
			wantErr:  false,
		},
		{
			name:     "URL without scheme",
			baseURL:  "example.com",
			apiKey:   "test-key",
			wantErr:  false,
		},
		{
			name:     "invalid URL",
			baseURL:  "://invalid",
			apiKey:   "test-key",
			wantErr:  true,
		},
		{
			name:     "empty base URL",
			baseURL:  "",
			apiKey:   "test-key",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.baseURL, tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client when no error expected")
			}
		})
	}
}

func TestNewClientWithOptions(t *testing.T) {
	client, err := NewClient("https://example.com", "test-key",
		WithTimeout(60*time.Second),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client.httpClient.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", client.httpClient.Timeout)
	}
}

func TestClientFields(t *testing.T) {
	baseURL := "https://example.com"
	apiKey := "test-key"

	client, err := NewClient(baseURL, apiKey)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, client.baseURL)
	}

	if client.apiKey != apiKey {
		t.Errorf("Expected apiKey %s, got %s", apiKey, client.apiKey)
	}

	if client.httpClient == nil {
		t.Error("Expected httpClient to be initialized")
	}
}