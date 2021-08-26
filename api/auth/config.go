package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type ConfigModel struct {
	Port        string
	GraphQL     GraphQLConfigModel
	JWT         JWTConfigModel
	Auth        AuthConfigModel
	Redis       RedisConfigModel
	RedirectURL *url.URL
}

type GraphQLConfigModel struct {
	URL    string
	Secret string
}

type JWTConfigModel struct {
	AuthTTL    int
	RefreshTTL int
	Secret     string
}

type ExternalAuthConfigModel struct {
	ID          string
	Secret      string
	RedirectURL string
}

type AuthConfigModel struct {
	Google ExternalAuthConfigModel
}

type RedisConfigModel struct {
	Host     string
	Port     string
	Password string
	DB       int
}

var envMap = map[string]func(*ConfigModel, string) (*ConfigModel, error){
	"PORT": func(m *ConfigModel, value string) (*ConfigModel, error) {
		m.Port = value
		return m, nil
	},
	"REDIRECT_URL": func(m *ConfigModel, value string) (*ConfigModel, error) {
		var err error
		m.RedirectURL, err = url.Parse(value)
		return m, err
	},
	"AUTH_GOOGLE_ID": func(m *ConfigModel, value string) (*ConfigModel, error) {
		m.Auth.Google.ID = value
		return m, nil
	},
	"AUTH_GOOGLE_SECRET": func(m *ConfigModel, value string) (*ConfigModel, error) {
		m.Auth.Google.Secret = value
		return m, nil
	},
	"AUTH_GOOGLE_REDIRECT_URL": func(m *ConfigModel, value string) (*ConfigModel, error) {
		m.Auth.Google.RedirectURL = value
		return m, nil
	},
	"JWT_SECRET": func(m *ConfigModel, value string) (*ConfigModel, error) {
		m.JWT.Secret = value
		return m, nil
	},
	"JWT_AUTH_TTL": func(m *ConfigModel, value string) (*ConfigModel, error) {
		ivalue, err := strconv.Atoi(value)
		if err != nil {
			return m, err
		}

		m.JWT.AuthTTL = ivalue
		return m, nil
	},
	"JWT_REFRESH_TTL": func(m *ConfigModel, value string) (*ConfigModel, error) {
		ivalue, err := strconv.Atoi(value)
		if err != nil {
			return m, err
		}

		m.JWT.RefreshTTL = ivalue
		return m, nil
	},
	"GRAPHQL_SECRET": func(m *ConfigModel, value string) (*ConfigModel, error) {
		m.GraphQL.Secret = value
		return m, nil
	},
	"GRAPHQL_URL": func(m *ConfigModel, value string) (*ConfigModel, error) {
		m.GraphQL.URL = value
		return m, nil
	},
	"REDIS_HOST": func(m *ConfigModel, value string) (*ConfigModel, error) {
		m.Redis.Host = value
		return m, nil
	},
	"REDIS_PORT": func(m *ConfigModel, value string) (*ConfigModel, error) {
		m.Redis.Port = value
		return m, nil
	},
	"REDIS_PASSWORD": func(m *ConfigModel, value string) (*ConfigModel, error) {
		m.Redis.Password = value
		return m, nil
	},
	"REDIS_DB": func(m *ConfigModel, value string) (*ConfigModel, error) {
		ivalue, err := strconv.Atoi(value)
		if err != nil {
			return m, err
		}

		m.Redis.DB = ivalue
		return m, nil
	},
}

var optionalVars = []string{
	"PORT",
	"AUTH_GITHUB_ID",
	"AUTH_GITHUB_SECRET",
	"REDIS_HOST",
	"REDIS_PORT",
	"REDIS_DB",
	"REDIS_PASSWORD",
}

var defaults = map[string]string{
	"PORT":       "8080",
	"REDIS_HOST": "redis",
	"REDIS_PORT": "6379",
}

func isOptional(envVar string) bool {
	for _, v := range optionalVars {
		if v == envVar {
			return true
		}
	}

	return false
}

func CreateConfig() (*ConfigModel, error) {
	config := &ConfigModel{}
	var err error

	for envName, mapper := range envMap {
		value, envOK := os.LookupEnv(envName)
		if !envOK {
			if !isOptional(envName) {
				return nil, fmt.Errorf("Missing %s environment variable", envName)
			}

			if dft, ok := defaults[envName]; ok {
				value = dft
			} else {
				continue
			}
		}

		config, err = mapper(config, value)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}
