package registry

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

type ServiceRegistry struct {
	client *api.Client
}

func NewServiceRegistry(consulAddr string) (*ServiceRegistry, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ServiceRegistry{client: client}, nil
}

func (sr *ServiceRegistry) RegisterService(name, host string, port int) error {
	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", name, host, port),
		Name:    name,
		Address: host,
		Port:    port,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", host, port),
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	return sr.client.Agent().ServiceRegister(registration)
}
