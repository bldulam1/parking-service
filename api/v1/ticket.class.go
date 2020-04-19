package v1

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"time"
)

type Ticket struct {
	Id          uuid.UUID
	TimeEntry   time.Time
	TimeExit    time.Time
	Vehicle     string `json:"vehicle"`
	ParkingSlot string `json:"parkingSlot"`
}

func (t *Ticket) JSON() []byte {
	thisJson, err := json.Marshal(*t)
	if err != nil {
		log.Fatal(err)
	}
	return thisJson
}
