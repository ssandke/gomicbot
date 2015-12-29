package main

import (
	"os"
	"strings"
	"testing"
)

const tokenVar = "GOBOT_TOKEN"
const readersVar = "READERS"
const writersVar = "WRITERS"
const redisUrlVar = "REDIS_URL"
const sayingChanceVar = "SAYING_CHANCE"

type configTestSpec struct {
	// inputs
	tokenStr       string
	readers        string
	writers        string
	redisUrl       string
	sayingChance   string
	directCommands string

	// expected results
	errorSubstring         string
	expectedDirectCommands bool
	expectedSayingChance   float64
}

const validToken = "a:b"
const validReaders = "1,2,3"
const validWriters = "1,2"
const validRedisUrl = "r:foo" // for now
const validSayingChance = "0.3"
const sayingChanceVal = 0.3
const trueString = "true"
const falseString = "false"

var testSpecs = []configTestSpec{
	{validToken, validReaders, validWriters, validRedisUrl, validSayingChance, trueString, "", true, sayingChanceVal},
	{validToken, validReaders, validWriters, validRedisUrl, validSayingChance, falseString, "", false, sayingChanceVal},
	{validToken, validReaders, validWriters, validRedisUrl, validSayingChance, "", "", true, sayingChanceVal},
}

func TestSayingChanceNonNumeric(t *testing.T) {

	os.Setenv(sayingChanceVar, "This ain't no number")

	c, err := loadConfiguration()
	if c != nil {
		t.Error("config return wasn't nil")
	}
	if err == nil {
		t.Error("error was nil")
	}
	if !strings.Contains(err.Error(), "invalid syntax") {
		t.Error("error message didn't contain 'invalid syntax'")
	}
}

func TestConfigs(t *testing.T) {

	for _, test := range testSpecs {
		os.Setenv(TokenEnvName, test.tokenStr)
		os.Setenv(readersVar, test.readers)
		os.Setenv(writersVar, test.writers)
		os.Setenv(SayingChanceEnvName, test.sayingChance)
		os.Setenv(DirectCommandsOnlyEnvName, test.directCommands)

		c, err := loadConfiguration()

		if test.errorSubstring != "" {
			if c != nil {
				t.Error("config return wasn't nil")
			}
			if err == nil {
				t.Error("error was nil")
			}

			if !strings.Contains(err.Error(), test.errorSubstring) {
				t.Errorf("error message didn't contain '%s'", test.errorSubstring)
			}

		}

		if c == nil {
			t.Error("config was nil")
		}
		if err != nil {
			t.Error("error wasn't nil")
		}
		if test.expectedDirectCommands != c.directCommandsOnly {
			t.Errorf("DirectCommands was %v, expected %v\n%v", c.directCommandsOnly, test.expectedDirectCommands, test)
		}
		if test.expectedSayingChance != c.sayingChance {
			t.Errorf("Saying Chance was %v, expected %v", c.sayingChance, test.expectedSayingChance)
		}
	}
}
