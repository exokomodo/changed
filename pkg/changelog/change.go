package changelog

import "time"

type Change struct {
	ID        uint      `json:"id"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Actor     string    `json:"actor"`
	Service   string    `json:"service"`
	Details   string    `json:"details"`
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
