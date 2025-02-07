package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/sagikazarmark/modern-go-application/internal/platform/database"
	"github.com/sagikazarmark/modern-go-application/internal/platform/log"
	"github.com/sagikazarmark/modern-go-application/internal/platform/opencensus"
	"github.com/sagikazarmark/modern-go-application/internal/platform/redis"
	"github.com/sagikazarmark/modern-go-application/internal/platform/watermill"
)

// configuration holds any kind of configuration that comes from the outside world and
// is necessary for running the application.
type configuration struct {
	// Meaningful values are recommended (eg. production, development, staging, release/123, etc)
	Environment string

	// Turns on some debug functionality
	Debug bool

	// Timeout for graceful shutdown
	ShutdownTimeout time.Duration

	// Log configuration
	Log log.Config

	// Instrumentation configuration
	Instrumentation instrumentationConfig

	// OpenCensus configuration
	Opencensus struct {
		Exporter struct {
			Enabled                   bool
			opencensus.ExporterConfig `mapstructure:",squash"`
		}

		Trace opencensus.TraceConfig

		// Prometheus configuration
		Prometheus struct {
			Enabled bool
		}
	}

	// App configuration
	App struct {
		// HTTP server address
		// nolint: golint
		HttpAddr string

		// GRPC server address
		GrpcAddr string
	}

	// Database connection information
	Database database.Config

	// Redis configuration
	Redis redis.Config

	// Watermill configuration
	Watermill struct {
		RouterConfig watermill.RouterConfig
	}
}

// Validate validates the configuration.
func (c configuration) Validate() error {
	if c.Environment == "" {
		return errors.New("environment is required")
	}

	if err := c.Instrumentation.Validate(); err != nil {
		return err
	}

	if c.App.HttpAddr == "" {
		return errors.New("http app server address is required")
	}

	if c.App.GrpcAddr == "" {
		return errors.New("grpc app server address is required")
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	// Uncomment to enable redis config validation
	// if err := c.Redis.Validate(); err != nil {
	// 	return err
	// }

	return nil
}

// instrumentationConfig represents the instrumentation related configuration.
type instrumentationConfig struct {
	// Instrumentation HTTP server address
	Addr string
}

// Validate validates the configuration.
func (c instrumentationConfig) Validate() error {
	if c.Addr == "" {
		return errors.New("instrumentation http server address is required")
	}

	return nil
}

// configure configures some defaults in the Viper instance.
func configure(v *viper.Viper, p *pflag.FlagSet) {
	// Viper settings
	v.AddConfigPath(".")
	v.AddConfigPath(fmt.Sprintf("$%s_CONFIG_DIR/", strings.ToUpper(envPrefix)))

	// Environment variable settings
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	// Application constants
	v.Set("appName", appName)

	// Global configuration
	v.SetDefault("environment", "production")
	v.SetDefault("debug", false)
	v.SetDefault("shutdownTimeout", 15*time.Second)
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		v.SetDefault("no_color", true)
	}

	// Log configuration
	v.SetDefault("log.format", "json")
	v.SetDefault("log.level", "info")
	v.RegisterAlias("log.noColor", "no_color")

	// Instrumentation configuration
	p.String("instrumentation-addr", ":10000", "Instrumentation HTTP server address")
	_ = v.BindPFlag("instrumentation.addr", p.Lookup("instrumentation-addr"))
	v.SetDefault("instrumentation.addr", ":10000")

	// OpenCensus configuration
	v.SetDefault("opencensus.exporter.enabled", false)
	_ = v.BindEnv("opencensus.exporter.address")
	_ = v.BindEnv("opencensus.exporter.insecure")
	_ = v.BindEnv("opencensus.exporter.reconnectPeriod")
	v.SetDefault("opencensus.trace.sampling.sampler", "never")
	v.SetDefault("opencensus.prometheus.enabled", false)

	// App configuration
	p.String("http-addr", ":8000", "App HTTP server address")
	_ = v.BindPFlag("app.httpAddr", p.Lookup("http-addr"))
	v.SetDefault("app.httpAddr", ":8000")

	p.String("grpc-addr", ":8001", "App GRPC server address")
	_ = v.BindPFlag("app.grpcAddr", p.Lookup("grpc-addr"))
	v.SetDefault("app.grpcAddr", ":8001")

	// Database configuration
	_ = v.BindEnv("database.host")
	v.SetDefault("database.port", 3306)
	_ = v.BindEnv("database.user")
	_ = v.BindEnv("database.pass")
	_ = v.BindEnv("database.name")
	v.SetDefault("database.params", map[string]string{
		"collation": "utf8mb4_general_ci",
	})

	// Redis configuration
	_ = v.BindEnv("redis.host")
	v.SetDefault("redis.port", 6379)
	_ = v.BindEnv("redis.password")

	// Watermill configuration
	v.RegisterAlias("watermill.routerConfig.closeTimeout", "shutdownTimeout")
}
