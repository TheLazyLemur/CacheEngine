package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Command string

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
	CMDHas Command = "HAS"
	CMDDel Command = "DEL"
)

type Message struct {
	Cmd Command
	Key []byte
	Value []byte
	Ttl int64
}

func (m *Message) ToBytes() []byte {
	var cmd string
	switch m.Cmd {
	case CMDSet:
		cmd = fmt.Sprintf("%s %s %s %d", m.Cmd, m.Key, m.Value, m.Ttl)
	case CMDGet:
		cmd = fmt.Sprintf("%s %s", m.Cmd, m.Key)
	case CMDHas:
		cmd = fmt.Sprintf("%s %s", m.Cmd, m.Key)
	case CMDDel:
		cmd = fmt.Sprintf("%s %s", m.Cmd, m.Key)
	}

	return []byte(cmd)
}

func parseMessage(rawCmd []byte) (*Message, error) {
	parts := strings.Split(string(rawCmd), " ")
	if len(parts) == 0 {
		return nil, errors.New("invalid protocol format")
	}

	msg := &Message{
		Cmd: Command(parts[0]),
		Key: []byte(parts[1]),
	}

	if msg.Cmd == CMDSet {
		if len(parts) != 4 {
			return nil, errors.New("invalid SET command")
		}

		msg.Value = []byte(parts[2])

		ttl, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, errors.New("invalid ttl format")
		}

		msg.Ttl = int64(ttl)
	}

	if msg.Cmd == CMDGet {
		if len(parts) != 2 {
			return nil, errors.New("invalid GET command")
		}
	}

	if msg.Cmd == CMDHas {
		if len(parts) != 2 {
			return nil, errors.New("invalid HAS command")
		}
	}

	if msg.Cmd == CMDDel {
		if len(parts) != 2 {
			return nil, errors.New("invalid HAS command")
		}
	}

	return msg, nil
}
