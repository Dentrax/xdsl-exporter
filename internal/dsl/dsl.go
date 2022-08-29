package dsl

import (
	"fmt"

	"3e8.eu/go/dsl"

	"github.com/Dentrax/xdsl-exporter/internal/config"
)

func GetSupportedClients() []string {
	clientTypes := dsl.GetClientTypes()
	var result []string
	for _, clientType := range clientTypes {
		result = append(result, clientType.ClientDesc().Title)
	}
	return result
}

func New(cfg config.Config) (dsl.Client, error) {
	c, err := GenerateConfigFrom(cfg)
	if err != nil {
		return nil, fmt.Errorf("generate dsl config: %w", err)
	}

	client, err := dsl.NewClient(*c)
	if err != nil {
		return nil, fmt.Errorf("new dsl client: %w", err)
	}

	return client, nil
}
