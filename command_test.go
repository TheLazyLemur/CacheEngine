package main

import (
	"testing"
)

func TestToBytesForSetCommand(t *testing.T) {
	msg :=  "SET Foo Bar 25"

	expectedBytes := []byte(msg)

	message := &Message {
		Cmd: "SET",
		Key: []byte("Foo"),
		Value: []byte("Bar"),
		Ttl: 25,
	}

	actualBytes := message.ToBytes()

	for i := 0; i < len(actualBytes); i++ {
		if actualBytes[i] != expectedBytes[i] {
			t.Errorf("Expected %v, got %v", expectedBytes, actualBytes)
		}
	}

	if string(actualBytes) != msg {
		t.Errorf("Expected %v, got %v", msg, string(actualBytes))
	}
}

func TestToBytesForGetCommand(t *testing.T) {
	msg :=  "GET Foo"

	expectedBytes := []byte(msg)

	message := &Message {
		Cmd: "GET",
		Key: []byte("Foo"),
	}

	actualBytes := message.ToBytes()

	for i := 0; i < len(actualBytes); i++ {
		if actualBytes[i] != expectedBytes[i] {
			t.Errorf("Expected %v, got %v", expectedBytes, actualBytes)
		}
	}

	if string(actualBytes) != msg {
		t.Errorf("Expected %v, got %v", msg, string(actualBytes))
	}
}

func TestToBytesForDeleteCommand(t *testing.T) {
	msg :=  "DEL Foo"

	expectedBytes := []byte(msg)

	message := &Message {
		Cmd: "DEL",
		Key: []byte("Foo"),
	}

	actualBytes := message.ToBytes()

	for i := 0; i < len(actualBytes); i++ {
		if actualBytes[i] != expectedBytes[i] {
			t.Errorf("Expected %v, got %v", expectedBytes, actualBytes)
		}
	}

	if string(actualBytes) != msg {
		t.Errorf("Expected %v, got %v", msg, string(actualBytes))
	}
}

func TestToBytesForHasCommand(t *testing.T) {
	msg :=  "HAS Foo"

	expectedBytes := []byte(msg)

	message := &Message {
		Cmd: "HAS",
		Key: []byte("Foo"),
	}

	actualBytes := message.ToBytes()

	for i := 0; i < len(actualBytes); i++ {
		if actualBytes[i] != expectedBytes[i] {
			t.Errorf("Expected %v, got %v", expectedBytes, actualBytes)
		}
	}

	if string(actualBytes) != msg {
		t.Errorf("Expected %v, got %v", msg, string(actualBytes))
	}
}
