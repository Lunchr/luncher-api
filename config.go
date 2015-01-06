package main

import (
	"strconv"

	"github.com/deiwin/luncher-api/config"
)

var (
	domainProperty = config.NewRequiredEnvProperty("LUNCHER_DOMAIN")
	portProperty   = config.NewEnvProperty("LUNCHER_PORT", "8080")
)

type Config struct {
	Domain string
	Port   int
}

func NewConfig() (conf Config, err error) {
	port, err := strconv.Atoi(portProperty.DefaultValue())
	if err != nil {
		return
	}
	conf = Config{
		Domain: domainProperty.DefaultValue(),
		Port:   port,
	}
	return
}
