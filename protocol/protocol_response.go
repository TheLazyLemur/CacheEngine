package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
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

func ParseDelReponse(r io.Reader) (*ResponseDelete, error) {
	resp := &ResponseDelete{}
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

