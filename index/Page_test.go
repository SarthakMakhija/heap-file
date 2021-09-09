package index

import (
	"testing"
)

func TestGetsTheIndexForAKey(t *testing.T) {
	page := Page{
		keyValuePairs: []keyValuePair{
			{
				key:   []byte("A"),
				value: []byte("A-Value"),
			},
			{
				key:   []byte("B"),
				value: []byte("B-Value"),
			},
			{
				key:   []byte("C"),
				value: []byte("C-Value"),
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
				value: []byte("A-Value"),
			},
			{
				key:   []byte("B"),
				value: []byte("B-Value"),
			},
			{
				key:   []byte("C"),
				value: []byte("C-Value"),
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
				value: []byte("C-Value"),
			},
		},
	}
	_, _, found := page.Get([]byte("D"))

	if found != false {
		t.Fatalf("Expected A to not be found")
	}
}
