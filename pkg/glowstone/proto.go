package glowstone

import (
	"encoding/binary"
)

type Header []byte

const HeaderSize = 10

type MessageType uint16

const (
	PASS MessageType = iota
	REPL
)

func (h Header) Encode(mtype MessageType, id int32, len int32) {
	binary.BigEndian.PutUint16(h[0:2], uint16(mtype))
	binary.BigEndian.PutUint32(h[2:6], uint32(id))
	binary.BigEndian.PutUint32(h[6:10], uint32(len))
}

func (h Header) MessageType() MessageType {
	value := binary.BigEndian.Uint16(h[0:2])
	return MessageType(value)
}

func (h Header) ID() uint32 {
	return binary.BigEndian.Uint32(h[2:6])
}

func (h Header) Len() uint32 {
	return binary.BigEndian.Uint32(h[6:10])
}
