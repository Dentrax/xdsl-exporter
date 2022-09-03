package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/mitchellh/go-homedir"
)

type Config struct {
	ListenAddress       string
	MetricsPath         string
	KnownHostsPath      string
	TargetHost          string
	TargetPort          int
	TargetUser          string
	TargetPassword      string
	TargetSSHKeyPath    string
	TargetSSHPassphrase string
	TargetClient        string
}

func (c Config) Check() error {
	if c.ListenAddress == "" {
		return fmt.Errorf("listen address is empty")
	}

	if c.MetricsPath == "" {
		return fmt.Errorf("metrics path is empty")
	}

	if c.TargetHost == "" {
		return fmt.Errorf("target host is empty")
	}

	_, err := url.Parse(c.TargetHost)
	if err != nil {
		return fmt.Errorf("invalid target host: %w", err)
	}

	if c.TargetUser == "" {
		return fmt.Errorf("target user is empty")
	}

	if c.TargetPassword == "" && c.TargetSSHKeyPath == "" {
		return fmt.Errorf("no password or ssh key path provided")
	}

	return nil
}

func (c Config) ReadSSHKey() (string, error) {
	value, err := os.ReadFile(c.TargetSSHKeyPath)
	if err != nil {
		return "", fmt.Errorf("read ssh key: %w", err)
	}
	return string(value), nil
}

func (c Config) ReadKnownHosts() (string, error) {
	expanded, err := homedir.Expand(c.KnownHostsPath)
	if err != nil {
		return "", fmt.Errorf("get known_hosts: %w", err)
	}

	value, err := os.ReadFile(expanded)
	if err != nil {
		return "", fmt.Errorf("read known_hosts: %w", err)
	}
	return string(value), nil
}
