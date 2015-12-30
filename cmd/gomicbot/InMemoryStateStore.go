package main

import (
	"sync"

	"time"
)

type InMemoryStateStore struct {
	lock     sync.Mutex
	lastSeen map[string]time.Time
	sayings  map[string]bool
}

func (s *InMemoryStateStore) Initialize(config *Configuration) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	s.lastSeen = make(map[string]time.Time)
	s.sayings = make(map[string]bool)

	return nil
}

func (*InMemoryStateStore) Save() error {
	return nil
}

func (s *InMemoryStateStore) LoadSayings() ([]string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	sayings := make([]string, 0, len(s.sayings))
	for saying := range s.sayings {
		sayings = append(sayings, saying)
	}

	return sayings, nil
}

func (s *InMemoryStateStore) StoreSaying(saying string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.sayings[saying] = true

	return nil
}

func (s *InMemoryStateStore) RemoveSaying(saying string) (present bool, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	present = s.sayings[saying]
	delete(s.sayings, saying)

	err = nil
	return
}

func (s *InMemoryStateStore) UpdateLastSeen(user string, seen time.Time) (lastseen time.Time, err error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	lastSeen := s.lastSeen[user]
	s.lastSeen[user] = seen

	return lastSeen, nil
}
