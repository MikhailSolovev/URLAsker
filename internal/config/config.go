package config

import (
	"github.com/caarlos0/env/v6"
	"sync"
	"time"
)

// singleton config

var Config = Options{}
var once sync.Once

type Options struct {
	ServiceName string `env:"SERVICE_NAME" envDefault:"asker"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"debug"`
	// IsPretty - if it's true human-readable outcome generated, otherwise json
	IsPretty            bool          `env:"IS_PRETTY" envDefault:"true"`
	SrvAddr             string        `env:"SRV_ADDR" envDefault:":8080"`
	SrvWriteTimeout     time.Duration `env:"SRV_WRITE_TIMEOUT" envDefault:"15s"`
	SrvReadTimeout      time.Duration `env:"SRV_READ_TIMEOUT" envDefault:"15s"`
	IdleTimeout         time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
	GracefulTimeout     time.Duration `env:"GRACEFUL_TIMEOUT" envDefault:"15s"`
	PostgresUser        string        `env:"POSTGRES_USER"`
	PostgresPass        string        `env:"POSTGRES_PASS"`
	PostgresHost        string        `env:"POSTGRES_HOST"`
	PostgresPort        string        `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresDB          string        `env:"POSTGRES_DB"`
	PostgresConnTimeout time.Duration `env:"POSTGRES_CONN_TIMEOUT" envDefault:"5s"`
}

func New() *Options {
	var err error
	once.Do(func() {
		err = env.Parse(&Config)
	})

	if err != nil {
		return nil
	} else {
		return &Config
	}
}
