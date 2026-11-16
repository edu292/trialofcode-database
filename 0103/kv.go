package db0103

import "slices"

type KV struct {
	log Log
	mem map[string][]byte
}

func (kv *KV) Open() error {
	if err := kv.log.Open(); err != nil {
		return err
	}

	kv.mem = map[string][]byte{}
	var ent Entry
	for {
		eof, err := kv.log.Read(&ent)
		if err != nil {
			return err
		}

		if eof {
			break
		}

		key_str := string(ent.key)
		if ent.deleted {
			delete(kv.mem, key_str)
		} else {
			kv.mem[key_str] = ent.val
		}
	}

	return nil
}

func (kv *KV) Close() error { return kv.log.Close() }

func (kv *KV) Get(key []byte) (val []byte, ok bool, err error) {
	val, ok = kv.mem[string(key)]
	return
}

func (kv *KV) Set(key []byte, val []byte) (updated bool, err error) {
	key_str := string(key)
	prev, exist := kv.mem[key_str]
	updated = !exist || !slices.Equal(prev, val)
	if updated {
		if err := kv.log.Write(&Entry{key: key, val: val, deleted: false}); err != nil {
			return false, err
		}
	}
	kv.mem[key_str] = val

	return updated, nil
}

func (kv *KV) Del(key []byte) (deleted bool, err error) {
	key_str := string(key)
	if _, exist := kv.mem[key_str]; !exist {
		return false, nil
	}

	if err := kv.log.Write(&Entry{key: key, deleted: true}); err != nil {
		return false, err
	}

	delete(kv.mem, key_str)
	return true, nil
}

// QzBQWVJJOUhU https://trialofcode.org/
