package event

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

type EventType string

const (
	EVENT_CREATE EventType = "Create"
	EVENT_UPDATE EventType = "Update"
	EVENT_DELETE EventType = "Delete"
)

type Event struct {
	ID        string          `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	EventType EventType       `json:"eventType"`
	Payload   json.RawMessage `json:"payload"`
}

func NewEvent(eventType EventType, payload interface{}) (*Event, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("event creation failed for payload:%v", payload)
		return nil, err
	}

	return &Event{
		ID:        uuid.New().String(),
		Timestamp: time.Now().UTC(),
		EventType: eventType,
		Payload:   payloadBytes,
	}, nil
}

func (e *Event) String() string {
	body, err := json.Marshal(e)
	if err != nil {
		log.Warnf("event marshall error:%v", err)
	}
	return string(body)
}
