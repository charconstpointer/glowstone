package glowstone

import (
	"fmt"
	"strings"
)

//Value is a []byte alias for util type functions
type Value []byte

func (v Value) String() string {
	return string(v)
}

//ID represents identifier of maximum size of 64 bytes
type ID struct {
	Value Value
}

//NewID is
func NewID(value string) (*ID, error) {
	if len(value) > 64 {
		return nil, fmt.Errorf("ID can have maximum length of %d, provided value is larger", 64)
	}

	b := make([]byte, 64)
	copy(b, value)

	return &ID{
		Value: b,
	}, nil
}

//AddID appends ID to a given slice of data
func AddID(payload []byte, ID ID) []byte {
	payloadID := append(payload, ID.Value...)
	x := string(payloadID)
	println(x)
	return payloadID
}

//GetID returns encoded ID in a []byte
func GetID(payload []byte) string {
	b := payload[len(payload)-64:]
	ID := string(b)
	return fmt.Sprintf("%s", strings.TrimFunc(ID, func(r rune) bool {
		return r == 0
	}))
}

//RemoveID removes ID from given payload and returns identifier and cleaned payload
func RemoveID(payload []byte) ([]byte, error) {
	if len(payload) < 64 {
		return nil, fmt.Errorf("Provided payload is not small, it could not have ID encoded into it")
	}
	return payload[:len(payload)-64], nil
}
