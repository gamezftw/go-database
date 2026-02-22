package db0202

import (
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
	key = append(key, []byte(schema.Table)...)
	key = append([]byte(schema.Table), 0x00)
	for _, col := range schema.PKey {
		key = row[col].Encode(key)
	}
	return key
}

func (row Row) EncodeVal(schema *Schema) (val []byte) {
	for colIndex := range schema.Cols {
		if !slices.Contains(schema.PKey, colIndex) {
			val = row[colIndex].Encode(val)
		}
	}

	return val
}

func (row Row) DecodeKey(schema *Schema, key []byte) (err error) {
	key = key[len(schema.Table)+1:]
	for _, colIndex := range schema.PKey {
		row[colIndex].Type = schema.Cols[colIndex].Type
		key, err = row[colIndex].Decode(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (row Row) DecodeVal(schema *Schema, val []byte) (err error) {
	for colIndex := range schema.Cols {
		if !slices.Contains(schema.PKey, colIndex) {
			row[colIndex].Type = schema.Cols[colIndex].Type
			val, err = row[colIndex].Decode(val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// QzBQWVJJOUhU https://trialofcode.org/
