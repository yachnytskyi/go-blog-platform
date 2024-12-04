package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
)

const (
	yamlConfigPath        = "config/yaml/v1"
	environmentConfigPath = "config/environment"
	v1                    = "v1=%s\n"

	expectedLocation           = "pkg.dependency.factory.config.viper."
	expectedLogger             = "Zerolog"
	expectedPort               = "8080"
	expectedMongoDBDatabaseURI = "mongodb://root:root@localhost:27017/golang_mongodb"
	expectedAllowedOrigins     = "http://localhost:8080"
	expectedAllowedMethod      = "GET"

	writePermissions = 0755
	readPermissions  = 0644
	openFileError    = "open %s: no such file or directory"
	yamlParsingError = "While parsing config: yaml: line 1: did not find expected ',' or ']'"
)

func setupYamlFilePath() string {
	yamlContent := []byte(`
core:
  logger: Zerolog
  email: GoMail
  database: MongoDB
  delivery: Gin

security:
  cookie_secure: true
  http_only: true
  rate_limit: 5.0
  content_security_policy_header:
    key: "Content-Security-Policy"
    value: "default-src 'self'"
  allowed_http_methods:
    - GET
    - POST
    - PUT
    - PATCH
    - DELETE
  allowed_content_types:
    - application/json
    - application/grpc

mongodb:
  name: golang_mongodb
  uri: mongodb://root:root@localhost:27017/golang_mongodb

gin:
  port: 8080
  allow_origins: http://localhost:8080
  allow_credentials: true
  server_group: /api

grpc:
  server_url: 0.0.0.0:8081
`)

	err := os.MkdirAll(yamlConfigPath, writePermissions)
	if err != nil {
		return ""
	}

	err = os.WriteFile(constants.YamlConfigPath, yamlContent, readPermissions)
	if err != nil {
		return ""
	}

	return constants.YamlConfigPath
}

func setupInvalidYamlFilePath() {
	invalidYAMLContent := []byte(`invalid_yaml: [unterminated`)

	err := os.MkdirAll(yamlConfigPath, writePermissions)
	fmt.Println(err)

	err = os.WriteFile(constants.YamlConfigPath, invalidYAMLContent, readPermissions)
	fmt.Println(err)
}

func setupEmptyYamlFilePath() string {
	invalidYAMLContent := []byte(``)

	err := os.MkdirAll(yamlConfigPath, writePermissions)
	if err != nil {
		return ""
	}
	err = os.WriteFile(constants.YamlConfigPath, invalidYAMLContent, readPermissions)
	if err != nil {
		return ""
	}

	return constants.YamlConfigPath
}

func setupEnvFilePath() string {
	envFilePath := constants.EnvironmentsPath + constants.Environment
	envContent := []byte(fmt.Sprintf(v1, constants.ConfigPath))

	err := os.MkdirAll(environmentConfigPath, writePermissions)
	if err != nil {
		return ""
	}
	err = os.WriteFile(envFilePath, envContent, readPermissions)
	if err != nil {
		return ""
	}

	return envFilePath
}

func cleanupTestEnvironment() {
	os.RemoveAll("config")
}

func TestViperLoadYamlConfiguration(t *testing.T) {
	setupYamlFilePath()
	setupEnvFilePath()
	defer cleanupTestEnvironment()

	viper := config.NewViper()
	assert.NotNil(t, viper, test.EqualMessage)
	assert.Equal(t, expectedLogger, viper.Core.Logger, test.EqualMessage)
	assert.Equal(t, expectedPort, viper.Gin.Port, test.EqualMessage)
	assert.Equal(t, expectedMongoDBDatabaseURI, viper.MongoDB.URI, test.EqualMessage)
	assert.Equal(t, expectedAllowedOrigins, viper.Gin.AllowOrigins, test.EqualMessage)
	assert.Contains(t, viper.Security.AllowedHTTPMethods, expectedAllowedMethod, test.EqualMessage)
}

func TestViperWithoutEnvironment(t *testing.T) {
	notification := fmt.Sprintf(openFileError, constants.EnvironmentsPath+constants.Environment)
	expectedError := domain.NewInternalError(expectedLocation+"loadDefaultEnvironment", notification)

	defer func() {
		recover := recover()
		if recover != nil {
			assert.Equal(t, recover, expectedError, test.EqualMessage)
		}
	}()
	config.NewViper()
}

func TestViperLoadEnvironmentWithoutYamlConfig(t *testing.T) {
	setupEnvFilePath()
	defer cleanupTestEnvironment()

	notification := fmt.Sprintf(openFileError, constants.ConfigPath)
	expectedError := domain.NewInternalError(expectedLocation+"loadDefaultConfig", notification)
	defer func() {
		recover := recover()
		if recover != nil {
			assert.Equal(t, recover, expectedError, test.EqualMessage)
		}
	}()
	config.NewViper()
}

func TestViperUnmarshalInvalidYAML(t *testing.T) {
	setupInvalidYamlFilePath()
	setupEnvFilePath()
	defer cleanupTestEnvironment()

	expectedError := domain.NewInternalError(expectedLocation+"loadDefaultConfig", yamlParsingError)
	defer func() {
		recover := recover()
		if recover != nil {
			assert.Equal(t, recover, expectedError, test.EqualMessage)
		}
	}()
	config.NewViper()
}

func TestViperUnmarshalEmptyAML(t *testing.T) {
	setupEmptyYamlFilePath()
	setupEnvFilePath()
	defer cleanupTestEnvironment()

	expectedError := domain.NewInternalError(expectedLocation+"loadDefaultConfig", yamlParsingError)
	defer func() {
		recover := recover()
		if recover != nil {
			assert.Equal(t, recover, expectedError, test.EqualMessage)
		}
	}()
	config.NewViper()
}
