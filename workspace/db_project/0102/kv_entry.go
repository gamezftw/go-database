package db0102

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Entry struct {
	key []byte
	val []byte
}

func (ent *Entry) Encode() []byte {
	keyLen := uint32(len(ent.key))
	valLen := uint32(len(ent.val))

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, valLen)
	buf.Write(ent.key)
	buf.Write(ent.val)

	return buf.Bytes()
}

func (ent *Entry) Decode(r io.Reader) error {
	keyLenByte := make([]byte, 4)
	_, keyLenError := r.Read(keyLenByte)

	valLenByte := make([]byte, 4)
	_, valLenError := r.Read(valLenByte)

	keyData := make([]byte, binary.LittleEndian.Uint32(keyLenByte))
	_, keyDataError := r.Read(keyData)
	ent.key = keyData

	valData := make([]byte, binary.LittleEndian.Uint32(valLenByte))
	_, valDataError := r.Read(valData)
	ent.val = valData

	return errors.Join(keyLenError, valLenError, keyDataError, valDataError)
}

// QzBQWVJJOUhU https://trialofcode.org/
