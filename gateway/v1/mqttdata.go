package v1

const (
	EXTENSION_ID              = "id"
	EXTENSION_TOPIC           = "topic"
	EXTENSION_QOS             = "qos"
	EXTENSION_TARGET_ENTITYID = "entity_id"
)

type PublishData struct {
	TopicName        string `json:"topic_name"`
	Qos              int    `json:"qos"`
	Payload          []byte `json:"payload"`
	TargetEntityID   string `json:"target_entity_id,omitempty"`
	PacketIdentifier uint16 `json:"packet_identifier,omitempty"`
}

type UnSubData struct {
	TopicNames []string `json:"topic_names"`
}

type Topic struct {
	TopicName string `json:"topic_name"`
	Qos       int    `json:"qos"`
}

type SubData struct {
	TopicFilters []Topic `json:"topic_filters"`
}
