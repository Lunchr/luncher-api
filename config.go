package main

import (
	"strconv"

	"github.com/deiwin/gonfigure"
)

var (
	domainProperty = gonfigure.NewRequiredEnvProperty("LUNCHER_DOMAIN")
	portProperty   = gonfigure.NewEnvProperty("LUNCHER_PORT", "8080")
)

type Config struct {
	Domain string
	Port   int
}

func NewConfig() (conf Config, err error) {
	port, err := strconv.Atoi(portProperty.Value())
	if err != nil {
		return
	}
	conf = Config{
		Domain: domainProperty.Value(),
		Port:   port,
	}
	return
}
