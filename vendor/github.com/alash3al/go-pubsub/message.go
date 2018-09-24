package pubsub

// Message the message metadata
type Message struct {
	topic     string
	payload   interface{}
	createdAt int64
}

// GetTopic return the topic of the current message
func (m *Message) GetTopic() string {
	return m.topic
}

// GetPayload get the payload of the current message
func (m *Message) GetPayload() interface{} {
	return m.payload
}

// GetCreatedAt get the creation time of this message
func (m *Message) GetCreatedAt() int64 {
	return m.createdAt
}
