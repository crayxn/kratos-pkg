package logger

import (
	kc "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
)

func NewLogger(kc kc.Config) log.Logger {
	//use zap
	var conf ZapConfig
	if err := kc.Value("logger").Scan(&conf); err != nil {
		conf = ZapConfig{
			Level:  "info",
			Writer: STDOUT,
		}
	}
	return NewZapLogger(&conf)
}
