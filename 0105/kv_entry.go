package db0105

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"slices"
)

type Entry struct {
	key     []byte
	val     []byte
	deleted bool
}

func (ent *Entry) Encode() []byte {
	key_len := uint32(len(ent.key))
	val_len := uint32(len(ent.val))
	data := make([]byte, 4+4+4+1+key_len+val_len)
	binary.LittleEndian.PutUint32(data[4:8], key_len)
	copy(data[13:], ent.key)
	if ent.deleted {
		data[12] = 1
	} else {
		binary.LittleEndian.PutUint32(data[8:12], val_len)
		copy(data[13+key_len:], ent.val)
	}
	binary.LittleEndian.PutUint32(data[0:], crc32.ChecksumIEEE(data[4:]))

	return data
}

var ErrBadSum = errors.New("bad checksum")

func (ent *Entry) Decode(r io.Reader) error {
	header_len := 13
	data := make([]byte, header_len)
	if _, err := io.ReadFull(r, data); err != nil {
		return err
	}

	key_len := int(binary.LittleEndian.Uint32(data[4:8]))
	val_len := int(binary.LittleEndian.Uint32(data[8:12]))
	data = slices.Grow(data, key_len+val_len)[:header_len+key_len+val_len]
	if _, err := io.ReadFull(r, data[header_len:]); err != nil {
		return err
	}

	if computed := crc32.ChecksumIEEE(data[4:]); computed != binary.LittleEndian.Uint32(data[0:4]) {
		return ErrBadSum
	}

	ent.key = data[header_len : header_len+key_len]
	if data[12] != 0 {
		ent.deleted = true
		ent.val = nil
	} else {
		ent.deleted = false
		ent.val = data[header_len+key_len:]
	}
	return nil
}

// QzBQWVJJOUhU https://trialofcode.org/
