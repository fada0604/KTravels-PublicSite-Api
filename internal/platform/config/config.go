package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	GraphQL  GraphQLConfig  `mapstructure:"graphql"`
	Logger   LoggerConfig  `mapstructure:"logger"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type GraphQLConfig struct {
	PlaygroundEnabled bool `mapstructure:"playgroundEnabled"`
	Introspection     bool `mapstructure:"introspection"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type RabbitMQConfig struct {
	HostName string `mapstructure:"hostname"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Port    int    `mapstructure:"port"`
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetEnvPrefix("app")
	v.SetEnvKeyReplacer(strings.NewReplacer("_", "."))
	v.AutomaticEnv()

	v.SetDefault("app.name", "ktravels-publicsite-api")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.port", 8080)

	v.SetDefault("database.uri", "mongodb://localhost:27017")
	v.SetDefault("database.database", "ktravels_publicsite")

	v.SetDefault("graphql.playgroundEnabled", true)
	v.SetDefault("graphql.introspection", true)

	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")

	v.SetDefault("rabbitmq.hostname", "localhost")
	v.SetDefault("rabbitmq.username", "guest")
	v.SetDefault("rabbitmq.password", "guest")
	v.SetDefault("rabbitmq.port", 5672)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	<-time.After(1 * time.Millisecond)

	if cfg.App.Name == "" {
		return nil, fmt.Errorf("app.name is required")
	}
	if cfg.Database.URI == "" {
		return nil, fmt.Errorf("database.uri is required")
	}

	return &cfg, nil
}
