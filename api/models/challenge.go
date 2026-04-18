package models

import "time"

type Challenge struct {
	Id          uint64    `json:"id"`
	X           int       `json:"x"`
	Y           int       `json:"y"`
	MapID       uint64    `json:"map_id"`
	Map         Map       `json:"-"`
	GeneratedAt time.Time `json:"generated_at"`
}

func init() {
	Models = append(Models, Challenge{})
}
