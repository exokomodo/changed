package changelog

import "time"

type Change struct {
	ID        uint
	Timestamp time.Time
	Actor     string
	Service   string
	Details   string
}

func NewChange(actor, service, details string) Change {
	return Change{
		Timestamp: time.Now().UTC(),
		Actor:     actor,
		Service:   service,
		Details:   details,
	}
}

func NewChangeWithTimestamp(actor, service, details string, timestamp time.Time) Change {
	return Change{
		Timestamp: timestamp.UTC(),
		Actor:     actor,
		Service:   service,
		Details:   details,
	}
}
