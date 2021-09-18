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
	pageSize     int
}

func NewHeapFile(bufferPool *BufferPool, freePageList *FreePageList, options Options) *HeapFile {
	isAvailable, pageId := freePageList.allocateAndUpdate(1)
	if !isAvailable {
		panic("No free page available for allocation during creation of NewHeapFile")
	}
	return &HeapFile{
		bufferPool:   bufferPool,
		freePageList: freePageList,
		currentPage:  page.NewSlottedPage(pageId, options.PageSize),
		pageSize:     options.PageSize,
	}
}

func (heapFile *HeapFile) Put(tuple *tuple.Tuple) (tuple.TupleId, error) {
	marshalledTuple := tuple.MarshalBinary()
	if !heapFile.isCurrentSlottedPageLargeEnoughToHold(marshalledTuple) {
		heapFile.currentPage = heapFile.newCurrentPage()
	}
	tupleId := heapFile.currentPage.Put(marshalledTuple)
	err := heapFile.bufferPool.Write(heapFile.currentPage)
	return tupleId, err
}

func (heapFile *HeapFile) GetAt(slotNo int) *tuple.Tuple {
	return heapFile.currentPage.GetAt(slotNo)
}

func (heapFile *HeapFile) isCurrentSlottedPageLargeEnoughToHold(marshalledTuple tuple.MarshalledTuple) bool {
	return uint16(marshalledTuple.Size()) <= heapFile.currentPage.SizeAvailable()
}

func (heapFile *HeapFile) newCurrentPage() *page.SlottedPage {
	isAvailable, newPageId := heapFile.freePageList.allocateAndUpdate(1)
	if isAvailable {
		return page.NewSlottedPage(newPageId, heapFile.pageSize)
	} else {
		pageId, err := heapFile.bufferPool.Allocate(1) //might change
		if err != nil {
			panic(fmt.Sprintf("Error while allocating a page %v", err))
		}
		return page.NewSlottedPage(uint32(pageId), heapFile.pageSize)
	}
}
