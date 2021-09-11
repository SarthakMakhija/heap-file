package index

import (
	"bytes"
	"encoding/binary"
	"sort"
)

var littleEndian = binary.LittleEndian

const (
	LeafPage = 0x00
)

type Page struct {
	id            int
	keyValuePairs []KeyValuePair
	childPageIds  []int
}

func NewPage(id int) *Page {
	return &Page{
		id: id,
	}
}

func (page Page) Get(key []byte) (int, bool) {
	return page.binarySearch(key)
}

func (page Page) GetKeyValuePairAt(index int) KeyValuePair {
	return page.keyValuePairs[index]
}

func (page Page) MarshalBinary() []byte {
	buffer := make([]byte, 30) //will be replaced with page.Size()
	offset := 0

	writeLeafPageType := func() {
		buffer[offset] = LeafPage
		offset++
	}
	writeKeyValuePairCount := func() {
		littleEndian.PutUint16(buffer[offset:offset+2], uint16(len(page.keyValuePairs)))
		offset += 2
	}
	writeValueLength := func(length uint16) {
		littleEndian.PutUint16(buffer[offset:offset+2], length)
		offset += 2
	}
	writeKeyLength := func(length uint16) {
		littleEndian.PutUint16(buffer[offset:offset+2], length)
		offset += 2
	}
	writeValue := func(value []byte) {
		copy(buffer[offset:offset+len(value)], value)
		offset += len(value)
	}
	writeKey := func(key []byte) {
		copy(buffer[offset:], key)
		offset += len(key)
	}

	if page.isLeaf() {
		writeLeafPageType()
		writeKeyValuePairCount()

		for index := 0; index < len(page.keyValuePairs); index++ {
			keyValuePair := page.keyValuePairs[index]

			writeValueLength(uint16(len(keyValuePair.value)))
			writeValue(keyValuePair.value)
			writeKeyLength(uint16(len(keyValuePair.key)))
			writeKey(keyValuePair.key)
		}
	}
	return buffer
}

func (page *Page) UnMarshalBinary(buffer []byte) {
	offset := 1

	readKeyValuePairCount := func() int {
		keyValuePairCount := int(littleEndian.Uint16(buffer[offset : offset+2]))
		offset += 2
		return keyValuePairCount
	}
	readValue := func() []byte {
		valueSize := int(littleEndian.Uint16(buffer[offset : offset+2]))
		offset += 2

		value := make([]byte, valueSize)
		copy(value, buffer[offset:offset+valueSize])
		offset += valueSize
		return value
	}
	readKey := func() []byte {
		keySize := int(littleEndian.Uint16(buffer[offset : offset+2]))
		offset += 2

		key := make([]byte, keySize)
		copy(key, buffer[offset:offset+keySize])
		offset += keySize
		return key
	}

	if buffer[0]&LeafPage == 0 {
		keyValuePairCount := readKeyValuePairCount()
		for index := 0; index < keyValuePairCount; index++ {
			pair := KeyValuePair{}
			pair.value = readValue()
			pair.key = readKey()
			page.keyValuePairs = append(page.keyValuePairs, pair)
		}
	}
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
	return len(page.childPageIds) == 0
}

func (page *Page) insertAt(index int, keyValuePair KeyValuePair) {
	page.keyValuePairs = append(page.keyValuePairs, KeyValuePair{})
	copy(page.keyValuePairs[index+1:], page.keyValuePairs[index:])
	page.keyValuePairs[index] = keyValuePair
}

func (page *Page) insertChildAt(index int, childPage *Page) {
	page.childPageIds = append(page.childPageIds, 0)
	copy(page.childPageIds[index+1:], page.childPageIds[index:])
	page.childPageIds[index] = childPage.id
}

func (page *Page) split(parentPage *Page, siblingPage *Page, index int) error {
	if page.isLeaf() {

		siblingPage.keyValuePairs = make([]KeyValuePair, len(page.keyValuePairs)/2+1)   //may change later - len(page.keyValuePairs)
		copy(siblingPage.keyValuePairs, page.keyValuePairs[len(page.keyValuePairs)/2:]) //may change later
		page.keyValuePairs = page.keyValuePairs[:len(page.keyValuePairs)/2]

		parentPage.insertChildAt(index+1, siblingPage)
		parentPage.insertAt(index, siblingPage.keyValuePairs[0])
	}
	return nil
}
