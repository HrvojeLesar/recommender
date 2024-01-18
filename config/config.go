package config

import "os"

func LookupEnvVariableOrDefault(variable string, defaultVar string) string {
	envVar, found := os.LookupEnv(variable)
	if !found {
		return defaultVar
	}
	return envVar
}

type Config struct {
	Mongo MongoConfig
}

func New() Config {
	return Config{
		Mongo: NewMongoConfig(),
	}
}
