package config

var CONFIG_MODE ConfigMode = DEV_DOCKER

type ConfigMode int

const (
	DEV_DOCKER ConfigMode = iota
	DEV
	PRODUCTION
)
