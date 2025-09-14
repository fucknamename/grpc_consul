package registry

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

type ConsulRegistry struct {
	client *api.Client
}

func NewConsulRegistry() *ConsulRegistry {
	cfg := api.DefaultConfig()
	// cfg.Address = "127.0.0.1:8500"  // 默认就是这个，一般不需要改
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to connect to consul: %v", err)
	}
	return &ConsulRegistry{client: client}
}

// Register gRPC 服务到 Consul
func (r *ConsulRegistry) Register(name, id, ip string, port int) error {
	reg := &api.AgentServiceRegistration{
		ID:      id,   // 服务节点的名称
		Name:    name, // 服务名
		Address: ip,   // 或内网IP
		Port:    port,
		Tags:    []string{"v1", "in"}, // 版本, 国家/币种, 可以为空
		Check: &api.AgentServiceCheck{
			// HTTP:                           "http://127.0.0.1:8080",
			GRPC:                           fmt.Sprintf("%s:%d", ip, port), // grpc, 执行健康检查的地址, 会传到 Health.Check 函数中
			Interval:                       "10s",                          // 健康检查间隔
			Timeout:                        "5s",                           // 超时时间
			DeregisterCriticalServiceAfter: "1m",                           // 某个服务的健康检查一直处于 critical 状态超过指定时间，Consul 会自动把这个服务从注册中心里移除
		},
	}
	return r.client.Agent().ServiceRegister(reg)
}

// Deregister 服务
func (r *ConsulRegistry) Deregister(id string) error {
	return r.client.Agent().ServiceDeregister(id)
}

// Discover 服务列表
func (r *ConsulRegistry) Discover(name string) ([]*api.ServiceEntry, error) {
	services, _, err := r.client.Health().Service(name, "", true, nil)
	return services, err
}
