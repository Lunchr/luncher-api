/*
Package config helps creating configuration structs.

The intended usage would be a simple struct that calls Value() on
a fields initialization. E.g.

	var portProperty   = config.NewEnvProperty("PORT", "8080")
	// If the $DOMAIN env variable not set the configuration creation will fail with a fatal error
	var domainProperty = config.NewRequiredEnvProperty("DOMAIN")

	type Config struct {
		Port   string
		Domain string
	}

	func NewConfig() Config {
		return Config{
			Port:   portProperty.Value(),
			Domain: domainProperty.Value(),
		}
	}
*/
package config

import (
	"log"
	"os"
)

// Property can be used to fetch default values for configuration properties.
type Property interface {
	Value() string
}

// NewRequiredEnvProperty returns a Property that gets its value from
// the specified environment variable. Panics if the variable is not set.
func NewRequiredEnvProperty(envVariableName string) Property {
	return requiredEnvProperty{
		envVariableName: envVariableName,
	}
}

// NewEnvProperty returns a Property that gets its value from the
// specified environment variable. If the environment vatiable is not set
// the fallback value will be used instead
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

func (prop envProperty) Value() (defaultValue string) {
	defaultValue = os.Getenv(prop.envVariableName)
	if defaultValue == "" {
		defaultValue = prop.fallbackValue
	}
	return
}

type requiredEnvProperty struct {
	envVariableName string
}

func (prop requiredEnvProperty) Value() (defaultValue string) {
	defaultValue = os.Getenv(prop.envVariableName)
	if defaultValue == "" {
		log.Fatalf("Please set the %s environment variable", prop.envVariableName)
	}
	return
}
