package entity

import "time"

type Registrant struct {
	Id       string    `json:"id"`
	Address  string    `json:"address"`
	LastPing time.Time `json:"last_ping,omitempty"`
}
