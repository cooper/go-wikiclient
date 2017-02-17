// Copyright (c) 2017, Mitchell Cooper
package wikiclient

// a Message represents a wikiserver message. the Message struct is used both
// for outgoing messages (client requests) and incoming messages (replies from
// the wikiserver).

import (
	"encoding/json"
	"errors"
	"fmt"
)

type messageArgs map[string]interface{}

var idCounter uint

type Message struct {
	Command string      // message type
	Args    messageArgs // message arguments
	ID      uint        // message ID
}

// creates a new Message with an automatically-generated ID
func NewMessage(cmd string, args messageArgs) Message {
	idCounter++
	return NewMessageWithID(cmd, args, idCounter)
}

// creates a new Message with the specified ID
func NewMessageWithID(cmd string, args messageArgs, id uint) Message {
	return Message{cmd, args, id}
}

// creates a new Message from JSON data
// this is ugly, but I don't think there's a nicer way to do it since we use
// a JSON array and not an object?
func MessageFromJson(data []byte) (msg Message, err error) {
	if err = json.Unmarshal(data, &msg); err != nil {
		return msg, err
	}
	return msg, nil
}

// fetch an arguments as a string
func (msg Message) String(arg string) string {
	iface, ok := msg.Args[arg]
	if !ok {
		return ""
	}
	switch val := iface.(type) {
	case nil:
		return ""
	case string:
		return val
	}
	return fmt.Sprintf("%v", iface)
}

// allows JSON unmarshal into a message
func (msg Message) UnmarshalJSON(buf []byte) error {
	parts := []interface{}{&msg.Command, &msg.Args, &msg.ID}
	need := len(parts)
	if err := json.Unmarshal(buf, &parts); err != nil {
		return err
	}
	if len(parts) != need {
		return errors.New("Message must be a JSON array of length 3")
	}
	return nil
}

// translates the Message to JSON
func (msg Message) ToJson() []byte {
	ary := [...]interface{}{msg.Command, msg.Args, msg.ID}
	json, _ := json.Marshal(ary)
	return json
}
