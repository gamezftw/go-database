package db0102

import (
	"encoding/binary"
	"io"
)

type Entry struct {
	key []byte
	val []byte
}

func (ent *Entry) Encode() []byte {
	keyLen := len(ent.key)
	valLen := len(ent.val)
	data := make([]byte, 4+4+keyLen+valLen)
	binary.LittleEndian.PutUint32(data[0:4], uint32(keyLen))
	binary.LittleEndian.PutUint32(data[4:8], uint32(valLen))
	copy(data[8:], ent.key)
	copy(data[8+keyLen:], ent.val)
	return data
}

func (ent *Entry) Decode(r io.Reader) error {
	keyLenByte := make([]byte, 4)
	r.Read(keyLenByte)

	valLenByte := make([]byte, 4)
	r.Read(valLenByte)

	keyData := make([]byte, binary.LittleEndian.Uint32(keyLenByte))
	r.Read(keyData)
	ent.key = keyData

	valData := make([]byte, binary.LittleEndian.Uint32(valLenByte))
	r.Read(valData)
	ent.val = valData

	return nil
}

// QzBQWVJJOUhU https://trialofcode.org/
