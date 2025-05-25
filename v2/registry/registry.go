package registry

import (
	"context"
	kc "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/errors"
	kr "github.com/go-kratos/kratos/v2/registry"
)

type Registry interface {
	Register() kr.Registrar
	Discover() kr.Discovery
}

type MockRegistry struct{}

func (m MockRegistry) Register(ctx context.Context, service *kr.ServiceInstance) error {
	return nil
}

func (m MockRegistry) Deregister(ctx context.Context, service *kr.ServiceInstance) error {
	return nil
}

type MockDiscovery struct{}

func (m MockDiscovery) GetService(ctx context.Context, serviceName string) ([]*kr.ServiceInstance, error) {
	return nil, errors.New(-1, "no discovery", "no discovery")
}

func (m MockDiscovery) Watch(ctx context.Context, serviceName string) (kr.Watcher, error) {
	return nil, errors.New(-1, "no discovery", "no discovery")
}

type MockGovernance struct {
}

func (g *MockGovernance) Register() kr.Registrar {
	return &MockRegistry{}
}

func (g *MockGovernance) Discover() kr.Discovery {
	return &MockDiscovery{}
}

func NewRegistry(kc kc.Config) Registry {
	driver, _ := kc.Value("registry.driver").String()
	switch driver {
	case "nacos":
		return NewNacosNaming(kc)
	default:
		return &MockGovernance{}
	}
}

func NewRegistrar(r Registry) kr.Registrar {
	return r.Register()
}

func NewDiscovery(r Registry) kr.Discovery {
	return r.Discover()
}
