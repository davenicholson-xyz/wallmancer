package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	values map[string]any
}

func load(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	cfg := &Config{values: make(map[string]any)}

	for key, value := range raw {
		envKey := "WMCR_" + strings.ToUpper(key)
		if envVal, exists := os.LookupEnv(envKey); exists {
			switch value.(type) {
			case bool:
				if b, err := strconv.ParseBool(envVal); err == nil {
					value = b
				}
			case int:
				if i, err := strconv.Atoi(envVal); err == nil {
					value = i
				}
			case float64:
				if i, err := strconv.Atoi(envVal); err == nil {
					value = i
				}
			default:
				value = envVal
			}
		}
		cfg.values[key] = value
	}

	return cfg, nil
}

func (c *Config) GetString(key string) string {
	if val, ok := c.values[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

func (c *Config) GetInt(key string) int {
	if val, ok := c.values[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

func (c *Config) GetBool(key string) bool {
	if val, ok := c.values[key]; ok {
		switch v := val.(type) {
		case bool:
			return v
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				return b
			}
		}
	}
	return false
}

func (c *Config) GetStringWithDefault(key, defaultValue string) string {
	if val := c.GetString(key); val != "" {
		return val
	}
	return defaultValue
}

func (c *Config) GetIntWithDefault(key string, defaultValue int) int {
	if val := c.GetInt(key); val != 0 {
		return val
	}
	return defaultValue
}

func (c *Config) GetBoolWithDefault(key string, defaultValue bool) bool {
	if _, exists := c.values[key]; exists {
		return c.GetBool(key)
	}
	return defaultValue
}

func New(path string) (*Config, error) {
	cfg, err := load("config.yml")
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
