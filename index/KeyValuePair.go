package index

import "bytes"

type KeyValuePair struct {
	key   []byte
	value []byte
}

func (keyValuePair KeyValuePair) Equals(other KeyValuePair) bool {
	if bytes.Equal(keyValuePair.value, other.value) && bytes.Equal(keyValuePair.key, other.key) {
		return true
	}
	return false
}

func (keyValuePair KeyValuePair) PrettyValue() string {
	return string(keyValuePair.value)
}
