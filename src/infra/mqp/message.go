package mqp

import jsoniter "github.com/json-iterator/go"

type Message struct {
	Body []byte
}

func NewMessage(body []byte) *Message {
	return &Message{Body: body}
}

func (sel *Message) MarshalToBody(val interface{}) (err error) {
	sel.Body, err = jsoniter.Marshal(val)
	return
}

func (sel *Message) UnMarshalFromBody(val interface{}) error {
	return jsoniter.Unmarshal(sel.Body, val)
}

func (sel *Message) Copy() *Message {
	return &Message{
		Body: sel.Body,
	}
}
