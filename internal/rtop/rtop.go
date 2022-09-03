package rtop

import (
	"fmt"
	"github.com/Dentrax/xdsl-exporter/internal/config"
	"github.com/rapidloop/rtop/pkg/client"
)

func New(cfg config.Config) (*client.Client, error) {
	opts := []client.Option{
		client.WithUser(cfg.TargetUser),
		client.WithHost(cfg.TargetHost),
		client.WithPort(cfg.TargetPort),
		client.WithKeyPath(cfg.TargetSSHKeyPath),
		client.WithWorkers(2),
	}

	client, err := client.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("new rtop client: %w", err)
	}

	return client, nil
}
