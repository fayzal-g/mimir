// SPDX-License-Identifier: AGPL-3.0-only

package vault

import (
	"context"
	"errors"
	"flag"
	"fmt"

	hashivault "github.com/hashicorp/vault/api"
)

// Config for the Vault used to fetch secrets
type Config struct {
	Enabled bool `yaml:"enabled" category:"experimental"`

	URL       string `yaml:"url" category:"experimental"`
	Token     string `yaml:"token" category:"experimental"`
	MountPath string `yaml:"mount_path" category:"experimental"`
}

func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	f.BoolVar(&cfg.Enabled, "vault.enabled", false, "Enables fetching of keys and certificates from Vault")
	f.StringVar(&cfg.URL, "vault.url", "", "Location of the Vault server")
	f.StringVar(&cfg.Token, "vault.token", "", "Token used to authenticate with Vault")
	f.StringVar(&cfg.MountPath, "vault.mount-path", "", "Location of secrets engine within Vault")
}

type SecretsEngine interface {
	Get(ctx context.Context, path string) (*hashivault.KVSecret, error)
}

type Vault struct {
	KVStore SecretsEngine
}

func NewVault(cfg Config) (*Vault, error) {
	if cfg.URL == "" || cfg.Token == "" || cfg.MountPath == "" {
		return nil, errors.New("invalid vault configuration supplied")
	}

	config := hashivault.DefaultConfig()
	config.Address = cfg.URL

	client, err := hashivault.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(cfg.Token)
	vault := &Vault{
		KVStore: client.KVv2(cfg.MountPath),
	}

	return vault, nil
}

func (v *Vault) ReadSecret(path string) ([]byte, error) {
	secret, err := v.KVStore.Get(context.Background(), path)
	if err != nil {
		return nil, fmt.Errorf("unable to read secret from vault: %v", err)
	}

	data := []byte(secret.Data["value"].(string))
	return data, nil
}
