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

// fetch an argument as a string
func (msg Message) Get(arg string) string {
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
func (msg *Message) UnmarshalJSON(buf []byte) error {
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

// allows marshaling of messages
func (msg Message) MarshalJSON() ([]byte, error) {
	parts := [...]interface{}{msg.Command, msg.Args, msg.ID}
	return json.Marshal(parts)
}
