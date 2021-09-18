package heap_file

import (
	"fmt"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/page"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
)

type HeapFile struct {
	bufferPool   *BufferPool
	freePageList *FreePageList
	currentPage  *page.SlottedPage
}

func NewHeapFile(bufferPool *BufferPool, freePageList *FreePageList) *HeapFile {
	isAvailable, pageId := freePageList.allocateAndUpdate(1)
	if !isAvailable {
		panic("No free page available for allocation during creation of NewHeapFile")
	}
	slottedPage, err := bufferPool.Read(pageId)
	if err != nil {
		panic(fmt.Sprintf("Error while reading a page with page id %v", pageId))
	}
	return &HeapFile{
		bufferPool:   bufferPool,
		freePageList: freePageList,
		currentPage:  slottedPage,
	}
}

func (heapFile *HeapFile) Put(tuple *tuple.Tuple) tuple.TupleId {
	return heapFile.currentPage.Put(tuple)
}

func (heapFile *HeapFile) GetAt(slotNo int) *tuple.Tuple {
	return heapFile.currentPage.GetAt(slotNo)
}
