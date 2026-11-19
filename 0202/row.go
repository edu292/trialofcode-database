package db0202

import (
	"errors"
	"slices"
)

type Schema struct {
	Table string
	Cols  []Column
	PKey  []int // indexes of primary key columns
}

type Column struct {
	Name string
	Type CellType
}

type Row []Cell

func (schema *Schema) NewRow() Row {
	return make(Row, len(schema.Cols))
}

func (row Row) EncodeKey(schema *Schema) (key []byte) {
	key = append([]byte(schema.Table), 0x00)
	for _, pkeyIndex := range schema.PKey {
		cell := row[pkeyIndex]
		cell.Type = schema.Cols[pkeyIndex].Type
		key = cell.Encode(key)
	}

	return key
}

func (row Row) EncodeVal(schema *Schema) (val []byte) {
	for idx, cell := range row {
		if slices.Contains(schema.PKey, idx) {
			continue
		}
		cell.Type = schema.Cols[idx].Type

		val = cell.Encode(val)
	}

	return val
}

func (row Row) DecodeKey(schema *Schema, key []byte) (err error) {
	prefix := append([]byte(schema.Table), 0x00)
	if !slices.Equal(prefix, key[:len(prefix)]) {
		errors.New("mismatched table name")
	}

	key = key[len(prefix):]
	for _, pkeyIndex := range schema.PKey {
		cell := &row[pkeyIndex]
		cell.Type = schema.Cols[pkeyIndex].Type

		if key, err = cell.Decode(key); err != nil {
			return err
		}
	}

	return nil
}

func (row Row) DecodeVal(schema *Schema, val []byte) (err error) {
	for idx := range row {
		if slices.Contains(schema.PKey, idx) {
			continue
		}
		cell := &row[idx]
		cell.Type = schema.Cols[idx].Type

		if val, err = cell.Decode(val); err != nil {
			return err
		}
	}

	return nil
}

// QzBQWVJJOUhU https://trialofcode.org/
