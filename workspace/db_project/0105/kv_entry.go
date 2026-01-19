package db0105

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

type Entry struct {
	key     []byte
	val     []byte
	deleted bool
}

func (ent *Entry) Encode() []byte {
	keyLen := uint32(len(ent.key))
	valLen := uint32(len(ent.val))

	dataEnt := make([]byte, 4+4+1+keyLen+valLen)

	binary.LittleEndian.PutUint32(dataEnt[:4], keyLen)
	binary.LittleEndian.PutUint32(dataEnt[4:8], valLen)
	if ent.deleted {
		dataEnt[8] = 1
	}

	copy(dataEnt[9:9+keyLen], ent.key)
	copy(dataEnt[9+keyLen:], ent.val)

	data := make([]byte, 4+len(dataEnt))
	chk := crc32.ChecksumIEEE(dataEnt)
	binary.LittleEndian.PutUint32(data[:4], chk)

	copy(data[4:], dataEnt)

	return data
}

var ErrBadSum = errors.New("bad checksum")

func (ent *Entry) Decode(r io.Reader) error {

	header := make([]byte, 4+4+4+1)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return err
	}

	dataChk := binary.LittleEndian.Uint32(header[:4])
	keyLen := binary.LittleEndian.Uint32(header[4:8])
	valLen := binary.LittleEndian.Uint32(header[8:])
	deleted := header[12] == 1

	data := make([]byte, keyLen+valLen)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return err
	}
	key := data[:keyLen]
	val := data[keyLen:]

	dataForChecksum := make([]byte, 4+4+1+keyLen+valLen)
	binary.LittleEndian.PutUint32(dataForChecksum[0:4], keyLen)
	binary.LittleEndian.PutUint32(dataForChecksum[4:8], valLen)
	if deleted {
		dataForChecksum[8] = 1
	}
	copy(dataForChecksum[9:9+keyLen], key)
	copy(dataForChecksum[9+keyLen:], val)

	chk := crc32.ChecksumIEEE(dataForChecksum)

	if chk != dataChk {
		return ErrBadSum
	}

	ent.deleted = deleted
	ent.key = key
	if len(data) > int(keyLen) {
		ent.val = val
	}

	return nil
}

// QzBQWVJJOUhU https://trialofcode.org/
