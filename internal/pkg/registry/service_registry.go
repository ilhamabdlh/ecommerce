package registry

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/hashicorp/consul/api"
)

type ServiceRegistry struct {
	client    *api.Client
	serviceID string
}

func NewServiceRegistry(consulAddr string, serviceName string, servicePort int) (*ServiceRegistry, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()
	serviceID := fmt.Sprintf("%s-%s-%d", serviceName, hostname, servicePort)

	return &ServiceRegistry{
		client:    client,
		serviceID: serviceID,
	}, nil
}

func (sr *ServiceRegistry) Register(serviceName string, tags []string) error {
	ip, err := sr.getLocalIP()
	if err != nil {
		return err
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}

	reg := &api.AgentServiceRegistration{
		ID:      sr.serviceID,
		Name:    serviceName,
		Tags:    tags,
		Port:    port,
		Address: ip,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", ip, port),
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	return sr.client.Agent().ServiceRegister(reg)
}

func (sr *ServiceRegistry) Deregister() error {
	return sr.client.Agent().ServiceDeregister(sr.serviceID)
}

func (sr *ServiceRegistry) getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no IP address found")
}
