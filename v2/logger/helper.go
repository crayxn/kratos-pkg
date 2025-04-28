package logger

import (
	"fmt"
	kc "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
)

type LogHelper struct {
	log.Helper
	modules       map[string]*log.Helper
	config        kc.Config
	defaultConfig *ZapConfig
}

func NewLogHelper(logger log.Logger, config kc.Config) *LogHelper {
	return &LogHelper{
		Helper:        *log.NewHelper(log.With(logger, "module", "default")),
		modules:       map[string]*log.Helper{},
		config:        config,
		defaultConfig: nil,
	}
}

func (h *LogHelper) Module(module string) *log.Helper {
	if helper, ok := h.modules[module]; ok {
		return helper
	}
	//获取配置
	moduleConfig, err := h.config.Value(fmt.Sprintf("modules.%s", module)).String()
}

func (h *LogHelper) getConfigWithDefault(module string) *ZapConfig {
	if h.defaultConfig == nil {
		if err := h.config.Value("logger.default").Scan(h.defaultConfig); err != nil {
			h.defaultConfig = &ZapConfig{
				Level:  "info",
				Writer: STDOUT,
			}
		}
	}
	var conf ZapConfig
	if err := h.config.Value(fmt.Sprintf("logger.%s", module)).Scan(&conf); err == nil {
		return nil
	}
	//
}

func NewLogger(kc kc.Config) log.Logger {
	var conf ZapConfig
	if err := kc.Value("logger.default").Scan(&conf); err != nil {
		conf = ZapConfig{
			Level:  "info",
			Writer: STDOUT,
		}
	}
	return NewZapLogger(&conf)
}
