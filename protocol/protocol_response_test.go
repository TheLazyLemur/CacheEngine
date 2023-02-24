package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSetResponse(t *testing.T) {
	resp := &ResponseSet{
		Status: StatusNone,
	}
	reader := bytes.NewReader(resp.Bytes())
	actual, err := ParseSetReponse(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp, actual)

	resp = &ResponseSet{
		Status: StatusOK,
	}
	reader = bytes.NewReader(resp.Bytes())
	actual, err = ParseSetReponse(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp, actual)

	resp = &ResponseSet{
		Status: StatusError,
	}
	reader = bytes.NewReader(resp.Bytes())
	actual, err = ParseSetReponse(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp, actual)

	resp = &ResponseSet{
		Status: StatusKeyNotFound,
	}
	reader = bytes.NewReader(resp.Bytes())
	actual, err = ParseSetReponse(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp, actual)
}

func TestParseGetResponseStatusOk(t *testing.T) {
	resp := &ResponseGet{
		Status: StatusOK,
		Value:  []byte("value"),
	}
	reader := bytes.NewReader(resp.Bytes())
	x, err := ParseGetReponse(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp, x)

	resp = &ResponseGet{
		Status: StatusError,
		Value:  []byte(""),
	}
	reader = bytes.NewReader(resp.Bytes())
	x, err = ParseGetReponse(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp, x)

	resp = &ResponseGet{
		Status: StatusKeyNotFound,
		Value:  []byte(""),
	}
	reader = bytes.NewReader(resp.Bytes())
	x, err = ParseGetReponse(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp, x)

	resp = &ResponseGet{
		Status: StatusNone,
		Value:  []byte(""),
	}
	reader = bytes.NewReader(resp.Bytes())
	x, err = ParseGetReponse(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp, x)
}

func TestParseDelResponseStatusOk(t *testing.T) {
	resp := &ResponseDelete{
		Status: StatusOK,
	}
	reader := bytes.NewReader(resp.Bytes())
	x, err := ParseDelReponse(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp, x)

	resp = &ResponseDelete{
		Status: StatusError,
	}
	reader = bytes.NewReader(resp.Bytes())
	x, err = ParseDelReponse(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp, x)

	resp = &ResponseDelete{
		Status: StatusKeyNotFound,
	}
	reader = bytes.NewReader(resp.Bytes())
	x, err = ParseDelReponse(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp, x)

	resp = &ResponseDelete{
		Status: StatusNone,
	}
	reader = bytes.NewReader(resp.Bytes())
	x, err = ParseDelReponse(reader)
	assert.NoError(t, err)
	assert.Equal(t, resp, x)
}
