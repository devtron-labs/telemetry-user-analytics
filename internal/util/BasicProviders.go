package util

import (
	"fmt"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
	"net/http"
)

var (
	// Logger is the defaut logger
	logger *zap.SugaredLogger
	//FIXME: remove this
	//defer Logger.Sync()
)

// Deprecated: instead calling this method inject logger from wire
func GetLogger() *zap.SugaredLogger {
	return logger
}

type SentryConfig struct {
	DSN           string `env:"DSN" envDefault:"https://0e57adc987494c3f99c56e8287475c20@sentry.io/1887839"`
	SentryEnv     string `env:"SENTRY_ENV" envDefault:"Staging"`
	SentryEnabled bool   `env:"SENTRY_ENABLED" envDefault:"false"`
}

func init() {
	cfg := &SentryConfig{}
	err := env.Parse(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := zap.NewProduction()
	if err != nil {
		panic("failed to create the default logger: " + err.Error())
	}
	if cfg.SentryEnabled {
		logger = l.Sugar() //modifyToSentryLogger(l, cfg.DSN, cfg.SentryEnv)
	} else {
		logger = l.Sugar()
	}
}

/*func modifyToSentryLogger(log *zap.Logger, DSN string, env string) *zap.SugaredLogger {
	cfg := zapsentry.Configuration{
		//when to send message to sentry
		Tags: map[string]string{
			"component": "system",
		},
		Level: zapcore.ErrorLevel,
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSNAndEnv(DSN, env))
	//core.Enabled(zapcore.ErrorLevel)
	//in case of err it will return noop core. so we can safely attach it
	if err != nil {
		log.Warn("failed to init zap", zap.Error(err))
	}
	return zapsentry.AttachCoreToLogger(core, log).Sugar()
}*/

func NewSugardLogger() *zap.SugaredLogger {
	return logger
}

func NewHttpClient() *http.Client {
	return http.DefaultClient
}
