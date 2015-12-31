package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	token              string
	redisUrl           string
	sayingChance       float64
	directCommandsOnly bool
	readers            map[int]bool
	writers            map[int]bool
}

const TokenEnvName = "GOBOT_TOKEN"
const RedisUrlEnvName = "REDIS_URL"
const SayingChanceEnvName = "SAYING_CHANCE"
const DirectCommandsOnlyEnvName = "DIRECT_COMMANDS_ONLY"
const ReadersEnvName = "READERS"
const WritersEnvName = "WRITERS"

const defaultDirectCommandOnly = "true"
const defaultSayingChance = "0.8"

func loadConfiguration() (*Configuration, error) {
	p := new(Configuration)
	p.token = os.Getenv(TokenEnvName)
	p.redisUrl = os.Getenv(RedisUrlEnvName)

	var err error
	p.readers, err = listToMap(os.Getenv(ReadersEnvName))
	if err != nil {
		return nil, err
	}

	p.writers, err = listToMap(os.Getenv(WritersEnvName))
	if err != nil {
		return nil, err
	}

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

func (c *Configuration) isReader(userid int) bool {
	return c.readers[userid]
}

func (c *Configuration) isWriter(userid int) bool {
	return c.writers[userid]
}

func listToMap(idList string) (result map[int]bool, err error) {
	result = make(map[int]bool)

	for _, id := range strings.Split(idList, ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			idNum, err := strconv.Atoi(id)
			if err != nil {
				break
			}
			result[idNum] = true
		}
	}
	return
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
	redisStr := c.redisUrl
	if redisStr != "" {
		redisStr = "<redacted>"
	}

	return fmt.Sprintf("[token: %s\n redis url: %s\n saying chance: %f\ndirect commands only: %v\nreaders: %v\nwriters: %v]",
		tokenStr, redisStr, c.sayingChance, c.directCommandsOnly, c.readers, c.writers)
}

func getEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	return value
}
