package main

import "testing"

func testHelper(t *testing.T, message string, command string, directOnly bool, isPersonal bool, argsAllowed bool, expectMatch bool, expectArgs string) {
	config := new(Configuration)
	config.directCommandsOnly = directOnly

	matched, args := checkCommandMatch(message, command, isPersonal, argsAllowed, config)

	if matched != expectMatch {
		t.Errorf("Matched %v != Expected Match %v", matched, expectMatch)
	}

	if args != expectArgs {
		t.Errorf("Returned Args '%v' != Expected Returned Args '%v'", args, expectArgs)
	}
}

func TestDirectOnlyBareCommandPresent(t *testing.T) {
	testHelper(t, "/ho", "/ho", true, false, false, false, "")
}

func TestDirectOnlyBareCommandPresentPersonalMessage(t *testing.T) {
	testHelper(t, "/ho", "/ho", true, true, false, true, "")
}

func TestBareCommandPresent(t *testing.T) {
	testHelper(t, "/ho", "/ho", false, false, false, true, "")
}

func TestBareCommandPresentWithArgs(t *testing.T) {
	testHelper(t, "/nose knows", "/nose", false, false, true, true, "knows")
}

func TestBareCommandPresentWithArgsAndSpace(t *testing.T) {
	testHelper(t, "/nose   knows yo   ", "/nose", false, false, true, true, "knows yo")
}

func TestBareCommandNotPresent(t *testing.T) {
	testHelper(t, "/yo", "/ho", false, false, false, false, "")
}
