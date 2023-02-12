package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)


type Command byte

const(
	CMDNone Command = iota
	CmdSet
	CmdGet
	CmdDel
)

type Status byte
func (s Status) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusError:
		return "ERROR"
	case StatusKeyNotFound:
		return "KEY_NOT_FOUND"
	default:
		return "UNKNOWN"
	}
}

const (
	StatusNone Status = iota
	StatusOK
	StatusError
	StatusKeyNotFound
)

type ResponseDelete struct {
	Status
}

func (r *ResponseDelete) Bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, r.Status)
	
	return buf.Bytes()
}

type ResponseSet struct {
	Status Status
}

func (r *ResponseSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, r.Status)
	
	return buf.Bytes()
}


type ResponseGet struct {
	Status
	Value []byte
}

func (r *ResponseGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, r.Status)
	
	valueLen := int32(len(r.Value))
	_ = binary.Write(buf, binary.LittleEndian, valueLen)
	_ = binary.Write(buf, binary.LittleEndian, r.Value)

	return buf.Bytes()
}

func ParseSetReponse(r io.Reader) (*ResponseSet, error) {
	resp := &ResponseSet{}
	err := binary.Read(r, binary.LittleEndian, &resp.Status)
	return resp, err
}

func ParseGetReponse(r io.Reader) (*ResponseGet, error) {
	resp := &ResponseGet{}
	_ = binary.Read(r, binary.LittleEndian, &resp.Status)

	var valueLen int32
	_ = binary.Read(r, binary.LittleEndian, &valueLen)

	resp.Value = make([]byte, valueLen)
	_ = binary.Read(r, binary.LittleEndian, &resp.Value)

	return resp, nil
}

type CommandSet struct {
	Key []byte
	Value []byte
	TTL int
}

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, CmdSet)

	keyLen := int32(len(c.Key))
	_ = binary.Write(buf, binary.LittleEndian, keyLen)
	_ = binary.Write(buf, binary.LittleEndian, c.Key)

	valLen := int32(len(c.Value))
	_ = binary.Write(buf, binary.LittleEndian, valLen)
	_ = binary.Write(buf, binary.LittleEndian, c.Value)

	_ = binary.Write(buf, binary.LittleEndian, int32(c.TTL))

	return buf.Bytes()
}

type CommandGet struct {
	Key []byte
}

func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, CmdGet)

	keyLen := int32(len(c.Key))
	_ = binary.Write(buf, binary.LittleEndian, keyLen)
	_ = binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()
}

type CommandDel struct {
	Key []byte
}

func (c *CommandDel) Bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, CmdDel)

	keyLen := int32(len(c.Key))
	_ = binary.Write(buf, binary.LittleEndian, keyLen)
	_ = binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()
}

func ParseCommand(r io.Reader) (any, error) {
	var cmd Command
	if err := binary.Read(r, binary.LittleEndian, &cmd); err != nil {
		return nil, err
	}

	switch cmd {
	case CmdSet:	
		return parseSetCommand(r), nil
	case CmdGet:	
		return parseGetCommand(r), nil
	case CmdDel:	
		return parseDelCommand(r), nil
	default:
		return nil, fmt.Errorf("unknown command: %d", cmd)
	}
}

func parseSetCommand(r io.Reader) *CommandSet {
	cmd := &CommandSet{}

	var keyLen int32
	_ = binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	_ = binary.Read(r, binary.LittleEndian, &cmd.Key)

	var valueLen int32
	_ = binary.Read(r, binary.LittleEndian, &valueLen)
	cmd.Value = make([]byte, valueLen)
	_ = binary.Read(r, binary.LittleEndian, &cmd.Value)

	var ttl int32
	_ = binary.Read(r, binary.LittleEndian, &ttl)
	cmd.TTL = int(ttl)

	return cmd
}

func parseGetCommand(r io.Reader) *CommandGet {
	cmd := &CommandGet{}

	var keyLen int32
	_ = binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	_ = binary.Read(r, binary.LittleEndian, &cmd.Key)

	return cmd
}

func parseDelCommand(r io.Reader) *CommandDel {
	cmd := &CommandDel{}

	var keyLen int32
	_ = binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	_ = binary.Read(r, binary.LittleEndian, &cmd.Key)

	return cmd
}
