package v1

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"time"
)

type Ticket struct {
	Id          uuid.UUID `json:"id"`
	TimeEntry   time.Time `json:"timeEntry"`
	TimeExit    time.Time `json:"timeExit"`
	Vehicle     string    `json:"vehicle"`
	ParkingSlot string    `json:"parkingSlot"`
}

func (t *Ticket) JSON() []byte {
	thisJson, err := json.Marshal(*t)
	if err != nil {
		log.Fatal(err)
	}
	return thisJson
}
