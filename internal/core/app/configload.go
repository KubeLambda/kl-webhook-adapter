package app

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
)

// This prefix is used for environment variables to override config file values
const envVarPrefix = "WEBHOOK_ADAPTER"

const baseDir = "./configs"

func LoadConfig(deployment string) *Config {
	var cfg *Config = nil
	opt := loadOptions{
		EnvVarPrefix: envVarPrefix,
		Deployment:   deployment,
		BaseDir:      baseDir,
	}
	cfg = &Config{}
	err := loadFromYaml(opt, cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration file: %v", err))
	}
	cfg.Deployment = deployment
	return cfg
}

type loadOptions struct {
	EnvVarPrefix string
	Deployment   string
	BaseDir      string
}

func loadFromYaml(opt loadOptions, cfg *Config) error {
	v := viper.New()
	v.SetEnvPrefix(opt.EnvVarPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetConfigName(opt.Deployment)
	v.SetConfigType("yaml")
	v.AddConfigPath(opt.BaseDir)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading configuration file: %w", err)
	}
	err := v.Unmarshal(cfg, decoderWithEnvVariablesSupport())
	if err != nil {
		return fmt.Errorf("error parsing configuration file: %w", err)
	}
	return nil
}

// decoderWithEnvVariablesSupport allows us to resolve values containing environment variables,
// e.g. in yaml file you should be able to use constructions such as:
//
//	file_path: "${HOME}/notes.txt"
func decoderWithEnvVariablesSupport() func(c *mapstructure.DecoderConfig) {
	return func(c *mapstructure.DecoderConfig) {
		c.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			c.DecodeHook,
			mapstructure.StringToSliceHookFunc(","),
			replaceEnvVarsHookFunc,
		)
	}
}

func replaceEnvVars(value string) string {
	return os.ExpandEnv(value)
}

func replaceEnvVarsHookFunc(
	f reflect.Type,
	t reflect.Type,
	data interface{},
) (interface{}, error) {
	if f.Kind() != reflect.String || t.Kind() != reflect.String {
		return data, nil
	}

	return replaceEnvVars(data.(string)), nil
}
