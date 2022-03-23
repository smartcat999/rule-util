package checkpoint

import (
	"encoding/json"
	"github.com/tkeel-io/rule-util/stream/types"
)

var (
	CodecPbMessage codecPbMessage
)

type codecPbMessage struct {
}

func (v codecPbMessage) Encode(message Message) ([]byte, error) {
	return v.encodeStructured(message)
}

func (v codecPbMessage) Decode(body []byte) (Message, error) {
	// only structured is supported as of v0.2
	return v.decodeStructured(body)
}

func (v codecPbMessage) encodeStructured(message Message) ([]byte, error) {
	//	body, err := codec.JsonEncodeV02(e)
	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (v codecPbMessage) decodeStructured(body []byte) (Message, error) {
	event := types.NewMessage()
	err := json.Unmarshal(body, event)
	return event, err
}
