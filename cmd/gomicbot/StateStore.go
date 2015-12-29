package main

import "time"

type StateStore interface {
	Initialize(config *Configuration) error

	Save() error

	LoadSayings() ([]string, error)
	StoreSaying(saying string) error
	RemoveSaying(saying string) (present bool, err error)

	UpdateLastSeen(user string, seen time.Time) (lastseen time.Time, err error)
}
