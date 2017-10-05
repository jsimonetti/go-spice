package spice

import (
	"fmt"
	"sync"
)

type SessionID [4]uint8

// sessionTable holds a mapping of SessionID and destination node address
// map[sessionid]address
type sessionTable struct {
	lock    sync.Mutex
	entries map[SessionID]*sessionEntry
}

type sessionEntry struct {
	address string
	count   int
}

func (s *sessionTable) Lookup(session SessionID) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.entries[session]
	return ok
}

func (s *sessionTable) Add(session SessionID, destination string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.entries[session]; !ok {
		s.entries[session] = &sessionEntry{address: destination, count: 1}
	}
	return
}

func (s *sessionTable) Connect(session SessionID) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.entries[session]; !ok {
		return "", fmt.Errorf("no such session in table")
	}
	s.entries[session].count = s.entries[session].count + 1
	return s.entries[session].address, nil
}

func (s *sessionTable) Disconnect(session SessionID) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.entries[session]; !ok {
		return
	}
	s.entries[session].count = s.entries[session].count - 1
	if s.entries[session].count < 1 {
		delete(s.entries, session)
	}
	return
}
