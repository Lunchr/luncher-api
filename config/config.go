/*
Package config helps creating configuration structs.

The intended usage would be a simple struct that calls DefaultValue() on
a fields initialization. E.g.

  var portProperty = config.NewEnvProperty("PORT", "8080")

  type Config struct {
    Port  string
  }

  func NewConfig() *Config {
    return &Config {
      Port: portProperty.DefaultValue()
    }
  }
*/
package config

import (
	"log"
	"os"
)

// Property can be used to define configuration properties with default values.
type Property interface {
	DefaultValue() string
}

// NewRequiredEnvProperty returns a Property that gets its default value from
// the specified environment variable. Panics if the variable is not set.
func NewRequiredEnvProperty(envVariableName string) Property {
	return requiredEnvProperty{
		envVariableName: envVariableName,
	}
}

// NewEnvProperty returns a Property that gets its default value from the
// specified environment variable. If the environment vatiable is not set
// the fallback value will be set as default instead
func NewEnvProperty(envVariableName string, fallbackValue string) Property {
	return envProperty{
		envVariableName: envVariableName,
		fallbackValue:   fallbackValue,
	}
}

type envProperty struct {
	envVariableName string
	fallbackValue   string
}

func (prop envProperty) DefaultValue() (defaultValue string) {
	defaultValue = os.Getenv(prop.envVariableName)
	if defaultValue == "" {
		defaultValue = prop.fallbackValue
	}
	return
}

type requiredEnvProperty struct {
	envVariableName string
}

func (prop requiredEnvProperty) DefaultValue() (defaultValue string) {
	defaultValue = os.Getenv(prop.envVariableName)
	if defaultValue == "" {
		log.Fatalf("Please set the %s environment variable", prop.envVariableName)
	}
	return
}
