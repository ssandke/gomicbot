package main

import (
	"sync"
	"time"
)

type InMemoryStateStore struct {
	lock     sync.Mutex
	lastSeen map[string]time.Time
	sayings  map[string]struct{}
}

func (s *InMemoryStateStore) Initialize(config *Configuration) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	s.lastSeen = make(map[string]time.Time)
	s.sayings = make(map[string]struct{})

	return nil
}

func (*InMemoryStateStore) Save() error {
	return nil
}

func (*InMemoryStateStore) LoadSayings() ([]string, error) {
	return nil, nil
}

func (*InMemoryStateStore) StoreSaying(saying string) error {
	return nil
}

func (*InMemoryStateStore) RemoveSaying(saying string) (present bool, err error) {
	return false, nil
}

func (s *InMemoryStateStore) UpdateLastSeen(user string, seen time.Time) (lastseen time.Time, err error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	lastSeen := s.lastSeen[user]
	s.lastSeen[user] = seen

	return lastSeen, nil
}
