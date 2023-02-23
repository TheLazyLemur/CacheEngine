package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSetResponseStatusNone(t *testing.T) {
	resp := &ResponseSet{
		Status: StatusNone,
	}

	expeced := []byte{0}
	actual := resp.Bytes()
	assert.Equal(t, expeced, actual)
}

func TestParseSetResponseStatusOk(t *testing.T) {
	resp := &ResponseSet{
		Status: StatusOK,
	}

	expeced := []byte{1}
	actual := resp.Bytes()
	assert.Equal(t, expeced, actual)
}

func TestParseSetResponseStatusErr(t *testing.T) {
	resp := &ResponseSet{
		Status: StatusError,
	}

	expeced := []byte{2}
	actual := resp.Bytes()
	assert.Equal(t, expeced, actual)
}

func TestParseSetResponseStatusKeyNotFound(t *testing.T) {
	resp := &ResponseSet{
		Status: StatusKeyNotFound,
	}

	expeced := []byte{3}
	actual := resp.Bytes()
	assert.Equal(t, expeced, actual)
}
