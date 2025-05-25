package discovery

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/tuandq2112/go-microservices/shared/logger"
	"go.uber.org/zap"
)

type ConsulRegistrar struct {
	client    *api.Client
	serviceID string
	logger    logger.Logger
}

type ServiceRegistrationConfig struct {
	ID      string
	Name    string
	Address string
	Port    int
	Tags    []string
}

func NewConsulRegistrar() (*ConsulRegistrar, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &ConsulRegistrar{
		client: client,
		logger: logger.GetLogger(),
	}, nil
}


func (c *ConsulRegistrar) RegisterService(cfg ServiceRegistrationConfig) error {
	registration := &api.AgentServiceRegistration{
		ID:      cfg.ID,
		Name:    cfg.Name,
		Address: cfg.Address,
		Port:    cfg.Port,
		Tags:    cfg.Tags,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d/%s", cfg.Address, cfg.Port, "grpc.health.v1.Health"),
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",	
		},
	}

	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	c.serviceID = cfg.ID
	c.logger.Info("Service registered with Consul", zap.String("service_id", cfg.ID))
	return nil
}

func (c *ConsulRegistrar) DeregisterService() error {
	if c.serviceID == "" {
		return fmt.Errorf("no service registered to deregister")
	}

	err := c.client.Agent().ServiceDeregister(c.serviceID)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	c.logger.Info("Service deregistered from Consul", zap.String("service_id", c.serviceID))
	return nil
}
