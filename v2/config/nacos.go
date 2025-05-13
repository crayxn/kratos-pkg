package config

import (
	nc "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	kc "github.com/go-kratos/kratos/v2/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosConfig struct {
	Host        string `json:"host"`
	Port        int32  `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	GroupName   string `json:"group_name"`
	NamespaceId string `json:"namespace_id"`
	DataId      string `json:"data_id"`
	LogLevel    string `json:"log_level"`
	LogDir      string `json:"log_dir"`
	CacheDir    string `json:"cache_dir"`
}

// NewNacosConfigSource 创建Nacos配置源
func NewNacosConfigSource(kc kc.Config) kc.Source {
	var conf NacosConfig
	if err := kc.Value("remote_config.nacos").Scan(&conf); err != nil {
		panic(err)
	}
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
