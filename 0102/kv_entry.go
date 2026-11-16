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
	key_len := uint32(len(ent.key))
	val_len := uint32(len(ent.val))
	data := make([]byte, 4+4+key_len+val_len)
	binary.LittleEndian.PutUint32(data[0:4], key_len)
	binary.LittleEndian.PutUint32(data[4:8], val_len)
	copy(data[8:], ent.key)
	copy(data[8+key_len:], ent.val)

	return data
}

func (ent *Entry) Decode(r io.Reader) error {
	var len_buf [4]byte
	if _, err := io.ReadFull(r, len_buf[:]); err != nil {
		return err
	}
	key_len := binary.LittleEndian.Uint32(len_buf[:])

	if _, err := io.ReadFull(r, len_buf[:]); err != nil {
		return err
	}
	val_len := binary.LittleEndian.Uint32(len_buf[:])

	key := make([]byte, key_len)
	if _, err := io.ReadFull(r, key); err != nil {
		return err
	}

	val := make([]byte, val_len)
	if _, err := io.ReadFull(r, val); err != nil {
		return err
	}

	ent.key = key
	ent.val = val
	return nil
}

// QzBQWVJJOUhU https://trialofcode.org/
