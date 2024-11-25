package discovery

import (
	"ecommerce/warehouse-service/utils"
	"fmt"

	"github.com/hashicorp/consul/api"
)

type ServiceDiscovery struct {
	client *api.Client
}

func NewServiceDiscovery(consulAddr string) (*ServiceDiscovery, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ServiceDiscovery{client: client}, nil
}

func (sd *ServiceDiscovery) Register(serviceName, serviceID string, port int) error {
	registration := &api.AgentServiceRegistration{
		ID:   serviceID,
		Name: serviceName,
		Port: port,
		Tags: []string{"warehouse", "microservice"},
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://localhost:%d/health", port),
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	err := sd.client.Agent().ServiceRegister(registration)
	if err != nil {
		utils.Logger.Errorf("Failed to register service: %v", err)
		return err
	}

	utils.Logger.Infof("Successfully registered service: %s", serviceID)
	return nil
}

func (sd *ServiceDiscovery) Deregister(serviceID string) error {
	err := sd.client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		utils.Logger.Errorf("Failed to deregister service: %v", err)
		return err
	}

	utils.Logger.Infof("Successfully deregistered service: %s", serviceID)
	return nil
}

func (sd *ServiceDiscovery) GetService(serviceName string) ([]*api.ServiceEntry, error) {
	services, _, err := sd.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		utils.Logger.Errorf("Failed to get service: %v", err)
		return nil, err
	}
	return services, nil
}
