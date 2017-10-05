package spice

type SessionID [4]uint8

// sessionTable holds a mapping of SessionID and destination node address
// map[sessionid]address
type sessionTable map[SessionID]*sessionEntry

type sessionEntry struct {
	address string
	count   int
}
