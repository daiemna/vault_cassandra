package vault

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
)

//Client for vault cassandra test.
type Client struct {
	inner *api.Client
}

// Payload is a type that holds data to write to vault
type Payload map[string]interface{}

//Write uses logical-api of vault to write to path, the given data.
func (c *Client) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	return c.inner.Logical().Write(path, data)
}

//PutPolicy writes poicy to vault
func (c *Client) PutPolicy(name, rules string) error {
	return c.inner.Sys().PutPolicy(name, rules)
}

func (c *Client) Read(path string) (*api.Secret, error) {
	return c.inner.Logical().Read(path)
}

// NewClientFromToken creates a vault client from token
func NewClientFromToken(token string, config *api.Config) *Client {
	if config == nil {
		config = api.DefaultConfig()
	}

	client, err := api.NewClient(config)
	// Final error return should be in the `return` section at the end
	if err != nil {
		msg := "Failed to create a vault client with error: %s"
		log.Debug().Msgf(msg, err.Error())
		return nil
	}
	client.SetToken(token)

	return &Client{
		inner: client,
	}
}

//NewRootClient creates client with root token.
func NewRootClient() (*Client, error) {
	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		return nil, fmt.Errorf("cannot find VAULT_TOKEN in environment")
	}
	client := NewClientFromToken(vaultToken, nil)
	return client, nil
}
