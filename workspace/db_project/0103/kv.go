package db0103

import (
	"bytes"
)

type KV struct {
	log Log
	mem map[string][]byte
}

func (kv *KV) Open() error {
	if err := kv.log.Open(); err != nil {
		return err
	}
	kv.mem = map[string][]byte{}

	for {
		ent := Entry{}
		eof, err := kv.log.Read(&ent)

		if eof {
			return nil
		}
		if err != nil {
			return err
		}

		keyString := string(ent.key)
		if !ent.deleted {
			kv.mem[string(ent.key)] = ent.val
		} else {
			delete(kv.mem, keyString)
		}
	}
}

func (kv *KV) Close() error { return kv.log.Close() }

func (kv *KV) Get(key []byte) (val []byte, ok bool, err error) {
	val, ok = kv.mem[string(key)]
	return val, ok, nil
}

func (kv *KV) Set(key []byte, val []byte) (updated bool, err error) {
	keyString := string(key)
	currentValue, exists := kv.mem[keyString]
	isTheSameValue := bytes.Equal(currentValue, val)
	shouldUpdate := !exists && !isTheSameValue

	if shouldUpdate {
		kv.mem[keyString] = val
		if err := kv.log.Write(&Entry{key: key, val: val, deleted: false}); err != nil {
			return false, err
		}
	}

	return shouldUpdate, nil
}

func (kv *KV) Del(key []byte) (deleted bool, err error) {
	keyString := string(key)

	if _, exists := kv.mem[keyString]; exists {
		deleted = true
		delete(kv.mem, keyString)
		if err := kv.log.Write(&Entry{key: key, val: nil, deleted: true}); err != nil {
			return false, err
		}
	}

	return deleted, nil
}

// QzBQWVJJOUhU https://trialofcode.org/
