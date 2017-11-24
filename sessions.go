package spice

import (
	"fmt"
	"sync"
)

// sessionTable holds a mapping of SessionID and destination node address
// map[sessionid]address
type sessionTable struct {
	lock    sync.Mutex
	entries map[uint32]*sessionEntry
}

type sessionEntry struct {
	address string
	count   int
	otp     string
}

func (s *sessionTable) Lookup(session uint32) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.entries[session]
	return ok
}

func (s *sessionTable) OTP(session uint32) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.entries[session]; ok {
		return s.entries[session].otp
	}
	return ""
}

func (s *sessionTable) Add(session uint32, destination string, otp string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.entries[session]; !ok {
		s.entries[session] = &sessionEntry{address: destination, count: 1, otp: otp}
	}
	return
}

func (s *sessionTable) Connect(session uint32) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.entries[session]; !ok {
		return "", fmt.Errorf("no such session in table")
	}
	s.entries[session].count = s.entries[session].count + 1
	return s.entries[session].address, nil
}

func (s *sessionTable) Disconnect(session uint32) {
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
