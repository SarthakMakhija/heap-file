package heap_file

import (
	"fmt"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/page"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
)

type HeapFile struct {
	bufferPool      *BufferPool
	freePageList    *FreePageList
	currentPage     *page.SlottedPage
	pageSize        int
	tupleDescriptor tuple.TupleDescriptor
}

func NewHeapFile(bufferPool *BufferPool, freePageList *FreePageList, options DbOptions) *HeapFile {
	isAvailable, pageId := freePageList.allocateAndUpdate(1)
	if !isAvailable {
		panic("No free page available for allocation during creation of NewHeapFile")
	}
	return &HeapFile{
		bufferPool:      bufferPool,
		freePageList:    freePageList,
		currentPage:     page.NewSlottedPage(pageId, options.PageSize, options.TupleDescriptor),
		pageSize:        options.PageSize,
		tupleDescriptor: options.TupleDescriptor,
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

func (heapFile *HeapFile) GetBy(tupleId tuple.TupleId) *tuple.Tuple {
	slottedPage, err := heapFile.bufferPool.Read(tupleId.PageId)
	if err != nil {
		return nil
	}
	return slottedPage.GetAt(tupleId.SlotNo)
}

func (heapFile *HeapFile) isCurrentSlottedPageLargeEnoughToHold(marshalledTuple tuple.MarshalledTuple) bool {
	return uint16(marshalledTuple.Size()) <= heapFile.currentPage.SizeAvailable()
}

func (heapFile *HeapFile) newCurrentPage() *page.SlottedPage {
	isAvailable, newPageId := heapFile.freePageList.allocateAndUpdate(1)
	if isAvailable {
		return page.NewSlottedPage(newPageId, heapFile.pageSize, heapFile.tupleDescriptor)
	} else {
		pageId, err := heapFile.bufferPool.Allocate(1) //might change
		if err != nil {
			panic(fmt.Sprintf("Error while allocating a page %v", err))
		}
		return page.NewSlottedPage(uint32(pageId), heapFile.pageSize, heapFile.tupleDescriptor)
	}
}
