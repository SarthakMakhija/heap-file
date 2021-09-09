package index

import (
	"testing"
)

func TestGetsTheIndexForAKey(t *testing.T) {
	page := Page{
		keyValuePairs: []KeyValuePair{
			{
				key:   []byte("A"),
				value: uint64(100),
			},
			{
				key:   []byte("B"),
				value: uint64(200),
			},
			{
				key:   []byte("C"),
				value: uint64(300),
			},
		},
	}
	expectedIndex := 0
	_, index, _ := page.Get([]byte("A"))

	if index != expectedIndex {
		t.Fatalf("Expected index of searched key A to be %v, received %v", expectedIndex, index)
	}
}

func TestReturnsTrueIfKeyIsPresentInThePage(t *testing.T) {
	page := Page{
		keyValuePairs: []KeyValuePair{
			{
				key:   []byte("A"),
				value: uint64(100),
			},
			{
				key:   []byte("B"),
				value: uint64(200),
			},
			{
				key:   []byte("C"),
				value: uint64(300),
			},
		},
	}
	_, _, found := page.Get([]byte("B"))

	if found != true {
		t.Fatalf("Expected A to be found")
	}
}

func TestReturnsFalseIfKeyIsNotPresentInThePage(t *testing.T) {
	page := Page{
		keyValuePairs: []KeyValuePair{
			{
				key:   []byte("C"),
				value: uint64(300),
			},
		},
	}
	_, _, found := page.Get([]byte("D"))

	if found != false {
		t.Fatalf("Expected A to not be found")
	}
}

func TestUnMarshalsAPageWithKeyValuePairCountAs1(t *testing.T) {
	page := Page{
		keyValuePairs: []KeyValuePair{
			{
				key:   []byte("C"),
				value: uint64(300),
			},
		},
	}
	bytes := page.MarshalBinary()

	newPage := &Page{}
	newPage.UnMarshalBinary(bytes)

	keyValuePairCount := len(newPage.keyValuePairs)
	if keyValuePairCount != 1 {
		t.Fatalf("Expected keyValuePairCount to be 1, received %v", keyValuePairCount)
	}
}

func TestUnMarshalsAPageWithKey(t *testing.T) {
	page := Page{
		keyValuePairs: []KeyValuePair{
			{
				key:   []byte("C"),
				value: uint64(300),
			},
		},
	}
	bytes := page.MarshalBinary()

	newPage := &Page{}
	newPage.UnMarshalBinary(bytes)

	key := string(newPage.keyValuePairs[0].key)
	if key != "C" {
		t.Fatalf("Expected key to be C, received %v", key)
	}
}

func TestUnMarshalsAPageWithValue(t *testing.T) {
	page := Page{
		keyValuePairs: []KeyValuePair{
			{
				key:   []byte("C"),
				value: uint64(300),
			},
		},
	}
	bytes := page.MarshalBinary()

	newPage := &Page{}
	newPage.UnMarshalBinary(bytes)

	value := newPage.keyValuePairs[0].value
	if value != 300 {
		t.Fatalf("Expected value to be 300, received %v", value)
	}
}

func TestUnMarshalsAPageWithMultipleKeyValuePairs(t *testing.T) {
	page := Page{
		keyValuePairs: []KeyValuePair{
			{
				key:   []byte("A"),
				value: uint64(100),
			},
			{
				key:   []byte("B"),
				value: uint64(200),
			},
		},
	}
	bytes := page.MarshalBinary()

	newPage := &Page{}
	newPage.UnMarshalBinary(bytes)

	keyValuePairCount := len(newPage.keyValuePairs)
	if keyValuePairCount != 2 {
		t.Fatalf("Expected keyValuePairCount to be 2, received %v", keyValuePairCount)
	}

	expectedFirstKeyValuePair := page.keyValuePairs[0]
	firstKeyValuePair := newPage.keyValuePairs[0]

	if !expectedFirstKeyValuePair.Equals(firstKeyValuePair) {
		t.Fatalf("Expected first key value pair to be %v, received %v", expectedFirstKeyValuePair, firstKeyValuePair)
	}

	expectedSecondKeyValuePair := page.keyValuePairs[1]
	secondKeyValuePair := newPage.keyValuePairs[1]

	if !expectedSecondKeyValuePair.Equals(secondKeyValuePair) {
		t.Fatalf("Expected second key value pair to be %v, received %v", expectedSecondKeyValuePair, secondKeyValuePair)
	}
}
