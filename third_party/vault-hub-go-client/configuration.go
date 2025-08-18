package openapi

import (
	"context"
	"net/http"
)

type contextKey string

var (
	ContextAPIKeys                  = contextKey("apiKeys")
	ContextServerIndex              = contextKey("serverIndex")
	ContextOperationServerIndices   = contextKey("serverOperationIndices")
	ContextServerVariables          = contextKey("serverVariables")
	ContextOperationServerVariables = contextKey("serverOperationServerVariables")
)

// APIKey provides API key based authentication to a request passed via context using ContextAPIKeys
type APIKey struct {
	Key    string
	Prefix string
}

// ServerConfiguration stores the information about a server
type ServerConfiguration struct {
	URL         string
	Description string
}

// ServerConfigurations stores multiple ServerConfiguration items
type ServerConfigurations []ServerConfiguration

// Configuration stores the configuration of the API client
type Configuration struct {
	Servers    ServerConfigurations
	HTTPClient *http.Client
}

// NewConfiguration returns a new Configuration object
func NewConfiguration() *Configuration {
	return &Configuration{
		Servers: ServerConfigurations{{URL: ""}},
	}
}

// ServerURLWithContext returns a server URL based on configuration
func (c *Configuration) ServerURLWithContext(ctx context.Context, endpoint string) (string, error) {
	if len(c.Servers) == 0 {
		return "", nil
	}
	return c.Servers[0].URL, nil
}

// AddDefaultHeader is a no-op in this minimal shim
func (c *Configuration) AddDefaultHeader(key, value string) {}

package openapi

import (
	"context"
	"net/http"
)

type contextKey string

var (
	ContextAPIKeys                 = contextKey("apiKeys")
	ContextServerIndex             = contextKey("serverIndex")
	ContextOperationServerIndices  = contextKey("serverOperationIndices")
	ContextServerVariables         = contextKey("serverVariables")
	ContextOperationServerVariables = contextKey("serverOperationServerVariables")
)

type APIKey struct {
	Key    string
	Prefix string
}

type ServerConfiguration struct {
	URL         string
	Description string
}

type ServerConfigurations []ServerConfiguration

type Configuration struct {
	Servers    ServerConfigurations
	HTTPClient *http.Client
}

func NewConfiguration() *Configuration {
	return &Configuration{Servers: ServerConfigurations{{URL: ""}}}
}

func (c *Configuration) ServerURLWithContext(ctx context.Context, endpoint string) (string, error) {
	// We only use the first server URL in this minimal shim
	if len(c.Servers) == 0 {
		return "", nil
	}
	return c.Servers[0].URL, nil
}

func (c *Configuration) AddDefaultHeader(key, value string) {}
