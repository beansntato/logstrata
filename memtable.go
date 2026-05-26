package logstrata

import (
	"time"

	"github.com/huandu/skiplist"
)

type MemTable struct {
	data skiplist.SkipList
	size int64
}

func New(size int64) *MemTable {
	return &MemTable{data: *skiplist.New(skiplist.String), size: size}
}

func (mt *MemTable) Put(key string, value []byte) {
	sizeChange := int64(len(value))
	existingEntry := mt.data.Get(key)

	if existingEntry != nil {
		mt.size -= int64(len(existingEntry.Value.(*LSMEntry).Value))
	} else {
		sizeChange += int64(len(key))
	}

	entry := getLSMEntry(key, &value, Command_PUT)
	mt.data.Set(key, entry)
	mt.size += sizeChange
}

func (mt *MemTable) Delete(key string) {
	existingEntry := mt.data.Get(key)

	if existingEntry != nil {
		mt.size -= int64(len(existingEntry.Value.(*LSMEntry).Value))
	} else {
		mt.size += int64(len(key))
	}

	mt.data.Set(key, getLSMEntry(key, nil, Command_DELETE))
}

func (mt *MemTable) Get(key string) *LSMEntry {
	value := mt.data.Get(key)

	if value == nil {
		return nil
	}

	return value.Value.(*LSMEntry)
}

func (mt *MemTable) RangeScan(startKey string, endKey string) []*LSMEntry {
	var results []*LSMEntry

	iter := mt.data.Find(startKey)
	for iter != nil {
		if iter.Element().Key().(string) > endKey {
			break
		}

		results = append(results, iter.Value.(*LSMEntry))
		iter = iter.Next()
	}

	return results
}

func (mt *MemTable) GetEntries() []*LSMEntry {
	var results []*LSMEntry
	iter := mt.data.Front()
	for iter != nil {
		results = append(results, iter.Value.(*LSMEntry))
		iter = iter.Next()
	}

	return results
}

func getLSMEntry(key string, value *[]byte, command Command) *LSMEntry {
	entry := &LSMEntry{
		Key:       key,
		Command:   command,
		Timestamp: time.Now().UnixNano(),
	}
	if value != nil {
		entry.Value = *value
	}
	return entry
}
