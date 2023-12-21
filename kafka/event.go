package kafka

import (
	"fmt"
	"time"

	common_utils "github.com/dispenal/go-common/utils"
	"github.com/google/uuid"
)

// EventType is the type of any event, used as its unique identifier.
type EventType string

// Event is an internal representation of an event, returned when the Aggregate
// uses NewEvent to create a new event. The events loaded from the db is
// represented by each DBs internal event type, implementing Event.
type Event struct {
	EventID   string
	EventType EventType
	Version   uint64
	Data      []byte
	Metadata  []byte
	Timestamp time.Time
}

// NewEvent creates a new event, with the given aggregateID, eventType and data.
// The eventID is generated automatically, and the version is set to 0.
func NewEvent(eventType EventType, data []byte, metadata ...[]byte) *Event {
	if len(metadata) == 0 {
		metadata = append(metadata, []byte("{}"))
	}
	return &Event{
		EventID:   uuid.New().String(),
		EventType: eventType,
		Version:   0,
		Data:      data,
		Metadata:  metadata[0],
		Timestamp: time.Now(),
	}
}

func (e *Event) GetEventID() string {
	return e.EventID
}

func (e *Event) GetEventType() EventType {
	return e.EventType
}

func (e *Event) GetVersion() uint64 {
	return e.Version
}

func (e *Event) GetData() []byte {
	return e.Data
}

func (e *Event) GetMetadata() []byte {
	return e.Metadata
}

func (e *Event) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetJsonData json unmarshal data attached to the Event.
func (e *Event) GetJsonData(data any) error {
	return common_utils.Unmarshal(e.GetData(), data)
}

// GetJsonMetadata unmarshal app-specific metadata serialized as json for the Event.
func (e *Event) GetJsonMetadata(metaData any) error {
	return common_utils.Unmarshal(e.GetMetadata(), metaData)
}

func (e *Event) SetMetadata(metaData any) error {

	metaDataBytes, err := common_utils.Marshal(metaData)
	if err != nil {
		return err
	}

	e.Metadata = metaDataBytes
	return nil
}

func (e *Event) String() string {
	return fmt.Sprintf("(Event) EventID: %s, Version: %d, EventType: %s, Metadata: %s, TimeStamp: %s",
		e.EventID,
		e.Version,
		e.EventType,
		string(e.Metadata),
		e.Timestamp.UTC().String(),
	)
}
