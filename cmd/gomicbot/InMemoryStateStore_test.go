package main

import (
	"testing"
	"time"
)

func makeTestStateStoreConfig() *Configuration {
	return new(Configuration)
}

const testStateStore_User1 = "steve"
const testStateStore_User2 = "bob"

func failOnError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func assertEqual(expected interface{}, actual interface{}, message string, t *testing.T) {
	if expected != actual {
		t.Logf("FAIL %s - Expected '%v', got '%v'", message, expected, actual)
		t.FailNow()
	}
}

func assertTrue(condition bool, message string, t *testing.T) {
	if !condition {
		t.Logf("FAIL %s - condition not satisfied", message)
		t.FailNow()
	}
}

func assertFalse(condition bool, message string, t *testing.T) {
	if condition {
		t.Logf("FAIL %s - condition unexpectedly satisfied", message)
		t.FailNow()
	}
}

func TestUpdateLastSeen(t *testing.T) {
	s := new(InMemoryStateStore)

	err := s.Initialize(makeTestStateStoreConfig())
	failOnError(t, err)

	time1 := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	time2 := time.Date(2009, time.November, 11, 23, 0, 0, 0, time.UTC)

	lastSeen, err := s.UpdateLastSeen(testStateStore_User1, time1)
	failOnError(t, err)

	assertTrue(lastSeen.IsZero(), "Last seen should be zero timeval on first sight", t)

	lastSeen, err = s.UpdateLastSeen(testStateStore_User1, time2)
	failOnError(t, err)
	assertEqual(time1, lastSeen, "Last seen doesn't match", t)

	lastSeen, err = s.UpdateLastSeen(testStateStore_User2, time2)
	failOnError(t, err)
	assertTrue(lastSeen.IsZero(), "Last seen should be zero timeval on first sight", t)

	lastSeen, err = s.UpdateLastSeen(testStateStore_User2, time1)
	failOnError(t, err)
	assertEqual(time2, lastSeen, "Last seen doesn't match", t)
}

func TestSayingsStorage(t *testing.T) {
	s := new(InMemoryStateStore)

	err := s.Initialize(makeTestStateStoreConfig())
	failOnError(t, err)

	sayings, err := s.LoadSayings()
	failOnError(t, err)

	assertEqual(0, len(sayings), "Sayings data should be empty to start", t)

	saying1 := "The quick brown fox jumped over the lazy dog."
	saying2 := "End of line."
	saying3 := "Ain't here."

	s.StoreSaying(saying1)

	sayings, err = s.LoadSayings()
	failOnError(t, err)

	assertEqual(1, len(sayings), "After first addition", t)
	assertTrue(sayings[0] == saying1, "saying in list doesn't match", t)

	s.StoreSaying(saying2)

	sayings, err = s.LoadSayings()
	failOnError(t, err)

	assertEqual(2, len(sayings), "After second addition", t)
	assertTrue((sayings[0] == saying1 && sayings[1] == saying2) || (sayings[0] == saying2 && sayings[1] == saying1), "sayings in list don't match", t)

	present, err := s.RemoveSaying(saying3)
	failOnError(t, err)
	assertFalse(present, "removing nonexistent saying", t)

	sayings, err = s.LoadSayings()
	failOnError(t, err)

	assertEqual(2, len(sayings), "After removal of nonexistent saying", t)
	assertTrue((sayings[0] == saying1 && sayings[1] == saying2) || (sayings[0] == saying2 && sayings[1] == saying1), "sayings in list don't match", t)

	present, err = s.RemoveSaying(saying1)
	failOnError(t, err)
	assertTrue(present, "removing existing saying", t)

	sayings, err = s.LoadSayings()
	failOnError(t, err)

	assertEqual(1, len(sayings), "After first removal", t)
	assertTrue(sayings[0] == saying2, "saying in list doesn't match", t)

}
