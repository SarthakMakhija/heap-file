package index

import "testing"

func TestReturnsTrueGivenKeyValuePairsAreEqual(t *testing.T) {
	firstKeyValuePair := KeyValuePair{
		key:   []byte("A"),
		value: uint64(100),
	}

	secondKeyValuePair := KeyValuePair{
		key:   []byte("A"),
		value: uint64(100),
	}

	if !firstKeyValuePair.Equals(secondKeyValuePair) {
		t.Fatalf("Expected key value pairs to be equals")
	}
}

func TestReturnsFalseGivenKeyValuePairsAreNotEqualByKey(t *testing.T) {
	firstKeyValuePair := KeyValuePair{
		key:   []byte("A"),
		value: uint64(100),
	}

	secondKeyValuePair := KeyValuePair{
		key:   []byte("B"),
		value: uint64(100),
	}

	if firstKeyValuePair.Equals(secondKeyValuePair) {
		t.Fatalf("Expected key value pairs to not be equal")
	}
}

func TestReturnsFalseGivenKeyValuePairsAreNotEqualByValue(t *testing.T) {
	firstKeyValuePair := KeyValuePair{
		key:   []byte("A"),
		value: uint64(100),
	}

	secondKeyValuePair := KeyValuePair{
		key:   []byte("A"),
		value: uint64(200),
	}

	if firstKeyValuePair.Equals(secondKeyValuePair) {
		t.Fatalf("Expected key value pairs to not be equal")
	}
}
