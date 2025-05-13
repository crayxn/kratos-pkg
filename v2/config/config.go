package config

import (
	kc "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

func NewRemoteConfigSource(conf kc.Config) kc.Source {
	driver, err := conf.Value("remote_config.driver").String()
	if err != nil || driver == "" {
		driver = "noop" //default nacos
	}
	if driver == "nacos" {
		return NewNacosConfigSource(conf)
	}
	return nil
}

func New(path string) kc.Config {
	basic := kc.New(
		kc.WithSource(
			file.NewSource(path),
		),
	)
	// close
	defer func(basic kc.Config) {
		err := basic.Close()
		if err != nil {
			panic(err)
		}
	}(basic)
	// load
	if err := basic.Load(); err != nil {
		panic(err)
	}
	// disable remote
	if res, err := basic.Value("remote_config.enable").Bool(); !res || err != nil {
		return basic
	}

	return kc.New(
		kc.WithSource(
			//local
			file.NewSource(path),
			//remote
			NewRemoteConfigSource(basic),
			//... other
		),
	)
}
