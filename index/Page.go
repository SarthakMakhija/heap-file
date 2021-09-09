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
	keyValuePairs []keyValuePair
}

type keyValuePair struct {
	key   []byte
	value uint64
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

func (page Page) MarshalBinary() []byte {
	buffer := make([]byte, 400) //will be replaced with page.Size()
	offset := 0

	writeLeafPageType := func() {
		buffer[offset] = LeafPage
		offset++
	}
	writeKeyValuePairCount := func() {
		littleEndian.PutUint16(buffer[offset:offset+2], uint16(len(page.keyValuePairs)))
		offset += 2
	}
	writeValue := func(value uint64) {
		littleEndian.PutUint64(buffer[offset:offset+8], value)
		offset += 8
	}
	writeKeyLength := func(length uint16) {
		littleEndian.PutUint16(buffer[offset:offset+2], length)
		offset += 2
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
	readValue := func() uint64 {
		value := littleEndian.Uint64(buffer[offset : offset+8])
		offset += 8
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
			pair := keyValuePair{}
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
	return true //for now this is true
}
