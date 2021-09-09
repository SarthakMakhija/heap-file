package index

import "bytes"

type KeyValuePair struct {
	key   []byte
	value uint64
}

func (keyValuePair KeyValuePair) Equals(other KeyValuePair) bool {
	if keyValuePair.value == other.value && bytes.Equal(keyValuePair.key, other.key) {
		return true
	}
	return false
}
