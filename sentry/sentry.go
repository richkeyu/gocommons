package sentry

import (
	"fmt"
	"time"

	"github.com/richkeyu/gocommons/config"
	"github.com/getsentry/sentry-go"
)

const (
	configName = "sentry"
)

type Config struct {
	Dsn string `json:"dsn" yaml:"dsn"`
}

var conf Config

// Init 初始化
// 在入口出添加代码
// sentry.Init()
// defer sentry.Flush(2 * time.Second)
// https://docs.sentry.io/platforms/go/
func Init() {
	err := config.Load(configName, &conf)
	if err != nil {
		panic(fmt.Sprintf("load config fail: %s %s", configName, err))
	}

	if len(conf.Dsn) == 0 {
		panic(fmt.Sprintf("dsn empty: %s", configName))
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn: conf.Dsn,
		//EnableTracing: true,
		// Specify a fixed sample rate:
		// We recommend adjusting this value in production
		//TracesSampleRate: 1.0,
		// Or provide a custom sampler:
		//TracesSampler: sentry.TracesSamplerFunc(func(ctx sentry.SamplingContext) sentry.Sampled {
		//	return sentry.SampledTrue
		//}),
	})
	if err != nil {
		panic(fmt.Sprintf("sentry.Init %s", err))
	}
}

func Flush(duration time.Duration) {
	sentry.Flush(duration)
}

func CaptureMessage(message string) *sentry.EventID {
	return sentry.CaptureMessage(message)
}
