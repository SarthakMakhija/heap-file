package index

import (
	"reflect"
	"testing"
)

func TestGetsTheIndexForAKey(t *testing.T) {
	page := Page{
		keyValuePairs: []KeyValuePair{
			{key: []byte("A")},
			{key: []byte("B")},
			{key: []byte("C")},
		},
	}
	expectedIndex := 0
	index, _ := page.Get([]byte("A"))

	if index != expectedIndex {
		t.Fatalf("Expected index of searched key A to be %v, received %v", expectedIndex, index)
	}
}

func TestReturnsTrueIfKeyIsPresentInThePage(t *testing.T) {
	page := Page{
		keyValuePairs: []KeyValuePair{
			{key: []byte("A")},
			{key: []byte("B")},
			{key: []byte("C")},
		},
	}
	_, found := page.Get([]byte("B"))

	if found != true {
		t.Fatalf("Expected A to be found")
	}
}

func TestReturnsFalseIfKeyIsNotPresentInThePage(t *testing.T) {
	page := Page{
		keyValuePairs: []KeyValuePair{
			{key: []byte("C")},
		},
	}
	_, found := page.Get([]byte("D"))

	if found != false {
		t.Fatalf("Expected A to not be found")
	}
}

func TestUnMarshalsAPageWithKeyValuePairCountAs1(t *testing.T) {
	page := Page{
		keyValuePairs: []KeyValuePair{
			{
				key:   []byte("C"),
				value: []byte("Storage"),
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
				value: []byte("Storage"),
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
				value: []byte("Storage"),
			},
		},
	}
	bytes := page.MarshalBinary()

	newPage := &Page{}
	newPage.UnMarshalBinary(bytes)

	value := newPage.keyValuePairs[0].PrettyValue()
	if value != "Storage" {
		t.Fatalf("Expected value to be Storage, received %v", value)
	}
}

func TestUnMarshalsAPageWithMultipleKeyValuePairs(t *testing.T) {
	page := Page{
		keyValuePairs: []KeyValuePair{
			{
				key:   []byte("A"),
				value: []byte("Database"),
			},
			{
				key:   []byte("B"),
				value: []byte("Storage"),
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

func TestInsertsAtAnIndexInAPage(t *testing.T) {
	page := &Page{
		keyValuePairs: []KeyValuePair{
			{key: []byte("A"), value: []byte("Database")},
			{key: []byte("C"), value: []byte("Storage")},
			{key: []byte("F"), value: []byte("Systems")},
		},
	}
	page.insertAt(2, KeyValuePair{key: []byte("D"), value: []byte("Operating")})
	expected := []KeyValuePair{
		{key: []byte("A"), value: []byte("Database")},
		{key: []byte("C"), value: []byte("Storage")},
		{key: []byte("D"), value: []byte("Operating")},
		{key: []byte("F"), value: []byte("Systems")},
	}

	pageKeyValuePairs := page.keyValuePairs
	if !reflect.DeepEqual(expected, pageKeyValuePairs) {
		t.Fatalf("Expected Key value pairs to be %v, received %v", expected, pageKeyValuePairs)
	}
}

func TestInsertsChildPageAtAnIndex(t *testing.T) {
	page := &Page{
		childPageIds: []int{8, 10, 14},
	}
	childPage := NewPage(11)
	expected := []int{8, 10, 11, 14}
	page.insertChildAt(2, childPage)

	actualChildPageId := page.childPageIds
	if !reflect.DeepEqual(expected, actualChildPageId) {
		t.Fatalf("Expected child page ids to be %v, received %v", expected, actualChildPageId)
	}
}

func TestSplitsALeafPageWithKeyValuePairs(t *testing.T) {
	page := &Page{
		id:            0,
		keyValuePairs: []KeyValuePair{{key: []byte("A"), value: []byte("Database")}, {key: []byte("B"), value: []byte("Systems")}},
	}
	parentPage := NewPage(100)
	parentPage.childPageIds = []int{0}
	siblingPage := NewPage(200)

	_ = page.split(parentPage, siblingPage, 0)

	keyValuePairsAfterSplit := page.keyValuePairs
	expected := []KeyValuePair{{key: []byte("A"), value: []byte("Database")}}

	if !reflect.DeepEqual(expected, keyValuePairsAfterSplit) {
		t.Fatalf("Expected key value pairs in the page after split to be %v, received %v", expected, keyValuePairsAfterSplit)
	}
}

func TestSplitsALeafPageWithKeyValuePairsInParent(t *testing.T) {
	page := &Page{
		id:            0,
		keyValuePairs: []KeyValuePair{{key: []byte("A"), value: []byte("Database")}, {key: []byte("B"), value: []byte("Systems")}},
	}
	parentPage := NewPage(100)
	parentPage.childPageIds = []int{0}
	siblingPage := NewPage(200)

	_ = page.split(parentPage, siblingPage, 0)

	keyValuePairsAfterSplit := parentPage.keyValuePairs
	expected := []KeyValuePair{{key: []byte("B")}}

	if !reflect.DeepEqual(expected, keyValuePairsAfterSplit) {
		t.Fatalf("Expected key value pairs in the parent page after split to be %v, received %v", expected, keyValuePairsAfterSplit)
	}
}

func TestSplitsALeafPageWithKeyValuePairsInSibling(t *testing.T) {
	t.Skip()

	page := &Page{
		id:            0,
		keyValuePairs: []KeyValuePair{{key: []byte("A"), value: []byte("Database")}, {key: []byte("B"), value: []byte("Systems")}},
	}
	parentPage := NewPage(100)
	parentPage.childPageIds = []int{0}
	siblingPage := NewPage(200)

	_ = page.split(parentPage, siblingPage, 0)

	keyValuePairsAfterSplit := siblingPage.keyValuePairs
	expected := []KeyValuePair{{key: []byte("B"), value: []byte("Systems")}}

	if !reflect.DeepEqual(expected, keyValuePairsAfterSplit) {
		t.Fatalf("Expected key value pairs in the parent page after split to be %v, received %v", expected, keyValuePairsAfterSplit)
	}
}
