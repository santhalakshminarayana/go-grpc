// Package config inits env config from a file and environment
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var cfg *config

type config struct {
	ServiceName          string `config:"SERVICE_NAME"`
	GrpcReflectionEnable bool   `config:"GRPC_REFLECTION_ENABLE"`
}

func GetConfig() *config {
	if cfg != nil {
		return cfg
	}

	// Initialize config
	InitConfig("")
	return cfg
}

func (c *config) validate() error {
	return nil
}

func loadConfig(viperConfig *viper.Viper) error {
	// Read from environment
	// viperConfig.AutomaticEnv()

	// Initalize new config
	cfg = &config{
		ServiceName:          viperConfig.GetString("SERVICE_NAME"),
		GrpcReflectionEnable: viperConfig.GetBool("GRPC_REFLECTION_ENABLE"),
	}

	// Validate config
	if err := cfg.validate(); err != nil {
		return fmt.Errorf("error in validating Config: %w", err)
	}

	return nil
}

func InitConfig(envFilePath string) {
	var err error

	defer func() {
		if err != nil {
			panic(fmt.Sprintf("Error initalizing Config: %v", err))
		}
	}()

	// New Viper
	viperConfig := viper.New()

	// Set and Read env file config
	if envFilePath != "" {
		viperConfig.SetConfigFile(envFilePath)
	}

	if err = viperConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = fmt.Errorf("error reading env file. No file found: %w", err)
			return
		}
		err = fmt.Errorf("error reading env file config: %w", err)
		return
	}

	// Load config from both local file and environment
	if err = loadConfig(viperConfig); err != nil {
		err = fmt.Errorf("error loading config: %w", err)
		return
	}
}
