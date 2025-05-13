package registry

import (
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	kc "github.com/go-kratos/kratos/v2/config"
	kr "github.com/go-kratos/kratos/v2/registry"
)

type NacosConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	UserName  string `json:"user"`
	Password  string `json:"password"`
	Namespace string `json:"namespace"`
	Group     string `json:"group"`
}

func NewNacosNaming(kc kc.Config) Registry {
	var conf NacosConfig
	if err := kc.Value("registry.nacos").Scan(&conf); err != nil {
		panic(err)
	}
	//创建nacos 客户端
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: []constant.ServerConfig{*constant.NewServerConfig(conf.Host, uint64(conf.Port))},
			ClientConfig: &constant.ClientConfig{
				NamespaceId: conf.Namespace,
				Username:    conf.UserName,
				Password:    conf.Password,
			},
		},
	)
	if err != nil {
		panic(err)
	}
	client := nacos.New(
		namingClient,
		nacos.WithGroup(func() string {
			if conf.Group == "" {
				return constant.DEFAULT_GROUP
			}
			return conf.Group
		}()),
	)
	return &NacosNaming{
		client,
	}
}

type NacosNaming struct {
	client *nacos.Registry
}

func (n NacosNaming) Register() kr.Registrar {
	return n.client
}

func (n NacosNaming) Discover() kr.Discovery {
	return n.client
}
