package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSetCommand(t *testing.T) {
	cmd := &CommandSet{
		Key: []byte("Foo"),
		Value: []byte("Bar"),
		TTL: 2,
	}

	r := bytes.NewReader(cmd.Bytes())
	pcmd := ParseCommand(r)

	assert.Equal(t, cmd, pcmd)
}

func TestParseGetCommand(t *testing.T) {
	cmd := &CommandGet{
		Key: []byte("Foo"),
	}

	r := bytes.NewReader(cmd.Bytes())
	pcmd := ParseCommand(r)

	assert.Equal(t, cmd, pcmd)
}

func TestParseDelCommand(t *testing.T) {
	cmd := &CommandDel{
		Key: []byte("Foo"),
	}

	r := bytes.NewReader(cmd.Bytes())
	pcmd := ParseCommand(r)

	assert.Equal(t, cmd, pcmd)
}
