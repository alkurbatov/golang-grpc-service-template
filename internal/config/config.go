// Package config implements service configuration layer.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func getSeconds(v *viper.Viper, key string) time.Duration {
	return time.Duration(v.GetInt(key)) * time.Second
}

func formatTimeout(timeout time.Duration) string {
	return fmt.Sprintf("%.fs", timeout.Seconds())
}

type Config struct {
	// General
	GRPCAddress string
	HTTPAddress string

	PrestopTimeout time.Duration
	GRPCTimeout    time.Duration

	// Logging
	LogLevel string
	LogJSON  bool

	// Other
	ProfAddress    string
	GRPCReflection bool
}

// New creates application config by reading environment variables.
func New() *Config {
	v := viper.New()

	v.SetDefault("GRPC_LISTEN_ADDRESS", "0.0.0.0:50051")
	v.SetDefault("HTTP_LISTEN_ADDRESS", "0.0.0.0:8888")
	v.SetDefault("GRACEFUL_SHUTDOWN_PRESTOP_TIMEOUT_SEC", 10)
	v.SetDefault("GRACEFUL_SHUTDOWN_GRPC_TIMEOUT_SEC", 200)

	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("JSON_LOGGING", true)

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	return &Config{
		GRPCAddress:    v.GetString("GRPC_LISTEN_ADDRESS"),
		HTTPAddress:    v.GetString("HTTP_LISTEN_ADDRESS"),
		PrestopTimeout: getSeconds(v, "GRACEFUL_SHUTDOWN_PRESTOP_TIMEOUT_SEC"),
		GRPCTimeout:    getSeconds(v, "GRACEFUL_SHUTDOWN_GRPC_TIMEOUT_SEC"),
		LogLevel:       v.GetString("LOG_LEVEL"),
		LogJSON:        v.GetBool("JSON_LOGGING"),
		ProfAddress:    v.GetString("PROF_LISTEN_ADDRESS"),
		GRPCReflection: v.GetBool("GRPC_REFLECTION"),
	}
}

func (c *Config) String() string {
	var sb strings.Builder

	sb.WriteString("General:\n")
	fmt.Fprintf(&sb, "\tGRPC_LISTEN_ADDRESS: %s\n", c.GRPCAddress)
	fmt.Fprintf(&sb, "\tHTTP_LISTEN_ADDRESS: %s\n", c.HTTPAddress)
	fmt.Fprintf(&sb, "\tGRACEFUL_SHUTDOWN_PRESTOP_TIMEOUT_SEC: %s\n", formatTimeout(c.PrestopTimeout))
	fmt.Fprintf(&sb, "\tGRACEFUL_SHUTDOWN_GRPC_TIMEOUT_SEC: %s\n", formatTimeout(c.GRPCTimeout))

	sb.WriteString("Logging:\n")
	fmt.Fprintf(&sb, "\tLOG_LEVEL: %s\n", c.LogLevel)
	fmt.Fprintf(&sb, "\tJSON_LOGGING: %t\n", c.LogJSON)

	sb.WriteString("Other:\n")
	fmt.Fprintf(&sb, "\tPROF_LISTEN_ADDRESS: %s\n", c.ProfAddress)
	fmt.Fprintf(&sb, "\tGRPC_REFLECTION: %t\n", c.GRPCReflection)

	return sb.String()
}
