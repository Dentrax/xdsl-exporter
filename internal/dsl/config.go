package dsl

import (
	"fmt"

	"3e8.eu/go/dsl"

	"github.com/Dentrax/xdsl-exporter/internal/config"
)

func GenerateConfigFrom(cfg config.Config) (*dsl.Config, error) {
	if err := cfg.Check(); err != nil {
		return nil, fmt.Errorf("config check: %w", err)
	}

	client := dsl.ClientType(cfg.TargetClient)
	if !client.IsValid() {
		return nil, fmt.Errorf("invalid client type: %s: alloweds: %s", client, GetSupportedClients())
	}

	sshKey, err := cfg.ReadSSHKey()
	if err != nil {
		return nil, err
	}

	knownHosts, err := cfg.ReadKnownHosts()
	if err != nil {
		return nil, err
	}

	return &dsl.Config{
		Type:            client,
		Host:            cfg.TargetHost,
		User:            cfg.TargetUser,
		AuthPassword:    getAuthPassword(cfg.TargetPassword),
		AuthPrivateKeys: getAuthPrivateKeys(sshKey, cfg.TargetSSHPassphrase),
		KnownHosts:      knownHosts,
		Options:         nil,
	}, nil
}

func getAuthPassword(password string) dsl.PasswordCallback {
	if password == "" {
		return nil
	}
	return dsl.Password(password)
}

func getAuthPrivateKeys(sshKey, sshKeyPassphrase string) dsl.PrivateKeysCallback {
	getKeys := func(sshKey string) func() ([]string, error) {
		if sshKey == "" {
			return nil
		}
		return func() ([]string, error) {
			return []string{sshKey}, nil
		}
	}

	getPassphrase := func(sshKeyPassphrase string) func(fingerprint string) (string, error) {
		if sshKeyPassphrase == "" {
			return nil
		}
		return func(fingerprint string) (string, error) {
			return sshKeyPassphrase, nil
		}
	}

	return dsl.PrivateKeysCallback{
		Keys:       getKeys(sshKey),
		Passphrase: getPassphrase(sshKeyPassphrase),
	}
}
