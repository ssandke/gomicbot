package main

import (
	"errors"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/ssandke/gomicbot/Godeps/_workspace/src/github.com/mediocregopher/radix.v2/redis"
)

type RedisStateStore struct {
	lock    sync.Mutex
	sayings map[string]bool
	client  *redis.Client
}

const sayingsRedisKey = "BasicSayings"

func (s *RedisStateStore) Initialize(config *Configuration) (err error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	client, err := s.connectUrl(config.redisUrl)
	if err != nil {
		return
	}
	s.client = client

	err = s.load()

	return
}

func (*RedisStateStore) Save() error {
	return nil
}

func (s *RedisStateStore) LoadSayings() ([]string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	sayings := make([]string, 0, len(s.sayings))
	for saying := range s.sayings {
		sayings = append(sayings, saying)
	}

	return sayings, nil
}

func (s *RedisStateStore) StoreSaying(saying string) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	resp := s.client.Cmd("SADD", sayingsRedisKey, saying)

	if err = resp.Err; err != nil {
		return
	}

	s.sayings[saying] = true

	return
}

func (s *RedisStateStore) RemoveSaying(saying string) (present bool, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	err = nil
	present = false

	resp := s.client.Cmd("SREM", sayingsRedisKey, saying)

	if err = resp.Err; err != nil {
		return
	}

	present = s.sayings[saying]
	delete(s.sayings, saying)

	return
}

func (s *RedisStateStore) UpdateLastSeen(user string, seen time.Time) (lastseen time.Time, err error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	// lastSeen := s.lastSeen[user]
	// s.lastSeen[user] = seen

	lastSeen := time.Now()
	return lastSeen, nil
}

func (s *RedisStateStore) load() (err error) {
	err = nil

	s.sayings = make(map[string]bool)

	resp := s.client.Cmd("SMEMBERS", sayingsRedisKey)

	err = resp.Err
	if err != nil {
		return
	}

	results, err := resp.List()
	if err != nil {
		return
	}

	for _, saying := range results {
		if saying != "" {
			s.sayings[saying] = true
		}
	}

	log.Printf("Loaded %d sayings from Redis.\n", len(s.sayings))

	return
}

func (s *RedisStateStore) connectUrl(redisUrl string) (client *redis.Client, err error) {
	client = nil
	err = nil

	var uri *url.URL

	// Parse the URL
	uri, err = url.Parse(redisUrl)
	if err != nil {
		return
	}

	if uri.Scheme != "redis" {
		err = errors.New("unexpected scheme in URL: " + redisUrl)
		return
	}

	// Open the connection
	log.Println("Connecting to Redis server at " + uri.Host)
	client, err = redis.Dial("tcp", uri.Host)
	if err != nil {
		return
	}

	// AUTH as needed
	if uri.User != nil {
		password, set := uri.User.Password()
		if set {
			log.Println("Authenticating redis connection.")
			result := client.Cmd("AUTH", password)
			if err = result.Err; err != nil {
				client.Close()
				client = nil
				return
			}
		}
	}

	return
}
