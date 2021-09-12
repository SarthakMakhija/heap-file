package index

import (
	"bytes"
	"encoding/binary"
	"sort"
)

var littleEndian = binary.LittleEndian
var (
	pageTypeSize          = 1
	keyValuePairCountSize = 2
	keyLengthSize         = 2
	valueLengthSize       = 2
	childPageIdSize       = 4
	childPageId0Size      = childPageIdSize
	leafPageMetaSize      = pageTypeSize + keyValuePairCountSize
	nonLeafPageMetaSize   = pageTypeSize + keyValuePairCountSize + childPageId0Size
)

const (
	LeafPage    = uint8(0x0)
	NonLeafPage = uint8(0x01)
)

type Page struct {
	id            int
	keyValuePairs []KeyValuePair
	childPageIds  []int
	dirty         bool
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
	buffer := make([]byte, page.size())
	offset := 0

	writePageType := func(pageType byte) {
		buffer[offset] = pageType
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
	writeChildPageId := func(childPageId uint32) {
		littleEndian.PutUint32(buffer[offset:offset+4], childPageId)
		offset += 4
	}

	if page.isLeaf() {
		writePageType(LeafPage)
		writeKeyValuePairCount()

		for index := 0; index < len(page.keyValuePairs); index++ {
			keyValuePair := page.keyValuePairs[index]

			writeValueLength(uint16(len(keyValuePair.value)))
			writeValue(keyValuePair.value)
			writeKeyLength(uint16(len(keyValuePair.key)))
			writeKey(keyValuePair.key)
		}
	} else {
		writePageType(NonLeafPage)
		writeKeyValuePairCount()
		writeChildPageId(uint32(page.childPageIds[0]))

		for index := 0; index < len(page.keyValuePairs); index++ {
			keyValuePair := page.keyValuePairs[index]

			writeChildPageId(uint32(page.childPageIds[index+1]))
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
	readChildPageId := func() int {
		pageId := int(littleEndian.Uint32(buffer[offset : offset+4]))
		offset += 4
		return pageId
	}

	if buffer[0]&NonLeafPage == 0 {
		keyValuePairCount := readKeyValuePairCount()
		for index := 0; index < keyValuePairCount; index++ {
			pair := KeyValuePair{}
			pair.value = readValue()
			pair.key = readKey()
			page.keyValuePairs = append(page.keyValuePairs, pair)
		}
	} else {
		keyValuePairCount := readKeyValuePairCount()
		page.childPageIds = append(page.childPageIds, readChildPageId())

		for index := 0; index < keyValuePairCount; index++ {
			childPageId := readChildPageId()
			key := readKey()

			page.childPageIds = append(page.childPageIds, childPageId)
			page.keyValuePairs = append(page.keyValuePairs, KeyValuePair{key: key})
		}
	}
}

func (page Page) size() int {
	size := 0
	if page.isLeaf() {
		size = leafPageMetaSize
		for _, keyValuePair := range page.keyValuePairs {
			size = size + keyLengthSize + valueLengthSize + len(keyValuePair.key) + len(keyValuePair.value)
		}
	} else {
		size = nonLeafPageMetaSize
		for _, keyValuePair := range page.keyValuePairs {
			size = size + keyLengthSize + len(keyValuePair.key) + childPageIdSize
		}
	}
	return size
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
	page.MarkDirty()
	page.keyValuePairs = append(page.keyValuePairs, KeyValuePair{})

	copy(page.keyValuePairs[index+1:], page.keyValuePairs[index:])
	if page.isLeaf() {
		page.keyValuePairs[index] = keyValuePair
	} else {
		page.keyValuePairs[index] = KeyValuePair{key: keyValuePair.key}
	}
}

func (page *Page) insertChildAt(index int, childPage *Page) {
	page.MarkDirty()

	page.childPageIds = append(page.childPageIds, 0)
	copy(page.childPageIds[index+1:], page.childPageIds[index:])
	page.childPageIds[index] = childPage.id
}

func (page *Page) split(parentPage *Page, siblingPage *Page, index int) error {
	page.MarkDirty()
	parentPage.MarkDirty()
	siblingPage.MarkDirty()

	if page.isLeaf() {
		pageKeyValuePairs := page.AllKeyValuePairs()
		siblingPage.keyValuePairs = append(siblingPage.keyValuePairs, page.keyValuePairs[len(pageKeyValuePairs)/2:]...)
		page.keyValuePairs = page.keyValuePairs[:len(pageKeyValuePairs)/2]

		parentPage.insertChildAt(index+1, siblingPage)
		parentPage.insertAt(index, siblingPage.keyValuePairs[0])
	} else {
		parentKey := page.keyValuePairs[len(page.AllKeyValuePairs())/2]

		siblingPage.keyValuePairs = append(siblingPage.keyValuePairs, page.keyValuePairs[0:len(page.AllKeyValuePairs())/2]...)
		page.keyValuePairs = page.keyValuePairs[len(page.AllKeyValuePairs())/2:]

		siblingPage.childPageIds = append(siblingPage.childPageIds, page.childPageIds[:len(page.AllKeyValuePairs())/2]...)
		page.childPageIds = page.childPageIds[len(page.AllKeyValuePairs())/2:]

		parentPage.insertChildAt(index, siblingPage)
		parentPage.insertAt(index, parentKey)
	}
	return nil
}

func (page *Page) AllKeyValuePairs() []KeyValuePair {
	return page.keyValuePairs
}

func (page *Page) MarkDirty() {
	page.dirty = true
}

func (page *Page) ClearDirty() {
	page.dirty = false
}

func (page *Page) IsDirty() bool {
	return page.dirty
}
