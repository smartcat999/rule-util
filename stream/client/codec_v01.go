package client

import (
	"encoding/json"
	"github.com/tkeel-io/rule-util/stream/types"
)

type CodecV02 struct {
}

var _ Codec = (*CodecV02)(nil)

func (v CodecV02) Encode(e types.Message) ([]byte, error) {
	return v.encodeStructured(e)
}

func (v CodecV02) Decode(data []byte) (types.Message, error) {
	// only structured is supported as of v0.2
	return v.decodeStructured(data)
}

func (v CodecV02) encodeStructured(e types.Message) ([]byte, error) {
	//	body, err := codec.JsonEncodeV02(e)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (v CodecV02) decodeStructured(body []byte) (types.Message, error) {
	event := types.NewMessage()
	err := json.Unmarshal(body, event)
	return event, err
}
