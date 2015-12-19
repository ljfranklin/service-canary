package adapters

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ljfranklin/service-canary/config"
	"github.com/pivotal-golang/lager"

	"gopkg.in/redis.v3"
)

type RedisAdapter struct {
	logger        lager.Logger
	name          string
	serviceConfig *config.ServiceConfig
	client        *redis.Client
}

type redisConfig struct {
	Credentials redisCredentialsConfig `json:"credentials"`
}

type redisCredentialsConfig struct {
	Hostname string `json:"hostname"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

func NewRedisAdapter(serviceConfig *config.ServiceConfig, logger lager.Logger) *RedisAdapter {
	return &RedisAdapter{
		logger:        logger,
		name:          serviceConfig.Name,
		serviceConfig: serviceConfig,
	}
}

func (a RedisAdapter) Name() string {
	return a.name
}

func (a *RedisAdapter) Setup() error {
	a.logger.Info("Setting up redis adapter...")

	var configProperties redisConfig
	if err := json.Unmarshal(a.serviceConfig.ConfigJSON, &configProperties); err != nil {
		return fmt.Errorf("Failed to parse config for %s: %s", a.Name(), err.Error())
	}

	creds := configProperties.Credentials
	a.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", creds.Hostname, creds.Port),
		Password: creds.Password,
		DB:       0, // use default DB
	})

	_, err := a.client.Ping().Result()
	if err != nil {
		return fmt.Errorf("Error setting up redis client: %s", err.Error())
	}

	return nil
}

func (a *RedisAdapter) Run() error {
	a.logger.Info("Running redis adapter...")

	key := "service_canary_test"
	val := time.Now().Unix()
	err := a.client.Set(key, string(val), 0).Err()
	if err != nil {
		return fmt.Errorf("Failed to set redis key: %s", err.Error())
	}

	returnedVal, err := a.client.Get(key).Result()
	if err == redis.Nil || err != nil {
		return fmt.Errorf("Could not find inserted key '%s': %s", key, err.Error())
	}

	if returnedVal != string(val) {
		return fmt.Errorf("Returned value '%s' did not match expected '%s'", returnedVal, val)
	}
	return nil
}
