package main

import (
	"fmt"
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
	expectedReaders        []int
	expectedWriters        []int
}

const validToken = "a:b"
const validReaders = "1,2,3"
const validWriters = "1,2"
const validRedisUrl = "r:foo" // for now
const validSayingChance = "0.3"
const sayingChanceVal = 0.3
const envTrue = "true"
const envFalse = "false"

var testSpecs = []configTestSpec{
	{validToken, validReaders, validWriters, validRedisUrl, validSayingChance, envTrue, "", true, sayingChanceVal, []int{1, 2, 3}, []int{1, 2}},
	{validToken, validReaders, "", validRedisUrl, validSayingChance, envFalse, "", false, sayingChanceVal, []int{1, 2, 3}, []int{}},
	{validToken, "", validWriters, validRedisUrl, validSayingChance, "", "", true, sayingChanceVal, []int{}, []int{1, 2}},
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
		assertEqual(len(test.expectedReaders), len(c.readers), "reader count", t)
		for _, id := range test.expectedReaders {
			assertTrue(c.isReader(id), fmt.Sprintf("reader %d", id), t)
		}

		assertEqual(len(test.expectedWriters), len(c.writers), "writer count", t)
		for _, id := range test.expectedWriters {
			assertTrue(c.isWriter(id), fmt.Sprintf("writer %d", id), t)
		}
	}
}

func TestReadersWriterParsing(t *testing.T) {

	idMap, err := listToMap("")
	failOnError(t, err)
	assertEqual(0, len(idMap), "empty string should be empty map", t)

	idMap, err = listToMap("1,2,3")
	failOnError(t, err)
	assertEqual(3, len(idMap), "wrong element count", t)
	assertTrue(idMap[1], "1", t)
	assertTrue(idMap[2], "2", t)
	assertTrue(idMap[3], "3", t)

	idMap, err = listToMap(" 2  ,2   , 3 ")
	failOnError(t, err)
	assertEqual(2, len(idMap), "wrong element count", t)
	assertTrue(idMap[2], "2", t)
	assertTrue(idMap[3], "3", t)
}
