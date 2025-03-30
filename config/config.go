package config

import (
	"fmt"
	"maps"
	"reflect"
)

type ConfigValue[T any] struct {
	Value   T
	Default T
}

func NewConfigValue[T comparable](value, defaultValue T) (ConfigValue[T], error) {
	var zero T
	if reflect.DeepEqual(defaultValue, zero) {
		return ConfigValue[T]{}, fmt.Errorf(`default value "%v" is empty`, defaultValue)
	}
	return ConfigValue[T]{Value: value, Default: defaultValue}, nil
}

var singletonCfg *GConfig

type GConfig struct {
	values map[string]ConfigValue[any]
}

func NewConfig(inputValue map[string]ConfigValue[any]) *GConfig {
	cfg := &GConfig{
		values: make(map[string]ConfigValue[any]),
	}
	cfg.SetDefault()
	cfg.Copy(inputValue)

	singletonCfg = cfg

	return singletonCfg
}

func (cfg *GConfig) SetDefault() {
	cfg.values["log_level"] = ConfigValue[any]{"", "debug"}
	cfg.values["log_folder"] = ConfigValue[any]{"", "log"}
	cfg.values["log_max_size"] = ConfigValue[any]{0, 1024}
	cfg.values["log_max_backup"] = ConfigValue[any]{0, 1024}
	cfg.values["log_max_age"] = ConfigValue[any]{0, 28}
	cfg.values["log_compress"] = ConfigValue[any]{false, true}
}

func (cfg *GConfig) Copy(inputValue map[string]ConfigValue[any]) {
	if inputValue == nil {
		return
	}
	maps.Copy(cfg.values, inputValue)
}

func Get[T any](key string) T {
	if singletonCfg == nil {
		NewConfig(nil)
	}

	var returnValue T

	anyValue, ok := singletonCfg.values[key]
	if !ok {
		return returnValue
	}

	value, _ := anyValue.Value.(T)
	defaultValue, _ := anyValue.Default.(T)

	if reflect.DeepEqual(value, returnValue) {
		return defaultValue
	}

	return value
}
