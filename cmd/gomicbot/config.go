package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Configuration struct {
	token        string
	redisUrl     string
	sayingChance float64
}

func loadConfiguration() (*Configuration, error) {
	p := new(Configuration)
	p.token = os.Getenv("GOBOT_TOKEN")
	p.redisUrl = os.Getenv("REDIS_URL")

	var err error
	p.sayingChance, err = strconv.ParseFloat(os.Getenv("SAYING_CHANCE"), 64)
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
