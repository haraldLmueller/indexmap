package indexmap

import "encoding/json"

func (imap *IndexMap[K, V]) MarshalJSON() ([]byte, error) {
	imap.lock.RLock()
	defer imap.lock.RUnlock()

	return json.Marshal(imap.primaryIndex.inner)
}

func (imap *IndexMap[K, V]) UnmarshalJSON(data []byte) error {
	imap.lock.Lock()
	defer imap.lock.Unlock()

	if err := json.Unmarshal(data, &imap.primaryIndex.inner); err != nil {
		return err
	}

	// the values went into the primary index directly, the secondary indexes
	// and the sorted view still have to be built
	imap.primaryIndex.iterate(func(_ K, value *V) {
		// don't use Insert() that locks the already locked map (dead lock)
		imap.insert(value)
	})

	return nil
}
