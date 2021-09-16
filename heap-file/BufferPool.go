package heap_file

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/page"
	"os"
)

type BufferPool struct {
	file      *os.File
	pageSize  int
	pageCount int
}

func NewBufferPool(file *os.File, pageSize int) *BufferPool {
	bufferPool := &BufferPool{
		file: file,
	}
	bufferPool.pageSize = pageSize
	bufferPool.pageCount = bufferPool.numberOfPages()

	return bufferPool
}

func (bufferPool *BufferPool) Allocate(pages int) (int, error) {
	nextPageId := bufferPool.pageCount
	targetSize := bufferPool.fileSize() + (pages * bufferPool.pageSize)
	if err := bufferPool.file.Truncate(int64(targetSize)); err != nil {
		return 0, err
	}
	bufferPool.pageCount = bufferPool.numberOfPages()
	return nextPageId, nil
}

func (bufferPool BufferPool) Read(pageId uint32) (*page.SlottedPage, error) {
	buffer := make([]byte, bufferPool.pageSize)
	_, err := bufferPool.file.ReadAt(buffer, bufferPool.offsetOf(pageId))
	if err != nil {
		return nil, err
	}
	return page.NewReadonlySlottedPageFrom(buffer), nil
}

func (bufferPool *BufferPool) Write(page *page.SlottedPage) error {
	_, err := bufferPool.file.WriteAt(page.Buffer(), bufferPool.offsetOf(page.PageId()))
	if err != nil {
		return err
	}
	return nil
}

func (bufferPool BufferPool) ContainsZeroPages() bool {
	return bufferPool.pageCount == 0
}

func (bufferPool *BufferPool) Close() error {
	return bufferPool.file.Close()
}

func (bufferPool BufferPool) offsetOf(pageId uint32) int64 {
	return int64(uint32(bufferPool.pageSize) * pageId)
}

func (bufferPool BufferPool) numberOfPages() int {
	return bufferPool.fileSize() / bufferPool.pageSize
}

func (bufferPool BufferPool) fileSize() int {
	stat, err := bufferPool.file.Stat()
	if err != nil {
		return 0 //Handle later
	}
	return int(stat.Size())
}
