package checkpoint

import (
	"fmt"
	"github.com/tkeel-io/rule-util/stream/types"
	"time"
)

// Event stream event, abstracts different types of events
// - Message
// - Checkpoint
type Event interface {
	string() string
}

// Message represents a message
type PacketMessage struct {
	Timestamp time.Time
	Message   Message
}

func (m *PacketMessage) string() string {
	return fmt.Sprintf("Message %s:%s:%s", m.Message.PacketIdentifier(), m.Message.Domain(), m.Message.(PublishMessage).Topic())
}

// Checkpoint
type Checkpoint struct {
	CheckpointID types.CheckpointID
	Timestamp    time.Time
}

func (n *Checkpoint) string() string {
	return "checkpoint"
}
