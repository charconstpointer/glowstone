package glowstone

import (
	"strings"
	"testing"
)

func TestAddIdLength(t *testing.T) {
	payload := make([]byte, 10)
	value := "someid"
	ID, _ := NewID(value)

	payloadLen := len(payload)
	payloadWithID := AddID(payload, *ID)
	got := len(payloadWithID)
	expected := payloadLen + 64
	if got != expected {
		t.Errorf("Expected new payload to be of length %d but instead got len of %d", expected, got)
	}
}

func TestAddIdContent(t *testing.T) {
	payload := make([]byte, 10)
	value := "someid"
	ID, _ := NewID(value)
	payloadWithID := AddID(payload, *ID)
	expected := true
	got := strings.Contains(string(payloadWithID), ID.Value.String())

	if got != expected {
		t.Errorf("Expected payload to contain provided ID")
	}
}

func TestGetID(t *testing.T) {
	value := "identifier"
	ID, _ := NewID(value)
	payload := AddID(make([]byte, 10), *ID)
	expected := value
	got := GetID(payload)
	if expected != got {
		t.Errorf("Expected extracted ID to be %s instead got %s", expected, got)
	}
}

func TestNewID(t *testing.T) {
	value := "foobar"
	got, err := NewID(value)
	if err != nil {
		t.Errorf("Expected error to be nil")
	}
	if got == nil {
		t.Errorf("Expected ID to no be nil")
		return
	}
	expected := value
	gotValue := GetID(got.Value)
	if gotValue != expected {
		t.Errorf("Expected ID value to be %s instead got %s", expected, gotValue)
	}
}
