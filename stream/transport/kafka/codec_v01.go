package kafka

import (
	"encoding/json"
	"github.com/tkeel-io/rule-util/stream/types"
)

type CodecV02 struct {
}

var _ types.Codec = (*CodecV02)(nil)

func (v CodecV02) Encode(e interface{}) ([]byte, error) {
	return v.encodeStructured(e)
}

func (v CodecV02) encodeStructured(e interface{}) ([]byte, error) {
	//	body, err := codec.JsonEncodeV02(e)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//func (v CodecV02) decodeStructured(body []byte) (types.Message, error) {
//	event := types.NewMessage()
//	err := json.Unmarshal(body, event)
//	return event, err
//}
