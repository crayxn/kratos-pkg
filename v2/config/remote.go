package config

import (
	nc "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	kc "github.com/go-kratos/kratos/v2/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosConfig struct {
	Host        string `json:"host,omitempty"`
	Port        int32  `json:"port,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	GroupName   string `json:"group_name,omitempty"`
	NamespaceId string `json:"namespace_id,omitempty"`
	DataId      string `json:"data_id,omitempty"`
	LogLevel    string `json:"log_level,omitempty"`
	LogDir      string `json:"log_dir,omitempty"`
	CacheDir    string `json:"cache_dir,omitempty"`
}

func NewRemoteConfigSource(conf kc.Config) kc.Source {
	driver, err := conf.Value("remote_config.driver").String()
	if err != nil || driver == "" {
		driver = "nacos" //default nacos
	}
	if driver == "nacos" {
		var c NacosConfig
		if err := conf.Value("remote_config.nacos").Scan(&c); err != nil {
			panic(err)
		}
		return NewNacosConfigSource(&c)
	}
	return nil
}

// NewNacosConfigSource 创建Nacos配置源
func NewNacosConfigSource(conf *NacosConfig) kc.Source {
	//config default
	if conf.LogLevel == "" {
		conf.LogLevel = "error"
	}
	if conf.LogDir == "" {
		conf.LogDir = "./runtime/log/nacos"
	}
	if conf.CacheDir == "" {
		conf.CacheDir = "./runtime/cache/nacos"
	}
	//create client
	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ServerConfigs: []constant.ServerConfig{*constant.NewServerConfig(conf.Host, uint64(conf.Port))},
		ClientConfig: &constant.ClientConfig{
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              conf.LogDir,
			CacheDir:            conf.CacheDir,
			LogLevel:            conf.LogLevel,
			NamespaceId:         conf.NamespaceId,
			Username:            conf.Username,
			Password:            conf.Password,
		},
	})
	if err != nil {
		panic(err)
	}
	return nc.NewConfigSource(client, nc.WithGroup(conf.GroupName), nc.WithDataID(conf.DataId))
}
