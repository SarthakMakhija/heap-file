package index

import (
	"testing"
)

func TestGetsTheIndexForAKey(t *testing.T) {
	page := Page{
		keyValuePairs: []keyValuePair{
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
		keyValuePairs: []keyValuePair{
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
		keyValuePairs: []keyValuePair{
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
		keyValuePairs: []keyValuePair{
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
		keyValuePairs: []keyValuePair{
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
		keyValuePairs: []keyValuePair{
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
		keyValuePairs: []keyValuePair{
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

	key := string(newPage.keyValuePairs[0].key)
	if key != "A" {
		t.Fatalf("Expected first key to be A, received %v", key)
	}

	value := newPage.keyValuePairs[0].value
	if value != 100 {
		t.Fatalf("Expected first value to be 100, received %v", value)
	}

	key = string(newPage.keyValuePairs[1].key)
	if key != "B" {
		t.Fatalf("Expected second key to be B, received %v", key)
	}

	value = newPage.keyValuePairs[1].value
	if value != 200 {
		t.Fatalf("Expected second value to be 200, received %v", value)
	}
}
