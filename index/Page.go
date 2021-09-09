package index

import (
	"bytes"
	"sort"
)

type Page struct {
	id            int
	keyValuePairs []keyValuePair
}

type keyValuePair struct {
	key   []byte
	value []byte
}

func NewPage(id int) *Page {
	return &Page{
		id: id,
	}
}

func (page Page) Get(key []byte) (Page, int, bool) {
	index, found := page.binarySearch(key)
	if page.isLeaf() {
		return page, index, found
	}
	//else :: handle non leaf and found, which means load the child page..
	return Page{}, 0, false
}

func (page Page) binarySearch(key []byte) (int, bool) {
	index := sort.Search(len(page.keyValuePairs), func(index int) bool {
		if bytes.Compare(key, page.keyValuePairs[index].key) < 0 {
			return true
		}
		return false
	})
	if index > 0 && bytes.Compare(page.keyValuePairs[index-1].key, key) == 0 {
		return index - 1, true
	}
	return index, false
}

func (page Page) isLeaf() bool {
	return true //for now this is true
}
