package consul

import (
	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/hashicorp/consul/api"
)

func GetConsulConfiguration(address string, path string) (config.Config, error) {
	consulClient, err := api.NewClient(&api.Config{
		Address: address,
	})
	if err != nil {
		return nil, err
	}
	cs, err := consul.New(consulClient, consul.WithPath(path))
	if err != nil {
		return nil, err
	}
	c := config.New(config.WithSource(cs))
	if err = c.Load(); err != nil {
		return nil, err
	}
	return c, nil
}
