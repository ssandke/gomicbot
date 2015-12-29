package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Configuration struct {
	token              string
	redisUrl           string
	sayingChance       float64
	directCommandsOnly bool
}

const TokenEnvName = "GOBOT_TOKEN"
const RedisUrlEnvName = "REDIS_URL"
const SayingChanceEnvName = "SAYING_CHANCE"
const DirectCommandsOnlyEnvName = "DIRECT_COMMANDS_ONLY"

const defaultDirectCommandOnly = "true"
const defaultSayingChance = "0.8"

func loadConfiguration() (*Configuration, error) {
	p := new(Configuration)
	p.token = os.Getenv(TokenEnvName)
	p.redisUrl = os.Getenv(RedisUrlEnvName)

	var err error
	p.sayingChance, err = strconv.ParseFloat(getEnvWithDefault(SayingChanceEnvName, defaultSayingChance), 64)
	if err != nil {
		return nil, err
	}

	p.directCommandsOnly, err = strconv.ParseBool(getEnvWithDefault(DirectCommandsOnlyEnvName, defaultDirectCommandOnly))
	if err != nil {
		return nil, err
	}

	return p, p.validate()
}

func (c *Configuration) validate() error {
	if c.sayingChance < 0.0 || c.sayingChance > 1.0 {
		return errors.New(fmt.Sprintf("Saying chance (%f) must be between 0.0 and 1.0 inclusive", c.sayingChance))
	}
	return nil
}

func (c *Configuration) String() string {
	tokenStr := c.token
	if tokenStr != "" {
		tokenStr = "<redacted>"
	}

	return fmt.Sprintf("[token: %s, redis url: %s, saying chance: %f]", tokenStr, c.redisUrl, c.sayingChance)
}

func getEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	return value
}
