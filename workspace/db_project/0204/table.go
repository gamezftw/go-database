package db0204

type DB struct {
	KV KV
}

func (db *DB) Open() error  { return db.KV.Open() }
func (db *DB) Close() error { return db.KV.Close() }

func (db *DB) Select(schema *Schema, row Row) (ok bool, err error) {
	key := row.EncodeKey(schema)
	val, ok, err := db.KV.Get(key)
	if err != nil {
		return false, err
	}
	if ok {
		err = row.DecodeVal(schema, val)
	}
	return ok, err
}

func (db *DB) Insert(schema *Schema, row Row) (updated bool, err error) {
	return db.setWithMode(schema, row, ModeInsert)
}

func (db *DB) Upsert(schema *Schema, row Row) (updated bool, err error) {
	return db.setWithMode(schema, row, ModeUpsert)
}

func (db *DB) Update(schema *Schema, row Row) (updated bool, err error) {
	return db.setWithMode(schema, row, ModeUpdate)
}

func (db *DB) setWithMode(schema *Schema, row Row, mode UpdateMode) (updated bool, err error) {
	key := row.EncodeKey(schema)
	val := row.EncodeVal(schema)
	return db.KV.SetEx(key, val, mode)
}

func (db *DB) Delete(schema *Schema, row Row) (deleted bool, err error) {
	key := row.EncodeKey(schema)
	return db.KV.Del(key)
}

// QzBQWVJJOUhU https://trialofcode.org/
