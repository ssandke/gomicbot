package main

import (
	"fmt"
	"strings"
)

const BotName = "gomicbot"

func __internalCheckCommandMatch(message string, command string, argsAllowed bool) (matched bool, args string) {

	matched = false
	args = ""

	if argsAllowed {
		if strings.HasPrefix(message, command) {
			matched = true
			args = strings.TrimSpace(message[len(command):])
		}
	} else if message == command {
		matched = true
	}
	return
}

func checkCommandMatch(message string, command string, argsAllowed bool, config *Configuration) (matched bool, args string) {
	directCommand := fmt.Sprintf("%s@%s", command, BotName)

	matched, args = __internalCheckCommandMatch(message, directCommand, argsAllowed)
	if matched {
		return
	}

	if !config.directCommandsOnly {
		matched, args = __internalCheckCommandMatch(message, command, argsAllowed)
		if matched {
			return
		}
	}

	args = ""
	return
}
